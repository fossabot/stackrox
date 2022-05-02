// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	"github.com/stackrox/rox/central/globaldb"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableTestgrandchild1Stmt holds the create statement for table `testgrandchild1`.
	CreateTableTestgrandchild1Stmt = &postgres.CreateStmts{
		Table: `
               create table if not exists testgrandchild1 (
                   Id varchar,
                   ParentId varchar,
                   ChildId varchar,
                   Val varchar,
                   serialized bytea,
                   PRIMARY KEY(Id),
                   CONSTRAINT fk_parent_table_0 FOREIGN KEY (ParentId) REFERENCES testchild1(Id) ON DELETE CASCADE
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}

	// Testgrandchild1Schema is the go schema for table `testgrandchild1`.
	Testgrandchild1Schema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("testgrandchild1")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestGrandChild1)(nil)), "testgrandchild1")
		referencedSchemas := map[string]*walker.Schema{
			"storage.TestChild1":       Testchild1Schema,
			"storage.TestGGrandChild1": Testggrandchild1Schema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(64), "testgrandchild1", (*storage.TestGrandChild1)(nil)))
		globaldb.RegisterTable(schema)
		return schema
	}()
)
