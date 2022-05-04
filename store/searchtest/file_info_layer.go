// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package searchtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cjdelisle/matterfoss-server/v6/model"
	"github.com/cjdelisle/matterfoss-server/v6/store"
)

var searchFileInfoStoreTests = []searchTest{
	{
		Name: "Should be able to search posts including results from DMs",
		Fn:   testFileInfoSearchFileInfosIncludingDMs,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to search posts using pagination",
		Fn:   testFileInfoSearchFileInfosWithPagination,
		Tags: []string{EngineElasticSearch, EngineBleve},
	},
	{
		Name: "Should be able to search for exact phrases in quotes",
		Fn:   testFileInfoSearchExactPhraseInQuotes,
		Tags: []string{EnginePostgres, EngineMySql, EngineElasticSearch},
	},
	{
		Name: "Should be able to search for email addresses with or without quotes",
		Fn:   testFileInfoSearchEmailAddresses,
		Tags: []string{EngineElasticSearch},
	},
	{
		Name: "Should be able to search when markdown underscores are applied",
		Fn:   testFileInfoSearchMarkdownUnderscores,
		Tags: []string{EnginePostgres, EngineElasticSearch},
	},
	{
		Name: "Should be able to search for non-latin words",
		Fn:   testFileInfoSearchNonLatinWords,
		Tags: []string{EngineElasticSearch},
	},
	{
		Name: "Should be able to search for alternative spellings of words",
		Fn:   testFileInfoSearchAlternativeSpellings,
		Tags: []string{EngineElasticSearch},
	},
	{
		Name: "Should be able to search for alternative spellings of words with and without accents",
		Fn:   testFileInfoSearchAlternativeSpellingsAccents,
		Tags: []string{EngineElasticSearch},
	},
	{
		Name: "Should be able to search or exclude messages written by a specific user",
		Fn:   testFileInfoSearchOrExcludeFileInfosBySpecificUser,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to search or exclude messages written in a specific channel",
		Fn:   testFileInfoSearchOrExcludeFileInfosInChannel,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to search or exclude messages written in a DM or GM",
		Fn:   testFileInfoSearchOrExcludeFileInfosInDMGM,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to search or exclude files by extensions",
		Fn:   testFileInfoSearchOrExcludeByExtensions,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to filter messages written after a specific date",
		Fn:   testFileInfoFilterFilesAfterSpecificDate,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to filter messages written before a specific date",
		Fn:   testFileInfoFilterFilesBeforeSpecificDate,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to filter messages written on a specific date",
		Fn:   testFileInfoFilterFilesInSpecificDate,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to exclude messages that contain a search term",
		Fn:   testFileInfoFilterFilesWithATerm,
		Tags: []string{EngineMySql, EnginePostgres},
	},
	{
		Name: "Should be able to search using boolean operators",
		Fn:   testFileInfoSearchUsingBooleanOperators,
		Tags: []string{EngineMySql, EnginePostgres, EngineElasticSearch},
	},
	{
		Name: "Should be able to search with combined filters",
		Fn:   testFileInfoSearchUsingCombinedFilters,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should be able to ignore stop words",
		Fn:   testFileInfoSearchIgnoringStopWords,
		Tags: []string{EngineMySql, EngineElasticSearch},
	},
	{
		Name: "Should support search stemming",
		Fn:   testFileInfoSupportStemming,
		Tags: []string{EnginePostgres, EngineElasticSearch},
	},
	{
		Name: "Should support search with wildcards",
		Fn:   testFileInfoSupportWildcards,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should not support search with preceding wildcards",
		Fn:   testFileInfoNotSupportPrecedingWildcards,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should discard a wildcard if it's not placed immediately by text",
		Fn:   testFileInfoSearchDiscardWildcardAlone,
		Tags: []string{EngineAll},
	},
	{
		Name: "Should support terms with dash",
		Fn:   testFileInfoSupportTermsWithDash,
		Tags: []string{EngineAll},
		Skip: true,
	},
	{
		Name: "Should support terms with underscore",
		Fn:   testFileInfoSupportTermsWithUnderscore,
		Tags: []string{EngineMySql, EngineElasticSearch},
	},
	{
		Name: "Should be able to search in deleted/archived channels",
		Fn:   testFileInfoSearchInDeletedOrArchivedChannels,
		Tags: []string{EngineMySql, EnginePostgres},
	},
	{
		Name:        "Should be able to search terms with dashes",
		Fn:          testFileInfoSearchTermsWithDashes,
		Tags:        []string{EngineAll},
		Skip:        true,
		SkipMessage: "Not working",
	},
	{
		Name: "Should be able to search terms with dots",
		Fn:   testFileInfoSearchTermsWithDots,
		Tags: []string{EnginePostgres, EngineElasticSearch},
	},
	{
		Name: "Should be able to search terms with underscores",
		Fn:   testFileInfoSearchTermsWithUnderscores,
		Tags: []string{EngineMySql, EngineElasticSearch},
	},
	{
		Name: "Should be able to combine stemming and wildcards",
		Fn:   testFileInfoSupportStemmingAndWildcards,
		Tags: []string{EngineElasticSearch},
	},
	{
		Name: "Should support wildcard outside quotes",
		Fn:   testFileInfoSupportWildcardOutsideQuotes,
		Tags: []string{EngineElasticSearch},
	},
	{
		Name: "Should not support slash as character separator",
		Fn:   testFileInfoSlashShouldNotBeCharSeparator,
		Tags: []string{EngineMySql, EngineElasticSearch},
	},
	{
		Name: "Should be able to search emails without quoting them",
		Fn:   testFileInfoSearchEmailsWithoutQuotes,
		Tags: []string{EngineElasticSearch},
	},
}

