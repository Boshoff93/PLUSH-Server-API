package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
)

func getFollowCounts(w http.ResponseWriter, r *http.Request){
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
    var follow_counts FollowCounts
    var follow_id string
    var countFollowing int
    var countFollowers int

    iterFollowing := session.Query("SELECT follow_id FROM following WHERE user_id = ?", user.User_Id).Iter()
    for iterFollowing.Scan(&follow_id) {
      countFollowing++
    }
    follow_counts.FollowingCount = countFollowing

    iterFollowers := session.Query("SELECT follow_id FROM followers WHERE user_id = ?", user.User_Id).Iter()
    for iterFollowers.Scan(&follow_id) {
      countFollowers++
    }
    follow_counts.FollowerCount = countFollowers

    fmt.Println(follow_counts.FollowingCount)
    fmt.Println(follow_counts.FollowerCount)
    json.NewEncoder(w).Encode(follow_counts)
    finished <- true
  }()
  <- finished

}
