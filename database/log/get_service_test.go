// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestLog_Engine_GetLogForService(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetRepoID(1)
	_log.SetBuildID(1)
	_log.SetServiceID(1)
	_log.SetData([]byte{})

	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "repo_id", "service_id", "step_id", "initstep_id", "data"}).
		AddRow(1, 1, 1, 1, 0, 0, []byte{})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "logs" WHERE service_id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateLog(_log)
	if err != nil {
		t.Errorf("unable to create test log for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Log
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _log,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _log,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetLogForService(_service)

			if test.failure {
				if err == nil {
					t.Errorf("GetLogForService for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetLogForService for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetLogForService for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
