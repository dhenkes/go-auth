package main

type User struct {
	Uuid       string `json:"uuid"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	PasswordDB string `json:"-"`
	Created    int    `json:"created"`
	Removed    int    `json:"removed"`
}

type Token struct {
	Uuid    string `json:"uuid"`
	User    string `json:"user"`
	Token   string `json:"token"`
	Created int    `json:"created"`
	Expires int    `json:"expires"`
	Removed int    `json:"removed"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
