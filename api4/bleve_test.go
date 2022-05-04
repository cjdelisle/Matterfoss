// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"testing"

	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/stretchr/testify/require"
)

func TestBlevePurgeIndexes(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	t.Run("as system user", func(t *testing.T) {
		resp, err := th.Client.PurgeBleveIndexes()
		require.Error(t, err)
		CheckForbiddenStatus(t, resp)
	})

	t.Run("as system user with write experimental permission", func(t *testing.T) {
		th.AddPermissionToRole(model.PermissionPurgeBleveIndexes.Id, model.SystemUserRoleId)
		defer th.RemovePermissionFromRole(model.PermissionSysconsoleWriteExperimental.Id, model.SystemUserRoleId)
		resp, err := th.Client.PurgeBleveIndexes()
		require.NoError(t, err)
		CheckOKStatus(t, resp)
	})

	t.Run("as system admin", func(t *testing.T) {
		resp, err := th.SystemAdminClient.PurgeBleveIndexes()
		require.NoError(t, err)
		CheckOKStatus(t, resp)
	})

	t.Run("as restricted system admin", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ExperimentalSettings.RestrictSystemAdmin = true })

		resp, err := th.SystemAdminClient.PurgeBleveIndexes()
		require.Error(t, err)
		CheckForbiddenStatus(t, resp)
	})
}
