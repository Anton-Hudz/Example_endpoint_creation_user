package main

import (
	"fmt"
	"log"
	"net/http"
)

type Users []User
type User struct {
	Id              string
	Login           string
	FirstName       string
	LastName        string
	Password        string
	ConfirmPassword string
}

func main() {
	http.HandleFunc("/random", GetUserDataHandler)
	port := ":8080"
	fmt.Println("Server listen on port:", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var u User
		u.Login = r.URL.Query().Get("login")
		u.FirstName = r.URL.Query().Get("firstName")
		u.LastName = r.URL.Query().Get("lastName")
		u.Password = r.URL.Query().Get("password")
		u.ConfirmPassword = r.URL.Query().Get("confirmPassword")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Your Login is: %v.\n Your first name is %v.\n Your lastname name is %v.\n ", u.Login, u.FirstName, u.LastName)))
		fmt.Println(u.Password)

	case http.MethodGet:
	case http.MethodPut:
	case http.MethodDelete:
	default:
	}
}
