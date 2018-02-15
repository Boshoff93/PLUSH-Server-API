package main

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
)

func searchUsers(w http.ResponseWriter, r *http.Request){

  session := getSession()
  defer session.Close()

  var search SearchUser
  params:= mux.Vars(r)
  search.Search = params["like_name"]

  finished := make(chan bool)
  go func() {
    var searchedUsers SearchedUsers
    var user_id string
    var email string
    var fullname string
    itr := session.Query("SELECT email, fullname, user_id FROM users_by_email_like_fullname WHERE fullname Like ? LIMIT 50","%" + search.Search + "%").Iter()
    for itr.Scan(&email,&fullname, &user_id) {
		    searchedUsers.User_Ids = append(searchedUsers.User_Ids, user_id)
        searchedUsers.Emails = append(searchedUsers.Emails, email)
        searchedUsers.Fullnames = append(searchedUsers.Fullnames ,fullname)
      }
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(searchedUsers)
      finished <- true
  }()
  <- finished
}
