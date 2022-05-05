// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package slashcommands

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cjdelisle/matterfoss-server/v6/model"
)

func TestMeProviderDoCommand(t *testing.T) {
	th := setup(t)
	defer th.tearDown()

	mp := MeProvider{}

	msg := "hello"

	resp := mp.DoCommand(th.App, th.Context, &model.CommandArgs{}, msg)

	assert.Equal(t, model.CommandResponseTypeInChannel, resp.ResponseType)
	assert.Equal(t, model.PostTypeMe, resp.Type)
	assert.Equal(t, "*"+msg+"*", resp.Text)
}
