package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
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
    itr := session.Query("SELECT display_name, user_id FROM users_by_id_like_display_name WHERE display_name Like ? LIMIT 10","%" + search.Search + "%").Iter()
    for itr.Scan(&display_name, &user_id) {
		  searchedUsers.User_Ids = append(searchedUsers.User_Ids, user_id)
      searchedUsers.Display_Names = append(searchedUsers.Display_Names ,display_name)
    }

    for _, id := range searchedUsers.User_Ids {
      var pp_name string
      if err := session.Query("SELECT pp_name FROM profile_picture_names WHERE user_id = ?",id).Scan(&pp_name); err != nil {
        fmt.Println("Could not find profile picture name, error: " + err.Error() )
        searchedUsers.Pp_Names = append(searchedUsers.Pp_Names, "empty")
      } else {
        searchedUsers.Pp_Names = append(searchedUsers.Pp_Names, pp_name)
      }
    }
    json.NewEncoder(w).Encode(searchedUsers)
    finished <- true
  }()
  <- finished
}
