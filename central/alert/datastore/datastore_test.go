package datastore

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	indexMocks "github.com/stackrox/rox/central/alert/datastore/internal/index/mocks"
	searchMocks "github.com/stackrox/rox/central/alert/datastore/internal/search/mocks"
	storeMocks "github.com/stackrox/rox/central/alert/datastore/internal/store/mocks"
	_ "github.com/stackrox/rox/central/alert/mappings"
	"github.com/stackrox/rox/central/alerttest"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/alert/convert"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stretchr/testify/suite"
)

var (
	errFake = errors.New("fake error")
)

func TestAlertDataStore(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(alertDataStoreTestSuite))
}

type alertDataStoreTestSuite struct {
	suite.Suite

	hasReadCtx  context.Context
	hasWriteCtx context.Context

	dataStore DataStore
	storage   *storeMocks.MockStore
	indexer   *indexMocks.MockIndexer
	searcher  *searchMocks.MockSearcher

	mockCtrl *gomock.Controller
}

func (s *alertDataStoreTestSuite) SetupTest() {
	s.hasReadCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
			sac.ResourceScopeKeys(resources.Alert)))
	s.hasWriteCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Alert)))

	s.mockCtrl = gomock.NewController(s.T())
	s.storage = storeMocks.NewMockStore(s.mockCtrl)
	s.storage.EXPECT().GetKeysToIndex(gomock.Any()).Return(nil, nil)

	s.indexer = indexMocks.NewMockIndexer(s.mockCtrl)
	s.indexer.EXPECT().NeedsInitialIndexing().Return(false, nil)

	s.searcher = searchMocks.NewMockSearcher(s.mockCtrl)

	var err error
	s.dataStore, err = New(s.storage, s.indexer, s.searcher)
	s.Require().NoError(err)
}

func (s *alertDataStoreTestSuite) TestSearchAlerts() {
	s.searcher.EXPECT().SearchAlerts(s.hasReadCtx, &v1.Query{}).Return([]*v1.SearchResult{{Id: alerttest.FakeAlertID}}, errFake)

	result, err := s.dataStore.SearchAlerts(s.hasReadCtx, &v1.Query{})

	s.Equal(errFake, err)
	s.Equal([]*v1.SearchResult{{Id: alerttest.FakeAlertID}}, result)
}

func (s *alertDataStoreTestSuite) TestSearchRawAlerts() {
	s.searcher.EXPECT().SearchRawAlerts(s.hasReadCtx, &v1.Query{}).Return([]*storage.Alert{{Id: alerttest.FakeAlertID}}, errFake)

	result, err := s.dataStore.SearchRawAlerts(s.hasReadCtx, &v1.Query{})

	s.Equal(errFake, err)
	s.Equal([]*storage.Alert{{Id: alerttest.FakeAlertID}}, result)
}

func (s *alertDataStoreTestSuite) TestSearchListAlerts() {
	s.searcher.EXPECT().SearchListAlerts(s.hasReadCtx, &v1.Query{}).Return(alerttest.NewFakeListAlertSlice(), errFake)

	result, err := s.dataStore.SearchListAlerts(s.hasReadCtx, &v1.Query{})

	s.Equal(errFake, err)
	s.Equal(alerttest.NewFakeListAlertSlice(), result)
}

func (s *alertDataStoreTestSuite) TestCountAlerts_Success() {
	expectedQ := search.NewQueryBuilder().AddStrings(search.ViolationState, storage.ViolationState_ACTIVE.String()).ProtoQuery()
	s.searcher.EXPECT().Count(s.hasReadCtx, expectedQ).Return(1, nil)

	result, err := s.dataStore.CountAlerts(s.hasReadCtx)

	s.NoError(err)
	s.Equal(1, result)
}

func (s *alertDataStoreTestSuite) TestCountAlerts_Error() {
	expectedQ := search.NewQueryBuilder().AddStrings(search.ViolationState, storage.ViolationState_ACTIVE.String()).ProtoQuery()
	s.searcher.EXPECT().Count(s.hasReadCtx, expectedQ).Return(0, errFake)

	_, err := s.dataStore.CountAlerts(s.hasReadCtx)

	s.Equal(errFake, err)
}

func (s *alertDataStoreTestSuite) TestAddAlert() {
	fakeAlert := alerttest.NewFakeAlert()
	s.storage.EXPECT().Upsert(gomock.Any(), fakeAlert).Return(nil)
	s.indexer.EXPECT().AddListAlert(fillSortHelperFields(convert.AlertToListAlert(alerttest.NewFakeAlert()))).Return(errFake)

	// We don't expect AckKeysIndexed, since the error returned from the above call will prevent this.
	err := s.dataStore.UpsertAlert(s.hasWriteCtx, alerttest.NewFakeAlert())

	s.Equal(errFake, err)
}

