// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

const (
	baseTable  = "processbaselines"
	countStmt  = "SELECT COUNT(*) FROM processbaselines"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM processbaselines WHERE Id = $1)"

	getStmt     = "SELECT serialized FROM processbaselines WHERE Id = $1"
	deleteStmt  = "DELETE FROM processbaselines WHERE Id = $1"
	walkStmt    = "SELECT serialized FROM processbaselines"
	getIDsStmt  = "SELECT Id FROM processbaselines"
	getManyStmt = "SELECT serialized FROM processbaselines WHERE Id = ANY($1::text[])"

	deleteManyStmt = "DELETE FROM processbaselines WHERE Id = ANY($1::text[])"
)

func init() {
	globaldb.RegisterTable(baseTable, "ProcessBaseline")
}

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (*storage.ProcessBaseline, bool, error)
	Upsert(ctx context.Context, obj *storage.ProcessBaseline) error
	UpsertMany(ctx context.Context, objs []*storage.ProcessBaseline) error
	Delete(ctx context.Context, id string) error
	GetIDs(ctx context.Context) ([]string, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.ProcessBaseline, []int, error)
	DeleteMany(ctx context.Context, ids []string) error

	Walk(ctx context.Context, fn func(obj *storage.ProcessBaseline) error) error

	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

func createTableProcessbaselines(ctx context.Context, db *pgxpool.Pool) {
	table := `
create table if not exists processbaselines (
    Id varchar,
    Key_DeploymentId varchar,
    Key_ContainerName varchar,
    Key_ClusterId varchar,
    Key_Namespace varchar,
    Created timestamp,
    UserLockedTimestamp timestamp,
    StackRoxLockedTimestamp timestamp,
    LastUpdate timestamp,
    serialized bytea,
    PRIMARY KEY(Id)
)
`

	_, err := db.Exec(ctx, table)
	if err != nil {
		panic("error creating table: " + table)
	}

	indexes := []string{}
	for _, index := range indexes {
		if _, err := db.Exec(ctx, index); err != nil {
			panic(err)
		}
	}

	createTableProcessbaselinesElements(ctx, db)
	createTableProcessbaselinesElementGraveyard(ctx, db)
}

func createTableProcessbaselinesElements(ctx context.Context, db *pgxpool.Pool) {
	table := `
create table if not exists processbaselines_Elements (
    processbaselines_Id varchar,
    idx integer,
    Element_ProcessName varchar,
    Auto bool,
    PRIMARY KEY(processbaselines_Id, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (processbaselines_Id) REFERENCES processbaselines(Id) ON DELETE CASCADE
)
`

	_, err := db.Exec(ctx, table)
	if err != nil {
		panic("error creating table: " + table)
	}

	indexes := []string{

		"create index if not exists processbaselinesElements_idx on processbaselines_Elements using btree(idx)",
	}
	for _, index := range indexes {
		if _, err := db.Exec(ctx, index); err != nil {
			panic(err)
		}
	}

}

func createTableProcessbaselinesElementGraveyard(ctx context.Context, db *pgxpool.Pool) {
	table := `
create table if not exists processbaselines_ElementGraveyard (
    processbaselines_Id varchar,
    idx integer,
    Element_ProcessName varchar,
    Auto bool,
    PRIMARY KEY(processbaselines_Id, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (processbaselines_Id) REFERENCES processbaselines(Id) ON DELETE CASCADE
)
`

	_, err := db.Exec(ctx, table)
	if err != nil {
		panic("error creating table: " + table)
	}

	indexes := []string{

		"create index if not exists processbaselinesElementGraveyard_idx on processbaselines_ElementGraveyard using btree(idx)",
	}
	for _, index := range indexes {
		if _, err := db.Exec(ctx, index); err != nil {
			panic(err)
		}
	}

}

