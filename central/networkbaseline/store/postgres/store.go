// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"bytes"
	"context"
	"reflect"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
)

const (
	countStmt  = "SELECT COUNT(*) FROM NetworkBaseline"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM NetworkBaseline WHERE DeploymentId = $1)"

	getStmt        = "SELECT serialized FROM NetworkBaseline WHERE DeploymentId = $1"
	deleteStmt     = "DELETE FROM NetworkBaseline WHERE DeploymentId = $1"
	walkStmt       = "SELECT serialized FROM NetworkBaseline"
	getIDsStmt     = "SELECT DeploymentId FROM NetworkBaseline"
	getManyStmt    = "SELECT serialized FROM NetworkBaseline WHERE DeploymentId = ANY($1::text[])"
	deleteManyStmt = "DELETE FROM NetworkBaseline WHERE DeploymentId = ANY($1::text[])"
)

var (
	log = logging.LoggerForModule()

	table = "NetworkBaseline"

	marshaler = &jsonpb.Marshaler{EnumsAsInts: true, EmitDefaults: true}
)

type Store interface {
	Count() (int, error)
	Exists(deploymentId string) (bool, error)
	Get(deploymentId string) (*storage.NetworkBaseline, bool, error)
	Upsert(obj *storage.NetworkBaseline) error
	UpsertMany(objs []*storage.NetworkBaseline) error
	Delete(deploymentId string) error
	GetIDs() ([]string, error)
	GetMany(ids []string) ([]*storage.NetworkBaseline, []int, error)
	DeleteMany(ids []string) error

	Walk(fn func(obj *storage.NetworkBaseline) error) error
	AckKeysIndexed(keys ...string) error
	GetKeysToIndex() ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

const (
	batchInsertTemplate = "<no value>"
)

// New returns a new Store instance using the provided sql instance.
func New(db *pgxpool.Pool) Store {
	globaldb.RegisterTable(table, "NetworkBaseline")

	for _, table := range []string{
		"create table if not exists NetworkBaseline(serialized jsonb not null, DeploymentId varchar, PRIMARY KEY (DeploymentId));",
		"create table if not exists NetworkBaseline_Peers(parent_DeploymentId varchar not null, idx integer not null, Entity_Info_Desc_ExternalSource_Default bool, PRIMARY KEY (parent_DeploymentId, idx), CONSTRAINT fk_parent_table FOREIGN KEY (parent_DeploymentId) REFERENCES NetworkBaseline(DeploymentId) ON DELETE CASCADE);",
		"create table if not exists NetworkBaseline_ForbiddenPeers(parent_DeploymentId varchar not null, idx integer not null, Entity_Info_Desc_ExternalSource_Default bool, PRIMARY KEY (parent_DeploymentId, idx), CONSTRAINT fk_parent_table FOREIGN KEY (parent_DeploymentId) REFERENCES NetworkBaseline(DeploymentId) ON DELETE CASCADE);",
	} {
		_, err := db.Exec(context.Background(), table)
		if err != nil {
			panic("error creating table: " + table)
		}
	}

	//
	return &storeImpl{
		db: db,
	}
	//
}

// Count returns the number of objects in the store
func (s *storeImpl) Count() (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "NetworkBaseline")

	row := s.db.QueryRow(context.Background(), countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(deploymentId string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "NetworkBaseline")

	row := s.db.QueryRow(context.Background(), existsStmt, deploymentId)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, nilNoRows(err)
	}
	return exists, nil
}

