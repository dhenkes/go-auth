package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (db *Database) createUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		resCouldNotReadBody(w)
		return
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		resCouldNotParseBody(w)
		return
	}

	var rowsFound int
	err = db.DB.QueryRow("SELECT count(0) FROM users WHERE username = $1", strings.ToLower(user.Username)).Scan(&rowsFound)
	if err != nil {
		resCouldNotInsertIntoDB(w)
		return
	}

	if rowsFound > 0 {
		resCouldNotInsertIntoDB(w)
		return
	}

	user.PasswordDB, err = generateBcryptHash(user.Password)
	if err != nil {
		resCouldNotHashPassword(w)
		return
	}

	user.Created = int(time.Now().Unix())
	user.Uuid, err = newUUID()
	if err != nil {
		resCouldNotGenerateUUID(w)
		return
	}

	stmt, err := db.DB.Prepare("INSERT INTO users(uuid, username, password, created) VALUES($1, $2, $3, $4)")
	if err != nil {
		resCouldNotPrepareStmt(w)
		return
	}

	_, err = stmt.Exec(user.Uuid, strings.ToLower(user.Username), user.PasswordDB, user.Created)
	stmt.Close()
	if err != nil {
		resCouldNotInsertIntoDB(w)
		return
	}

	user.Password = ""
	user.PasswordDB = ""

	resData(w, user)
}

func (db *Database) getAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := db.DB.Query("SELECT uuid, username, created, removed FROM users WHERE removed = 0")
	defer rows.Close()
	if err != nil {
		// No rows found
	}

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Uuid, &user.Username, &user.Created, &user.Removed)
		if err != nil {
			break
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		// No rows found
	}

	resData(w, users)
}

func (db *Database) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := User{}
	err := db.DB.QueryRow("SELECT uuid, username, created, removed FROM users WHERE uuid = $1 AND removed = 0", ps.ByName("id")).Scan(&user.Uuid, &user.Username, &user.Created, &user.Removed)
	if err != nil {
		resNoRowFound(w)
		return
	}

	resData(w, user)
}
