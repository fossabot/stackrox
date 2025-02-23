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
	// CreateTableImageComponentsStmt holds the create statement for table `image_components`.
	CreateTableImageComponentsStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists image_components (
                   Id varchar,
                   Name varchar,
                   Version varchar,
                   Source integer,
                   RiskScore numeric,
                   TopCvss numeric,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}

	// ImageComponentsSchema is the go schema for table `image_components`.
	ImageComponentsSchema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("image_components")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ImageComponent)(nil)), "image_components")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_IMAGE_COMPONENTS, "image_components", (*storage.ImageComponent)(nil)))
		globaldb.RegisterTable(schema)
		return schema
	}()
)
