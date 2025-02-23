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
	// CreateTableServiceaccountsStmt holds the create statement for table `serviceaccounts`.
	CreateTableServiceaccountsStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists serviceaccounts (
                   Id varchar,
                   Name varchar,
                   Namespace varchar,
                   ClusterName varchar,
                   ClusterId varchar,
                   Labels jsonb,
                   Annotations jsonb,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}

	// ServiceaccountsSchema is the go schema for table `serviceaccounts`.
	ServiceaccountsSchema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("serviceaccounts")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ServiceAccount)(nil)), "serviceaccounts")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_SERVICE_ACCOUNTS, "serviceaccounts", (*storage.ServiceAccount)(nil)))
		globaldb.RegisterTable(schema)
		return schema
	}()
)
