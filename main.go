package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"unicode"
)

type User struct {
	Id              string `json:"id"`
	Email           string `json:"email"`
	FullName        string `json:"fullName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

var Users = []User{}
var u User

var (
	problemConfirmPasswordErr = errors.New("problem with confirming password")
	problemLenPasswordErr     = errors.New("problem with length of password")
	problemSymbPasswordErr    = errors.New("problem with correct password symbols")
	problemSymbEmailErr       = errors.New("problem with correct email symbols")
	problemLenEmailErr        = errors.New("problem with length of email")
	problemSymbFullNameErr    = errors.New("problem with correct full name symbols")
	problemLenFullNameErr     = errors.New("problem with length of full name")
)

func main() {
	http.HandleFunc("/users", OperationUserHandler)
	port := ":8080"
	fmt.Println("Server listen on port:", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func OperationUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var customerResponse error

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validationPassword(&u); err != nil {
			customerResponse = err
		}
		if err := validationEmail(&u); err != nil {
			customerResponse = err
		}
		if err := validationFullName(&u); err != nil {
			customerResponse = err
		}
		if customerResponse != nil {
			w.WriteHeader(http.StatusBadRequest) //пока одна ошибка
			errorText := map[string]string{"error": customerResponse.Error()}
			err := json.NewEncoder(w).Encode(errorText)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			return
		}

		u.Id = uuid.New().String()

		successfullSignUpText := map[string]string{"you have successfully created a user with login": u.Email}
		err = json.NewEncoder(w).Encode(successfullSignUpText)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("Successfuly created user")

	case http.MethodGet:
	case http.MethodPut:
	}
}

func validationPassword(u *User) error {
	const (
		minPasswordLen = 8
		maxPasswordLen = 256
	)
	var correctPasswordSymb = true

	for i := 0; i < len(u.Password); i++ {
		if u.Password[i] > unicode.MaxASCII {
			correctPasswordSymb = false
		}
	}

	if correctPasswordSymb == false {
		return problemSymbPasswordErr
	}

	if u.Password != u.ConfirmPassword {
		return problemConfirmPasswordErr
	}

	if len(u.Password) < minPasswordLen || len(u.Password) > maxPasswordLen {
		return problemLenPasswordErr
	}

	return nil
}

func validationEmail(u *User) error {
	const (
		minEmailLen = 4
		maxEmailLen = 256
	)
	var correctEmailSymb = true

	if len(u.Email) < minEmailLen || len(u.Email) > maxEmailLen {
		return problemLenEmailErr
	}

	for i := 0; i < len(u.Email); i++ {
		if u.Email[i] > unicode.MaxASCII {
			correctEmailSymb = false
		}
	}

	if !strings.Contains(u.Email, "@") {
		correctEmailSymb = false
	}

	if correctEmailSymb == false {
		return problemSymbEmailErr
	}

	return nil
}
func validationFullName(u *User) error {
	const minFullNamelLen = 1
	var correctFullNameSymb = true

	for i := 0; i < len(u.FullName); i++ {
		if u.FullName[i] > unicode.MaxASCII {
			correctFullNameSymb = false
		}
	}
	if correctFullNameSymb == false {
		return problemSymbFullNameErr
	}

	if len(u.FullName) < minFullNamelLen {
		return problemLenFullNameErr
	}

	return nil
}

//----------------------------------------------------------------------------------

// As a user, I want to register in the system. To register I provide an email (login), fullname and password.

// Requirements for password validation:
// - min 8
// - max 256
// - ASCI symbols

// Requirements for login(email) validation:
// - max 256 symbols
// - must be @
// - TODO: search typical email requirements
// - must be unique

// Requirements for Fullname:
// - min 2/3 symbols

// User must have ID (UUID)
// CreatedAt column must be in the database (and maybe UpdatedAt)