func TestSearchFileInfoStore(t *testing.T, s store.Store, testEngine *SearchTestEngine) {
	th := &SearchTestHelper{
		Store: s,
	}
	err := th.SetupBasicFixtures()
	require.NoError(t, err)
	defer th.CleanFixtures()

	runTestSearch(t, testEngine, searchFileInfoStoreTests, th)
}

func testFileInfoSearchFileInfosIncludingDMs(t *testing.T, th *SearchTestHelper) {
	direct, err := th.createDirectChannel(th.Team.Id, "direct-"+th.Team.Id, "direct-"+th.Team.Id, []*model.User{th.User, th.User2})
	require.NoError(t, err)
	defer th.deleteChannel(direct)

	post, err := th.createPost(th.User.Id, direct.Id, "dm test", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	post2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "dm test", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "dm test filename", "dm contenttest filename", "jpg", "image/jpeg", 0, 1)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "dm other filename", "dm other filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "channel test filename", "channel contenttest filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("by-name", func(t *testing.T) {
		params := &model.SearchParams{Terms: "test"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("by-content", func(t *testing.T) {
		params := &model.SearchParams{Terms: "contenttest"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoSearchFileInfosWithPagination(t *testing.T, th *SearchTestHelper) {
	direct, err := th.createDirectChannel(th.Team.Id, "direct", "direct", []*model.User{th.User, th.User2})
	require.NoError(t, err)
	defer th.deleteChannel(direct)

	post, err := th.createPost(th.User.Id, direct.Id, "dm test", "", model.PostTypeDefault, 10000, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	post2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "dm test", "", model.PostTypeDefault, 20000, false)
	require.NoError(t, err)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "dm test filename", "dm contenttest filename", "jpg", "image/jpeg", 10000, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "dm other filename", "dm other filename", "jpg", "image/jpeg", 20000, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "channel test filename", "channel contenttest filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("by-name", func(t *testing.T) {
		params := &model.SearchParams{Terms: "test"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 1)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)

		results, err = th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 1, 1)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("by-content", func(t *testing.T) {
		params := &model.SearchParams{Terms: "contenttest"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 1)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)

		results, err = th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 1, 1)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
}

func testFileInfoSearchExactPhraseInQuotes(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "channel test 1 2 3 filename", "channel content test 1 2 3 filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "channel test 123 filename", "channel content test 123 filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("by-name", func(t *testing.T) {
		params := &model.SearchParams{Terms: "\"channel test 1 2 3\""}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("by-content", func(t *testing.T) {
		params := &model.SearchParams{Terms: "\"channel content test 1 2 3\""}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
}

func testFileInfoSearchEmailAddresses(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "test email test@test.com", "test email test@content.com", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "test email test2@test.com", "test email test2@content.com", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("by-name", func(t *testing.T) {
		t.Run("Should search email addresses enclosed by quotes", func(t *testing.T) {
			params := &model.SearchParams{Terms: "\"test@test.com\""}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})

		t.Run("Should search email addresses without quotes", func(t *testing.T) {
			params := &model.SearchParams{Terms: "test@test.com"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})
	})
	t.Run("by-content", func(t *testing.T) {
		t.Run("Should search email addresses enclosed by quotes", func(t *testing.T) {
			params := &model.SearchParams{Terms: "\"test@content.com\""}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})

		t.Run("Should search email addresses without quotes", func(t *testing.T) {
			params := &model.SearchParams{Terms: "test@content.com"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})
	})
}

func testFileInfoSearchMarkdownUnderscores(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "_start middle end_ _another_", "_start middle end_ _another_", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should search the start inside the markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "start"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should search a word in the middle of the markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "middle"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should search in the end of the markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "end"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should search inside markdown underscore", func(t *testing.T) {
		params := &model.SearchParams{Terms: "another"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
}

func testFileInfoSearchNonLatinWords(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	t.Run("Should be able to search chinese words", func(t *testing.T) {
		p1, err := th.createFileInfo(th.User.Id, post.Id, "你好", "你好", "jpg", "image/jpeg", 0, 0)
		require.NoError(t, err)
		p2, err := th.createFileInfo(th.User.Id, post.Id, "你", "你", "jpg", "image/jpeg", 0, 0)
		require.NoError(t, err)
		defer th.deleteUserFileInfos(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "你"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		})
		t.Run("Should search two words", func(t *testing.T) {
			params := &model.SearchParams{Terms: "你好"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})
		t.Run("Should search with wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "你*"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 2)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
			th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		})
	})
	t.Run("Should be able to search cyrillic words", func(t *testing.T) {
		p1, err := th.createFileInfo(th.User.Id, post.Id, "слово test", "слово test", "jpg", "image/jpeg", 0, 0)
		require.NoError(t, err)
		defer th.deleteUserFileInfos(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "слово"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})
		t.Run("Should search using wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "слов*"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})
	})

	t.Run("Should be able to search japanese words", func(t *testing.T) {
		p1, err := th.createFileInfo(th.User.Id, post.Id, "本", "本", "jpg", "image/jpeg", 0, 0)
		require.NoError(t, err)
		p2, err := th.createFileInfo(th.User.Id, post.Id, "本木", "本木", "jpg", "image/jpeg", 0, 0)
		require.NoError(t, err)
		defer th.deleteUserFileInfos(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "本"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 2)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
			th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		})
		t.Run("Should search two words", func(t *testing.T) {
			params := &model.SearchParams{Terms: "本木"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		})
		t.Run("Should search with wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "本*"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 2)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
			th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		})
	})

	t.Run("Should be able to search korean words", func(t *testing.T) {
		p1, err := th.createFileInfo(th.User.Id, post.Id, "불", "불", "jpg", "image/jpeg", 0, 0)
		require.NoError(t, err)
		p2, err := th.createFileInfo(th.User.Id, post.Id, "불다", "불다", "jpg", "image/jpeg", 0, 0)
		require.NoError(t, err)
		defer th.deleteUserFileInfos(th.User.Id)

		t.Run("Should search one word", func(t *testing.T) {
			params := &model.SearchParams{Terms: "불"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		})
		t.Run("Should search two words", func(t *testing.T) {
			params := &model.SearchParams{Terms: "불다"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 1)
			th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		})
		t.Run("Should search with wildcard", func(t *testing.T) {
			params := &model.SearchParams{Terms: "불*"}
			results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
			require.NoError(t, err)

			require.Len(t, results.FileInfos, 2)
			th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
			th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		})
	})
}

func testFileInfoSearchAlternativeSpellings(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "Straße test", "Straße test", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "Strasse test", "Strasse test", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	params := &model.SearchParams{Terms: "Straße"}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 2)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)

	params = &model.SearchParams{Terms: "Strasse"}
	results, err = th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 2)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
}

func testFileInfoSearchAlternativeSpellingsAccents(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "café", "café", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "café", "café", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	params := &model.SearchParams{Terms: "café"}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 2)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)

	params = &model.SearchParams{Terms: "café"}
	results, err = th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 2)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)

	params = &model.SearchParams{Terms: "cafe"}
	results, err = th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 0)
}

func testFileInfoSearchOrExcludeFileInfosBySpecificUser(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "test fromuser filename", "test fromuser filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User2.Id, post.Id, "test fromuser filename", "test fromuser filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)
	defer th.deleteUserFileInfos(th.User2.Id)

	params := &model.SearchParams{Terms: "fromuser", FromUsers: []string{th.User.Id}}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 1)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
}

