package consts

//Version is the bot version.
//This number should be incremented for every release
const Version = "0.0.1-g α"

//InitSQLString is the initialization query, to run once at the bot startup
const InitSQLString = `
CREATE TABLE IF NOT EXISTS 'Users' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Nickname'  TEXT NOT NULL UNIQUE,
	'Biography'  TEXT,
	'Status'  INTEGER DEFAULT 0,
	'LastSeen'  TEXT DEFAULT '0000-00-00 00:00:00'
);

CREATE TABLE IF NOT EXISTS 'Groups' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Title'  TEXT NOT NULL,
	'Ref'	TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS 'Channels' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Group'  INTEGER NOT NULL,
	'Name'	TEXT NOT NULL,
	'Ref'	TEXT NOT NULL,	
	FOREIGN KEY('Group') REFERENCES Groups('ID')
);

CREATE TABLE IF NOT EXISTS 'Permissions' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'User'	INTEGER NOT NULL,
	'Group'	INTEGER NOT NULL,
	'Permissions'  INTEGER DEFAULT 0,
	FOREIGN KEY('User') REFERENCES Users('ID'),
	FOREIGN KEY('Group') REFERENCES Groups('ID'),
	CONSTRAINT con_perm_user_group_perm_unique UNIQUE ('User','Group','Permissions')
);

CREATE TABLE IF NOT EXISTS 'Lists' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'Name'  TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS 'Bookmarks' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'User'  INTEGER NOT NULL,
	'Group'  INTEGER NOT NULL,
	'MessageID'	INTEGER NOT NULL,
	'Alias'	TEXT,
	FOREIGN KEY('User') REFERENCES Users('ID'),
	FOREIGN KEY('Group') REFERENCES Groups('ID'),
	CONSTRAINT con_bookm_user_group_unique UNIQUE ('User','Group')

);

CREATE TABLE IF NOT EXISTS 'Subscriptions' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'List'  INTEGER NOT NULL,
	'User'  INTEGER NOT NULL,
	FOREIGN KEY('User') REFERENCES Users('ID'),
	FOREIGN KEY('List') REFERENCES 'Lists'('ID'),
	CONSTRAINT con_subs_user_list_unique UNIQUE ('User','List')
);

CREATE TABLE IF NOT EXISTS 'MessageCount' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'User'  INTEGER NOT NULL,
	'Group'  INTEGER NOT NULL,
	'MessageCount'  INTEGER NOT NULL,
	FOREIGN KEY('User') REFERENCES Users('ID'),
	FOREIGN KEY('Group') REFERENCES Groups('ID'),
	CONSTRAINT con_msgcoubt_user_group_unique UNIQUE ('User','Group')
);

CREATE TABLE IF NOT EXISTS 'Strings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL,
	'Value'	TEXT DEFAULT 'Not implemented',
	'Locale'	TEXT DEFAULT 'it',
	'Group'	INTEGER NOT NULL,
	FOREIGN KEY('Group') REFERENCES Groups('ID'),
	CONSTRAINT con_strings_key_group_locale_unique UNIQUE ('Key','Group','Locale')
);

CREATE TABLE IF NOT EXISTS 'Settings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL,
	'Value'	TEXT DEFAULT 'Not implemented',
	'Group'	INTEGER NOT NULL,
	FOREIGN KEY('Group') REFERENCES Groups('ID'),
	CONSTRAINT con_setting_key_group_unique UNIQUE ('Key','Group')
);

CREATE TABLE IF NOT EXISTS 'BotSettings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL UNIQUE,
	'Value'	TEXT DEFAULT 'Not implemented'
);


`

//User status constants
//UserStatusActive Is assigned when the user is created and it's immediately active
//UserStatusWaitingForBio Is assigned when the user wants to edit its biography
//UserStatusWaitingForList Is assigned when the user wants to create a new list, and as long as it's setted, it will create new lists
//UserStatusBanned Is assigned when the user is banned. Once the user is banned, the bot will not consider anymore its commands
const (
	UserStatusActive = iota
	UserStatusWaitingForBio
	UserStatusWaitingForList
	UserStatusBanned = -99
)

//UserPermissionAdmin is the admin privilege. It allows to do admin stuff
//TODO determine what admins can do and what not.
//UserPermissionCanAddAdmins is the privilege that allows an admin to add another one
//UserPermissionCanRemoveAdmins is the privilege that allows an admin to remove another one - dangerous!
//UserPermissionCanForwardToChannel is the privilege that allows an user to forward to a channel
// a message linking a message she's replying to - only for supergroups
//UserPermissionCanCreateList is the privilege that allows an user to create a list
//UserPermissionCanRemoveList is the privilege that allows an user to remove a list
const (
	UserPermissionAdmin = 1 << iota
	UserPermissionCanAddAdmins
	UserPermissionCanRemoveAdmins
	UserPermissionCanForwardToChannel
	UserPermissionCanCreateList
	UserPermissionCanRemoveList
)
