package main

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "fmt"
)

func getUserViewByEmail(w http.ResponseWriter, r *http.Request){

  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)
  var user User
  user.Email = params["email"]

  finished := make(chan bool)
  go func() {
    var firstname string
    var lastname string
    var user_id string

    if err := session.Query("SELECT firstname, lastname, user_id FROM users_by_email WHERE email = ?",user.Email).Scan(&firstname, &lastname, &user_id); err != nil {
      fmt.Println("User Does Not Exist: user_by_email");
      return
    }
    user.User_Id = user_id
    user.Firstname = firstname
    user.Lastname = lastname
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
    finished <- true
  }()
  <- finished
}
