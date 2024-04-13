package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// create 2 APIs.
// POST:
// Send the username and password as a json request. Save it in the DB.
// GET:
// Take the username as a query param, if a user existing with the name give the user
//  details as response. If the query param is empty give all user details available in the db.

var DB *gorm.DB

func Setupdb() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	log.Println("(db) Open Completed")
	if err != nil {
		log.Printf("(db) sqlite Error in Opening %v", err)
		return
	}
	log.Println("Migrateing Database")
	db.AutoMigrate(
		&User{},
	)
	DB = db

}

type User struct {
	gorm.Model
	UserName string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(user User) string {
	resp := DB.Create(&user)
	if resp.Error != nil {
		return "Error in Saving Data " + resp.Error.Error()
	} else {
		return "Data Added Successfully"
	}
}

func GetUserData(username string) (User, error) {
	var user User
	resp := DB.Where("user_name = ?", username).Find(user)
	if resp.Error != nil {
		return User{}, resp.Error
	} else {
		return user, nil
	}
}

func main() {
	go Setupdb()
	r := mux.NewRouter()
	r.HandleFunc("/adduser", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		data := r.Body
		var user User
		_ = json.NewDecoder(data).Decode(&user)

		RESP := CreateUser(user)
		log.Println(RESP)
	}).Methods("POST")

	r.HandleFunc("/getuser", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := r.URL.Query().Get("username")
		resp, err := GetUserData(params)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
		} else {
			jsondata, _ := json.Marshal(resp)
			w.Write(jsondata)
		}
	}).Methods("GET")

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Print(err)
	}

}
