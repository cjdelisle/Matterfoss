// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"github.com/cjdelisle/matterfoss-server/v6/app"
	"github.com/cjdelisle/matterfoss-server/v6/app/request"
	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/shared/i18n"
)

type PurposeProvider struct {
}

const (
	CmdPurpose = "purpose"
)

func init() {
	app.RegisterCommandProvider(&PurposeProvider{})
}

func (*PurposeProvider) GetTrigger() string {
	return CmdPurpose
}

func (*PurposeProvider) GetCommand(a *app.App, T i18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CmdPurpose,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_channel_purpose.desc"),
		AutoCompleteHint: T("api.command_channel_purpose.hint"),
		DisplayName:      T("api.command_channel_purpose.name"),
	}
}

func (*PurposeProvider) DoCommand(a *app.App, c *request.Context, args *model.CommandArgs, message string) *model.CommandResponse {
	channel, err := a.GetChannel(args.ChannelId)
	if err != nil {
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_purpose.channel.app_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	switch channel.Type {
	case model.ChannelTypeOpen:
		if !a.HasPermissionToChannel(args.UserId, args.ChannelId, model.PermissionManagePublicChannelProperties) {
			return &model.CommandResponse{
				Text:         args.T("api.command_channel_purpose.permission.app_error"),
				ResponseType: model.CommandResponseTypeEphemeral,
			}
		}
	case model.ChannelTypePrivate:
		if !a.HasPermissionToChannel(args.UserId, args.ChannelId, model.PermissionManagePrivateChannelProperties) {
			return &model.CommandResponse{
				Text:         args.T("api.command_channel_purpose.permission.app_error"),
				ResponseType: model.CommandResponseTypeEphemeral,
			}
		}
	default:
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_purpose.direct_group.app_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	if message == "" {
		return &model.CommandResponse{
			Text:         args.T("api.command_channel_purpose.message.app_error"),
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	patch := &model.ChannelPatch{
		Purpose: new(string),
	}
	*patch.Purpose = message

	_, err = a.PatchChannel(c, channel, patch, args.UserId)
	if err != nil {
		text := args.T("api.command_channel_purpose.update_channel.app_error")
		if err.Id == "model.channel.is_valid.purpose.app_error" {
			text = args.T("api.command_channel_purpose.update_channel.max_length", map[string]interface{}{
				"MaxLength": model.ChannelPurposeMaxRunes,
			})
		}

		return &model.CommandResponse{
			Text:         text,
			ResponseType: model.CommandResponseTypeEphemeral,
		}
	}

	return &model.CommandResponse{}
}
