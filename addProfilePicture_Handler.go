package main

import (
  b64 "encoding/base64"
  "encoding/json"
  "fmt"
  "net/http"
  "strings"
)

func addProfilePicture(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  var blob Blob
  if err := json.NewDecoder(r.Body).Decode(&blob); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }
  // Spliiting up data:image/jpeg;base64,/9j/ffdgfd...
  imageParts := strings.Split(blob.Data, ",")
  htmlEmbed := imageParts[0]
  string64 := imageParts[1]
  sDec, err := b64.StdEncoding.DecodeString(string64)
  finished := make(chan bool)
  if err != nil {
    fmt.Println("Could not decode profile picture, error: " + err.Error())
    return
  }

  go func() {
    if err := session.Query("INSERT INTO profile_pictures (user_id, html_embed, image_blob) VALUES (?,?,?)",blob.User_Id, htmlEmbed, sDec).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    finished <- true
    json.NewEncoder(w).Encode(blob)
  }()
  <- finished
}
