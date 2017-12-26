package models

import "fmt"

type User struct {
	ID   int64
	Name string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s>", u.ID, u.Name)
}
