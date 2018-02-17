package main

import (
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "net/http"
)

func getUserViewByUserId(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)
  var user User
  user.User_Id = params["user_id"]
  fmt.Println(user.User_Id)

  finished := make(chan bool)
  go func() {
    var firstname string
    var lastname string
    var user_id string
    if err := session.Query("SELECT firstname, lastname, user_id FROM users_by_id WHERE user_id = ?",user.User_Id).Scan(&firstname, &lastname, &user_id); err != nil {
      fmt.Println(err.Error())
      json.NewEncoder(w).Encode("User does not exist")
      finished <- true
      return
    }

    user.User_Id = user_id
    user.Firstname = firstname
    user.Lastname = lastname
    json.NewEncoder(w).Encode(user)
    finished <- true
  }()
  <- finished
}
