package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
)

func getPostsLikesAndDislikes(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)
  var user Post_User_Id
  user.User_Id = params["user_id"]


  finished := make(chan bool)
  var posts_likes_dislikes Posts_Likes_Dislikes
  go func() {
    posts_likes_dislikes = getPosts_Likes_Dislikes(session, posts_likes_dislikes, user)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts_likes_dislikes)
    finished <- true
  }()
  <- finished
}
