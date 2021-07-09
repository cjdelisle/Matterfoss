// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"github.com/cjdelisle/matterfoss-server/v5/app/plugin_api_tests"
	"github.com/cjdelisle/matterfoss-server/v5/model"
	"github.com/cjdelisle/matterfoss-server/v5/plugin"
)

type MyPlugin struct {
	plugin.MatterfossPlugin
	configuration plugin_api_tests.BasicConfig
}

func (p *MyPlugin) OnConfigurationChange() error {
	if err := p.API.LoadPluginConfiguration(&p.configuration); err != nil {
		return err
	}
	return nil
}

func (p *MyPlugin) MessageWillBePosted(_ *plugin.Context, _ *model.Post) (*model.Post, string) {

	channels, err := p.API.SearchChannels(p.configuration.BasicTeamID, p.configuration.BasicChannelName)
	if err != nil {
		return nil, err.Error()
	}
	if len(channels) != 1 {
		return nil, "Returned invalid number of channels"
	}

	channels, err = p.API.SearchChannels("invalidid", p.configuration.BasicChannelName)
	if err != nil {
		return nil, err.Error()
	}
	if len(channels) != 0 {
		return nil, "Returned invalid number of channels"
	}

	return nil, "OK"
}

func main() {
	plugin.ClientMain(&MyPlugin{})
}
