// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cjdelisle/matterfoss-server/v6/app/request"
	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/shared/mlog"
)

func (a *App) getSysAdminsEmailRecipients() ([]*model.User, *model.AppError) {
	userOptions := &model.UserGetOptions{
		Page:     0,
		PerPage:  100,
		Role:     model.SystemAdminRoleId,
		Inactive: false,
	}
	return a.GetUsers(userOptions)
}

// SendAdminUpgradeRequestEmail takes the username of user trying to alert admins and then applies rate limit of n (number of admins) emails per user per day
// before sending the emails.
func (a *App) SendAdminUpgradeRequestEmail(username string, subscription *model.Subscription, action string) *model.AppError {
	if a.Srv().License() == nil || (a.Srv().License() != nil && !*a.Srv().License().Features.Cloud) {
		return nil
	}

	if subscription != nil && subscription.IsPaidTier == "true" {
		return nil
	}

	year, month, day := time.Now().Date()
	key := fmt.Sprintf("%s-%d-%s-%d", action, day, month, year)

	if a.Srv().EmailService.GetPerDayEmailRateLimiter() == nil {
		return model.NewAppError("app.SendAdminUpgradeRequestEmail", "app.email.no_rate_limiter.app_error", nil, fmt.Sprintf("for key=%s", key), http.StatusInternalServerError)
	}

	// rate limit based on combination of date and action as key
	rateLimited, result, err := a.Srv().EmailService.GetPerDayEmailRateLimiter().RateLimit(key, 1)
	if err != nil {
		return model.NewAppError("app.SendAdminUpgradeRequestEmail", "app.email.setup_rate_limiter.app_error", nil, fmt.Sprintf("for key=%s, error=%v", key, err), http.StatusInternalServerError)
	}

	if rateLimited {
		return model.NewAppError("app.SendAdminUpgradeRequestEmail",
			"app.email.rate_limit_exceeded.app_error", map[string]interface{}{"RetryAfter": result.RetryAfter.String(), "ResetAfter": result.ResetAfter.String()},
			fmt.Sprintf("key=%s, retry_after_secs=%f, reset_after_secs=%f",
				key, result.RetryAfter.Seconds(), result.ResetAfter.Seconds()),
			http.StatusRequestEntityTooLarge)
	}

	sysAdmins, e := a.getSysAdminsEmailRecipients()
	if e != nil {
		return e
	}

	// we want to at least have one email sent out to an admin
	countNotOks := 0

	for admin := range sysAdmins {
		ok, err := a.Srv().EmailService.SendUpgradeEmail(username, sysAdmins[admin].Email, sysAdmins[admin].Locale, *a.Config().ServiceSettings.SiteURL, action)
		if !ok || err != nil {
			a.Log().Error("Error sending upgrade request email", mlog.Err(err))
			countNotOks++
		}
	}

	// if not even one admin got an email, we consider that this operation errored
	if countNotOks == len(sysAdmins) {
		return model.NewAppError("app.SendAdminUpgradeRequestEmail", "app.user.send_emails.app_error", nil, "", http.StatusInternalServerError)
	}

	return nil
}

