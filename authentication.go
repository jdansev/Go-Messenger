package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("alpacas-are-just-humpless-camels-i-swear")

// Register : create a user account
func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	usr := r.FormValue("username")
	pwd := r.FormValue("password")

	if usr == "" {
		http.Error(w, "400 - no username!", http.StatusBadRequest)
		return
	}

	if pwd == "" {
		http.Error(w, "400 - no password!", http.StatusBadRequest)
		return
	}

	if findUserByName(usr) != nil {
		http.Error(w, "400 - user already exists!", http.StatusBadRequest)
		return
	}

	newUser := createUser(usr, pwd)
	json.NewEncoder(w).Encode(newUser)
}

// Login : login to an account
func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	usr := r.FormValue("username")
	pwd := r.FormValue("password")

	if usr == "" {
		http.Error(w, "400 - no username!", http.StatusBadRequest)
		return
	}

	if pwd == "" {
		http.Error(w, "400 - no password!", http.StatusBadRequest)
		return
	}

	u := findUserByName(usr)

	if u == nil {
		http.Error(w, "400 - user doesn't exists!", http.StatusBadRequest)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	if err != nil {
		http.Error(w, "403 - password doesn't match!", http.StatusForbidden)
		return
	}

	token := u.generateToken()

	js, _ := json.Marshal(map[string]string{
		"username": u.Username,
		"id": u.ID,
		"token": token,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
