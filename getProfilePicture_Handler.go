package main

import (
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "net/http"
)


func getProfilePicture(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)
  var user User
  user.User_Id = params["user_id"]

  finished := make(chan bool)
  go func() {
     var user_id string

    var pp_name string
    if err := session.Query("SELECT * FROM profile_picture_names WHERE user_id = ?",user.User_Id).Scan(&user_id, &pp_name); err != nil {
      fmt.Println("Could not get profile picture, error: " + err.Error() )
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
    }
    var blob Blob
    blob.Pp_Name = pp_name
    blob.User_Id = user_id
    json.NewEncoder(w).Encode(blob)
    finished <- true
  }()
  <- finished
}
