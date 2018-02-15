package main

import (
  "fmt"
  "net/http"
  "encoding/json"
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
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
    finished <- true
  }()
  <- finished
}
