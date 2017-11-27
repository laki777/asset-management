package main

import (
	"asset-management/boltfunc"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func findUserOnly(user string) bool {
	p, error := boltfunc.AdminGetPerson(user, "adminpeople")
	if error != nil {
		return false
	} else if p != nil {
		return true
	}
	return false
}

func hashPassword(pass string) string {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func signup(w http.ResponseWriter, req *http.Request) {

	boltfunc.Open("../database/bolt.db")
	defer boltfunc.Close()

	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	pd := pageData{
		Title: "Sign Up Page",
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	confirmpass := req.FormValue("confirmpass")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	// role := "user"

	if password != confirmpass {
		http.Redirect(w, req, "/signup", 301)
		fmt.Println("Password does not match the confirm password.")
	} else {

		if (username != "") && findUserOnly(username) {
			http.Redirect(w, req, "/signup", 301)
			fmt.Printf("User with (%s) username already exists!", username)
		} else if (username != "") && (password != "") {

			k, _ := boltfunc.AdminGetFirstPerson()

			role := k
			input := []*boltfunc.AdminPerson{
				{username, hashPassword(password), firstname, lastname, role},
			}

			for _, p := range input {
				p.AdminSave("adminpeople")
			}
			http.Redirect(w, req, "/", 301)
		}
	}
	err := tpl.ExecuteTemplate(w, "signup.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
