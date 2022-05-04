// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package email

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/shared/mail"
)

func TestCondenseSiteURL(t *testing.T) {
	require.Equal(t, "", condenseSiteURL(""))
	require.Equal(t, "mattermost.com", condenseSiteURL("mattermost.com"))
	require.Equal(t, "mattermost.com", condenseSiteURL("mattermost.com/"))
	require.Equal(t, "chat.mattermost.com", condenseSiteURL("chat.mattermost.com"))
	require.Equal(t, "chat.mattermost.com", condenseSiteURL("chat.mattermost.com/"))
	require.Equal(t, "mattermost.com/subpath", condenseSiteURL("mattermost.com/subpath"))
	require.Equal(t, "mattermost.com/subpath", condenseSiteURL("mattermost.com/subpath/"))
	require.Equal(t, "chat.mattermost.com/subpath", condenseSiteURL("chat.mattermost.com/subpath"))
	require.Equal(t, "chat.mattermost.com/subpath", condenseSiteURL("chat.mattermost.com/subpath/"))

	require.Equal(t, "mattermost.com:8080", condenseSiteURL("http://mattermost.com:8080"))
	require.Equal(t, "mattermost.com:8080", condenseSiteURL("http://mattermost.com:8080/"))
	require.Equal(t, "chat.mattermost.com:8080", condenseSiteURL("http://chat.mattermost.com:8080"))
	require.Equal(t, "chat.mattermost.com:8080", condenseSiteURL("http://chat.mattermost.com:8080/"))
	require.Equal(t, "mattermost.com:8080/subpath", condenseSiteURL("http://mattermost.com:8080/subpath"))
	require.Equal(t, "mattermost.com:8080/subpath", condenseSiteURL("http://mattermost.com:8080/subpath/"))
	require.Equal(t, "chat.mattermost.com:8080/subpath", condenseSiteURL("http://chat.mattermost.com:8080/subpath"))
	require.Equal(t, "chat.mattermost.com:8080/subpath", condenseSiteURL("http://chat.mattermost.com:8080/subpath/"))
}

func TestSendInviteEmails(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	th.ConfigureInbucketMail()

	emailTo := "test@example.com"

	retrieveEmail := func(t *testing.T) mail.JSONMessageInbucket {
		t.Helper()
		var resultsMailbox mail.JSONMessageHeaderInbucket
		err2 := mail.RetryInbucket(5, func() error {
			var err error
			resultsMailbox, err = mail.GetMailBox(emailTo)
			return err
		})
		if err2 != nil {
			t.Skipf("No email was received, maybe due load on the server: %v", err2)
		}
		require.Len(t, resultsMailbox, 1)
		require.Contains(t, resultsMailbox[0].To[0], emailTo, "Wrong To: recipient")
		resultsEmail, err := mail.GetMessageFromMailbox(emailTo, resultsMailbox[0].ID)
		require.NoError(t, err, "Could not get message from mailbox")
		return resultsEmail
	}

	verifyMailbox := func(t *testing.T) {
		t.Helper()
		email := retrieveEmail(t)
		require.Contains(t, email.Body.HTML, "http://testserver", "Wrong received message %s", email.Body.Text)
		require.Contains(t, email.Body.HTML, "test-user", "Wrong received message %s", email.Body.Text)
		require.Contains(t, email.Body.Text, "http://testserver", "Wrong received message %s", email.Body.Text)
		require.Contains(t, email.Body.Text, "test-user", "Wrong received message %s", email.Body.Text)
	}

	th.UpdateConfig(func(cfg *model.Config) {
		*cfg.ServiceSettings.EnableEmailInvitations = true
		*cfg.EmailSettings.SendEmailNotifications = false
	})
	t.Run("SendInviteEmails", func(t *testing.T) {
		mail.DeleteMailBox(emailTo)

		err := th.service.SendInviteEmails(th.BasicTeam, "test-user", th.BasicUser.Id, []string{emailTo}, "http://testserver", nil, false)
		require.NoError(t, err)

		verifyMailbox(t)
	})

	t.Run("SendInviteEmails can return error when SMTP connection fails", func(t *testing.T) {
		originalPort := *th.service.config().EmailSettings.SMTPPort
		th.UpdateConfig(func(cfg *model.Config) {
			os.Setenv("MM_EMAILSETTINGS_SMTPPORT", "5432")
			*cfg.EmailSettings.SMTPPort = "5432"
		})
		defer th.UpdateConfig(func(cfg *model.Config) {
			os.Setenv("MM_EMAILSETTINGS_SMTPPORT", originalPort)
			*cfg.EmailSettings.SMTPPort = originalPort
		})

		err := th.service.SendInviteEmails(th.BasicTeam, "test-user", th.BasicUser.Id, []string{emailTo}, "http://testserver", nil, true)
		require.Error(t, err)

		err = th.service.SendInviteEmails(th.BasicTeam, "test-user", th.BasicUser.Id, []string{emailTo}, "http://testserver", nil, false)
		require.NoError(t, err)
	})

	t.Run("SendGuestInviteEmails", func(t *testing.T) {
		mail.DeleteMailBox(emailTo)

		err := th.service.SendGuestInviteEmails(
			th.BasicTeam,
			[]*model.Channel{th.BasicChannel},
			"test-user",
			th.BasicUser.Id,
			nil,
			[]string{emailTo},
			"http://testserver",
			"hello world",
			false,
		)
		require.NoError(t, err)

		verifyMailbox(t)
	})

	t.Run("SendGuestInviteEmail can return error when SMTP connection fails", func(t *testing.T) {
		originalPort := *th.service.config().EmailSettings.SMTPPort
		th.UpdateConfig(func(cfg *model.Config) {
			os.Setenv("MM_EMAILSETTINGS_SMTPPORT", "5432")
			*cfg.EmailSettings.SMTPPort = "5432"
		})
		defer th.UpdateConfig(func(cfg *model.Config) {
			os.Setenv("MM_EMAILSETTINGS_SMTPPORT", originalPort)
			*cfg.EmailSettings.SMTPPort = originalPort
		})

		err := th.service.SendGuestInviteEmails(
			th.BasicTeam,
			[]*model.Channel{th.BasicChannel},
			"test-user",
			th.BasicUser.Id,
			nil,
			[]string{emailTo},
			"http://testserver",
			"hello world",
			false,
		)
		require.NoError(t, err)

		err = th.service.SendGuestInviteEmails(
			th.BasicTeam,
			[]*model.Channel{th.BasicChannel},
			"test-user",
			th.BasicUser.Id,
			nil,
			[]string{emailTo},
			"http://testserver",
			"hello world",
			true,
		)
		require.Error(t, err)

	})

	t.Run("SendGuestInviteEmails should sanitize HTML input", func(t *testing.T) {
		mail.DeleteMailBox(emailTo)

		message := `<a href="http://testserver">sanitized message</a>`
		err := th.service.SendGuestInviteEmails(
			th.BasicTeam,
			[]*model.Channel{th.BasicChannel},
			"test-user",
			th.BasicUser.Id,
			nil,
			[]string{emailTo},
			"http://testserver",
			message,
			false,
		)
		require.NoError(t, err)

		email := retrieveEmail(t)
		require.NotContains(t, email.Body.HTML, message)
		require.Contains(t, email.Body.HTML, "sanitized message")
		require.Contains(t, email.Body.Text, "sanitized message")
	})
}

