// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/shared/i18n"
)

func TestMsgProvider(t *testing.T) {
	th := setup(t).initBasic()
	defer th.tearDown()

	team := th.createTeam()
	th.linkUserToTeam(th.BasicUser, team)
	cmd := &msgProvider{}

	th.removePermissionFromRole(model.PermissionCreateDirectChannel.Id, model.SystemUserRoleId)

	// Check without permission to create a DM channel.
	resp := cmd.DoCommand(th.App, th.Context, &model.CommandArgs{
		T:       i18n.IdentityTfunc(),
		SiteURL: "http://test.url",
		TeamId:  team.Id,
		UserId:  th.BasicUser.Id,
	}, "@"+th.BasicUser2.Username+" hello")

	channelName := model.GetDMNameFromIds(th.BasicUser.Id, th.BasicUser2.Id)
	assert.Equal(t, "api.command_msg.permission.app_error", resp.Text)
	assert.Equal(t, "", resp.GotoLocation)

	th.addPermissionToRole(model.PermissionCreateDirectChannel.Id, model.SystemUserRoleId)

	// Check with permission to create a DM channel.
	resp = cmd.DoCommand(th.App, th.Context, &model.CommandArgs{
		T:       i18n.IdentityTfunc(),
		SiteURL: "http://test.url",
		TeamId:  team.Id,
		UserId:  th.BasicUser.Id,
	}, "@"+th.BasicUser2.Username+" hello")

	assert.Equal(t, "", resp.Text)
	assert.Equal(t, "http://test.url/"+team.Name+"/channels/"+channelName, resp.GotoLocation)

	// Check without permission to post to an existing DM channel.
	resp = cmd.DoCommand(th.App, th.Context, &model.CommandArgs{
		T:       i18n.IdentityTfunc(),
		SiteURL: "http://test.url",
		TeamId:  team.Id,
		UserId:  th.BasicUser.Id,
	}, "@"+th.BasicUser2.Username+" hello")

	assert.Equal(t, "", resp.Text)
	assert.Equal(t, "http://test.url/"+team.Name+"/channels/"+channelName, resp.GotoLocation)

	// Check that a guest user cannot message a user who is not in a channel/team with him
	guest := th.createGuest()
	user := th.createUser()

	th.linkUserToTeam(user, team)
	th.linkUserToTeam(guest, th.BasicTeam)
	th.addUserToChannel(guest, th.BasicChannel)

	resp = cmd.DoCommand(th.App, th.Context, &model.CommandArgs{
		T:       i18n.IdentityTfunc(),
		SiteURL: "http://test.url",
		TeamId:  th.BasicTeam.Id,
		UserId:  guest.Id,
	}, "@"+user.Username+" hello")

	assert.Equal(t, "api.command_msg.missing.app_error", resp.Text)
	assert.Equal(t, "", resp.GotoLocation)

	// Check that a guest user can message a user who is in a channel/team with him
	th.linkUserToTeam(user, th.BasicTeam)
	th.addUserToChannel(user, th.BasicChannel)

	resp = cmd.DoCommand(th.App, th.Context, &model.CommandArgs{
		T:       i18n.IdentityTfunc(),
		SiteURL: "http://test.url",
		TeamId:  th.BasicTeam.Id,
		UserId:  guest.Id,
	}, "@"+user.Username+" hello")

	channelName = model.GetDMNameFromIds(guest.Id, user.Id)

	assert.Equal(t, "api.channel.create_channel.direct_channel.guest_restricted_error", resp.Text)
	assert.Equal(t, "", resp.GotoLocation)
}
