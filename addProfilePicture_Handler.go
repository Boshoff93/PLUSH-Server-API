package main

import (
  "github.com/mitchellh/mapstructure"
  "fmt"
  b64 "encoding/base64"
  "strings"
)

func addProfilePicture(client *Client, data interface{}){
  var blob Blob
  err := mapstructure.Decode(data, &blob)
  if err != nil {
    client.send <- Message{"error", "could not decode addProfilePicture"}
    return
  }
  // Spliiting up data:image/jpeg;base64,/9j/ffdgfd...
  imageParts := strings.Split(blob.Data, ",")
  htmlEmbed := imageParts[0]
  string64 := imageParts[1]
  sDec, err := b64.StdEncoding.DecodeString(string64)
  if err != nil {
    client.send <- Message{"error", "could not decode Base64"}
    return
  }

  go func() {
    if err := client.session.Query("INSERT INTO profile_pictures (user_id, html_embed, image_blob) VALUES (?,?,?)",blob.User_Id, htmlEmbed, sDec).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    client.send <- Message{"profile picture add", blob}
  }()

}
