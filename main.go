package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type User struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	CreatedAT       string `json:"createAT"`
}

var (
	errProblemConfirmPassword = errors.New("problem with confirming password")
	errProblemLenPassword     = errors.New("problem with length of password")
	errProblemSymbPassword    = errors.New("problem with correct password symbols")
	errProblemSymbEmail       = errors.New("problem with correct email symbols")
	errProblemLenEmail        = errors.New("problem with length of email")
	errProblemSymbFirstName   = errors.New("problem with correct first name symbols")
	errProblemLenFirstName    = errors.New("problem with length of first name")
	errProblemSymbLastName    = errors.New("problem with correct last name symbols")
	errProblemLenLastName     = errors.New("problem with length of last name")
	errProblemOpenDB          = errors.New("problem with opening data base")
	errProblemUniqueEmail     = errors.New("problem with unique email adress")
)

const (
	host          = "localhost"
	port          = 5432
	user          = "postgres"
	password      = "123456"
	dbname        = "users"
	webServerPort = ":8080"
)

func main() {

	http.HandleFunc("/users", CreateUserHandler)
	fmt.Println("Server listen on port:", webServerPort)
	if err := http.ListenAndServe(webServerPort, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		const dateTmplate = "2006-1-2 15:4:5"
		var u User

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Println("Failed to unmarshal request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = u.validate()
		if err != nil {
			respondWithError(w, http.StatusNotAcceptable, err)

			return
		}

		u.ID = uuid.New().String()
		u.CreatedAT = time.Now().Format(dateTmplate)

		if err := createUserDB(&u); err != nil {
			respondWithError(w, http.StatusNotAcceptable, err)
			return
		}

		successfullSignUpText := ([]byte(fmt.Sprintf(`{"you have successfully created a user with login":"%v"}`, u.Email)))
		w.Write(successfullSignUpText)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func createUserDB(u *User) error {
	db, err := connectDB()
	if err != nil {
		return err
	}

	defer db.Close()

	insertDynStmt := `insert into "users"("id", "email", "firstname", "lastname", "password", "createdat") values($1,$2,$3,$4,$5,$6);`
	_, err = db.Exec(insertDynStmt, u.ID, u.Email, u.FirstName, u.LastName, u.Password, u.CreatedAT)
	if err != nil {
		return errProblemUniqueEmail
	}

	return nil
}

func connectDB() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, errProblemOpenDB
	}

	return db, err
}

func respondWithError(w http.ResponseWriter, statusCode int, customerResponse error) {
	w.WriteHeader(statusCode)
	errorText := ([]byte(fmt.Sprintf(`{"error":"%v"}`, customerResponse.Error())))
	w.Write(errorText)
	log.Println(customerResponse.Error())
}

func (u *User) validate() error {
	if err := validatePassword(u); err != nil {
		return err
	}
	if err := validateEmail(u); err != nil {
		return err
	}

	if err := validateFirstName(u); err != nil {
		return err
	}

	if err := validateLastName(u); err != nil {
		return err
	}

	return nil
}

func validatePassword(u *User) error {
	const (
		minPasswordLen = 8
		maxPasswordLen = 256
	)

	for i := 0; i < len(u.Password); i++ {
		if u.Password[i] > unicode.MaxASCII {
			return errProblemSymbPassword
		}
	}

	if u.Password != u.ConfirmPassword {
		return errProblemConfirmPassword
	}

	if len(u.Password) < minPasswordLen || len(u.Password) > maxPasswordLen {
		return errProblemLenPassword
	}

	return nil
}

func validateEmail(u *User) error {
	const (
		minEmailLen = 4
		maxEmailLen = 256
	)

	if len(u.Email) < minEmailLen || len(u.Email) > maxEmailLen {
		return errProblemLenEmail
	}

	for i := 0; i < len(u.Email); i++ {
		if u.Email[i] > unicode.MaxASCII {
			return errProblemSymbEmail
		}
	}

	if !strings.Contains(u.Email, "@") {
		return errProblemSymbEmail
	}

	return nil
}

func validateFirstName(u *User) error {
	const minFirstNameLen = 1

	for i := 0; i < len(u.FirstName); i++ {
		if u.FirstName[i] > unicode.MaxASCII {
			return errProblemSymbFirstName
		}
	}

	if len(u.FirstName) < minFirstNameLen {
		return errProblemLenFirstName
	}

	return nil
}

func validateLastName(u *User) error {
	const minLastNameLen = 1

	for i := 0; i < len(u.LastName); i++ {
		if u.LastName[i] > unicode.MaxASCII {
			return errProblemSymbLastName
		}
	}

	if len(u.LastName) < minLastNameLen {
		return errProblemLenLastName
	}

	return nil
}
