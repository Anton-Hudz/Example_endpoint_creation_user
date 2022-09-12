package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var Users = []User{}

type User struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func main() {
	// mux := http.NewServeMux()
	http.HandleFunc("/random", UserCreateHandler)
	port := ":8080"
	fmt.Println("Server listen on port:", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var u User
		//--------------------------------------------------
		// u.Login = r.URL.Query().Get("login")
		// u.FirstName = r.URL.Query().Get("firstName")
		// u.LastName = r.URL.Query().Get("lastName")
		// u.Password = r.URL.Query().Get("password")
		// u.ConfirmPassword = r.URL.Query().Get("confirmPassword")
		// w.WriteHeader(http.StatusCreated)
		// w.Write([]byte(fmt.Sprintf("Your Login is: %v.\n Your first name is %v.\n Your lastname name is %v.\n ", u.Login, u.FirstName, u.LastName)))
		// if u.Password != u.ConfirmPassword {
		// 	w.Write([]byte(fmt.Sprintf("your password and password confirmation do't match")))
		// 	fmt.Println("Problem with password")
		// 	return
		// }
		// fmt.Println(u.Password)
		// fmt.Println(u.ConfirmPassword)

		//-------------------------------------------------

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if u.Password != u.ConfirmPassword {
			w.Write([]byte(fmt.Sprintf("your password and password confirmation do't match")))
			fmt.Println("Problem with password")
			return
		}
		fmt.Fprintf(w, "Your login is: %s, first name is: %s, last name is: %s", u.Login, u.FirstName, u.LastName)

	case http.MethodGet:

	}
}
