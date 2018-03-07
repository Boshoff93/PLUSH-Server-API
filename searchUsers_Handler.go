package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
  b64 "encoding/base64"
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
      var user_id string
      var htmlEmbed string
      var string64 []byte
      if err := session.Query("SELECT * FROM profile_pictures WHERE user_id = ?",id).Scan(&user_id, &htmlEmbed, &string64); err != nil {
        fmt.Println("Could not get profile picture, error: " + err.Error() )
        json.NewEncoder(w).Encode(Error{Error: err.Error()})
        finished <- true
        return
      }
      encodedString := b64.StdEncoding.EncodeToString(string64)
      //Constructing html base64 embeded image
      var base64EmbededImage = htmlEmbed + "," + encodedString
      searchedUsers.Avatars = append(searchedUsers.Avatars ,base64EmbededImage)
    }
      json.NewEncoder(w).Encode(searchedUsers)
      finished <- true
  }()
  <- finished
}
