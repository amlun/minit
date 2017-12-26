package services

import (
	"log"
	"github.com/amlun/minit/models"
)

func GetAllUsers() []models.User {
	var users []models.User
	_, err := db.Query(&users, `SELECT * FROM users`)
	if err != nil {
		log.Printf("GetAllUsers with error: %v", err)
		return nil
	}
	return users
}

func AddUser(user *models.User) error {
	if err := db.Insert(user); err != nil {
		log.Printf("AddUser with error: %v", err)
		return err
	}
	return nil
}

func GetUserRelations(ownerID int64) []models.Relationship {
	var relations []models.Relationship
	_, err := db.Query(&relations, `SELECT * FROM relationships WHERE owner_id = ?`, ownerID)
	if err != nil {
		log.Printf("GetAllUsers with error: %v", err)
		return nil
	}
	return relations
}

func AddUserRelation(ownerID, userID int64, state string) (*models.Relationship, error) {
	var err error
	var r models.Relationship
	var r1 = &models.Relationship{OwnerID: ownerID, UserID: userID, State: state}
	var r2 = &models.Relationship{OwnerID: userID, UserID: ownerID, ReState: state}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	// owner_id like user_id
	_, err = tx.Model(r1).OnConflict("(owner_id, user_id) DO UPDATE").Set("state = EXCLUDED.state").Insert()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// user_id is liked by owner_id
	_, err = tx.Model(r2).OnConflict("(owner_id, user_id) DO UPDATE").Set("re_state = EXCLUDED.re_state").Insert()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	_, err = db.QueryOne(&r, `SELECT * FROM relationships WHERE owner_id = ? AND user_id = ?`, ownerID, userID)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
