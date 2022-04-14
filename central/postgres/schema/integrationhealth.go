// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"github.com/stackrox/rox/pkg/postgres"
)

var (
	// CreateTableIntegrationhealthStmt holds the create statement for table `Integrationhealth`.
	CreateTableIntegrationhealthStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists integrationhealth (
                   Id varchar,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}
)
