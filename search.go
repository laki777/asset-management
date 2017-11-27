package main

import (
	"asset-management/boltfunc"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func search(w http.ResponseWriter, req *http.Request) {

	boltfunc.Open("../database/bolt.db")
	defer boltfunc.Close()

	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	pd := pageData{
		Title: "Search Page",
	}

	boltfunc.List("people")

	// var results []*boltfunc.Person
	// var error_res error

	// if results, error_res = search_query(req.FormValue("search")); error_res != nil {
	// 	http.Error(w, error_res.Error(), http.StatusInternalServerError)
	// }

	err := tpl.ExecuteTemplate(w, "search.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// func search_query(query string) ([]*boltfunc.Person, error) {

// 	for _, id := range []string{"query"} {
// 		p, err := boltfunc.GetPerson("people", id)
// 		if err != nil {
// 			return []boltfunc.Person{}, err
// 		}
// 		fmt.Println(p)
// 	}

// }
