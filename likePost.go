package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)
func likePost(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  var post_user_id Post_User_Id
  if err := json.NewDecoder(r.Body).Decode(&post_user_id); err != nil {
          http.Error(w, err.Error(), 400)
          return
  }
  finished := make(chan bool)
  var posts_likes_dislikes Posts_Likes_Dislikes
  go func() {
    var post_id string
    var user_id string
    var like int
    var dislike int

    if err := session.Query("SELECT post_id, user_id, like, dislike FROM posts_likes_dislikes WHERE post_id = ? AND user_id= ?",
                                  post_user_id.Post_Id,post_user_id.User_Id).Scan(&post_id, &user_id, &like, &dislike); err != nil {
        if err := session.Query("INSERT INTO posts_likes_dislikes (post_id, user_id, like, dislike) VALUES (?,?,?,?)",
                                              post_user_id.Post_Id, post_user_id.User_Id , 1, 0).Exec(); err != nil {
          fmt.Println(err.Error());
          json.NewEncoder(w).Encode(Error{Error: err.Error()})
          finished <- true
          return
        }

        posts_likes_dislikes = getPosts_Likes_Dislikes(session, posts_likes_dislikes, post_user_id)
        json.NewEncoder(w).Encode(posts_likes_dislikes)
        finished <- true
        return
    }

    if(like == 0) {
      if err := session.Query("UPDATE posts_likes_dislikes SET like=?, dislike=? WHERE post_id=? AND user_id=?",1,0,post_user_id.Post_Id, post_user_id.User_Id).Exec(); err != nil {
        fmt.Println(err.Error());
        json.NewEncoder(w).Encode(Error{Error: err.Error()})
        finished <- true
        return
      }
    } else {
      if err := session.Query("UPDATE posts_likes_dislikes SET like=?, dislike=? WHERE post_id=? AND user_id=?",0,0,post_user_id.Post_Id, post_user_id.User_Id).Exec(); err != nil {
        fmt.Println(post_user_id)
        fmt.Println(err.Error());
        json.NewEncoder(w).Encode(Error{Error: err.Error()})
        finished <- true
        return
      }
    }

    posts_likes_dislikes = getPosts_Likes_Dislikes(session, posts_likes_dislikes, post_user_id)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts_likes_dislikes)
    finished <- true
  }()
  <- finished
}
