package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
  "strings"
)

func getFollowersAndFollowings(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params := mux.Vars(r)
  var user User
  user.User_Id = params["user_id"]
  finished := make(chan bool)

  go func() {
    var followers_and_followings FollowersAndFollowings
    var follow_id string
    var follow_display_name string
    var follow_pp_name string
    var countFollowing int
    var countFollowers int

    iterFollowing := session.Query("SELECT follow_id FROM following WHERE user_id = ?", user.User_Id).Iter()
    for iterFollowing.Scan(&follow_id) {
      countFollowing++
      followers_and_followings.Following_Ids = append(followers_and_followings.Following_Ids, follow_id);
      session.Query("SELECT display_name FROM users_by_id WHERE user_id = ?", follow_id).Scan(&follow_display_name)
      followers_and_followings.Following_Display_Names = append(followers_and_followings.Following_Display_Names, follow_display_name)
      if err := session.Query("SELECT pp_name FROM profile_picture_names WHERE user_id = ?", follow_id).Scan(&follow_pp_name); err != nil {
        followers_and_followings.Following_Pp_Names = append(followers_and_followings.Following_Pp_Names, "empty")
      } else {
        followers_and_followings.Following_Pp_Names = append(followers_and_followings.Following_Pp_Names, follow_pp_name)
      }
    }


    iterFollowers := session.Query("SELECT follow_id FROM followers WHERE user_id = ?", user.User_Id).Iter()
    for iterFollowers.Scan(&follow_id) {
      countFollowers++
      followers_and_followings.Follower_Ids = append(followers_and_followings.Follower_Ids, follow_id);
      session.Query("SELECT display_name FROM users_by_id WHERE user_id = ?", follow_id).Scan(&follow_display_name)
      followers_and_followings.Follower_Display_Names = append(followers_and_followings.Follower_Display_Names, follow_display_name)
      if err := session.Query("SELECT pp_name FROM profile_picture_names WHERE user_id = ?", follow_id).Scan(&follow_pp_name) ; err != nil {
        followers_and_followings.Follower_Pp_Names = append(followers_and_followings.Follower_Pp_Names, "empty")
      } else {
        followers_and_followings.Follower_Pp_Names = append(followers_and_followings.Follower_Pp_Names, follow_pp_name)
      }
    }

    followers_and_followings.FollowingCount = countFollowing;
    followers_and_followings.FollowerCount = countFollowers;

    //This needs to be improved, bubble sort is not an efficient solution
    var done bool = false;
    for {
        done = true;
        for i := 0; i < len(followers_and_followings.Following_Display_Names) - 1; i++ {
          if(strings.Compare(strings.ToLower(followers_and_followings.Following_Display_Names[i]), strings.ToLower(followers_and_followings.Following_Display_Names[i+1])) > 0) {
            done = false;
            followers_and_followings.Following_Display_Names[i], followers_and_followings.Following_Display_Names[i+1] = followers_and_followings.Following_Display_Names[i+1], followers_and_followings.Following_Display_Names[i]
            followers_and_followings.Following_Pp_Names[i], followers_and_followings.Following_Pp_Names[i+1] = followers_and_followings.Following_Pp_Names[i+1], followers_and_followings.Following_Pp_Names[i]
            followers_and_followings.Following_Ids[i], followers_and_followings.Following_Ids[i+1] = followers_and_followings.Following_Ids[i+1], followers_and_followings.Following_Ids[i]
          }
        }
        if(done == true) {
          break
        }
    }

    //This needs to be improved, bubble sort is not an efficient solution
    done = false;
    for {
        done = true;
        for i := 0; i < len(followers_and_followings.Follower_Display_Names) - 1; i++ {
          if(strings.Compare(strings.ToLower(followers_and_followings.Follower_Display_Names[i]), strings.ToLower(followers_and_followings.Follower_Display_Names[i+1])) > 0 ) {
            done = false;
            followers_and_followings.Follower_Display_Names[i], followers_and_followings.Follower_Display_Names[i+1] = followers_and_followings.Follower_Display_Names[i+1], followers_and_followings.Follower_Display_Names[i]
            followers_and_followings.Follower_Pp_Names[i], followers_and_followings.Follower_Pp_Names[i+1] = followers_and_followings.Follower_Pp_Names[i+1], followers_and_followings.Follower_Pp_Names[i]
            followers_and_followings.Follower_Ids[i], followers_and_followings.Follower_Ids[i+1] = followers_and_followings.Follower_Ids[i+1], followers_and_followings.Follower_Ids[i]
          }
        }
        if(done == true) {
          break
        }
    }

    fmt.Println(followers_and_followings)

    json.NewEncoder(w).Encode(followers_and_followings)
    finished <- true
  }()
  <- finished

}
