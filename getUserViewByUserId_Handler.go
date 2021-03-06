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
    var display_name string
    var user_id string
    if err := session.Query("SELECT display_name, user_id FROM users_by_id WHERE user_id = ?",user.User_Id).Scan(&display_name, &user_id); err != nil {
      fmt.Println(err.Error())
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
    }
    user.User_Id = user_id
    user.Display_Name = display_name
    json.NewEncoder(w).Encode(user)
    finished <- true
  }()
  <- finished
}
