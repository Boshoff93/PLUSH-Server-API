package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)


func deletePost(w http.ResponseWriter, r *http.Request){
  var post Post
  session := getSession()
  defer session.Close()

  if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }

  finished := make(chan bool)
  go func() {
    if err := session.Query("DELETE FROM posts WHERE user_id = ? AND post_id = ?",post.User_Id, post.Post_Id).Exec(); err != nil {
      fmt.Println(err.Error());
      json.NewEncoder(w).Encode(Error{Error: err.Error()})
      finished <- true
      return
    }
    json.NewEncoder(w).Encode(post)
    finished <- true
  }()
  <- finished
}
