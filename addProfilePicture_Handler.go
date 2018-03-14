package main

import (
  //b64 "encoding/base64"
  "encoding/json"
  "fmt"
  "net/http"
  //"strings"
)

func addProfilePicture(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  var blob Blob
  if err := json.NewDecoder(r.Body).Decode(&blob); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }
  finished := make(chan bool)
  var pp_name string
  pp_name = blob.User_Id + "_pp_picture"
  go func() {
    if err := session.Query("INSERT INTO profile_picture_names (user_id, pp_name) VALUES (?,?)",blob.User_Id, pp_name).Exec(); err != nil {
      fmt.Println(err.Error());
      finished <- true
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
    }
    blob.Pp_Name = pp_name
    finished <- true
    json.NewEncoder(w).Encode(blob)
  }()
  <- finished
}