func (s *alertDataStoreTestSuite) TestAddAlertWhenTheIndexerFails() {
	fakeAlert := alerttest.NewFakeAlert()
	s.storage.EXPECT().Upsert(gomock.Any(), fakeAlert).Return(errFake)

	// No AckKeysIndexed call due to error on upsert.
	err := s.dataStore.UpsertAlert(s.hasWriteCtx, alerttest.NewFakeAlert())

	s.Equal(errFake, err)
}

func (s *alertDataStoreTestSuite) TestMarkAlertStale() {
	fakeAlert := alerttest.NewFakeAlert()

	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlert, true, nil)
	s.storage.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
	s.indexer.EXPECT().AddListAlert(gomock.Any()).Return(nil)
	s.storage.EXPECT().AckKeysIndexed(gomock.Any(), fakeAlert.GetId()).Times(1).Return(nil)

	err := s.dataStore.MarkAlertStale(s.hasWriteCtx, alerttest.FakeAlertID)
	s.NoError(err)

	s.Equal(storage.ViolationState_RESOLVED, fakeAlert.GetState())
}

func (s *alertDataStoreTestSuite) TestMarkAlertStaleWhenStorageFails() {
	fakeAlert := alerttest.NewFakeAlert()

	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlert, false, errFake)

	err := s.dataStore.MarkAlertStale(s.hasWriteCtx, alerttest.FakeAlertID)

	s.Equal(errFake, err)
}

func (s *alertDataStoreTestSuite) TestMarkAlertStaleWhenTheAlertWasNotFoundInStorage() {
	fakeAlert := alerttest.NewFakeAlert()

	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlert, false, nil)

	err := s.dataStore.MarkAlertStale(s.hasWriteCtx, alerttest.FakeAlertID)

	s.EqualError(err, fmt.Sprintf("alert with id '%s' does not exist", alerttest.FakeAlertID))
}

func (s *alertDataStoreTestSuite) TestKeyIndexing() {
	fakeAlert := alerttest.NewFakeAlert()

	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlert, false, nil)

	err := s.dataStore.MarkAlertStale(s.hasWriteCtx, alerttest.FakeAlertID)

	s.EqualError(err, fmt.Sprintf("alert with id '%s' does not exist", alerttest.FakeAlertID))
}

func TestAlertDataStoreWithSAC(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(alertDataStoreWithSACTestSuite))
}

type alertDataStoreWithSACTestSuite struct {
	suite.Suite

	hasNoneCtx  context.Context
	hasReadCtx  context.Context
	hasWriteCtx context.Context

	dataStore DataStore
	storage   *storeMocks.MockStore
	indexer   *indexMocks.MockIndexer
	searcher  *searchMocks.MockSearcher

	mockCtrl *gomock.Controller
}

func (s *alertDataStoreWithSACTestSuite) SetupTest() {
	s.hasNoneCtx = sac.WithGlobalAccessScopeChecker(context.Background(), sac.DenyAllAccessScopeChecker())
	s.hasReadCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
			sac.ResourceScopeKeys(resources.Alert)))
	s.hasWriteCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Alert)))

	s.mockCtrl = gomock.NewController(s.T())
	s.storage = storeMocks.NewMockStore(s.mockCtrl)
	s.storage.EXPECT().GetKeysToIndex(gomock.Any()).Return(nil, nil)

	s.indexer = indexMocks.NewMockIndexer(s.mockCtrl)
	s.indexer.EXPECT().NeedsInitialIndexing().Return(false, nil)
	s.searcher = searchMocks.NewMockSearcher(s.mockCtrl)
	var err error
	s.dataStore, err = New(s.storage, s.indexer, s.searcher)
	s.NoError(err)
}

func (s *alertDataStoreWithSACTestSuite) TestAddAlertEnforced() {
	s.storage.EXPECT().Upsert(gomock.Any(), alerttest.NewFakeAlert()).Times(0)
	s.indexer.EXPECT().AddListAlert(convert.AlertToListAlert(alerttest.NewFakeAlert())).Times(0)

	err := s.dataStore.UpsertAlert(s.hasReadCtx, alerttest.NewFakeAlert())

	s.ErrorIs(err, sac.ErrResourceAccessDenied)
}

