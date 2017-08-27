package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func loadConfig() map[string]string {
	variables := [4]string{"host", "username", "password", "database"}
	config := make(map[string]string)

	var undefined []string
	for _, v := range variables {
		config[v] = os.Getenv(v)
		if len(config[v]) == 0 {
			undefined = append(undefined, v)
		}
	}

	if len(undefined) > 0 {
		fmt.Printf("Not all environment variables are set.\n")
		fmt.Printf("Please enter the value(s) for the following variable(s).\n")

		scanner := bufio.NewScanner(os.Stdin)
		for _, v := range undefined {
			fmt.Printf("%s: ", v)
			scanner.Scan()
			config[v] = scanner.Text()
		}
	}

	return config
}

func DatabaseInit(config map[string]string) (Database, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		config["username"], config["password"],
		config["host"], config["database"],
	)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return Database{}, err
	}

	err = db.Ping()
	if err != nil {
		return Database{}, err
	}

	return Database{DB: db}, nil
}

func main() {
	config := loadConfig()
	db, err := DatabaseInit(config)
	if err != nil {
		fmt.Println("ERROR: Could not connect to the database")
		fmt.Println(err)
		return
	}

	defer db.DB.Close()

	router := httprouter.New()

	router.POST("/users", db.createUser)
	router.GET("/users", db.getAllUsers)
	router.GET("/users/:id", db.getUser)

	router.POST("/tokens", db.createToken)

	log.Fatal(http.ListenAndServe(":8080", router).Error())
}
