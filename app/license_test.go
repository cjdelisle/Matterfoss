// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cjdelisle/matterfoss-server/v6/model"
)

func TestLoadLicense(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	th.App.Srv().LoadLicense()
	require.Nil(t, th.App.Srv().License(), "shouldn't have a valid license")
}

func TestSaveLicense(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	b1 := []byte("junk")

	_, err := th.App.Srv().SaveLicense(b1)
	require.NotNil(t, err, "shouldn't have saved license")
}

func TestRemoveLicense(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	err := th.App.Srv().RemoveLicense()
	require.Nil(t, err, "should have removed license")
}

func TestSetLicense(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	l1 := &model.License{}
	l1.Features = &model.Features{}
	l1.Customer = &model.Customer{}
	l1.StartsAt = model.GetMillis() - 1000
	l1.ExpiresAt = model.GetMillis() + 100000
	ok := th.App.Srv().SetLicense(l1)
	require.True(t, ok, "license should have worked")

	l3 := &model.License{}
	l3.Features = &model.Features{}
	l3.Customer = &model.Customer{}
	l3.StartsAt = model.GetMillis() + 10000
	l3.ExpiresAt = model.GetMillis() + 100000
	ok = th.App.Srv().SetLicense(l3)
	require.True(t, ok, "license should have passed")
}

func TestGetSanitizedClientLicense(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	setLicense(th, nil)

	m := th.App.Srv().GetSanitizedClientLicense()

	_, ok := m["Name"]
	assert.False(t, ok)
	_, ok = m["SkuName"]
	assert.False(t, ok)
	_, ok = m["SkuShortName"]
	assert.False(t, ok)
}

func TestGenerateRenewalToken(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	t.Run("test invalid token", func(t *testing.T) {
		_, err := th.App.Srv().renewalTokenValid("badtoken", "")
		var vErr *jwt.ValidationError
		require.True(t, errors.As(err, &vErr))
	})

	t.Run("renewal token generated correctly", func(t *testing.T) {
		setLicense(th, nil)
		token, appErr := th.App.Srv().GenerateRenewalToken(JWTDefaultTokenExpiration)
		require.Nil(t, appErr)
		require.NotEmpty(t, token)
		defer th.App.Srv().Store.System().PermanentDeleteByName(model.SystemLicenseRenewalToken)

		customerEmail := th.App.Srv().License().Customer.Email
		validToken, err := th.App.Srv().renewalTokenValid(token, customerEmail)
		require.NoError(t, err)
		require.True(t, validToken)
	})

	t.Run("only one token should be active", func(t *testing.T) {
		setLicense(th, nil)
		token, appErr := th.App.Srv().GenerateRenewalToken(JWTDefaultTokenExpiration)
		require.Nil(t, appErr)
		require.NotEmpty(t, token)
		defer th.App.Srv().Store.System().PermanentDeleteByName(model.SystemLicenseRenewalToken)

		newToken, appErr := th.App.Srv().GenerateRenewalToken(JWTDefaultTokenExpiration)
		require.Nil(t, appErr)
		require.Equal(t, token, newToken)
	})

	t.Run("return error if there is no active license", func(t *testing.T) {
		th.App.Srv().SetLicense(nil)
		_, appErr := th.App.Srv().GenerateRenewalToken(JWTDefaultTokenExpiration)
		require.NotNil(t, appErr)
	})

	t.Run("return another token if the license owner change", func(t *testing.T) {
		setLicense(th, nil)
		token, appErr := th.App.Srv().GenerateRenewalToken(JWTDefaultTokenExpiration)
		require.Nil(t, appErr)
		require.NotEmpty(t, token)
		defer th.App.Srv().Store.System().PermanentDeleteByName(model.SystemLicenseRenewalToken)
		setLicense(th, &model.Customer{
			Name:  "another customer",
			Email: "another@example.com",
		})
		newToken, appErr := th.App.Srv().GenerateRenewalToken(JWTDefaultTokenExpiration)
		require.Nil(t, appErr)
		require.NotEqual(t, token, newToken)
	})

	t.Run("return another token if the active one has expired", func(t *testing.T) {
		setLicense(th, nil)
		token, appErr := th.App.Srv().GenerateRenewalToken(1 * time.Second)
		require.Nil(t, appErr)
		require.NotEmpty(t, token)
		defer th.App.Srv().Store.System().PermanentDeleteByName(model.SystemLicenseRenewalToken)
		// The small time unit for expiration we're using is seconds
		time.Sleep(1 * time.Second)
		newToken, appErr := th.App.Srv().GenerateRenewalToken(JWTDefaultTokenExpiration)
		require.Nil(t, appErr)
		require.NotEqual(t, token, newToken)
	})

}

func setLicense(th *TestHelper, customer *model.Customer) {
	l1 := &model.License{}
	l1.Features = &model.Features{}
	if customer != nil {
		l1.Customer = customer
	} else {
		l1.Customer = &model.Customer{}
		l1.Customer.Name = "TestName"
		l1.Customer.Email = "test@example.com"
	}
	l1.SkuName = "SKU NAME"
	l1.SkuShortName = "SKU SHORT NAME"
	l1.StartsAt = model.GetMillis() - 1000
	l1.ExpiresAt = model.GetMillis() + 100000
	th.App.Srv().SetLicense(l1)
}