func testFileInfoSearchOrExcludeFileInfosInChannel(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "test fromuser filename", "test fromuser filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post2.Id, "test fromuser filename", "test fromuser filename", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)
	defer th.deleteUserFileInfos(th.User2.Id)

	params := &model.SearchParams{Terms: "fromuser", InChannels: []string{th.ChannelBasic.Id}}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 1)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
}

func testFileInfoSearchOrExcludeFileInfosInDMGM(t *testing.T, th *SearchTestHelper) {
	direct, err := th.createDirectChannel(th.Team.Id, "direct", "direct", []*model.User{th.User, th.User2})
	require.NoError(t, err)
	defer th.deleteChannel(direct)

	group, err := th.createGroupChannel(th.Team.Id, "test group", []*model.User{th.User, th.User2})
	require.NoError(t, err)
	defer th.deleteChannel(group)

	post1, err := th.createPost(th.User.Id, direct.Id, "test fromuser", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	post2, err := th.createPost(th.User2.Id, group.Id, "test fromuser 2", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)
	defer th.deleteUserPosts(th.User2.Id)

	p1, err := th.createFileInfo(th.User.Id, post1.Id, "test fromuser", "test fromuser", "jpg", "image/jpg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User2.Id, post2.Id, "test fromuser 2", "test fromuser 2", "jpg", "image/jpg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)
	defer th.deleteUserFileInfos(th.User2.Id)

	t.Run("Should be able to search in both DM and GM channels", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "fromuser",
			InChannels: []string{direct.Id, group.Id},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("Should be able to search only in DM channel", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "fromuser",
			InChannels: []string{direct.Id},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should be able to search only in GM channel", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "fromuser",
			InChannels: []string{group.Id},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoSearchOrExcludeByExtensions(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "test", "test", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "test", "test", "png", "image/png", 0, 0)
	require.NoError(t, err)
	p3, err := th.createFileInfo(th.User.Id, post.Id, "test", "test", "bmp", "image/bmp", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Search by one extension", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "test",
			InChannels: []string{th.ChannelBasic.Id},
			Extensions: []string{"jpg"},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search by multiple extensions", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "test",
			InChannels: []string{th.ChannelBasic.Id},
			Extensions: []string{"jpg", "bmp"},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Search excluding one extension", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:              "test",
			InChannels:         []string{th.ChannelBasic.Id},
			ExcludedExtensions: []string{"jpg"},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Search excluding multiple extensions", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:              "test",
			InChannels:         []string{th.ChannelBasic.Id},
			ExcludedExtensions: []string{"jpg", "bmp"},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoFilterFilesInSpecificDate(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	creationDate := model.GetMillisForTime(time.Date(2020, 03, 22, 12, 0, 0, 0, time.UTC))
	p1, err := th.createFileInfo(th.User.Id, post1.Id, "test in specific date", "test in specific date", "jpg", "image/jpeg", creationDate, 0)
	require.NoError(t, err)
	creationDate2 := model.GetMillisForTime(time.Date(2020, 03, 23, 0, 0, 0, 0, time.UTC))
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "test in the present", "test in the present", "jpg", "image/jpeg", creationDate2, 0)
	require.NoError(t, err)
	creationDate3 := model.GetMillisForTime(time.Date(2020, 03, 21, 23, 59, 59, 0, time.UTC))
	p3, err := th.createFileInfo(th.User.Id, post1.Id, "test in the present", "test in the present", "jpg", "image/jpeg", creationDate3, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should be able to search posts on date", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:  "test",
			OnDate: "2020-03-22",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
	t.Run("Should be able to exclude posts on date", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:        "test",
			ExcludedDate: "2020-03-22",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})
}

func testFileInfoFilterFilesBeforeSpecificDate(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	creationDate := model.GetMillisForTime(time.Date(2020, 03, 01, 12, 0, 0, 0, time.UTC))
	p1, err := th.createFileInfo(th.User.Id, post1.Id, "test in specific date", "test in specific date", "jpg", "image/jpeg", creationDate, 0)
	require.NoError(t, err)
	creationDate2 := model.GetMillisForTime(time.Date(2020, 03, 22, 23, 59, 59, 0, time.UTC))
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "test in specific date 2", "test in specific date 2", "jpg", "image/jpeg", creationDate2, 0)
	require.NoError(t, err)
	creationDate3 := model.GetMillisForTime(time.Date(2020, 03, 26, 16, 55, 0, 0, time.UTC))
	p3, err := th.createFileInfo(th.User.Id, post1.Id, "test in the present", "test in the present", "jpg", "image/jpeg", creationDate3, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should be able to search posts before a date", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "test",
			BeforeDate: "2020-03-23",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("Should be able to exclude posts before a date", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:              "test",
			ExcludedBeforeDate: "2020-03-23",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})
}

func testFileInfoFilterFilesAfterSpecificDate(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	creationDate := model.GetMillisForTime(time.Date(2020, 03, 01, 12, 0, 0, 0, time.UTC))
	p1, err := th.createFileInfo(th.User.Id, post1.Id, "test in specific date", "test in specific date", "jpg", "image/jpeg", creationDate, 0)
	require.NoError(t, err)
	creationDate2 := model.GetMillisForTime(time.Date(2020, 03, 22, 23, 59, 59, 0, time.UTC))
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "test in specific date 2", "test in specific date 2", "jpg", "image/jpeg", creationDate2, 0)
	require.NoError(t, err)
	creationDate3 := model.GetMillisForTime(time.Date(2020, 03, 26, 16, 55, 0, 0, time.UTC))
	p3, err := th.createFileInfo(th.User.Id, post1.Id, "test in the present", "test in the present", "jpg", "image/jpeg", creationDate3, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should be able to search posts after a date", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:     "test",
			AfterDate: "2020-03-23",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Should be able to exclude posts after a date", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:             "test",
			ExcludedAfterDate: "2020-03-23",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoFilterFilesWithATerm(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post1.Id, "one two three", "one two three", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "one four five six", "one four five six", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post1.Id, "one seven eight nine", "one seven eight nine", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should exclude terms", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:         "one",
			ExcludedTerms: "five eight",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should exclude quoted terms", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:         "one",
			ExcludedTerms: "\"eight nine\"",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoSearchUsingBooleanOperators(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "one two three message", "one two three message", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "two messages", "two messages", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "another message", "another message", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should search posts using OR operator", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:   "one two",
			OrTerms: true,
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("Should search posts using AND operator", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:   "one two",
			OrTerms: false,
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
}

func testFileInfoSearchUsingCombinedFilters(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	creationDate := model.GetMillisForTime(time.Date(2020, 03, 01, 12, 0, 0, 0, time.UTC))
	p1, err := th.createFileInfo(th.User.Id, post2.Id, "one two three message", "one two three message", "jpg", "image/jpeg", creationDate, 0)
	require.NoError(t, err)
	creationDate2 := model.GetMillisForTime(time.Date(2020, 03, 10, 12, 0, 0, 0, time.UTC))
	p2, err := th.createFileInfo(th.User2.Id, post2.Id, "two messages", "two messages", "jpg", "image/jpeg", creationDate2, 0)
	require.NoError(t, err)
	creationDate3 := model.GetMillisForTime(time.Date(2020, 03, 20, 12, 0, 0, 0, time.UTC))
	p3, err := th.createFileInfo(th.User.Id, post1.Id, "two another message", "two another message", "jpg", "image/jpeg", creationDate3, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)
	defer th.deleteUserFileInfos(th.User2.Id)

	t.Run("Should search combining from user and in channel filters", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:      "two",
			FromUsers:  []string{th.User2.Id},
			InChannels: []string{th.ChannelPrivate.Id},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("Should search combining excluding users and in channel filters", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:         "two",
			ExcludedUsers: []string{th.User2.Id},
			InChannels:    []string{th.ChannelPrivate.Id},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should search combining excluding dates and in channel filters", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:              "two",
			ExcludedBeforeDate: "2020-03-09",
			ExcludedAfterDate:  "2020-03-11",
			InChannels:         []string{th.ChannelPrivate.Id},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
	t.Run("Should search combining excluding dates and in channel filters", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:            "two",
			AfterDate:        "2020-03-11",
			ExcludedChannels: []string{th.ChannelPrivate.Id},
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})
}

func testFileInfoSearchIgnoringStopWords(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "the search for a bunch of stop words", "the search for a bunch of stop words", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "the objective is to avoid a bunch of stop words", "the objective is to avoid a bunch of stop words", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p3, err := th.createFileInfo(th.User.Id, post.Id, "in the a on to where you", "in the a on to where you", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p4, err := th.createFileInfo(th.User.Id, post.Id, "where is the car?", "where is the car?", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should avoid stop word 'the'", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "the search",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should avoid stop word 'a'", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "a avoid",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("Should avoid stop word 'in'", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "in where you",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Should avoid stop words 'where', 'is' and 'the'", func(t *testing.T) {
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{{Terms: "is the car"}}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)
		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p4.Id, results.FileInfos)
	})

	t.Run("Should remove all terms and return empty list", func(t *testing.T) {
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{{Terms: "is the"}}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)
		require.Empty(t, results.FileInfos)
	})
}

func testFileInfoSupportStemming(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "search post", "search post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "searching post", "searching post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "another post", "another post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	params := &model.SearchParams{
		Terms: "search",
	}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 2)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
}

func testFileInfoSupportWildcards(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "search post", "search post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "searching", "searching", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "another post", "another post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Simple wildcard-only search", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "search*",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("Wildcard search with another term placed after", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "sear* post",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
}

func testFileInfoNotSupportPrecedingWildcards(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	_, err = th.createFileInfo(th.User.Id, post.Id, "search post", "search post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "searching post", "searching post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "another post", "another post", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	params := &model.SearchParams{
		Terms: "*earch",
	}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 0)
}

func testFileInfoSearchDiscardWildcardAlone(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "qwerty", "qwerty", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "qwertyjkl", "qwertyjkl", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	params := &model.SearchParams{
		Terms: "qwerty *",
	}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 1)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
}

func testFileInfoSupportTermsWithDash(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "search term-with-dash", "search term-with-dash", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "searching term with dash", "searching term with dash", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should search terms with dash", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "term-with-dash",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should search terms with dash using quotes", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "\"term-with-dash\"",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
}

func testFileInfoSupportTermsWithUnderscore(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "search term_with_underscore", "search term_with_underscore", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "searching term with underscore", "searching term with underscore", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should search terms with underscore", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "term_with_underscore",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Should search terms with underscore using quotes", func(t *testing.T) {
		params := &model.SearchParams{
			Terms: "\"term_with_underscore\"",
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})
}

func testFileInfoSearchInDeletedOrArchivedChannels(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelDeleted.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)
	post2, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)
	post3, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post1.Id, "message in deleted channel", "message in deleted channel", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "message in regular channel", "message in regular channel", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p3, err := th.createFileInfo(th.User.Id, post3.Id, "message in private channel", "message in private channel", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Doesn't include posts in deleted channels", func(t *testing.T) {
		params := &model.SearchParams{Terms: "message", IncludeDeletedChannels: false}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Include posts in deleted channels", func(t *testing.T) {
		params := &model.SearchParams{Terms: "message", IncludeDeletedChannels: true}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 3)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Include posts in deleted channels using multiple terms", func(t *testing.T) {
		params := &model.SearchParams{Terms: "message channel", IncludeDeletedChannels: true}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 3)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Include posts in deleted channels using multiple OR terms", func(t *testing.T) {
		params := &model.SearchParams{
			Terms:                  "message channel",
			IncludeDeletedChannels: true,
			OrTerms:                true,
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 3)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("All IncludeDeletedChannels params should have same value if multiple SearchParams provided", func(t *testing.T) {
		params1 := &model.SearchParams{
			Terms:                  "message channel",
			IncludeDeletedChannels: true,
		}
		params2 := &model.SearchParams{
			Terms:                  "#hashtag",
			IncludeDeletedChannels: false,
		}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params1, params2}, th.User.Id, th.Team.Id, 0, 20)
		require.Nil(t, results)
		require.Error(t, err)
	})
}