func (s *alertDataStoreWithSACTestSuite) TestMarkAlertStaleEnforced() {
	fakeAlert := alerttest.NewFakeAlert()

	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlert, true, nil)
	s.storage.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(0)
	s.indexer.EXPECT().AddListAlert(gomock.Any()).Times(0)

	err := s.dataStore.MarkAlertStale(s.hasReadCtx, alerttest.FakeAlertID)
	s.ErrorIs(err, sac.ErrResourceAccessDenied)

	s.Equal(storage.ViolationState_ACTIVE, fakeAlert.GetState())
}

func (s *alertDataStoreWithSACTestSuite) TestAddAlertTagsAllowed() {
	fakeAlertWithNoTags := alerttest.NewFakeAlert()
	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlertWithNoTags, true, nil)
	fakeAlertWithTwoTags := alerttest.NewFakeAlertWithTwoTags()
	s.storage.EXPECT().Upsert(gomock.Any(), fakeAlertWithTwoTags).Return(nil)
	s.indexer.EXPECT().AddListAlert(fillSortHelperFields(convert.AlertToListAlert(fakeAlertWithTwoTags))).Return(nil)
	s.storage.EXPECT().AckKeysIndexed(gomock.Any(), fakeAlertWithTwoTags.GetId()).Return(nil)
	expectedResponse := alerttest.NewFakeTwoTags()

	response, err := s.dataStore.AddAlertTags(s.hasWriteCtx, alerttest.FakeAlertID, alerttest.NewFakeTwoTags())
	s.NoError(err)
	s.Equal(expectedResponse, response)
}

func (s *alertDataStoreWithSACTestSuite) TestAddAlertTagsAllowed2() {
	fakeAlertWithTwoTags := alerttest.NewFakeAlertWithTwoTags()
	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlertWithTwoTags, true, nil)
	fakeAlertWithThreeTags := alerttest.NewFakeAlertWithThreeTags()
	s.storage.EXPECT().Upsert(gomock.Any(), fakeAlertWithThreeTags).Return(nil)
	s.indexer.EXPECT().AddListAlert(fillSortHelperFields(convert.AlertToListAlert(fakeAlertWithThreeTags))).Return(nil)
	s.storage.EXPECT().AckKeysIndexed(gomock.Any(), fakeAlertWithThreeTags.GetId()).Return(nil)
	expectedResponse := alerttest.NewFakeThreeTags()

	response, err := s.dataStore.AddAlertTags(s.hasWriteCtx, alerttest.FakeAlertID, alerttest.NewFakeTwoTagsHasOverlap())
	s.NoError(err)
	s.Equal(expectedResponse, response)
}

func (s *alertDataStoreWithSACTestSuite) TestAddAlertTagsEnforced() {
	fakeAlert := alerttest.NewFakeAlert()
	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlert, true, nil)

	_, err := s.dataStore.AddAlertTags(s.hasReadCtx, alerttest.FakeAlertID, alerttest.NewFakeTwoTags())
	s.ErrorIs(err, sac.ErrResourceAccessDenied)
}

func (s *alertDataStoreWithSACTestSuite) TestRemoveAlertTagsAllowed() {
	fakeAlertWithTwoTags := alerttest.NewFakeAlertWithTwoTags()
	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlertWithTwoTags, true, nil)
	fakeAlertWithNoTags := alerttest.NewFakeAlert()
	s.storage.EXPECT().Upsert(gomock.Any(), fakeAlertWithNoTags).Return(nil)
	s.indexer.EXPECT().AddListAlert(fillSortHelperFields(convert.AlertToListAlert(fakeAlertWithNoTags))).Return(nil)
	s.storage.EXPECT().AckKeysIndexed(gomock.Any(), fakeAlertWithNoTags.GetId()).Return(nil)

	err := s.dataStore.RemoveAlertTags(s.hasWriteCtx, alerttest.FakeAlertID, alerttest.NewFakeTwoTags())
	s.NoError(err)
}

func (s *alertDataStoreWithSACTestSuite) TestRemoveAlertTagsAllowed2() {
	fakeAlertWithThreeTags := alerttest.NewFakeAlertWithThreeTags()
	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlertWithThreeTags, true, nil)
	fakeAlertWithOneTag := alerttest.NewFakeAlertWithOneTag()
	s.storage.EXPECT().Upsert(gomock.Any(), fakeAlertWithOneTag).Return(nil)
	s.indexer.EXPECT().AddListAlert(fillSortHelperFields(convert.AlertToListAlert(fakeAlertWithOneTag))).Return(nil)
	s.storage.EXPECT().AckKeysIndexed(gomock.Any(), fakeAlertWithOneTag.GetId()).Return(nil)

	err := s.dataStore.RemoveAlertTags(s.hasWriteCtx, alerttest.FakeAlertID, alerttest.NewFakeTwoTags())
	s.NoError(err)
}

