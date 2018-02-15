package main

import (
  "net/http"
  "encoding/json"
  "fmt"
)

func addUser(w http.ResponseWriter, r *http.Request){

  session := getSession()
  defer session.Close()

  var user User
  if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }
  fmt.Println(user)

  finished := make(chan bool)
  go func() {
    var email string
    if err := session.Query("SELECT email FROM users_by_email WHERE email = ?",user.Email).Scan(&email); err != nil {

      if err := session.Query("INSERT INTO users_by_email (email, user_id, created_at, firstname, lastname) VALUES (?,?,?,?,?)",
                                      user.Email, user.User_Id, user.Created_At, user.Firstname, user.Lastname).Exec(); err != nil {
        fmt.Println("Could not add user to users_by_email, error: " + err.Error())
        return
      }

      if err := session.Query("INSERT INTO users_by_id (email, user_id, created_at, firstname, lastname) VALUES (?,?,?,?,?)",
                                      user.Email, user.User_Id, user.Created_At, user.Firstname, user.Lastname).Exec(); err != nil {
        fmt.Println("Could not add user to users_by_id, error: " + err.Error())
        return
      }

      fullname := user.Firstname + " " + user.Lastname
      if err := session.Query("INSERT INTO users_by_email_like_fullname (fullname, email, user_id, created_at) VALUES (?,?,?,?)",
                                      fullname, user.Email, user.User_Id, user.Created_At).Exec(); err != nil {
        fmt.Println("Could not add user to users_by_email_like_fullname, error: " + err.Error())
        return
      }

      var userAdded User
      userAdded.Firstname = user.Firstname
      userAdded.Lastname = user.Lastname
      userAdded.Email = user.Email
      userAdded.User_Id = user.User_Id
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(userAdded)
    }
    finished <- true
  }()
  <- finished
}
