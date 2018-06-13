package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
)

func getDisplayName(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)
  var user User
  user.User_Id = params["user_id"]


  fmt.Println("AWE")
  finished := make(chan bool)
  go func() {

    if err := session.Query("SELECT display_name FROM users_by_id WHERE user_id = ?",user.User_Id).Scan(&user.Display_Name); err != nil {
      fmt.Println("Could not get display name, error: " + err.Error())
      json.NewEncoder(w).Encode(Error{Error:"Could not get user display name"})
      finished <- true
      return
    }
    json.NewEncoder(w).Encode(user)
    finished <- true
  }()
  <- finished

}
