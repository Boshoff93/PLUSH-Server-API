package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
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
    var display_name string
    itr := session.Query("SELECT display_name, user_id FROM users_by_id_like_display_name WHERE display_name Like ? LIMIT 50","%" + search.Search + "%").Iter()
    for itr.Scan(&display_name, &user_id) {
		    searchedUsers.User_Ids = append(searchedUsers.User_Ids, user_id)
        searchedUsers.Display_Names = append(searchedUsers.Display_Names ,display_name)
      }
      json.NewEncoder(w).Encode(searchedUsers)
      finished <- true
  }()
  <- finished
}
