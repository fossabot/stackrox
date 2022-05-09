package undostore

import (
	"context"

	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/networkpolicies/datastore/internal/undostore/bolt"
	"github.com/stackrox/rox/central/networkpolicies/datastore/internal/undostore/postgres"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	undoStoreInstance     UndoStore
	undoStoreInstanceInit sync.Once
)

// Singleton returns the singleton instance of the undo store.
func Singleton() UndoStore {
	undoStoreInstanceInit.Do(func() {
		if features.PostgresDatastore.Enabled() {
			undoStoreInstance = postgres.New(context.TODO(), globaldb.GetPostgres())
		} else {
			undoStoreInstance = bolt.New(globaldb.GetGlobalDB())
		}
	})
	return undoStoreInstance
}
