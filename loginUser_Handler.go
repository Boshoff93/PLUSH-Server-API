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
  fmt.Println(tokenInfo.Audience)
  if(tokenInfo.Audience != "729356241272-e9plc3mof68m60u319bjsu46b3udfukv.apps.googleusercontent.com") {
    json.NewEncoder(w).Encode("access denied")
    return
  }
  fmt.Println("Access Granted")

  user.Email = tokenInfo.Email
  finished := make(chan bool)

  go func() {
     var email string
     var display_name string
     var user_id string
     if err := session.Query("SELECT email, display_name, user_id FROM users_by_id WHERE user_id = ?",
                             user.User_Id).Scan(&email, &display_name, &user_id); err != nil {

      if err := session.Query("INSERT INTO users_by_id (user_id, email, created_at, display_name) VALUES (?,?,?,?)",
                              user.User_Id, user.Email, user.Created_At, user.Display_Name).Exec(); err != nil {
        fmt.Println("Could not add user to users_by_id, error: " + err.Error())
        finished <- true
        return
      }

    
      var user_id string
      var email string
      var fullname string
      itr := session.Query("SELECT email, display_name, user_id FROM users_by_id_like_display_name WHERE display_name Like ? LIMIT 50","%Wi%").Iter()
      for itr.Scan(&email,&fullname, &user_id) {
  		    fmt.Println("ssssssssssssssssssssssssssssss=-----------------" + email + fullname + user_id)
        }



      if err := session.Query("INSERT INTO users_by_id_like_display_name (user_id, email, created_at, display_name) VALUES (?,?,?,?)",
                              user.User_Id, user.Email, user.Created_At, user.Display_Name).Exec(); err != nil {
        fmt.Println("Could not add user to users_by_id_like_diplay_name, error: " + err.Error())
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
     userFound.Display_Name = display_name
     userFound.Email = email
     userFound.User_Id = user_id

     json.NewEncoder(w).Encode(userFound)
     finished <- true
   }()
    <- finished
}
