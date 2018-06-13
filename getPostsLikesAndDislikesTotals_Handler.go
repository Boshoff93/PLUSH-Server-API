package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "strings"
)

func getPostsLikesAndDislikesTotals(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)
  var posts_likes_dislikes Posts_Likes_Dislikes
  var post_ids_string = params["post_ids"]
  fmt.Println(post_ids_string)

  var post_id string
  var user_id string
  var like int
  var dislike int

  var likesTotal int = 0;
  var dislikesTotal int = 0;

  if(post_ids_string != "null") {
    posts_likes_dislikes.Post_Ids = strings.Split(post_ids_string, ","); // for one name
    finished := make(chan bool)
    go func() {
      for _, element := range posts_likes_dislikes.Post_Ids {

        itr := session.Query("SELECT * FROM posts_likes_dislikes WHERE post_id = ?",element).Iter()
        for itr.Scan(&post_id, &user_id, &dislike, &like) {
          if(like == 1) {
            likesTotal++;
          }

          if(dislike == 1) {
            dislikesTotal++;
          }
        }
        posts_likes_dislikes.TotalLikes = append(posts_likes_dislikes.TotalLikes, likesTotal);
        posts_likes_dislikes.TotalDislikes = append(posts_likes_dislikes.TotalDislikes, dislikesTotal);
        likesTotal = 0;
        dislikesTotal = 0;
      }
      json.NewEncoder(w).Encode(posts_likes_dislikes)
      finished <- true
    }()
    <- finished
    } else {
     var empty []string
     json.NewEncoder(w).Encode(empty)
   }

}
