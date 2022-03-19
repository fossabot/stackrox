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
	baseTable  = "networkentity"
	countStmt  = "SELECT COUNT(*) FROM networkentity"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM networkentity WHERE Info_Id = $1)"

	getStmt     = "SELECT serialized FROM networkentity WHERE Info_Id = $1"
	deleteStmt  = "DELETE FROM networkentity WHERE Info_Id = $1"
	walkStmt    = "SELECT serialized FROM networkentity"
	getIDsStmt  = "SELECT Info_Id FROM networkentity"
	getManyStmt = "SELECT serialized FROM networkentity WHERE Info_Id = ANY($1::text[])"

	deleteManyStmt = "DELETE FROM networkentity WHERE Info_Id = ANY($1::text[])"
)

func init() {
	globaldb.RegisterTable(baseTable, "NetworkEntity")
}

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, infoId string) (bool, error)
	Get(ctx context.Context, infoId string) (*storage.NetworkEntity, bool, error)
	Upsert(ctx context.Context, obj *storage.NetworkEntity) error
	UpsertMany(ctx context.Context, objs []*storage.NetworkEntity) error
	Delete(ctx context.Context, infoId string) error
	GetIDs(ctx context.Context) ([]string, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.NetworkEntity, []int, error)
	DeleteMany(ctx context.Context, ids []string) error

	Walk(ctx context.Context, fn func(obj *storage.NetworkEntity) error) error

	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

