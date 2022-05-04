// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"testing"

	"github.com/cjdelisle/matterfoss-server/v6/store/storetest"
)

func TestTokensStore(t *testing.T) {
	StoreTest(t, storetest.TestTokensStore)
}
