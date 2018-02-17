package main

import (
  b64 "encoding/base64"
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
    var htmlEmbed string
    var string64 []byte
    if err := session.Query("SELECT * FROM profile_pictures WHERE user_id = ?",user.User_Id).Scan(&user_id, &htmlEmbed, &string64); err != nil {
      fmt.Println("Could not get profile picture, error: " + err.Error() )
      json.NewEncoder(w).Encode("")
      finished <- true
      return
    }
    encodedString := b64.StdEncoding.EncodeToString(string64)
    //Constructing html base64 embeded image
    var base64EmbededImage = htmlEmbed + "," + encodedString
    var blob Blob
    blob.Data = base64EmbededImage
    blob.User_Id = user_id
    json.NewEncoder(w).Encode(blob)
    finished <- true
  }()
  <- finished
}
