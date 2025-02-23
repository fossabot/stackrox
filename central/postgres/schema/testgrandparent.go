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
	// CreateTableTestgrandparentStmt holds the create statement for table `testgrandparent`.
	CreateTableTestgrandparentStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists testgrandparent (
                   Id varchar,
                   Val varchar,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes: []string{},
		Children: []*postgres.CreateStmts{
			&postgres.CreateStmts{
				Table: `
               create table if not exists testgrandparent_Embedded (
                   testgrandparent_Id varchar,
                   idx integer,
                   Val varchar,
                   PRIMARY KEY(testgrandparent_Id, idx),
                   CONSTRAINT fk_parent_table_0 FOREIGN KEY (testgrandparent_Id) REFERENCES testgrandparent(Id) ON DELETE CASCADE
               )
               `,
				Indexes: []string{
					"create index if not exists testgrandparentEmbedded_idx on testgrandparent_Embedded using btree(idx)",
				},
				Children: []*postgres.CreateStmts{
					&postgres.CreateStmts{
						Table: `
               create table if not exists testgrandparent_Embedded_Embedded2 (
                   testgrandparent_Id varchar,
                   testgrandparent_Embedded_idx integer,
                   idx integer,
                   Val varchar,
                   PRIMARY KEY(testgrandparent_Id, testgrandparent_Embedded_idx, idx),
                   CONSTRAINT fk_parent_table_0 FOREIGN KEY (testgrandparent_Id, testgrandparent_Embedded_idx) REFERENCES testgrandparent_Embedded(testgrandparent_Id, idx) ON DELETE CASCADE
               )
               `,
						Indexes: []string{
							"create index if not exists testgrandparentEmbeddedEmbedded2_idx on testgrandparent_Embedded_Embedded2 using btree(idx)",
						},
						Children: []*postgres.CreateStmts{},
					},
				},
			},
		},
	}

	// TestgrandparentSchema is the go schema for table `testgrandparent`.
	TestgrandparentSchema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("testgrandparent")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestGrandparent)(nil)), "testgrandparent")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(61), "testgrandparent", (*storage.TestGrandparent)(nil)))
		globaldb.RegisterTable(schema)
		return schema
	}()
)
