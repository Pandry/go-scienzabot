package database

/*
TODO: Can be done later - channels are not implemented/required yet

CREATE TABLE IF NOT EXISTS 'Channels' (
	'ID'  INTEGER NOT NULL PRIMARY KEY,
	'Group'  INTEGER NOT NULL,
	'Name'	TEXT NOT NULL,
	'Ref'	TEXT NOT NULL,
	FOREIGN KEY('Group') REFERENCES Groups('ID'),
	CONSTRAINT con_channels_channel_group__unique UNIQUE ('ID','Group')
);
*/
