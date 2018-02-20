package main

import (
  // "encoding/json"
  // "fmt"
  // "github.com/gorilla/mux"
  // "net/http"
)

//func findUser(w http.ResponseWriter, r *http.Request){
  // session := getSession()
  // defer session.Close()
  //
  // params:= mux.Vars(r)
  // var user User
  // user.Email = params["email"]
  //
  // finished := make(chan bool)
  // go func() {
  //   var email string
  //   var firstname string
  //   var lastname string
  //   var user_id string
  //   if err := session.Query("SELECT email, firstname, lastname, user_id FROM users_by_email WHERE email = ?",
  //                                   user.Email).Scan(&email, &firstname, &lastname, &user_id); err != nil {
  //     fmt.Println(err.Error())
  //     json.NewEncoder(w).Encode("access denied")
  //     finished <- true
  //     return
  //   }
  //     var userFound User
  //     userFound.Firstname = firstname
  //     userFound.Lastname = lastname
  //     userFound.Email = email
  //     userFound.User_Id = user_id
  //
  //     json.NewEncoder(w).Encode(userFound)
  //     finished <- true
  // }()
  //  <- finished
//}
