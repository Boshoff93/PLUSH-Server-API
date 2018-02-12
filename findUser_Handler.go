package main

import (
  "github.com/mitchellh/mapstructure"
)

func findUser(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode findUser"}
    return
  }

  go func() {
    var email string
    var firstname string
    var lastname string
    var password string
    var user_id string
    if err := client.session.Query("SELECT email, firstname, lastname, password, user_id FROM users_by_email WHERE email = ?",
                                    user.Email).Scan(&email, &firstname, &lastname, &password, &user_id); err != nil {

      client.send <- Message{"account not found", ""}
      return
    }

    match := CheckPasswordHash(user.Password, password)
    if(match) {
      var userFound User
      userFound.Firstname = firstname
      userFound.Lastname = lastname
      userFound.Email = email
      userFound.User_Id = user_id
      client.send <- Message{"access granted", userFound}
    } else {
      client.send <- Message{"access denied", ""}
    }
  }()
}
