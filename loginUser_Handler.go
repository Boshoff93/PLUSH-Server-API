package main

import (

  "fmt"
  "google.golang.org/api/oauth2/v2"
  "net/http"
  "encoding/json"
  "github.com/dgrijalva/jwt-go"

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
    json.NewEncoder(w).Encode(Error{Error:"Access Denied"})
    return
  }
  fmt.Println("Access Granted")
  //Give uswer a token
  //**************************************************************************
  accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
   "username": user.Display_Name,
   "user_id": user.User_Id,
   })
  tokenString, err := accessToken.SignedString([]byte("MyFancySecret"))
  if err != nil {
   fmt.Println("Token signed string failed: " + err.Error())
   json.NewEncoder(w).Encode(Error{Error: err.Error()})
  }
  //**************************************************************************
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
        json.NewEncoder(w).Encode(Error{Error:"Could not add user to users_by_id"})
        finished <- true
        return
      }
      if err := session.Query("INSERT INTO users_by_id_like_display_name (user_id, email, created_at, display_name) VALUES (?,?,?,?)",
                              user.User_Id, user.Email, user.Created_At, user.Display_Name).Exec(); err != nil {
        fmt.Println("Could not add user to users_by_id_like_diplay_name, error: " + err.Error())
        json.NewEncoder(w).Encode(Error{Error:"Could not add user to users_by_id_like_diplay_name"})
        finished <- true
        return
      }
      user.Token = tokenString
      json.NewEncoder(w).Encode(user)
      finished <- true
      return
     }
     var userFound User
     userFound.Display_Name = display_name
     userFound.Email = email
     userFound.User_Id = user_id
     userFound.Token = tokenString
     json.NewEncoder(w).Encode(userFound)
     finished <- true
   }()
    <- finished
}
