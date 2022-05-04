// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main_test

import (
	"database/sql"
	"testing"

	"github.com/cjdelisle/matterfoss-server/v6/app/plugin_api_tests"
	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/plugin"
	"github.com/cjdelisle/matterfoss-server/v6/shared/driver"
	"github.com/cjdelisle/matterfoss-server/v6/store/sqlstore"
	"github.com/cjdelisle/matterfoss-server/v6/store/storetest"
)

type MyPlugin struct {
	plugin.MattermostPlugin
	config plugin_api_tests.BasicConfig
	t      *testing.T
}

func (p *MyPlugin) OnConfigurationChange() error {
	if err := p.API.LoadPluginConfiguration(&p.config); err != nil {
		return err
	}
	return nil
}

func (p *MyPlugin) MessageWillBePosted(_ *plugin.Context, _ *model.Post) (*model.Post, string) {
	settings := p.API.GetUnsanitizedConfig().SqlSettings
	settings.Trace = model.NewBool(false)
	store := sqlstore.New(settings, nil)
	store.GetMasterX().Close()

	for _, isMaster := range []bool{true, false} {
		handle := sql.OpenDB(driver.NewConnector(p.Driver, isMaster))
		store.SetMasterX(handle)

		wrapper := sqlstore.NewStoreTestWrapper(store)
		// Testing with a handful of stores
		storetest.TestPostStore(p.t, store, wrapper)
		storetest.TestUserStore(p.t, store, wrapper)
		storetest.TestTeamStore(p.t, store)
		storetest.TestChannelStore(p.t, store, wrapper)
		storetest.TestBotStore(p.t, store, wrapper)

		store.GetMasterX().Close()
	}

	// Use the API to instantiate the driver
	// And then run the full suite of tests.
	return nil, "OK"
}

// TestDBAPI is a test function which actually runs a plugin. The objective
// is to run the storetest suite from inside a plugin.
//
// The test runner compiles the test code to a binary, and runs it as a normal
// binary. But under the hood, a test runs.
func TestDBAPI(t *testing.T) {
	plugin.ClientMain(&MyPlugin{t: t})
}
