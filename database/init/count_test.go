// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInit_Engine_CountInits(t *testing.T) {
	// setup types
	_initOne := testInit()
	_initOne.SetID(1)
	_initOne.SetRepoID(1)
	_initOne.SetBuildID(1)
	_initOne.SetNumber(1)
	_initOne.SetReporter("Foobar Runtime")
	_initOne.SetName("foobar")
	_initOne.SetMimetype("text/plain")

	_initTwo := testInit()
	_initTwo.SetID(2)
	_initTwo.SetRepoID(1)
	_initTwo.SetBuildID(2)
	_initTwo.SetNumber(2)
	_initTwo.SetReporter("Foobar Runtime")
	_initTwo.SetName("foobar")
	_initTwo.SetMimetype("text/plain")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "inits"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateInit(_initOne)
	if err != nil {
		t.Errorf("unable to create test init for sqlite: %v", err)
	}

	err = _sqlite.CreateInit(_initTwo)
	if err != nil {
		t.Errorf("unable to create test init for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     2,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     2,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountInits()

			if test.failure {
				if err == nil {
					t.Errorf("CountInits for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountInits for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountInits for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
