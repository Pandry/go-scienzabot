package main

//TODO: explain how bitwise logic works

//HasPermission checks if the given user permission has the permission necessary to do an action (passed as permission in the 2nd arg)
func HasPermission(userPermission int, permissionToCheck int) bool {
	if userPermission&permissionToCheck == permissionToCheck {
		return true
	}
	return false
}

//SetPermission returns the new user permission value after assigning the new user.
//Needs the initial user permission and the permission we want to add
func SetPermission(userPermission int, permissionToSet int) int {
	return userPermission | permissionToSet
}

//RemovePermission returns a user's new permission, removing the permissiopn passed as 2nd arg
func RemovePermission(userPermission int, permissionToRemove int) int {
	return userPermission &^ permissionToRemove
}