func createTableNetworkentity(ctx context.Context, db *pgxpool.Pool) {
	table := `
create table if not exists networkentity (
    Info_Type integer,
    Info_Id varchar,
    Info_Deployment_Name varchar,
    Info_Deployment_Namespace varchar,
    Info_Deployment_Cluster varchar,
    Info_ExternalSource_Name varchar,
    Info_ExternalSource_Cidr varchar,
    Info_ExternalSource_Default bool,
    Scope_ClusterId varchar,
    serialized bytea,
    PRIMARY KEY(Info_Id)
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

	createTableNetworkentityListenPorts(ctx, db)
}

func createTableNetworkentityListenPorts(ctx context.Context, db *pgxpool.Pool) {
	table := `
create table if not exists networkentity_ListenPorts (
    networkentity_Info_Id varchar,
    idx integer,
    Port integer,
    L4Protocol integer,
    PRIMARY KEY(networkentity_Info_Id, idx),
    CONSTRAINT fk_parent_table FOREIGN KEY (networkentity_Info_Id) REFERENCES networkentity(Info_Id) ON DELETE CASCADE
)
`

	_, err := db.Exec(ctx, table)
	if err != nil {
		panic("error creating table: " + table)
	}

	indexes := []string{

		"create index if not exists networkentityListenPorts_idx on networkentity_ListenPorts using btree(idx)",
	}
	for _, index := range indexes {
		if _, err := db.Exec(ctx, index); err != nil {
			panic(err)
		}
	}

}

func insertIntoNetworkentity(ctx context.Context, tx pgx.Tx, obj *storage.NetworkEntity) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		obj.GetInfo().GetType(),
		obj.GetInfo().GetId(),
		obj.GetInfo().GetDeployment().GetName(),
		obj.GetInfo().GetDeployment().GetNamespace(),
		obj.GetInfo().GetDeployment().GetCluster(),
		obj.GetInfo().GetExternalSource().GetName(),
		obj.GetInfo().GetExternalSource().GetCidr(),
		obj.GetInfo().GetExternalSource().GetDefault(),
		obj.GetScope().GetClusterId(),
		serialized,
	}

	finalStr := "INSERT INTO networkentity (Info_Type, Info_Id, Info_Deployment_Name, Info_Deployment_Namespace, Info_Deployment_Cluster, Info_ExternalSource_Name, Info_ExternalSource_Cidr, Info_ExternalSource_Default, Scope_ClusterId, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT(Info_Id) DO UPDATE SET Info_Type = EXCLUDED.Info_Type, Info_Id = EXCLUDED.Info_Id, Info_Deployment_Name = EXCLUDED.Info_Deployment_Name, Info_Deployment_Namespace = EXCLUDED.Info_Deployment_Namespace, Info_Deployment_Cluster = EXCLUDED.Info_Deployment_Cluster, Info_ExternalSource_Name = EXCLUDED.Info_ExternalSource_Name, Info_ExternalSource_Cidr = EXCLUDED.Info_ExternalSource_Cidr, Info_ExternalSource_Default = EXCLUDED.Info_ExternalSource_Default, Scope_ClusterId = EXCLUDED.Scope_ClusterId, serialized = EXCLUDED.serialized"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetInfo().GetDeployment().GetListenPorts() {
		if err := insertIntoNetworkentityListenPorts(ctx, tx, child, obj.GetInfo().GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from networkentity_ListenPorts where networkentity_Info_Id = $1 AND idx >= $2"
	_, err = tx.Exec(ctx, query, obj.GetInfo().GetId(), len(obj.GetInfo().GetDeployment().GetListenPorts()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoNetworkentityListenPorts(ctx context.Context, tx pgx.Tx, obj *storage.NetworkEntityInfo_Deployment_ListenPort, networkentity_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		networkentity_Id,
		idx,
		obj.GetPort(),
		obj.GetL4Protocol(),
	}

	finalStr := "INSERT INTO networkentity_ListenPorts (networkentity_Info_Id, idx, Port, L4Protocol) VALUES($1, $2, $3, $4) ON CONFLICT(networkentity_Info_Id, idx) DO UPDATE SET networkentity_Info_Id = EXCLUDED.networkentity_Info_Id, idx = EXCLUDED.idx, Port = EXCLUDED.Port, L4Protocol = EXCLUDED.L4Protocol"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

// New returns a new Store instance using the provided sql instance.
func New(ctx context.Context, db *pgxpool.Pool) Store {
	createTableNetworkentity(ctx, db)

	return &storeImpl{
		db: db,
	}
}

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.NetworkEntity) error {
	conn, release := s.acquireConn(ctx, ops.Get, "NetworkEntity")
	defer release()

	for _, obj := range objs {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}

		if err := insertIntoNetworkentity(ctx, tx, obj); err != nil {
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

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.NetworkEntity) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "NetworkEntity")

	return s.upsert(ctx, obj)
}

func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.NetworkEntity) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "NetworkEntity")

	return s.upsert(ctx, objs...)
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "NetworkEntity")

	row := s.db.QueryRow(ctx, countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, infoId string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "NetworkEntity")

	row := s.db.QueryRow(ctx, existsStmt, infoId)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, pgutils.ErrNilIfNoRows(err)
	}
	return exists, nil
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, infoId string) (*storage.NetworkEntity, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "NetworkEntity")

	conn, release := s.acquireConn(ctx, ops.Get, "NetworkEntity")
	defer release()

	row := conn.QueryRow(ctx, getStmt, infoId)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.NetworkEntity
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
func (s *storeImpl) Delete(ctx context.Context, infoId string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "NetworkEntity")

	conn, release := s.acquireConn(ctx, ops.Remove, "NetworkEntity")
	defer release()

	if _, err := conn.Exec(ctx, deleteStmt, infoId); err != nil {
		return err
	}
	return nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs(ctx context.Context) ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "storage.NetworkEntityIDs")

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
func (s *storeImpl) GetMany(ctx context.Context, ids []string) ([]*storage.NetworkEntity, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "NetworkEntity")

	conn, release := s.acquireConn(ctx, ops.GetMany, "NetworkEntity")
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
	resultsByID := make(map[string]*storage.NetworkEntity)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		msg := &storage.NetworkEntity{}
		if err := proto.Unmarshal(data, msg); err != nil {
			return nil, nil, err
		}
		resultsByID[msg.GetInfo().GetId()] = msg
	}
	missingIndices := make([]int, 0, len(ids)-len(resultsByID))
	// It is important that the elems are populated in the same order as the input ids
	// slice, since some calling code relies on that to maintain order.
	elems := make([]*storage.NetworkEntity, 0, len(resultsByID))
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
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "NetworkEntity")

	conn, release := s.acquireConn(ctx, ops.RemoveMany, "NetworkEntity")
	defer release()
	if _, err := conn.Exec(ctx, deleteManyStmt, ids); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.NetworkEntity) error) error {
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
		var msg storage.NetworkEntity
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

func dropTableNetworkentity(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS networkentity CASCADE")
	dropTableNetworkentityListenPorts(ctx, db)

}

func dropTableNetworkentityListenPorts(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS networkentity_ListenPorts CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTableNetworkentity(ctx, db)
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
