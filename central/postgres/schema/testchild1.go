// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/central/globaldb"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableTestchild1Stmt holds the create statement for table `testchild1`.
	CreateTableTestchild1Stmt = &postgres.CreateStmts{
		Table: `
               create table if not exists testchild1 (
                   Id varchar,
                   Val varchar,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}

	// Testchild1Schema is the go schema for table `testchild1`.
	Testchild1Schema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("testchild1")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestChild1)(nil)), "testchild1")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(63), "testchild1", (*storage.TestChild1)(nil)))
		globaldb.RegisterTable(schema)
		return schema
	}()
)
