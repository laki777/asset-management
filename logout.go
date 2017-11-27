package main

import (
	"asset-management/boltfunc"
	"net/http"
)

func logout(w http.ResponseWriter, req *http.Request) {

	boltfunc.Open("../database/bolt.db")
	defer boltfunc.Close()

	if cleanSessions(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// clean up db Sessions
	// if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
	// 	go cleanSessions()
	// }

}
