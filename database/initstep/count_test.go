// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInitStep_Engine_CountInitSteps(t *testing.T) {
	// setup types
	_initStepOne := testInitStep()
	_initStepOne.SetID(1)
	_initStepOne.SetRepoID(1)
	_initStepOne.SetBuildID(1)
	_initStepOne.SetNumber(1)
	_initStepOne.SetReporter("Foobar Runtime")
	_initStepOne.SetName("foobar")
	_initStepOne.SetMimetype("text/plain")

	_initStepTwo := testInitStep()
	_initStepTwo.SetID(2)
	_initStepTwo.SetRepoID(1)
	_initStepTwo.SetBuildID(2)
	_initStepTwo.SetNumber(2)
	_initStepTwo.SetReporter("Foobar Runtime")
	_initStepTwo.SetName("foobar")
	_initStepTwo.SetMimetype("text/plain")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "initsteps"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateInitStep(_initStepOne)
	if err != nil {
		t.Errorf("unable to create test init step for sqlite: %v", err)
	}

	err = _sqlite.CreateInitStep(_initStepTwo)
	if err != nil {
		t.Errorf("unable to create test init step for sqlite: %v", err)
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
			got, err := test.database.CountInitSteps()

			if test.failure {
				if err == nil {
					t.Errorf("CountInitSteps for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountInitSteps for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountInitSteps for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}