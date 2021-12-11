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
	countStmt  = "SELECT COUNT(*) FROM Pod"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM Pod WHERE Id = $1)"

	getStmt        = "SELECT serialized FROM Pod WHERE Id = $1"
	deleteStmt     = "DELETE FROM Pod WHERE Id = $1"
	walkStmt       = "SELECT serialized FROM Pod"
	getIDsStmt     = "SELECT Id FROM Pod"
	getManyStmt    = "SELECT serialized FROM Pod WHERE Id = ANY($1::text[])"
	deleteManyStmt = "DELETE FROM Pod WHERE Id = ANY($1::text[])"
)

var (
	log = logging.LoggerForModule()

	table = "Pod"

	marshaler = &jsonpb.Marshaler{EnumsAsInts: true, EmitDefaults: true}
)

type Store interface {
	Count() (int, error)
	Exists(id string) (bool, error)
	Get(id string) (*storage.Pod, bool, error)
	Upsert(obj *storage.Pod) error
	UpsertMany(objs []*storage.Pod) error
	Delete(id string) error
	GetIDs() ([]string, error)
	GetMany(ids []string) ([]*storage.Pod, []int, error)
	DeleteMany(ids []string) error

	Walk(fn func(obj *storage.Pod) error) error
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
	globaldb.RegisterTable(table, "Pod")

	for _, table := range []string{
		"create table if not exists Pod(serialized jsonb not null, Id varchar, Name varchar, DeploymentId varchar, Namespace varchar, ClusterId varchar, PRIMARY KEY (Id));",
		"create table if not exists Pod_LiveInstances(parent_Id varchar not null, idx integer not null, ImageDigest varchar, PRIMARY KEY (parent_Id, idx), CONSTRAINT fk_parent_table FOREIGN KEY (parent_Id) REFERENCES Pod(Id) ON DELETE CASCADE);",
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
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "Pod")

	row := s.db.QueryRow(context.Background(), countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(id string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "Pod")

	row := s.db.QueryRow(context.Background(), existsStmt, id)
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
func (s *storeImpl) Get(id string) (*storage.Pod, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "Pod")

	conn, release := s.acquireConn(ops.Get, "Pod")
	defer release()

	row := conn.QueryRow(context.Background(), getStmt, id)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, nilNoRows(err)
	}

	var msg storage.Pod
	buf := bytes.NewBuffer(data)
	defer metrics.SetJSONPBOperationDurationTime(time.Now(), "Unmarshal", "Pod")
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
func (s *storeImpl) Upsert(obj0 *storage.Pod) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Add, "Pod")

	t := time.Now()
	serialized, err := marshaler.MarshalToString(obj0)
	if err != nil {
		return err
	}
	metrics.SetJSONPBOperationDurationTime(t, "Marshal", "Pod")
	conn, release := s.acquireConn(ops.Add, "Pod")
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

	localQuery := "insert into Pod(serialized, Id, Name, DeploymentId, Namespace, ClusterId) values($1, $2, $3, $4, $5, $6) on conflict(Id) do update set serialized = EXCLUDED.serialized, Id = EXCLUDED.Id, Name = EXCLUDED.Name, DeploymentId = EXCLUDED.DeploymentId, Namespace = EXCLUDED.Namespace, ClusterId = EXCLUDED.ClusterId"
	_, err = tx.Exec(context.Background(), localQuery, serialized, obj0.GetId(), obj0.GetName(), obj0.GetDeploymentId(), obj0.GetNamespace(), obj0.GetClusterId())
	if err != nil {
		return err
	}
	for idx1, obj1 := range obj0.GetLiveInstances() {
		localQuery := "insert into Pod_LiveInstances(parent_Id, idx, ImageDigest) values($1, $2, $3) on conflict(parent_Id, idx) do update set parent_Id = EXCLUDED.parent_Id, idx = EXCLUDED.idx, ImageDigest = EXCLUDED.ImageDigest"
		_, err := tx.Exec(context.Background(), localQuery, obj0.GetId(), idx1, obj1.GetImageDigest())
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(context.Background(), "delete from Pod_LiveInstances where parent_Id = $1 and idx >= $2", obj0.GetId(), len(obj0.GetLiveInstances()))
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
func (s *storeImpl) UpsertMany(objs []*storage.Pod) error {
	if len(objs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.AddMany, "Pod")
	for _, obj0 := range objs {
		t := time.Now()
		serialized, err := marshaler.MarshalToString(obj0)
		if err != nil {
			return err
		}
		metrics.SetJSONPBOperationDurationTime(t, "Marshal", "Pod")
		localQuery := "insert into Pod(serialized, Id, Name, DeploymentId, Namespace, ClusterId) values($1, $2, $3, $4, $5, $6) on conflict(Id) do update set serialized = EXCLUDED.serialized, Id = EXCLUDED.Id, Name = EXCLUDED.Name, DeploymentId = EXCLUDED.DeploymentId, Namespace = EXCLUDED.Namespace, ClusterId = EXCLUDED.ClusterId"
		batch.Queue(localQuery, serialized, obj0.GetId(), obj0.GetName(), obj0.GetDeploymentId(), obj0.GetNamespace(), obj0.GetClusterId())
		for idx1, obj1 := range obj0.GetLiveInstances() {
			localQuery := "insert into Pod_LiveInstances(parent_Id, idx, ImageDigest) values($1, $2, $3) on conflict(parent_Id, idx) do update set parent_Id = EXCLUDED.parent_Id, idx = EXCLUDED.idx, ImageDigest = EXCLUDED.ImageDigest"
			batch.Queue(localQuery, obj0.GetId(), idx1, obj1.GetImageDigest())
		}
		batch.Queue("delete from Pod_LiveInstances where parent_Id = $1 and idx >= $2", obj0.GetId(), len(obj0.GetLiveInstances()))

	}

	conn, release := s.acquireConn(ops.AddMany, "Pod")
	defer release()

	results := conn.SendBatch(context.Background(), batch)
	if err := results.Close(); err != nil {
		return err
	}
	return nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(id string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "Pod")

	conn, release := s.acquireConn(ops.Remove, "Pod")
	defer release()

	if _, err := conn.Exec(context.Background(), deleteStmt, id); err != nil {
		return err
	}
	return nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs() ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "PodIDs")

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
func (s *storeImpl) GetMany(ids []string) ([]*storage.Pod, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "Pod")

	conn, release := s.acquireConn(ops.GetMany, "Pod")
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
	elems := make([]*storage.Pod, 0, len(ids))
	foundSet := make(map[string]struct{})
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		var msg storage.Pod
		buf := bytes.NewBuffer(data)
		t := time.Now()
		if err := jsonpb.Unmarshal(buf, &msg); err != nil {
			return nil, nil, err
		}
		metrics.SetJSONPBOperationDurationTime(t, "Unmarshal", "Pod")
		foundSet[msg.GetId()] = struct{}{}
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
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "Pod")

	conn, release := s.acquireConn(ops.RemoveMany, "Pod")
	defer release()
	if _, err := conn.Exec(context.Background(), deleteManyStmt, ids); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(fn func(obj *storage.Pod) error) error {
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
		var msg storage.Pod
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
