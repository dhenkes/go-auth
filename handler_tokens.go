package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (db *Database) createToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	err = db.DB.QueryRow("SELECT uuid, password FROM users WHERE username = $1 AND removed = 0", strings.ToLower(user.Username)).Scan(&user.Uuid, &user.PasswordDB)
	if err != nil {
		user.Password = ""
		user.PasswordDB = ""

		resCouldNotGenerateToken(w)
		return
	}

	err = compareHashAndPassword(user.PasswordDB, user.Password)
	user.Password = ""
	user.PasswordDB = ""
	if err != nil {
		resCouldNotGenerateToken(w)
		return
	}

	token := Token{}
	token.Created = int(time.Now().Unix())
	token.Expires = token.Created + (60 * 60 * 24)
	token.User = user.Uuid
	token.Token, err = generateToken()
	if err != nil {
		resCouldNotGenerateToken(w)
		return
	}

	token.Uuid, err = newUUID()
	if err != nil {
		resCouldNotGenerateUUID(w)
		return
	}

	stmt, err := db.DB.Prepare("INSERT INTO tokens(uuid, user_uuid, token, created, expires) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		resCouldNotPrepareStmt(w)
		return
	}

	_, err = stmt.Exec(token.Uuid, token.User, token.Token, token.Created, token.Expires)
	stmt.Close()
	if err != nil {
		resCouldNotInsertIntoDB(w)
		return
	}

	resData(w, token)
}
