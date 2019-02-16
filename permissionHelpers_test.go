package main

import "testing"

func TestHasPermission(t *testing.T) {
	p := []struct {
		UserPermission  int
		PermissionValue int
		Result          bool
	}{
		{1, 1, true},
		{1, 2, false},
		{2, 1, false},
		{3, 1, true},
		{1, 3, false},
	}
	for _, e := range p {
		if e.Result != HasPermission(e.UserPermission, e.PermissionValue) {
			t.Errorf("HasPermission was incorrect with %d as UserPermission and %d as PermissionValue.", e.UserPermission, e.PermissionValue)
		}
	}
}

func TestSetPermission(t *testing.T) {
	p := []struct {
		UserPermission  int
		PermissionValue int
		Result          int
	}{
		{0, 1, 1},
		{1, 2, 3},
		{1, 1, 1},
		{1, 0, 1},
		{4, 1, 5},
	}
	for _, e := range p {
		res := SetPermission(e.UserPermission, e.PermissionValue)
		if e.Result != res {
			t.Errorf("SetPermission was incorrect with %d as UserPermission and %d as PermissionValue. Result was %d.", e.UserPermission, e.PermissionValue, res)
		}
	}
}

func TestRemovePermission(t *testing.T) {
	p := []struct {
		UserPermission  int
		PermissionValue int
		Result          int
	}{
		{0, 1, 0},
		{1, 2, 1},
		{3, 1, 2},
		{5, 0, 5},
		{5, 1, 4},
		{1, 1, 0},
	}
	for _, e := range p {
		res := RemovePermission(e.UserPermission, e.PermissionValue)
		if e.Result != res {
			t.Errorf("RemovePermission was incorrect with %d as UserPermission and %d as PermissionValue. Result was %d.", e.UserPermission, e.PermissionValue, res)
		}
	}
}
