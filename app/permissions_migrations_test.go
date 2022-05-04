// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cjdelisle/matterfoss-server/v6/model"
)

func TestApplyPermissionsMap(t *testing.T) {
	tt := []struct {
		Name           string
		RoleMap        map[string]map[string]bool
		TranslationMap permissionsMap
		ExpectedResult []string
	}{
		{
			"Split existing",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{On: permissionExists("test2"), Add: []string{"test4", "test5"}}},
			[]string{"test1", "test2", "test3", "test4", "test5"},
		},
		{
			"Remove existing",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{On: permissionExists("test2"), Remove: []string{"test2"}}},
			[]string{"test1", "test3"},
		},
		{
			"Rename existing",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{On: permissionExists("test2"), Add: []string{"test5"}, Remove: []string{"test2"}}},
			[]string{"test1", "test3", "test5"},
		},
		{
			"Remove when other not exists",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{On: permissionNotExists("test5"), Remove: []string{"test2"}}},
			[]string{"test1", "test3"},
		},
		{
			"Add when at least one exists",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:  permissionOr(permissionExists("test5"), permissionExists("test3")),
				Add: []string{"test4"},
			}},
			[]string{"test1", "test2", "test3", "test4"},
		},
		{
			"Add when all exists",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:  permissionAnd(permissionExists("test1"), permissionExists("test2")),
				Add: []string{"test4"},
			}},
			[]string{"test1", "test2", "test3", "test4"},
		},
		{
			"Not add when one in the and not exists",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:  permissionAnd(permissionExists("test1"), permissionExists("test5")),
				Add: []string{"test4"},
			}},
			[]string{"test1", "test2", "test3"},
		},
		{
			"Not Add when none on the or exists",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:  permissionOr(permissionExists("test7"), permissionExists("test9")),
				Add: []string{"test4"},
			}},
			[]string{"test1", "test2", "test3"},
		},
		{
			"When the role matches",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:  isRole("system_admin"),
				Add: []string{"test4"},
			}},
			[]string{"test1", "test2", "test3", "test4"},
		},
		{
			"When the role doesn't match",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:  isRole("system_user"),
				Add: []string{"test4"},
			}},
			[]string{"test1", "test2", "test3"},
		},
		{
			"Remove a permission conditional on another role having it, success case",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test3": true,
				},
				"other_role": {
					"test4": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:     onOtherRole("other_role", permissionExists("test4")),
				Remove: []string{"test1"},
			}},
			[]string{"test2", "test3"},
		},
		{
			"Remove a permission conditional on another role having it, failure case",
			map[string]map[string]bool{
				"system_admin": {
					"test1": true,
					"test2": true,
					"test4": true,
				},
				"other_role": {
					"test1": true,
				},
			},
			permissionsMap{permissionTransformation{
				On:     onOtherRole("other_role", permissionExists("test4")),
				Remove: []string{"test1"},
			}},
			[]string{"test1", "test2", "test4"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			result := applyPermissionsMap(&model.Role{Name: "system_admin"}, tc.RoleMap, tc.TranslationMap)
			sort.Strings(result)
			assert.Equal(t, tc.ExpectedResult, result)
		})
	}
}
