package consts

//This package is supposed to have all the values that remain constant.
// An exaple could be the initiazization query, or the admin permission value (Go does not provide enums)

//Version is the bot version.
//This number should be incremented for every release
//Also, this shoud be inserted in the database and be incremented from it
const Version = "0.3.00 α 210330135500"

//DefaultLocale identifies the default locale of the bot
const DefaultLocale = "it"

//TimeFormatString is the format the time is parsed to and from string
//The T and the Z was inserted because when pulling the string from the
//  database the chars was there and was influenciating the time parsing
//Also the UTC was added, since telegram returns a UTC UNIX timestamp
const TimeFormatString = "2006-01-02T15:04:05Z-0700UTC"

//MaximumInlineKeyboardRows is the maximum number of rows the inline keyboard can have
const MaximumInlineKeyboardRows = 7

//ListRegex is the expression that determines if a list name is valid or not
const ListRegex = "^[a-z\\-_]{1,30}$"

//MaximumMessageLength is the maximun length a message can have
//  this constant is used only in the send long message function extended in the bot itself
//  nad is implemented in the embtypes package
const MaximumMessageLength = 4000

//MaximumMessageLengthMargin is the margin given to the maximum length
//Basically it means MaximumMessageLength +- MaximumMessageLengthMargin
//  This because the bot tries to plit the message bsed on muktuple chars
//  and tries to find the most adapt char
const MaximumMessageLengthMargin = 90

//const MessageSplitter = []string{"\n\n\n", "\n\n", "\n", " "}

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
//TODO: determine what admins can do and what not.
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
	UserPermissionGroupAdmin
	UserPermissionListBanned
)

//ListPropertyGroupIndipendent is the property that defines a list that is available in any group
//ListPropertyGroupLocked is the property that defines a list that is joinable only via invite
const (
	ListPropertyGroupIndipendent = 1 << iota
	ListPropertyGroupLocked
)

//GroupActive is the property that tells that a group is active
//GroupBanned is the property that makes a group blocked and does not permit it to be enabled until it's unlocked
const (
	GroupActive = 1 << iota
	GroupBanned
)