func insertIntoProcessbaselines(ctx context.Context, tx pgx.Tx, obj *storage.ProcessBaseline) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		obj.GetId(),
		obj.GetKey().GetDeploymentId(),
		obj.GetKey().GetContainerName(),
		obj.GetKey().GetClusterId(),
		obj.GetKey().GetNamespace(),
		pgutils.NilOrStringTimestamp(obj.GetCreated()),
		pgutils.NilOrStringTimestamp(obj.GetUserLockedTimestamp()),
		pgutils.NilOrStringTimestamp(obj.GetStackRoxLockedTimestamp()),
		pgutils.NilOrStringTimestamp(obj.GetLastUpdate()),
		serialized,
	}

	finalStr := "INSERT INTO processbaselines (Id, Key_DeploymentId, Key_ContainerName, Key_ClusterId, Key_Namespace, Created, UserLockedTimestamp, StackRoxLockedTimestamp, LastUpdate, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Key_DeploymentId = EXCLUDED.Key_DeploymentId, Key_ContainerName = EXCLUDED.Key_ContainerName, Key_ClusterId = EXCLUDED.Key_ClusterId, Key_Namespace = EXCLUDED.Key_Namespace, Created = EXCLUDED.Created, UserLockedTimestamp = EXCLUDED.UserLockedTimestamp, StackRoxLockedTimestamp = EXCLUDED.StackRoxLockedTimestamp, LastUpdate = EXCLUDED.LastUpdate, serialized = EXCLUDED.serialized"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetElements() {
		if err := insertIntoProcessbaselinesElements(ctx, tx, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from processbaselines_Elements where processbaselines_Id = $1 AND idx >= $2"
	_, err = tx.Exec(ctx, query, obj.GetId(), len(obj.GetElements()))
	if err != nil {
		return err
	}
	for childIdx, child := range obj.GetElementGraveyard() {
		if err := insertIntoProcessbaselinesElementGraveyard(ctx, tx, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from processbaselines_ElementGraveyard where processbaselines_Id = $1 AND idx >= $2"
	_, err = tx.Exec(ctx, query, obj.GetId(), len(obj.GetElementGraveyard()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoProcessbaselinesElements(ctx context.Context, tx pgx.Tx, obj *storage.BaselineElement, processbaselines_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		processbaselines_Id,
		idx,
		obj.GetElement().GetProcessName(),
		obj.GetAuto(),
	}

	finalStr := "INSERT INTO processbaselines_Elements (processbaselines_Id, idx, Element_ProcessName, Auto) VALUES($1, $2, $3, $4) ON CONFLICT(processbaselines_Id, idx) DO UPDATE SET processbaselines_Id = EXCLUDED.processbaselines_Id, idx = EXCLUDED.idx, Element_ProcessName = EXCLUDED.Element_ProcessName, Auto = EXCLUDED.Auto"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func insertIntoProcessbaselinesElementGraveyard(ctx context.Context, tx pgx.Tx, obj *storage.BaselineElement, processbaselines_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		processbaselines_Id,
		idx,
		obj.GetElement().GetProcessName(),
		obj.GetAuto(),
	}

	finalStr := "INSERT INTO processbaselines_ElementGraveyard (processbaselines_Id, idx, Element_ProcessName, Auto) VALUES($1, $2, $3, $4) ON CONFLICT(processbaselines_Id, idx) DO UPDATE SET processbaselines_Id = EXCLUDED.processbaselines_Id, idx = EXCLUDED.idx, Element_ProcessName = EXCLUDED.Element_ProcessName, Auto = EXCLUDED.Auto"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

// New returns a new Store instance using the provided sql instance.
func New(ctx context.Context, db *pgxpool.Pool) Store {
	createTableProcessbaselines(ctx, db)

	return &storeImpl{
		db: db,
	}
}

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.ProcessBaseline) error {
	conn, release := s.acquireConn(ctx, ops.Get, "ProcessBaseline")
	defer release()

	for _, obj := range objs {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}

		if err := insertIntoProcessbaselines(ctx, tx, obj); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.ProcessBaseline) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "ProcessBaseline")

	return s.upsert(ctx, obj)
}

func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.ProcessBaseline) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "ProcessBaseline")

	return s.upsert(ctx, objs...)
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "ProcessBaseline")

	row := s.db.QueryRow(ctx, countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "ProcessBaseline")

	row := s.db.QueryRow(ctx, existsStmt, id)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, pgutils.ErrNilIfNoRows(err)
	}
	return exists, nil
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string) (*storage.ProcessBaseline, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "ProcessBaseline")

	conn, release := s.acquireConn(ctx, ops.Get, "ProcessBaseline")
	defer release()

	row := conn.QueryRow(ctx, getStmt, id)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.ProcessBaseline
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func (s *storeImpl) acquireConn(ctx context.Context, op ops.Op, typ string) (*pgxpool.Conn, func()) {
	defer metrics.SetAcquireDBConnDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		panic(err)
	}
	return conn, conn.Release
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(ctx context.Context, id string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "ProcessBaseline")

	conn, release := s.acquireConn(ctx, ops.Remove, "ProcessBaseline")
	defer release()

	if _, err := conn.Exec(ctx, deleteStmt, id); err != nil {
		return err
	}
	return nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs(ctx context.Context) ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "storage.ProcessBaselineIDs")

	rows, err := s.db.Query(ctx, getIDsStmt)
	if err != nil {
		return nil, pgutils.ErrNilIfNoRows(err)
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
func (s *storeImpl) GetMany(ctx context.Context, ids []string) ([]*storage.ProcessBaseline, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "ProcessBaseline")

	conn, release := s.acquireConn(ctx, ops.GetMany, "ProcessBaseline")
	defer release()

	rows, err := conn.Query(ctx, getManyStmt, ids)
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
	resultsByID := make(map[string]*storage.ProcessBaseline)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		msg := &storage.ProcessBaseline{}
		if err := proto.Unmarshal(data, msg); err != nil {
			return nil, nil, err
		}
		resultsByID[msg.GetId()] = msg
	}
	missingIndices := make([]int, 0, len(ids)-len(resultsByID))
	// It is important that the elems are populated in the same order as the input ids
	// slice, since some calling code relies on that to maintain order.
	elems := make([]*storage.ProcessBaseline, 0, len(resultsByID))
	for i, id := range ids {
		if result, ok := resultsByID[id]; !ok {
			missingIndices = append(missingIndices, i)
		} else {
			elems = append(elems, result)
		}
	}
	return elems, missingIndices, nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ctx context.Context, ids []string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "ProcessBaseline")

	conn, release := s.acquireConn(ctx, ops.RemoveMany, "ProcessBaseline")
	defer release()
	if _, err := conn.Exec(ctx, deleteManyStmt, ids); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.ProcessBaseline) error) error {
	rows, err := s.db.Query(ctx, walkStmt)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.ProcessBaseline
		if err := proto.Unmarshal(data, &msg); err != nil {
			return err
		}
		if err := fn(&msg); err != nil {
			return err
		}
	}
	return nil
}

//// Used for testing

func dropTableProcessbaselines(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS processbaselines CASCADE")
	dropTableProcessbaselinesElements(ctx, db)
	dropTableProcessbaselinesElementGraveyard(ctx, db)

}

func dropTableProcessbaselinesElements(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS processbaselines_Elements CASCADE")

}

func dropTableProcessbaselinesElementGraveyard(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS processbaselines_ElementGraveyard CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTableProcessbaselines(ctx, db)
}

//// Stubs for satisfying legacy interfaces

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(ctx context.Context, keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex(ctx context.Context) ([]string, error) {
	return nil, nil
}
