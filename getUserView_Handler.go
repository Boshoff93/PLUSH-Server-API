package main

import (
  "github.com/mitchellh/mapstructure"
  "fmt"
)

func getUserView(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode getUserView"}
    return
  }

  go func() {
    var firstname string
    var lastname string
    var user_id string

    if(user.User_Id != "" && user.Email == "") {
      if err := client.session.Query("SELECT firstname, lastname, user_id FROM users_by_id WHERE user_id = ?",user.User_Id).Scan(&firstname, &lastname, &user_id); err != nil {
        fmt.Println("User Does Not Exist: user_by_id");
        return
      }
    } else {
      if err := client.session.Query("SELECT firstname, lastname, user_id FROM users_by_email WHERE email = ?",user.Email).Scan(&firstname, &lastname, &user_id); err != nil {
        fmt.Println("User Does Not Exist: user_by_email");
        return
      }
    }
    user.User_Id = user_id
    user.Firstname = firstname
    user.Lastname = lastname
    user.Email = ""
    client.send <- Message{"user get", user}
  }()
}
