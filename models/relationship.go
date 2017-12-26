package models

type Relationship struct {
	ID      int64
	OwnerID int64
	UserID  int64
	State   string
	ReState string
}
