package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
)

func getFollowCounts(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params := mux.Vars(r)
  var user User
  user.User_Id = params["user_id"]
  finished := make(chan bool)

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
    json.NewEncoder(w).Encode(follow_counts)
    finished <- true
  }()
  <- finished

}
