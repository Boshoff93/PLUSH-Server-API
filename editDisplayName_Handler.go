package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)

func editDisplayName(w http.ResponseWriter, r *http.Request){
  fmt.Println("Got here now")
  session := getSession()
  defer session.Close()
  var displayName EditDisplayName

  if err := json.NewDecoder(r.Body).Decode(&displayName); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }

  finished := make(chan bool)
  go func() {
    var user User

    if err := session.Query("SELECT user_id, display_name, created_at, email FROM users_by_id WHERE user_id = ?",
                            displayName.User_Id).Scan(&user.User_Id, &user.Display_Name, &user.Created_At, &user.Email); err != nil {
      fmt.Println("Could not get user, error: " + err.Error())
      json.NewEncoder(w).Encode(Error{Error:"Could not get user to users_by_id"})
      finished <- true
      return
    }

    user.Display_Name = displayName.Display_Name

    if err := session.Query("DELETE FROM users_by_id WHERE user_id = ?",user.User_Id).Exec(); err != nil {
      fmt.Println(err.Error());
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
    }

    if err := session.Query("DELETE FROM users_by_id_like_display_name WHERE user_id = ?",user.User_Id).Exec(); err != nil {
      fmt.Println(err.Error());
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
    }

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

    finished <- true
    json.NewEncoder(w).Encode(displayName)
  }()
  <- finished
}
