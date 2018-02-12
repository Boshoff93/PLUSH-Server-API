package main

import (
  "github.com/mitchellh/mapstructure"
  b64 "encoding/base64"
)


func getProfilePicture(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode getProfilePicture"}
    return
  }
  go func() {
    var user_id string
    var htmlEmbed string
    var string64 []byte
    if err := client.session.Query("SELECT * FROM profile_pictures WHERE user_id = ?",user.User_Id).Scan(&user_id, &htmlEmbed, &string64); err != nil {
      client.send <- Message{"profile picture get", "" }
      return
    }
    encodedString := b64.StdEncoding.EncodeToString(string64)
    //Constructing html base64 embeded image
    var base64EmbededImage = htmlEmbed + "," + encodedString
    var blob Blob
    blob.Data = base64EmbededImage
    blob.User_Id = user_id
    client.send <- Message{"profile picture get", blob}
  }()

}
