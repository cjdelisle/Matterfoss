// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package plugin_test

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/plugin"
)

// configuration represents the configuration for this plugin as exposed via the Matterfoss
// server configuration.
type configuration struct {
	TeamName    string
	ChannelName string

	// channelID is resolved when the public configuration fields above change
	channelID string
}

// HelpPlugin implements the interface expected by the Matterfoss server to communicate
// between the server and plugin processes.
type HelpPlugin struct {
	plugin.MatterfossPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *HelpPlugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant.
func (p *HelpPlugin) setConfiguration(configuration *configuration) {
	// Replace the active configuration under lock.
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()
	p.configuration = configuration
}

// OnConfigurationChange updates the active configuration for this plugin under lock.
func (p *HelpPlugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Matterfoss server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	team, err := p.API.GetTeamByName(configuration.TeamName)
	if err != nil {
		return errors.Wrapf(err, "failed to find team %s", configuration.TeamName)
	}

	channel, err := p.API.GetChannelByName(configuration.ChannelName, team.Id, false)
	if err != nil {
		return errors.Wrapf(err, "failed to find channel %s", configuration.ChannelName)
	}

	configuration.channelID = channel.Id

	p.setConfiguration(configuration)

	return nil
}

// MessageHasBeenPosted automatically replies to posts that plea for help.
func (p *HelpPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	configuration := p.getConfiguration()

	// Ignore posts not in the configured channel
	if post.ChannelId != configuration.channelID {
		return
	}

	// Ignore posts this plugin made.
	if sentByPlugin, _ := post.GetProp("sent_by_plugin").(bool); sentByPlugin {
		return
	}

	// Ignore posts without a plea for help.
	if !strings.Contains(post.Message, "help") {
		return
	}

	p.API.SendEphemeralPost(post.UserId, &model.Post{
		ChannelId: configuration.channelID,
		Message:   "You asked for help? Checkout https://support.matterfoss.org/hc/en-us",
		Props: map[string]interface{}{
			"sent_by_plugin": true,
		},
	})
}

func Example_helpPlugin() {
	plugin.ClientMain(&HelpPlugin{})
}
