// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountHooksForRepo gets the count of hooks by repo ID from the database.
func (e *engine) CountHooksForRepo(r *library.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of hooks for repo %s from the database", r.GetFullName())

	// variable to store query results
	var h int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Where("repo_id = ?", r.GetID()).
		Count(&h).
		Error

	return h, err
}
