package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/amlun/minit/services"
	"io/ioutil"
	"github.com/amlun/minit/models"
	"strconv"
	"log"
)

var errorResponse = []byte("{\"message\": \"error\"}")

type UserPostRequest struct {
	Name string `json:"name"`
}

type RelationsPutRequest struct {
	State string `json:"state"`
}

type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RelationResponse struct {
	UserID int64  `json:"user_id"`
	State  string `json:"state"`
	Type   string `json:"type"`
}

// UsersGetHandler
// GET /users [List all users]
func UsersGetHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		body []byte
	)
	defer func() {
		writeBody(w, body, err)
	}()

	var u UserResponse
	var us []UserResponse

	users := services.GetAllUsers()
	for _, v := range users {
		u.ID = v.ID
		u.Name = v.Name
		u.Type = services.UserType
		us = append(us, u)
	}
	body, err = json.Marshal(us)
}

// UsersRelationshipsGetHandler
// GET /users/:user_id/relationships [List a users all relationships]
func UsersRelationshipsGetHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		body []byte
	)

	defer func() {
		writeBody(w, body, err)
	}()

	vars := mux.Vars(r)
	if userID, ok := vars["user_id"]; ok {
		id, err := strconv.ParseInt(userID, 10, 0)
		if err != nil {
			log.Printf("strconv.ParseInt(%v) with error: %v", userID, err)
			return
		}
		var rl RelationResponse
		var rls []RelationResponse

		relations := services.GetUserRelations(id)
		for _, v := range relations {
			rl.UserID = v.UserID
			if v.State == services.RelationshipStateLiked && v.ReState == services.RelationshipStateLiked {
				rl.State = services.RelationshipStateMatched
			} else {
				rl.State = v.State
			}
			rl.Type = services.UserType
			rls = append(rls, rl)
		}

		body, err = json.Marshal(rls)
	}
}

// UsersAddHandler
// POST /users [Create a user]
func UsersAddHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		body []byte
		user models.User
		u    UserResponse
	)

	defer func() {
		writeBody(w, body, err)
	}()

	var req UserPostRequest
	if err := parseRequest(r, &req); err != nil {
		return
	}

	user.Name = req.Name
	if err = services.AddUser(&user); err != nil {
		return
	}
	u.ID = user.ID
	u.Name = user.Name
	u.Type = services.UserType
	body, err = json.Marshal(u)
}

// UsersRelationshipsAddHandler
// PUT /users/:user_id/relationships/:other_user_id Create/update relationship state to another user.
func UsersRelationshipsAddHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		body    []byte
		ownerID int64
		userID  int64
		rl      *models.Relationship
		rlRes   RelationResponse
	)

	defer func() {
		writeBody(w, body, err)
	}()

	vars := mux.Vars(r)
	if ownerStr, ok := vars["owner_id"]; ok {
		ownerID, err = strconv.ParseInt(ownerStr, 10, 0)
		if err != nil {
			log.Printf("strconv.ParseInt(%v) with error: %v", ownerStr, err)
			return
		}
	}

	if userStr, ok := vars["user_id"]; ok {
		userID, err = strconv.ParseInt(userStr, 10, 0)
		if err != nil {
			log.Printf("strconv.ParseInt(%v) with error: %v", userStr, err)
			return
		}
	}

	var req RelationsPutRequest

	if err := parseRequest(r, &req); err != nil {
		return
	}

	if req.State != services.RelationshipStateLiked && req.State != services.RelationshipStateDisliked {
		return
	}

	if rl, err = services.AddUserRelation(ownerID, userID, req.State); err != nil {
		return
	}

	rlRes.Type = services.RelationshipType
	rlRes.UserID = rl.UserID
	if rl.State == services.RelationshipStateLiked && rl.ReState == services.RelationshipStateLiked {
		rlRes.State = services.RelationshipStateMatched
	} else {
		rlRes.State = rl.State
	}

	body, err = json.Marshal(rlRes)
}

// send the response
func writeBody(w http.ResponseWriter, body []byte, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil || body == nil {
		log.Printf("write body with error: %v or body is %v", err, body)
		w.Write(errorResponse)
	} else {
		w.Write(body)
	}
}

// parse request body (json string) to struct
func parseRequest(r *http.Request, v interface{}) error {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll(%v) with error: %v", r.Body, err)
		return err
	}
	if err := json.Unmarshal(req, v); err != nil {
		log.Printf("json.Unmarshal error: %v", err)
		return err
	}
	return nil
}
