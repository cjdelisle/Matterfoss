// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCondenseSiteURL(t *testing.T) {
	require.Equal(t, "", condenseSiteURL(""))
	require.Equal(t, "matterfoss.org", condenseSiteURL("matterfoss.org"))
	require.Equal(t, "matterfoss.org", condenseSiteURL("matterfoss.org/"))
	require.Equal(t, "chat.matterfoss.org", condenseSiteURL("chat.matterfoss.org"))
	require.Equal(t, "chat.matterfoss.org", condenseSiteURL("chat.matterfoss.org/"))
	require.Equal(t, "matterfoss.org/subpath", condenseSiteURL("matterfoss.org/subpath"))
	require.Equal(t, "matterfoss.org/subpath", condenseSiteURL("matterfoss.org/subpath/"))
	require.Equal(t, "chat.matterfoss.org/subpath", condenseSiteURL("chat.matterfoss.org/subpath"))
	require.Equal(t, "chat.matterfoss.org/subpath", condenseSiteURL("chat.matterfoss.org/subpath/"))

	require.Equal(t, "matterfoss.org:8080", condenseSiteURL("http://matterfoss.org:8080"))
	require.Equal(t, "matterfoss.org:8080", condenseSiteURL("http://matterfoss.org:8080/"))
	require.Equal(t, "chat.matterfoss.org:8080", condenseSiteURL("http://chat.matterfoss.org:8080"))
	require.Equal(t, "chat.matterfoss.org:8080", condenseSiteURL("http://chat.matterfoss.org:8080/"))
	require.Equal(t, "matterfoss.org:8080/subpath", condenseSiteURL("http://matterfoss.org:8080/subpath"))
	require.Equal(t, "matterfoss.org:8080/subpath", condenseSiteURL("http://matterfoss.org:8080/subpath/"))
	require.Equal(t, "chat.matterfoss.org:8080/subpath", condenseSiteURL("http://chat.matterfoss.org:8080/subpath"))
	require.Equal(t, "chat.matterfoss.org:8080/subpath", condenseSiteURL("http://chat.matterfoss.org:8080/subpath/"))
}
