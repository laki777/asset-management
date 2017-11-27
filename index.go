package main

import (
	"asset-management/boltfunc"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func finduser(user, pass string, w http.ResponseWriter, req *http.Request) bool {
	//Get a person from the database by their username.
	p, error := boltfunc.AdminGetPerson(user, "adminpeople")
	if error != nil {
		return false
	} else if p != nil {
		err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(pass))
		if err != nil {
			return false
		}
		setSession(user, w, req)
		return true
	}
	return false
}

func index(w http.ResponseWriter, req *http.Request) {

	boltfunc.Open("../database/bolt.db")
	defer boltfunc.Close()

	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/insert", http.StatusSeeOther)
		return
	}

	pd := pageData{
		Title: "Index Page",
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	if (username != "") && finduser(username, password, w, req) {
		http.Redirect(w, req, "/insert", 301)
	} else if (username != "") && (password != "") {
		http.Redirect(w, req, "/", 301)
		// http.Error(w, "Username and/or password do not match", http.StatusForbidden)
	}
	err := tpl.ExecuteTemplate(w, "index.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
