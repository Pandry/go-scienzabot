package consts

//InitSQLString is the initialization query, to run once at the bot startup
const InitSQLString = `
/*
The Users table is supposed to contain all the users subscribed to the bot
*/
CREATE TABLE IF NOT EXISTS 'Users' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Nickname'  TEXT UNIQUE,
	'Biography'  TEXT,
	'Status'  INTEGER NOT NULL DEFAULT 0,
	'LastSeen'  TIMESTAMP,
	'RegisterDate' TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

/*
The Groups table is supposed to contain all the groups the bot is added to.
Status refers if the bot was kicked form the group, not used atm
*/
CREATE TABLE IF NOT EXISTS 'Groups' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Title'  TEXT NOT NULL,
	'Ref'	TEXT NOT NULL,
	'Locale'	TEXT DEFAULT '` + DefaultLocale + `',
	'Status'	INTEGER NOT NULL DEFAULT 0
);

/*
The Channels table is supposed to contain the channels that admins may want to use to forward messages
from a group to a channel referring to a particular message
*/
CREATE TABLE IF NOT EXISTS 'Channels' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'GroupID'  INTEGER NOT NULL,
	'Name'	TEXT NOT NULL,
	'Ref'	TEXT NOT NULL,	
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_channels_channel_group__unique UNIQUE ('ID','GroupID')
);

/*
The Permissions table is supposed to contain the permissions for each user in each group.
*/
CREATE TABLE IF NOT EXISTS 'Permissions' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'UserID'	INTEGER NOT NULL,
	'GroupID'	INTEGER,
	'Permission' INTEGER DEFAULT 0,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_perm_user_group_perm_unique UNIQUE ('UserID','GroupID','Permission')
);

/*
The Lists table is suopposed to contain the lists where a user can subscribe to.
Such list should be group-dependent (if not specified otherwise on the status field, that shouold be based on a bit-based flag)
The status is not used yet
*/
CREATE TABLE IF NOT EXISTS 'Lists' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'Name'  TEXT NOT NULL UNIQUE,
	'GroupID'	INTEGER NOT NULL,
	'Properties'  INTEGER DEFAULT 0,
	'CreationDate' TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	'LatestInvocationDate' TIMESTAMP,
	FOREIGN KEY('GroupID') REFERENCES Groups('ID')
);

/*
The Bookmarks table is used to when a user wants to save a message for a future reference.
The bot will in fact save the group and the message, and will bind it to a user.
The bot will also save a copy of the message content (when possible).
Deletion of a row should be impossibilitated to a user
TODO: Remembere to check if the message still exists
*/
CREATE TABLE IF NOT EXISTS 'Bookmarks' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'UserID'  INTEGER NOT NULL,
	'GroupID'  INTEGER NOT NULL,
	'MessageID'	INTEGER NOT NULL,
	'Alias'	TEXT,
	'Status' INTEGER DEFAULT 0,
	'MessageContent' TEXT, 
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_bookm_user_group_unique UNIQUE ('UserID','GroupID')

);

/*
The Subscriptions table is used to subscribe a specific user to a "list" where he belongs
*/
CREATE TABLE IF NOT EXISTS 'Subscriptions' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'ListID'  INTEGER NOT NULL,
	'UserID'  INTEGER NOT NULL,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('ListID') REFERENCES 'Lists'('ID'),
	CONSTRAINT con_subs_user_list_unique UNIQUE ('UserID','ListID')
);

/*
The MessageCount table is used to count the message of each user in the various groups
This allows the bot to count the message of a specific user on a multitude of groups
*/
CREATE TABLE IF NOT EXISTS 'MessageCount' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'UserID'  INTEGER NOT NULL,
	'GroupID'  INTEGER NOT NULL,
	'MessageCount'  INTEGER NOT NULL,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_msgcoubt_user_group_unique UNIQUE ('UserID','GroupID')
);

/*
The Strings table will contain all the strings about a specific group 
(in fact it's group-dependent).
Such strings could be something like a welcome message, a help message an so on...
*/
CREATE TABLE IF NOT EXISTS 'Strings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL,
	'Value'	TEXT DEFAULT 'Not implemented',
	'Locale'	TEXT DEFAULT '` + DefaultLocale + `',
	'GroupID'	INTEGER NOT NULL,
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_strings_key_group_locale_unique UNIQUE ('Key','GroupID','Locale')
);

/*
The Settings table will contain all the settings about a specific group 
(in fact it's group-dependent).
An example could be the status (on/off) of a specific function of the bot
*/
CREATE TABLE IF NOT EXISTS 'Settings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL,
	'Value'	TEXT DEFAULT 'Not implemented',
	'GroupID'	INTEGER NOT NULL,
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_setting_key_group_unique UNIQUE ('Key','GroupID')
);

/*
The BotSettings table will contain all the settings of the bot
Such settings could be things like the default locale
TODO: evaluate UNIQUE on Key
TODO: evaluate removal of this table
*/
CREATE TABLE IF NOT EXISTS 'BotSettings' (
	'ID'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'	TEXT NOT NULL UNIQUE,
	'Value'	TEXT DEFAULT 'Not implemented'
);

/*
The BotStrings table will contain all the strings to be used from the bot, like the "cancel" text and so on...
As the constraint shows, there can only be a pair of key-locale per table(we can't have 2 way of saying the same thing
in the same language; which one should we take?) 
*/
CREATE TABLE IF NOT EXISTS 'BotStrings' (
	'ID'		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Key'		TEXT NOT NULL UNIQUE,
	'Value'		TEXT DEFAULT 'Not implemented',
	'Locale'	TEXT DEFAULT '` + DefaultLocale + `',
	CONSTRAINT con_botstrings_key_locale_unique UNIQUE ('Key','Locale')
);

/*
The Log table should contains the log of the bot, like a new user subscribed, a help message requested,
a non-working command received and so on.
Also crashes should be documented.
The Event field is supposed to be a human-readable value
The RelatedUser field is supposed to contain the user in the context who sent or triggered the action
The RelatedGroup field should contain the group where the action triggered
The Message field should be a human-readble string that should say what happened (moreless)
The UpdateValue field should contain the update raw string from telegram
The Error field should contain the string reported by the funcion that throwed the error (if any)
The Severity field should indicate he gravity of the event
The Date field is the date when the event occourred
*/
CREATE TABLE IF NOT EXISTS 'Log' (
	'ID'		INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Event'		TEXT NOT NULL,
	'RelatedUser'		INTEGER,
	'RelatedGroup'		INTEGER,
	'Message'		TEXT,
	'UpdateValue'		TEXT,
	'Error'		TEXT,
	'Severity'		INTEGER,
	'Date'	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

/*
The BotAdministrators table will contain the bot aministrator, which will manage the bot
*/
CREATE TABLE IF NOT EXISTS 'BotAdministrators' (
	'UserID'			TEXT NOT NULL PRIMARY KEY,
	'Permissions'	TEXT NOT NULL DEFAULT 0,
	FOREIGN KEY('UserID') REFERENCES Users('ID')
);

-- Inserting the default locale in DB
INSERT OR IGNORE INTO BotSettings (Key, Value ) VALUES ( "DefaultLocale", "'` + DefaultLocale + `'" );

-- Inserting Pandry and AndreaIdini as users
INSERT OR IGNORE INTO Users (ID, Nickname) VALUES (14092073, "Pandry"), (44917659, "AndreaIdini");

-- Inserting 							     Pandry's ID and  Idini's one as bot administrators
INSERT OR IGNORE INTO BotAdministrators (User) VALUES (14092073),     (44917659);


`
