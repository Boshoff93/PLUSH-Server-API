package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "time"
)

func getPosts(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)
  var user User
  user.User_Id = params["user_id"]

  finished := make(chan bool)
  go func() {
    var posts Posts
    var post_id string
    var content string
    var post_time time.Time
    var type_of_post int
    itr := session.Query("SELECT toTimeStamp(post_id), post_id, content, type FROM posts WHERE user_id = ?",user.User_Id).Iter()
    for itr.Scan(&post_time,&post_id, &content, &type_of_post) {
		    posts.Post_Ids = append(posts.Post_Ids, post_id)
        posts.Post_Times = append(posts.Post_Times, post_time)
        posts.Posts = append(posts.Posts ,content)
        posts.Types_Of_Posts = append(posts.Types_Of_Posts, type_of_post)
      }

    json.NewEncoder(w).Encode(posts)
    finished <- true
  }()
  <- finished

}
