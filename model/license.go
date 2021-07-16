// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"math"
)

const (
	EXPIRED_LICENSE_ERROR = "api.license.add_license.expired.app_error"
	INVALID_LICENSE_ERROR = "api.license.add_license.invalid.app_error"
	LICENSE_GRACE_PERIOD  = 1000 * 60 * 60 * 24 * 10 //10 days
	LICENSE_RENEWAL_LINK  = "https://matterfoss.org/renew/"
)

type LicenseRecord struct {
	Id       string `json:"id"`
	CreateAt int64  `json:"create_at"`
	Bytes    string `json:"-"`
}

type License struct {
	Id           string    `json:"id"`
	IssuedAt     int64     `json:"issued_at"`
	StartsAt     int64     `json:"starts_at"`
	ExpiresAt    int64     `json:"expires_at"`
	Customer     *Customer `json:"customer"`
	Features     *Features `json:"features"`
	SkuName      string    `json:"sku_name"`
	SkuShortName string    `json:"sku_short_name"`
}

type Customer struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
}

type TrialLicenseRequest struct {
	ServerID              string `json:"server_id"`
	Email                 string `json:"email"`
	Name                  string `json:"name"`
	SiteURL               string `json:"site_url"`
	SiteName              string `json:"site_name"`
	Users                 int    `json:"users"`
	TermsAccepted         bool   `json:"terms_accepted"`
	ReceiveEmailsAccepted bool   `json:"receive_emails_accepted"`
}

func (tlr *TrialLicenseRequest) ToJson() string {
	b, _ := json.Marshal(tlr)
	return string(b)
}

type Features struct {
	Users                     *int  `json:"users"`
	LDAP                      *bool `json:"ldap"`
	LDAPGroups                *bool `json:"ldap_groups"`
	MFA                       *bool `json:"mfa"`
	GoogleOAuth               *bool `json:"google_oauth"`
	Office365OAuth            *bool `json:"office365_oauth"`
	OpenId                    *bool `json:"openid"`
	Compliance                *bool `json:"compliance"`
	Cluster                   *bool `json:"cluster"`
	Metrics                   *bool `json:"metrics"`
	MHPNS                     *bool `json:"mhpns"`
	SAML                      *bool `json:"saml"`
	Elasticsearch             *bool `json:"elastic_search"`
	Announcement              *bool `json:"announcement"`
	ThemeManagement           *bool `json:"theme_management"`
	EmailNotificationContents *bool `json:"email_notification_contents"`
	DataRetention             *bool `json:"data_retention"`
	MessageExport             *bool `json:"message_export"`
	CustomPermissionsSchemes  *bool `json:"custom_permissions_schemes"`
	CustomTermsOfService      *bool `json:"custom_terms_of_service"`
	GuestAccounts             *bool `json:"guest_accounts"`
	GuestAccountsPermissions  *bool `json:"guest_accounts_permissions"`
	IDLoadedPushNotifications *bool `json:"id_loaded"`
	LockTeammateNameDisplay   *bool `json:"lock_teammate_name_display"`
	EnterprisePlugins         *bool `json:"enterprise_plugins"`
	AdvancedLogging           *bool `json:"advanced_logging"`
	Cloud                     *bool `json:"cloud"`
	SharedChannels            *bool `json:"shared_channels"`
	RemoteClusterService      *bool `json:"remote_cluster_service"`

	// after we enabled more features we'll need to control them with this
	FutureFeatures *bool `json:"future_features"`
}

func (f *Features) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ldap":                        *f.LDAP,
		"ldap_groups":                 *f.LDAPGroups,
		"mfa":                         *f.MFA,
		"google":                      *f.GoogleOAuth,
		"office365":                   *f.Office365OAuth,
		"openid":                      *f.OpenId,
		"compliance":                  *f.Compliance,
		"cluster":                     *f.Cluster,
		"metrics":                     *f.Metrics,
		"mhpns":                       *f.MHPNS,
		"saml":                        *f.SAML,
		"elastic_search":              *f.Elasticsearch,
		"email_notification_contents": *f.EmailNotificationContents,
		"data_retention":              *f.DataRetention,
		"message_export":              *f.MessageExport,
		"custom_permissions_schemes":  *f.CustomPermissionsSchemes,
		"guest_accounts":              *f.GuestAccounts,
		"guest_accounts_permissions":  *f.GuestAccountsPermissions,
		"id_loaded":                   *f.IDLoadedPushNotifications,
		"lock_teammate_name_display":  *f.LockTeammateNameDisplay,
		"enterprise_plugins":          *f.EnterprisePlugins,
		"advanced_logging":            *f.AdvancedLogging,
		"cloud":                       *f.Cloud,
		"shared_channels":             *f.SharedChannels,
		"remote_cluster_service":      *f.RemoteClusterService,
		"future":                      *f.FutureFeatures,
	}
}

func (f *Features) SetDefaults() {
	f.FutureFeatures = NewBool(true)
	f.Users = NewInt(math.MaxInt64)
	f.LDAP = NewBool(false)
	f.LDAPGroups = NewBool(true)
	f.MFA = NewBool(true)
	f.GoogleOAuth = NewBool(false)
	f.Office365OAuth = NewBool(false)
	f.OpenId = NewBool(true)
	f.Compliance = NewBool(true)
	f.Cluster = NewBool(true)
	f.Metrics = NewBool(true)
	f.MHPNS = NewBool(true)
	f.SAML = NewBool(false)
	f.Elasticsearch = NewBool(false)
	f.Announcement = NewBool(true)
	f.ThemeManagement = NewBool(true)
	f.EmailNotificationContents = NewBool(true)
	f.DataRetention = NewBool(true)
	f.MessageExport = NewBool(true)
	f.CustomPermissionsSchemes = NewBool(true)
	f.GuestAccounts = NewBool(true)
	f.GuestAccountsPermissions = NewBool(true)
	f.CustomTermsOfService = NewBool(true)
	f.IDLoadedPushNotifications = NewBool(true)
	f.LockTeammateNameDisplay = NewBool(true)
	f.EnterprisePlugins = NewBool(true)
	f.AdvancedLogging = NewBool(true)
	f.Cloud = NewBool(false)
	f.SharedChannels = NewBool(true)
	f.RemoteClusterService = NewBool(false)
}

func (l *License) IsExpired() bool {
	return false
}

func (l *License) IsPastGracePeriod() bool {
	return false
}

func (l *License) IsStarted() bool {
	return true
}

func (l *License) ToJson() string {
	b, _ := json.Marshal(l)
	return string(b)
}

// NewTestLicense returns a license that expires in the future and has the given features.
func NewTestLicense(features ...string) *License {
	ret := &License{
		ExpiresAt: GetMillis() + 90*24*60*60*1000,
		Customer:  &Customer{},
		Features:  &Features{},
	}
	ret.Features.SetDefaults()

	featureMap := map[string]bool{}
	for _, feature := range features {
		featureMap[feature] = true
	}
	featureJson, _ := json.Marshal(featureMap)
	json.Unmarshal(featureJson, &ret.Features)

	return ret
}

func LicenseFromJson(data io.Reader) *License {
	var o *License
	json.NewDecoder(data).Decode(&o)
	return o
}

func (lr *LicenseRecord) IsValid() *AppError {
	return nil
}

func (lr *LicenseRecord) PreSave() {
	lr.CreateAt = GetMillis()
}
