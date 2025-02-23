// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableNotifiersStmt holds the create statement for table `notifiers`.
	CreateTableNotifiersStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists notifiers (
                   Id varchar,
                   Name varchar UNIQUE,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}

	// NotifiersSchema is the go schema for table `notifiers`.
	NotifiersSchema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("notifiers")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.Notifier)(nil)), "notifiers")
		globaldb.RegisterTable(schema)
		return schema
	}()
)
