package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)


func deleteFollow(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  var ids IdFields
  if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }

  finished := make(chan bool)
  go func() {
    if err := session.Query("DELETE FROM following WHERE user_id = ? AND follow_id = ?",ids.User_Id, ids.Follow_Id).Exec(); err != nil {
      fmt.Println(err.Error());
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
    }
    if err := session.Query("DELETE FROM followers WHERE user_id = ? AND follow_id = ?",ids.Follow_Id, ids.User_Id).Exec(); err != nil {
      fmt.Println(err.Error());
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
    }
    json.NewEncoder(w).Encode(ids)
    finished <- true
  }()
  <- finished
}
