// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres workers table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
workers (
	id               SERIAL PRIMARY KEY,
	hostname         VARCHAR(250),
	address          VARCHAR(250),
	routes           VARCHAR(1000),
	active           BOOLEAN,
	last_checked_in  INTEGER,
	build_limit      INTEGER,
	UNIQUE(hostname)
);
`

	// CreateSqliteTable represents a query to create the Sqlite workers table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
workers (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	hostname        TEXT,
	address         TEXT,
	routes          TEXT,
	active          BOOLEAN,
	last_checked_in	INTEGER,
	build_limit     INTEGER,
	UNIQUE(hostname)
);
`
)

// CreateWorkerTable creates the workers table in the database.
func (e *engine) CreateWorkerTable(driver string) error {
	e.logger.Tracef("creating workers table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the workers table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the workers table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
