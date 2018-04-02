package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  //"time"
)

func getAllFollowPosts(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params := mux.Vars(r)
  var user User
  user.User_Id = params["user_id"]

  finished := make(chan bool)

  //Need to get all posts for users following group
  //Need to get all post times for users following group
  //Need to get profile pictures for users following group
  //Need to get all users display names for users following group

  //Need to get specified user_id's posts as well
  go func() {

    itr := session.Query("SELECT follow_id FROM following WHERE user_id = ?", user.User_Id).Iter()
    var follow_id_temp string
    for itr.Scan(&follow_id_temp) {
      //For every follow_id run a query to get the neccasary information and add to a array in struct
    }
    //Add ow user information
    //Sort each array

    // var posts Posts
    // var post_id string
    // var content string
    // var post_time time.Time
    // itr := session.Query("SELECT toTimeStamp(post_id), post_id, content FROM posts WHERE user_id = ?",user.User_Id).Iter()
    // for itr.Scan(&post_time,&post_id, &content) {
		//     posts.Post_Ids = append(posts.Post_Ids, post_id)
    //     posts.Post_Times = append(posts.Post_Times, post_time)
    //     posts.Posts = append(posts.Posts ,content)
    //   }

    json.NewEncoder(w).Encode("")
    finished <- true
  }()
  <- finished

}