func testFileInfoSearchTermsWithDashes(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "message with-dash-term", "message with-dash-term", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "message with dash term", "message with dash term", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Search for terms with dash", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with-dash-term"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for terms with quoted dash", func(t *testing.T) {
		params := &model.SearchParams{Terms: "\"with-dash-term\""}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for multiple terms with one having dash", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with-dash-term message"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for multiple OR terms with one having dash", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with-dash-term message", OrTerms: true}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoSearchTermsWithDots(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "message with.dots.term", "message with.dots.term", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "message with dots term", "message with dots term", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Search for terms with dots", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with.dots.term"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for terms with quoted dots", func(t *testing.T) {
		params := &model.SearchParams{Terms: "\"with.dots.term\""}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for multiple terms with one having dots", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with.dots.term message"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for multiple OR terms with one having dots", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with.dots.term message", OrTerms: true}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoSearchTermsWithUnderscores(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "message with_underscores_term", "message with_underscores_term", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post.Id, "message with underscores term", "message with underscores term", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Search for terms with underscores", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with_underscores_term"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for terms with quoted underscores", func(t *testing.T) {
		params := &model.SearchParams{Terms: "\"with_underscores_term\""}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for multiple terms with one having underscores", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with_underscores_term message"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
	})

	t.Run("Search for multiple OR terms with one having underscores", func(t *testing.T) {
		params := &model.SearchParams{Terms: "with_underscores_term message", OrTerms: true}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})
}

