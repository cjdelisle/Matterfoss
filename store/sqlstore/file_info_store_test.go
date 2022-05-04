// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"testing"

	"github.com/cjdelisle/matterfoss-server/v6/store/searchtest"
	"github.com/cjdelisle/matterfoss-server/v6/store/storetest"
)

func TestFileInfoStore(t *testing.T) {
	StoreTest(t, storetest.TestFileInfoStore)
}

func TestSearchFileInfoStore(t *testing.T) {
	StoreTestWithSearchTestEngine(t, searchtest.TestSearchFileInfoStore)
}
