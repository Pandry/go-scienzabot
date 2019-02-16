package main

// Bitwise operations introduction
//   Bitwise operations are operations that can be performed on a bit level
//   It basically works on 0 and 1, and it's highy efficent
//   It permits to use operators like AND, OR, NOT, XOR etc...
//   It is highly efficent
//   In this bot we are assigning to each bit of an integer a specific permission (e.g. administrator, moderator)
//
// HasPermission function
//   The HasPermission function does an AND bitwise operation with the current permission.
//   For instance, think to the "create list" permission as the bit vith the value of 1 in the 0010(2) number.
//   Another role could be the "Share to channel", with the value 1 in the 00000100(4) number.
//   A user could be both autorized to share to a channel and create a list, so it could have a permission value of 0110(6).
//   If we want to see if a user has a specific permission, we can do and AND operation with the permission of the user and the
//      permission we want to check.
//   The and operation compares 2 bit values, and returns 1, only if both the values are one, so here's a "demostration" :
//
//   | 0 0 1 0 | AND		(The permission we want to check)
//   | 0 1 1 0 |			(The permission of the user)
//   |---------|
//   | 0 0 1 0 |
//
//   Basically, the only way the operation can return a value, is if the user has the permission bit setted to 1.
//   In such case, the value the operation will return will be the value of the permission itself, otherwise it
//      will return 0.
//   Comparing the result, we can see if the the user has the permission we requested or not.

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