func testFileInfoSupportStemmingAndWildcards(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)

	defer th.deleteUserPosts(th.User.Id)
	p1, err := th.createFileInfo(th.User.Id, post1.Id, "approve", "approve", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "approved", "approved", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p3, err := th.createFileInfo(th.User.Id, post2.Id, "approvedz", "approvedz", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should stem appr", func(t *testing.T) {
		params := &model.SearchParams{Terms: "appr*"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 3)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})

	t.Run("Should stem approve", func(t *testing.T) {
		params := &model.SearchParams{Terms: "approve*"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p3.Id, results.FileInfos)
	})
}

func testFileInfoSupportWildcardOutsideQuotes(t *testing.T, th *SearchTestHelper) {
	post1, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)
	post2, err := th.createPost(th.User.Id, th.ChannelPrivate.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)

	p1, err := th.createFileInfo(th.User.Id, post1.Id, "hello world", "hello world", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	p2, err := th.createFileInfo(th.User.Id, post2.Id, "hell or heaven", "hell or heaven", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	t.Run("Should return results without quotes", func(t *testing.T) {
		params := &model.SearchParams{Terms: "hell*"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 2)
		th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

	t.Run("Should return just one result with quotes", func(t *testing.T) {
		params := &model.SearchParams{Terms: "\"hell\"*"}
		results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
		require.NoError(t, err)

		require.Len(t, results.FileInfos, 1)
		th.checkFileInfoInSearchResults(t, p2.Id, results.FileInfos)
	})

}

func testFileInfoSlashShouldNotBeCharSeparator(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "alpha/beta gamma, theta", "alpha/beta gamma, theta", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	params := &model.SearchParams{Terms: "gamma"}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 1)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)

	params = &model.SearchParams{Terms: "beta"}
	results, err = th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 1)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)

	params = &model.SearchParams{Terms: "alpha"}
	results, err = th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 1)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
}

func testFileInfoSearchEmailsWithoutQuotes(t *testing.T, th *SearchTestHelper) {
	post, err := th.createPost(th.User.Id, th.ChannelBasic.Id, "testmessage", "", model.PostTypeDefault, 0, false)
	require.NoError(t, err)
	defer th.deleteUserPosts(th.User.Id)

	p1, err := th.createFileInfo(th.User.Id, post.Id, "message test@test.com", "message test@test.com", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	_, err = th.createFileInfo(th.User.Id, post.Id, "message test2@test.com", "message test2@test.com", "jpg", "image/jpeg", 0, 0)
	require.NoError(t, err)
	defer th.deleteUserFileInfos(th.User.Id)

	params := &model.SearchParams{Terms: "test@test.com"}
	results, err := th.Store.FileInfo().Search([]*model.SearchParams{params}, th.User.Id, th.Team.Id, 0, 20)
	require.NoError(t, err)

	require.Len(t, results.FileInfos, 1)
	th.checkFileInfoInSearchResults(t, p1.Id, results.FileInfos)
}
