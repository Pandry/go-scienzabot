package database

import "time"

type User struct {
	ID        int
	Nickname  string
	Biography string
	Status    int
	LastSeen  time.Time
}

type Group struct {
	ID    int
	Title string
	Ref   string
}

type Channel struct {
	ID      int
	GroupID int
	Name    string
	Ref     string
}

type Permission struct {
	ID          int
	UserID      int
	GroupID     int
	Permissions int
}

type List struct {
	ID   int
	Name string
}

type Bookmark struct {
	ID        int
	UserID    int
	GroupID   int
	MessageID int
	Alias     string
}

type Subscription struct {
	ID     int
	ListID int
	UserID int
}

type MessageCount struct {
	ID           int
	UserID       int
	GroupID      int
	MessageCount int
}

type Strings struct {
	ID      int
	GroupID int
	Key     string
	Value   string
	Locale  string
}

type Settings struct {
	ID      int
	GroupID int
	Key     string
	Value   string
}

type BotSettings struct {
	ID    int
	Key   string
	Value string
}
