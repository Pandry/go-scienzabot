package database

//User represent the respective table in the database
type User struct {
	ID        int
	Nickname  string
	Biography string
	Status    int
	LastSeen  string
}

//Group represent the respective table in the database
type Group struct {
	ID     int
	Title  string
	Ref    string
	Locale string
}

//Channel represent the respective table in the database
type Channel struct {
	ID      int
	GroupID int
	Name    string
	Ref     string
}

//Permission represent the respective table in the database
type Permission struct {
	ID         int
	UserID     int
	GroupID    int
	Permission int
}

//List represent the respective table in the database
type List struct {
	ID               int
	Name             string
	GroupID          int
	GroupIndipendent int
	InviteOnly       int
}

//Bookmark represent the respective table in the database
type Bookmark struct {
	ID             int
	UserID         int
	GroupID        int
	MessageID      int
	Alias          string
	Status         int
	MessageContent string
}

//Subscription represent the respective table in the database
type Subscription struct {
	ID     int
	ListID int
	UserID int
}

//MessageCount represent the respective table in the database
type MessageCount struct {
	ID           int
	UserID       int
	GroupID      int
	MessageCount int
}

//Strings represent the respective table in the database
type Strings struct {
	ID      int
	Key     string
	Value   string
	Locale  string
	GroupID int
}

//Settings represent the respective table in the database
type Settings struct {
	ID      int
	Key     string
	Value   string
	GroupID int
}

//BotSettings represent the respective table in the database
type BotSettings struct {
	ID    int
	Key   string
	Value string
}

//BotString represent the respective table in the database
type BotString struct {
	ID     int
	Key    string
	Value  string
	Locale string
}

//Log represent the respective table in the database
type Log struct {
	ID             int
	Event          string
	RelatedUserID  int
	RelatedGroupID int
	Message        string
	UpdateValue    string
	Error          string
	Severity       int
	Date           string
}

//BotAdministrator represent the respective table in the database
type BotAdministrator struct {
	UserID      int
	Permissions int
}
