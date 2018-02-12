package main

import (
  "github.com/mitchellh/mapstructure"
)

func addUser(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode addUser"}
    return
  }
  hashedPassword, _ := HashPassword(user.Password)
  user.Password = hashedPassword

  go func() {
  var email string
  if err := client.session.Query("SELECT email FROM users_by_email WHERE email = ?",user.Email).Scan(&email); err != nil {

    if err := client.session.Query("INSERT INTO users_by_email (email, user_id, created_at, firstname, lastname, password) VALUES (?,?,?,?,?,?)",
                                    user.Email, user.User_Id, user.Created_At, user.Firstname, user.Lastname, user.Password).Exec(); err != nil {
      client.send <- Message{"error", "could add user to users_by_email"}
      return
    }

    if err := client.session.Query("INSERT INTO users_by_id (email, user_id, created_at, firstname, lastname, password) VALUES (?,?,?,?,?,?)",
                                    user.Email, user.User_Id, user.Created_At, user.Firstname, user.Lastname, user.Password).Exec(); err != nil {
      client.send <- Message{"error", "could add user to users_by_id"}
      return
    }
    var userAdded User
    userAdded.Firstname = user.Firstname
    userAdded.Lastname = user.Lastname
    userAdded.Email = user.Email
    userAdded.User_Id = user.User_Id
    client.send <- Message{"user add", userAdded}
    return
  }
    client.send <- Message{"email unavailible",""}
  }()

}
