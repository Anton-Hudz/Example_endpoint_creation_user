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

var u User

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
	// switch r.Method {
	// case http.MethodPost:
	var customerResponse error
	// if r.Method == http.MethodPost {
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// }

	// make function helper for printing error in log and sending to user
	if err := validationPassword(&u); err != nil {
		customerResponse = err
	}

	if err := validationEmail(&u); err != nil {
		customerResponse = err
	}

	if err := validationFirstName(&u); err != nil {
		customerResponse = err
	}

	if err := validationLastName(&u); err != nil {
		customerResponse = err
	}

	if customerResponse != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorText := map[string]string{"error": customerResponse.Error()}
		err := json.NewEncoder(w).Encode(errorText)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //тут можно прописать универсальную перем со статускодом, котор будет передаваться из валидации
			return
		}
		log.Println(customerResponse.Error())

		return
	}

	u.ID = uuid.New().String()
	u.CreatedAT = time.Now().Format("2006-1-2 15:4:5")

	if err := creationUserDB(&u); err != nil {
		customerResponse = err
	}

	if customerResponse != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorText := map[string]string{"error": customerResponse.Error()}
		err := json.NewEncoder(w).Encode(errorText)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //тут можно прописать универсальную перем со статускодом, котор будет передаваться из валидации
			return
		}
		log.Println(customerResponse.Error())

		return
	}

	successfullSignUpText := map[string]string{"you have successfully created a user with login": u.Email}
	err = json.NewEncoder(w).Encode(successfullSignUpText)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Successfuly created user")
}

func creationUserDB(u *User) error {
	db, err := connectionDB()
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

func connectionDB() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, errProblemOpenDB
	}

	return db, err
}

func validationPassword(u *User) error {
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

func validationEmail(u *User) error {
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

func validationFirstName(u *User) error {
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

func validationLastName(u *User) error {
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
