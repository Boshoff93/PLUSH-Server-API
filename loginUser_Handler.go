package main

import (

  "fmt"
  "google.golang.org/api/oauth2/v2"
  "net/http"
  "encoding/json"
)


func login(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()
  tokenId := r.Header.Get("Authorization")
  var user User

  if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }

  var httpClient = &http.Client{}
  oauth2Service, err := oauth2.New(httpClient)
  tokenInfoCall := oauth2Service.Tokeninfo()
  tokenInfoCall.IdToken(tokenId)

  tokenInfo, err := tokenInfoCall.Do()
  if err != nil {
    return
  }
  if(tokenInfo.Audience != "729356241272-g4gtvdrpvhsts6ogat3n8kv0ma1vidhm.apps.googleusercontent.com") {
    json.NewEncoder(w).Encode("access denied")
    return
  }
  fmt.Println("Access Granted")

  user.Email = tokenInfo.Email
  var googleId = tokenInfo.UserId
  finished := make(chan bool)

  go func() {
     var email string
     var firstname string
     var lastname string
     var user_id string
     if err := session.Query("SELECT email, firstname, lastname, user_id FROM users_by_google_id WHERE google_id = ?",
                             googleId).Scan(&email, &firstname, &lastname, &user_id); err != nil {

      if err := session.Query("INSERT INTO users_by_google_id (google_id, email, user_id, created_at, firstname, lastname) VALUES (?,?,?,?,?,?)",
                              googleId, user.Email, user.User_Id, user.Created_At, user.Firstname, user.Lastname).Exec(); err != nil {
        fmt.Println("Could not add user to users_by_google_id, error: " + err.Error())
        finished <- true
        return
      }
      fmt.Println("user created")
      json.NewEncoder(w).Encode(user)
      finished <- true
      return
     }
     fmt.Println("user found")
     var userFound User
     userFound.Firstname = firstname
     userFound.Lastname = lastname
     userFound.Email = email
     userFound.User_Id = user_id

     json.NewEncoder(w).Encode(userFound)
     finished <- true
   }()
    <- finished
}