func (s *alertDataStoreWithSACTestSuite) TestRemoveAlertTagsEnforced() {
	fakeAlertWithTwoTags := alerttest.NewFakeAlertWithTwoTags()
	s.storage.EXPECT().Get(gomock.Any(), alerttest.FakeAlertID).Return(fakeAlertWithTwoTags, true, nil)

	err := s.dataStore.RemoveAlertTags(s.hasReadCtx, alerttest.FakeAlertID, alerttest.NewFakeTwoTags())
	s.ErrorIs(err, sac.ErrResourceAccessDenied)
}

func TestAlertReindexSuite(t *testing.T) {
	suite.Run(t, new(AlertReindexSuite))
}

type AlertReindexSuite struct {
	suite.Suite

	storage  *storeMocks.MockStore
	indexer  *indexMocks.MockIndexer
	searcher *searchMocks.MockSearcher

	mockCtrl *gomock.Controller
}

func (suite *AlertReindexSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.storage = storeMocks.NewMockStore(suite.mockCtrl)
	suite.indexer = indexMocks.NewMockIndexer(suite.mockCtrl)
	suite.searcher = searchMocks.NewMockSearcher(suite.mockCtrl)
}

func (suite *AlertReindexSuite) TestReconciliationFullReindex() {
	suite.indexer.EXPECT().NeedsInitialIndexing().Return(true, nil)

	alert1 := convert.AlertToListAlert(fixtures.GetAlertWithID("A"))
	alert2 := convert.AlertToListAlert(fixtures.GetAlertWithID("B"))

	listAlerts := []*storage.ListAlert{alert1, alert2}

	suite.storage.EXPECT().GetIDs(gomock.Any()).Return([]string{"A", "B"}, nil)
	suite.storage.EXPECT().GetListAlerts(gomock.Any(), []string{"A", "B"}).Return(listAlerts, nil, nil)
	suite.indexer.EXPECT().AddListAlerts(listAlerts).Return(nil)

	suite.storage.EXPECT().GetKeysToIndex(gomock.Any()).Return([]string{"D", "E"}, nil)
	suite.storage.EXPECT().AckKeysIndexed(gomock.Any(), []string{"D", "E"}).Return(nil)

	suite.indexer.EXPECT().MarkInitialIndexingComplete().Return(nil)

	_, err := New(suite.storage, suite.indexer, suite.searcher)
	suite.NoError(err)
}

func (suite *AlertReindexSuite) TestReconciliationPartialReindex() {
	suite.storage.EXPECT().GetKeysToIndex(gomock.Any()).Return([]string{"A", "B", "C"}, nil)
	suite.indexer.EXPECT().NeedsInitialIndexing().Return(false, nil)

	alert1 := convert.AlertToListAlert(fixtures.GetAlertWithID("A"))
	alert2 := convert.AlertToListAlert(fixtures.GetAlertWithID("B"))
	alert3 := convert.AlertToListAlert(fixtures.GetAlertWithID("C"))

	listAlerts := []*storage.ListAlert{alert1, alert2, alert3}

	suite.storage.EXPECT().GetListAlerts(gomock.Any(), []string{"A", "B", "C"}).Return(listAlerts, nil, nil)
	suite.indexer.EXPECT().AddListAlerts(listAlerts).Return(nil)
	suite.storage.EXPECT().AckKeysIndexed(gomock.Any(), []string{"A", "B", "C"}).Return(nil)

	_, err := New(suite.storage, suite.indexer, suite.searcher)
	suite.NoError(err)

	// Make listAlerts just A,B so C should be deleted
	listAlerts = listAlerts[:1]
	suite.storage.EXPECT().GetKeysToIndex(gomock.Any()).Return([]string{"A", "B", "C"}, nil)
	suite.indexer.EXPECT().NeedsInitialIndexing().Return(false, nil)

	suite.storage.EXPECT().GetListAlerts(gomock.Any(), []string{"A", "B", "C"}).Return(listAlerts, []int{2}, nil)
	suite.indexer.EXPECT().AddListAlerts(listAlerts).Return(nil)
	suite.indexer.EXPECT().DeleteListAlerts([]string{"C"}).Return(nil)
	suite.storage.EXPECT().AckKeysIndexed(gomock.Any(), []string{"A", "B", "C"}).Return(nil)

	_, err = New(suite.storage, suite.indexer, suite.searcher)
	suite.NoError(err)
}