func (a *App) GetSubscriptionStats() (*model.SubscriptionStats, *model.AppError) {
	if a.Srv().License() == nil || !*a.Srv().License().Features.Cloud {
		return nil, model.NewAppError("app.GetSubscriptionStats", "api.cloud.license_error", nil, "", http.StatusInternalServerError)
	}

	subscription, appErr := a.Cloud().GetSubscription("")
	if appErr != nil {
		return nil, model.NewAppError("app.GetSubscriptionStats", "api.cloud.request_error", nil, appErr.Error(), http.StatusInternalServerError)
	}

	count, err := a.Srv().Store.User().Count(model.UserCountOptions{})
	if err != nil {
		return nil, model.NewAppError("app.GetSubscriptionStats", "app.user.get_total_users_count.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	cloudUserLimit := *a.Config().ExperimentalSettings.CloudUserLimit

	s := cloudUserLimit - count

	return &model.SubscriptionStats{
		RemainingSeats: int(s),
		IsPaidTier:     subscription.IsPaidTier,
	}, nil
}

func (a *App) CheckCloudAccountAtLimit() (bool, *model.AppError) {
	if a.Srv().License() == nil || (a.Srv().License() != nil && !*a.Srv().License().Features.Cloud) {
		// Not cloud instance, so no at limit checks
		return false, nil
	}

	stats, err := a.GetSubscriptionStats()
	if err != nil {
		return false, err
	}

	if stats.IsPaidTier == "true" {
		return false, nil
	}

	if stats.RemainingSeats < 1 {
		return true, nil
	}

	return false, nil
}

func (a *App) CheckAndSendUserLimitWarningEmails(c *request.Context) *model.AppError {
	return nil
}

func (a *App) SendPaymentFailedEmail(failedPayment *model.FailedPayment) *model.AppError {
	sysAdmins, err := a.getSysAdminsEmailRecipients()
	if err != nil {
		return err
	}

	for _, admin := range sysAdmins {
		_, err := a.Srv().EmailService.SendPaymentFailedEmail(admin.Email, admin.Locale, failedPayment, *a.Config().ServiceSettings.SiteURL)
		if err != nil {
			a.Log().Error("Error sending payment failed email", mlog.Err(err))
		}
	}
	return nil
}

// SendNoCardPaymentFailedEmail
func (a *App) SendNoCardPaymentFailedEmail() *model.AppError {
	sysAdmins, err := a.getSysAdminsEmailRecipients()
	if err != nil {
		return err
	}

	for _, admin := range sysAdmins {
		err := a.Srv().EmailService.SendNoCardPaymentFailedEmail(admin.Email, admin.Locale, *a.Config().ServiceSettings.SiteURL)
		if err != nil {
			a.Log().Error("Error sending payment failed email", mlog.Err(err))
		}
	}
	return nil
}

func (a *App) SendCloudTrialEndWarningEmail(trialEndDate, siteURL string) *model.AppError {
	sysAdmins, e := a.getSysAdminsEmailRecipients()
	if e != nil {
		return e
	}

	// we want to at least have one email sent out to an admin
	countNotOks := 0

	for admin := range sysAdmins {
		name := sysAdmins[admin].FirstName
		if name == "" {
			name = sysAdmins[admin].Username
		}
		err := a.Srv().EmailService.SendCloudTrialEndWarningEmail(sysAdmins[admin].Email, name, trialEndDate, sysAdmins[admin].Locale, siteURL)
		if err != nil {
			a.Log().Error("Error sending trial ending warning to", mlog.String("email", sysAdmins[admin].Email), mlog.Err(err))
			countNotOks++
		}
	}

	// if not even one admin got an email, we consider that this operation errored
	if countNotOks == len(sysAdmins) {
		return model.NewAppError("app.SendCloudTrialEndWarningEmail", "app.user.send_emails.app_error", nil, "", http.StatusInternalServerError)
	}

	return nil
}

func (a *App) SendCloudTrialEndedEmail() *model.AppError {
	sysAdmins, e := a.getSysAdminsEmailRecipients()
	if e != nil {
		return e
	}

	// we want to at least have one email sent out to an admin
	countNotOks := 0

	for admin := range sysAdmins {
		name := sysAdmins[admin].FirstName
		if name == "" {
			name = sysAdmins[admin].Username
		}

		err := a.Srv().EmailService.SendCloudTrialEndedEmail(sysAdmins[admin].Email, name, sysAdmins[admin].Locale, *a.Config().ServiceSettings.SiteURL)
		if err != nil {
			a.Log().Error("Error sending trial ended email to", mlog.String("email", sysAdmins[admin].Email), mlog.Err(err))
			countNotOks++
		}
	}

	// if not even one admin got an email, we consider that this operation errored
	if countNotOks == len(sysAdmins) {
		return model.NewAppError("app.SendCloudTrialEndedEmail", "app.user.send_emails.app_error", nil, "", http.StatusInternalServerError)
	}

	return nil
}
