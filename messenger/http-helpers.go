package main

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func validateToken(tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	return err == nil
}

func validateUserFromToken(tok string, w http.ResponseWriter) (*User, bool) {
	u := getUserFromToken(tok)
	if u == nil {
		http.Error(w, "500 - error fetching user!", http.StatusInternalServerError)
	}
	return u, u != nil
}

func validateURLToken(w http.ResponseWriter, r *http.Request) (string, bool) {
	var tokenParam []string
	var tok string
	var ok bool
	tokenParam, ok = r.URL.Query()["token"]
	if !ok || len(tokenParam) < 1 || tokenParam[0] == "undefined" {
		http.Error(w, "400 - token invalid!", http.StatusBadRequest)
		return tok, ok
	}
	tok = tokenParam[0]
	if ok = validateToken(tok); !ok {
		http.Error(w, "403 - you are not authorized to do this action!", http.StatusForbidden)
	}
	return tok, ok
}

func validateFormToken(w http.ResponseWriter, r *http.Request) (string, bool) {
	var tok string
	var ok bool
	tok = r.FormValue("token")
	ok = validateToken(tok)
	if !ok {
		http.Error(w, "403 - you are not authorized to do this action!", http.StatusForbidden)
	}
	return tok, ok
}

func validateUserIDFromPath(w http.ResponseWriter, r *http.Request) (*User, bool) {
	var u *User
	params := mux.Vars(r)
	if u = findUserByID(params["user_id"]); u == nil {
		http.Error(w, "400 - user doesn't exist!", http.StatusBadRequest)
	}
	return u, u != nil
}

func validateHubIDFromPath(w http.ResponseWriter, r *http.Request) (*Hub, bool) {
	var h *Hub
	params := mux.Vars(r)
	if h = getHub(params["hub_id"]); h == nil {
		http.Error(w, "400 - hub doesn't exist!", http.StatusBadRequest)
	}
	return h, h != nil
}

func validateUserIDFromForm(w http.ResponseWriter, r *http.Request) (*User, bool) {
	var fu *User // form user
	var fid string
	fid = r.FormValue("user_id")
	if fu = findUserByID(fid); fu == nil {
		http.Error(w, "400 - user not found!", http.StatusBadRequest)
	}
	return fu, fu != nil
}