func nilNoRows(err error) error {
	if err == pgx.ErrNoRows {
		return nil
	}
	return err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(deploymentId string) (*storage.NetworkBaseline, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "NetworkBaseline")

	conn, release := s.acquireConn(ops.Get, "NetworkBaseline")
	defer release()

	row := conn.QueryRow(context.Background(), getStmt, deploymentId)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, nilNoRows(err)
	}

	var msg storage.NetworkBaseline
	buf := bytes.NewBuffer(data)
	defer metrics.SetJSONPBOperationDurationTime(time.Now(), "Unmarshal", "NetworkBaseline")
	if err := jsonpb.Unmarshal(buf, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func convertEnumSliceToIntArray(i interface{}) []int32 {
	enumSlice := reflect.ValueOf(i)
	enumSliceLen := enumSlice.Len()
	resultSlice := make([]int32, 0, enumSliceLen)
	for i := 0; i < enumSlice.Len(); i++ {
		resultSlice = append(resultSlice, int32(enumSlice.Index(i).Int()))
	}
	return resultSlice
}

func nilOrStringTimestamp(t *types.Timestamp) *string {
	if t == nil {
		return nil
	}
	s := t.String()
	return &s
}

// Upsert inserts the object into the DB
func (s *storeImpl) Upsert(obj0 *storage.NetworkBaseline) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Add, "NetworkBaseline")

	t := time.Now()
	serialized, err := marshaler.MarshalToString(obj0)
	if err != nil {
		return err
	}
	metrics.SetJSONPBOperationDurationTime(t, "Marshal", "NetworkBaseline")
	conn, release := s.acquireConn(ops.Add, "NetworkBaseline")
	defer release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	doRollback := true
	defer func() {
		if doRollback {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				log.Errorf("error rolling backing: %v", err)
			}
		}
	}()

	localQuery := "insert into NetworkBaseline(serialized, DeploymentId) values($1, $2) on conflict(DeploymentId) do update set serialized = EXCLUDED.serialized, DeploymentId = EXCLUDED.DeploymentId"
	_, err = tx.Exec(context.Background(), localQuery, serialized, obj0.GetDeploymentId())
	if err != nil {
		return err
	}
	for idx1, obj1 := range obj0.GetPeers() {
		localQuery := "insert into NetworkBaseline_Peers(parent_DeploymentId, idx, Entity_Info_Desc_ExternalSource_Default) values($1, $2, $3) on conflict(parent_DeploymentId, idx) do update set parent_DeploymentId = EXCLUDED.parent_DeploymentId, idx = EXCLUDED.idx, Entity_Info_Desc_ExternalSource_Default = EXCLUDED.Entity_Info_Desc_ExternalSource_Default"
		_, err := tx.Exec(context.Background(), localQuery, obj0.GetDeploymentId(), idx1, obj1.GetEntity().GetInfo().GetExternalSource().GetDefault())
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(context.Background(), "delete from NetworkBaseline_Peers where parent_DeploymentId = $1 and idx >= $2", obj0.GetDeploymentId(), len(obj0.GetPeers()))
	if err != nil {
		return err
	}
	for idx1, obj1 := range obj0.GetForbiddenPeers() {
		localQuery := "insert into NetworkBaseline_ForbiddenPeers(parent_DeploymentId, idx, Entity_Info_Desc_ExternalSource_Default) values($1, $2, $3) on conflict(parent_DeploymentId, idx) do update set parent_DeploymentId = EXCLUDED.parent_DeploymentId, idx = EXCLUDED.idx, Entity_Info_Desc_ExternalSource_Default = EXCLUDED.Entity_Info_Desc_ExternalSource_Default"
		_, err := tx.Exec(context.Background(), localQuery, obj0.GetDeploymentId(), idx1, obj1.GetEntity().GetInfo().GetExternalSource().GetDefault())
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(context.Background(), "delete from NetworkBaseline_ForbiddenPeers where parent_DeploymentId = $1 and idx >= $2", obj0.GetDeploymentId(), len(obj0.GetForbiddenPeers()))
	if err != nil {
		return err
	}

	doRollback = false
	return tx.Commit(context.Background())
}

func (s *storeImpl) acquireConn(op ops.Op, typ string) (*pgxpool.Conn, func()) {
	defer metrics.SetAcquireDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	return conn, conn.Release
}

// UpsertMany batches objects into the DB
func (s *storeImpl) UpsertMany(objs []*storage.NetworkBaseline) error {
	if len(objs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.AddMany, "NetworkBaseline")
	for _, obj0 := range objs {
		t := time.Now()
		serialized, err := marshaler.MarshalToString(obj0)
		if err != nil {
			return err
		}
		metrics.SetJSONPBOperationDurationTime(t, "Marshal", "NetworkBaseline")
		localQuery := "insert into NetworkBaseline(serialized, DeploymentId) values($1, $2) on conflict(DeploymentId) do update set serialized = EXCLUDED.serialized, DeploymentId = EXCLUDED.DeploymentId"
		batch.Queue(localQuery, serialized, obj0.GetDeploymentId())
		for idx1, obj1 := range obj0.GetPeers() {
			localQuery := "insert into NetworkBaseline_Peers(parent_DeploymentId, idx, Entity_Info_Desc_ExternalSource_Default) values($1, $2, $3) on conflict(parent_DeploymentId, idx) do update set parent_DeploymentId = EXCLUDED.parent_DeploymentId, idx = EXCLUDED.idx, Entity_Info_Desc_ExternalSource_Default = EXCLUDED.Entity_Info_Desc_ExternalSource_Default"
			batch.Queue(localQuery, obj0.GetDeploymentId(), idx1, obj1.GetEntity().GetInfo().GetExternalSource().GetDefault())
		}
		batch.Queue("delete from NetworkBaseline_Peers where parent_DeploymentId = $1 and idx >= $2", obj0.GetDeploymentId(), len(obj0.GetPeers()))
		for idx1, obj1 := range obj0.GetForbiddenPeers() {
			localQuery := "insert into NetworkBaseline_ForbiddenPeers(parent_DeploymentId, idx, Entity_Info_Desc_ExternalSource_Default) values($1, $2, $3) on conflict(parent_DeploymentId, idx) do update set parent_DeploymentId = EXCLUDED.parent_DeploymentId, idx = EXCLUDED.idx, Entity_Info_Desc_ExternalSource_Default = EXCLUDED.Entity_Info_Desc_ExternalSource_Default"
			batch.Queue(localQuery, obj0.GetDeploymentId(), idx1, obj1.GetEntity().GetInfo().GetExternalSource().GetDefault())
		}
		batch.Queue("delete from NetworkBaseline_ForbiddenPeers where parent_DeploymentId = $1 and idx >= $2", obj0.GetDeploymentId(), len(obj0.GetForbiddenPeers()))

	}

	conn, release := s.acquireConn(ops.AddMany, "NetworkBaseline")
	defer release()

	results := conn.SendBatch(context.Background(), batch)
	if err := results.Close(); err != nil {
		return err
	}
	return nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(deploymentId string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "NetworkBaseline")

	conn, release := s.acquireConn(ops.Remove, "NetworkBaseline")
	defer release()

	if _, err := conn.Exec(context.Background(), deleteStmt, deploymentId); err != nil {
		return err
	}
	return nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs() ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "NetworkBaselineIDs")

	rows, err := s.db.Query(context.Background(), getIDsStmt)
	if err != nil {
		return nil, nilNoRows(err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (s *storeImpl) GetMany(ids []string) ([]*storage.NetworkBaseline, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "NetworkBaseline")

	conn, release := s.acquireConn(ops.GetMany, "NetworkBaseline")
	defer release()

	rows, err := conn.Query(context.Background(), getManyStmt, ids)
	if err != nil {
		if err == pgx.ErrNoRows {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	defer rows.Close()
	elems := make([]*storage.NetworkBaseline, 0, len(ids))
	foundSet := make(map[string]struct{})
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		var msg storage.NetworkBaseline
		buf := bytes.NewBuffer(data)
		t := time.Now()
		if err := jsonpb.Unmarshal(buf, &msg); err != nil {
			return nil, nil, err
		}
		metrics.SetJSONPBOperationDurationTime(t, "Unmarshal", "NetworkBaseline")
		foundSet[msg.GetDeploymentId()] = struct{}{}
		elems = append(elems, &msg)
	}
	missingIndices := make([]int, 0, len(ids)-len(foundSet))
	for i, id := range ids {
		if _, ok := foundSet[id]; !ok {
			missingIndices = append(missingIndices, i)
		}
	}
	return elems, missingIndices, nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ids []string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "NetworkBaseline")

	conn, release := s.acquireConn(ops.RemoveMany, "NetworkBaseline")
	defer release()
	if _, err := conn.Exec(context.Background(), deleteManyStmt, ids); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(fn func(obj *storage.NetworkBaseline) error) error {
	rows, err := s.db.Query(context.Background(), walkStmt)
	if err != nil {
		return nilNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.NetworkBaseline
		buf := bytes.NewReader(data)
		if err := jsonpb.Unmarshal(buf, &msg); err != nil {
			return err
		}
		return fn(&msg)
	}
	return nil
}

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex() ([]string, error) {
	return nil, nil
}
