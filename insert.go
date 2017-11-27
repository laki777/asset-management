package main

import (
	"asset-management/boltfunc"
	"log"
	"net/http"
)

func insert(w http.ResponseWriter, req *http.Request) {

	boltfunc.Open("../database/bolt.db")
	defer boltfunc.Close()

	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	username := req.FormValue("username")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	hardware := req.FormValue("hardware")
	software := req.FormValue("software")

	input := []*boltfunc.Person{
		{username, firstname, lastname, hardware, software},
	}

	pd := pageData{
		Title:      "Insert Data",
		FirstName:  firstname,
		SecondName: lastname,
	}

	for _, p := range input {
		p.Save("people")
	}

	// boltfunc.List("people")

	err := tpl.ExecuteTemplate(w, "insert.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
