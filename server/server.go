package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/codes4coffee/bag-tag-server/user"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type userLogic struct {
	users []user.User
}

func (users *userLogic) userFinderHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, uuid.New().String())
	vars := mux.Vars(r)
	userIdx := user.UserFinder(users.users, vars["id"])

	if userIdx != -1 {
		fmt.Fprintf(w, "User found. Name is: "+users.users[userIdx].Name)
	}
}

func (users *userLogic) userRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic("Oh shit")
	}
	fName := r.PostFormValue("first-name")
	passcode := r.PostFormValue("passcode")
	newUser := user.User{Id: uuid.New(), Name: fName, passcode: passcode}
	users.users = append(users.users, newUser)
	fmt.Fprintf(w, "Created ID: "+newUser.Id.String()+" for "+newUser.Name)
}

func (users *userLogic) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())

	err := r.ParseForm()

	if err != nil {
		panic("Error getting post data from login")
	}

	uName := r.PostFormValue("username")
	passcode := r.PostFormValue("passcode")
	userIdx := userFinder(users.users, uName)

	if userIdx != -1 {
		userObj := &users.users[userIdx] //Need to get the pointer so I can change the value on line 85
		if passcode == userObj.passcode {
			token := uuid.New().String()
			userObj.sessionToken = token
			userObj.tokenGeneratedAt = time.Now().Unix()
			fmt.Fprintf(w, token)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func (users *userLogic) getUserAccountHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-SESSION-TOKEN")
	userObj := findUserBySessionToken(token, users.users)
	if userObj == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	} else {
		returnedUserJson, err := json.Marshal(userObj)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, string(returnedUserJson))
	}
}

func main() {
	uLogic := &userLogic{users: make([]user, 3)}
	r := mux.NewRouter()
	r.HandleFunc("/found/{id}", uLogic.userFinderHandler).Methods("GET")
	r.HandleFunc("/register", uLogic.userRegistrationHandler).Methods("POST")
	r.HandleFunc("/login", uLogic.userLoginHandler).Methods("POST")
	r.HandleFunc("/account", uLogic.getUserAccountHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	log.Fatal(http.ListenAndServe(":8080", r))
}
