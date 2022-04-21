// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package plugintest_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/plugin"
	"github.com/cjdelisle/matterfoss-server/v6/plugin/plugintest"
)

type HelloUserPlugin struct {
	plugin.MattermostPlugin
}

func (p *HelloUserPlugin) ServeHTTP(context *plugin.Context, w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	user, err := p.API.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.API.LogError(err.Error())
		return
	}

	fmt.Fprintf(w, "Welcome back, %s!", user.Username)
}

func Example() {
	t := &testing.T{}
	user := &model.User{
		Id:       model.NewId(),
		Username: "billybob",
	}

	api := &plugintest.API{}
	api.On("GetUser", user.Id).Return(user, nil)
	defer api.AssertExpectations(t)

	p := &HelloUserPlugin{}
	p.SetAPI(api)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Mattermost-User-Id", user.Id)
	p.ServeHTTP(&plugin.Context{}, w, r)
	body, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, "Welcome back, billybob!", string(body))
}
