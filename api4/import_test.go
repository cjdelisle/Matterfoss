// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/utils/fileutils"
)

func TestListImports(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	testsDir, _ := fileutils.FindDir("tests")
	require.NotEmpty(t, testsDir)

	uploadNewImport := func(c *model.Client4, t *testing.T) string {
		file, err := os.Open(testsDir + "/import_test.zip")
		require.NoError(t, err)

		info, err := file.Stat()
		require.NoError(t, err)

		us := &model.UploadSession{
			Filename: info.Name(),
			FileSize: info.Size(),
			Type:     model.UploadTypeImport,
		}

		if c == th.LocalClient {
			us.UserId = model.UploadNoUserID
		}

		u, _, err := c.CreateUpload(us)
		require.NoError(t, err)
		require.NotNil(t, u)

		finfo, _, err := c.UploadData(u.Id, file)
		require.NoError(t, err)
		require.NotNil(t, finfo)

		return u.Id
	}

	t.Run("no permissions", func(t *testing.T) {
		imports, _, err := th.Client.ListImports()
		require.Error(t, err)
		CheckErrorID(t, err, "api.context.permissions.app_error")
		require.Nil(t, imports)
	})

	dataDir, found := fileutils.FindDir("data")
	require.True(t, found)

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, c *model.Client4) {
		imports, _, err := c.ListImports()
		require.NoError(t, err)
		require.Empty(t, imports)
	}, "no imports")

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, c *model.Client4) {
		id := uploadNewImport(c, t)
		id2 := uploadNewImport(c, t)

		importDir := filepath.Join(dataDir, "import")
		f, err := os.Create(filepath.Join(importDir, "import.zip.tmp"))
		require.NoError(t, err)
		f.Close()

		imports, _, err := c.ListImports()
		require.NoError(t, err)
		require.NotEmpty(t, imports)
		require.Len(t, imports, 2)
		require.Contains(t, imports, id+"_import_test.zip")
		require.Contains(t, imports, id2+"_import_test.zip")

		require.NoError(t, os.RemoveAll(importDir))
	}, "expected imports")

	th.TestForSystemAdminAndLocal(t, func(t *testing.T, c *model.Client4) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ImportSettings.Directory = "import_new" })
		defer th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ImportSettings.Directory = "import" })

		importDir := filepath.Join(dataDir, "import_new")

		imports, _, err := c.ListImports()
		require.NoError(t, err)
		require.Empty(t, imports)

		id := uploadNewImport(c, t)
		imports, _, err = c.ListImports()
		require.NoError(t, err)
		require.NotEmpty(t, imports)
		require.Len(t, imports, 1)
		require.Equal(t, id+"_import_test.zip", imports[0])

		require.NoError(t, os.RemoveAll(importDir))
	}, "change import directory")
}