func TestSendCloudTrialEndWarningEmail(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	th.ConfigureInbucketMail()

	emailTo := "testclouduser@example.com"
	emailToUsername := strings.Split(emailTo, "@")[0]

	t.Run("SendCloudTrialEndWarningEmail", func(t *testing.T) {
		verifyMailbox := func(t *testing.T) {
			t.Helper()

			var resultsMailbox mail.JSONMessageHeaderInbucket
			err2 := mail.RetryInbucket(5, func() error {
				var err error
				resultsMailbox, err = mail.GetMailBox(emailTo)
				return err
			})
			if err2 != nil {
				t.Skipf("No email was received, maybe due load on the server: %v", err2)
			}

			require.Len(t, resultsMailbox, 1)
			require.Contains(t, resultsMailbox[0].To[0], emailTo, "Wrong To: recipient")
			resultsEmail, err := mail.GetMessageFromMailbox(emailTo, resultsMailbox[0].ID)
			require.NoError(t, err, "Could not get message from mailbox")
			require.Contains(t, resultsEmail.Body.HTML, "http://testserver", "Wrong received message %s", resultsEmail.Body.Text)
			require.Contains(t, resultsEmail.Body.HTML, emailToUsername, "Wrong received message %s", resultsEmail.Body.Text)
			require.Contains(t, resultsEmail.Body.Text, "http://testserver", "Wrong received message %s", resultsEmail.Body.Text)
			require.Contains(t, resultsEmail.Body.Text, emailToUsername, "Wrong received message %s", resultsEmail.Body.Text)
			require.Contains(t, resultsEmail.Body.Text, "feedback-cloud@mattermost.com")
		}
		mail.DeleteMailBox(emailTo)

		err := th.service.SendCloudTrialEndWarningEmail(emailTo, emailToUsername, "June 23, 2200", th.BasicUser.Locale, "http://testserver")
		require.NoError(t, err)

		verifyMailbox(t)
	})
}

func TestSendCloudTrialEndedEmail(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()
	th.ConfigureInbucketMail()

	emailTo := "testclouduser@example.com"
	emailToUsername := strings.Split(emailTo, "@")[0]

	t.Run("SendCloudTrialEndedEmail", func(t *testing.T) {
		verifyMailbox := func(t *testing.T) {
			t.Helper()

			var resultsMailbox mail.JSONMessageHeaderInbucket
			err2 := mail.RetryInbucket(5, func() error {
				var err error
				resultsMailbox, err = mail.GetMailBox(emailTo)
				return err
			})
			if err2 != nil {
				t.Skipf("No email was received, maybe due load on the server: %v", err2)
			}

			require.Len(t, resultsMailbox, 1)
			require.Contains(t, resultsMailbox[0].To[0], emailTo, "Wrong To: recipient")
			resultsEmail, err := mail.GetMessageFromMailbox(emailTo, resultsMailbox[0].ID)
			require.NoError(t, err, "Could not get message from mailbox")
			require.Contains(t, resultsEmail.Body.Text, "your 14-day free trial of Mattermost Cloud Enterprise has ended today", "Wrong received message %s", resultsEmail.Body.Text)
			require.Contains(t, resultsEmail.Body.Text, "we will delete your Cloud workspace permanently", "Wrong received message %s", resultsEmail.Body.Text)
		}
		mail.DeleteMailBox(emailTo)

		err := th.service.SendCloudTrialEndedEmail(emailTo, emailToUsername, "June 23, 2200", th.BasicUser.Locale)
		require.NoError(t, err)

		verifyMailbox(t)
	})
}
