package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	Routers()
}

func Routers() {
	InitDB()
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	http.ListenAndServe(":9080", &CORSRouterDecorator{router})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	result, err := db.Query("SELECT id, first_name," + "last_name, email from users")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO users(first_name," + "last_name, email) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	first_name := keyVal["firstName"]
	last_name := keyVal["lastName"]
	email := keyVal["email"]
	_, err = stmt.Exec(first_name, last_name, email)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New user was created")
}

type User struct {
	ID string `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
}

var db *sql.DB
var err error

func InitDB() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/userdb")
	if err != nil {
		panic(err.Error())
	}
}

type CORSRouterDecorator struct {
	R *mux.Router
}

func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language," + " Content-Type, YourOwnHeader")
	}
	if req.Method == "OPTIONS" {
		return
	}

	c.R.ServeHTTP(rw, req)
}