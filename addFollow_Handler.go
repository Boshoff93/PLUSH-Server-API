package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)

func addFollow(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  var ids IdFields
  if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }

  finished := make(chan bool)
  go func() {
    if err := session.Query("INSERT INTO following (user_id, follower_id) VALUES (?,?)",ids.User_Id, ids.Follower_Id).Exec(); err != nil {
      fmt.Println(err.Error());
      finished <- true
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
    }
    if err := session.Query("INSERT INTO followers (user_id, follower_id) VALUES (?,?)", ids.Follower_Id, ids.User_Id).Exec(); err != nil {
      fmt.Println(err.Error());
      finished <- true
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
    }
    finished <- true
    json.NewEncoder(w).Encode(ids)
  }()
  <- finished
}
