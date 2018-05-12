package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "time"
  "fmt"
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
    var following_posts FollowingPosts
    itr := session.Query("SELECT follow_id FROM following WHERE user_id = ?", user.User_Id).Iter()
    var follow_id_temp string
    var content string
    var pp_name string
    var post_time time.Time

    for itr.Scan(&follow_id_temp) {
      //For every follow_id run a query to get the neccasary information and add to a array in struct
      var display_name string
      session.Query("SELECT display_name FROM users_by_id WHERE user_id = ?",follow_id_temp).Scan(&display_name)
      itrFolID := session.Query("SELECT toTimeStamp(post_id), content FROM posts WHERE user_id = ?",follow_id_temp).Iter()

      if err := session.Query("SELECT pp_name FROM profile_picture_names WHERE user_id = ?",follow_id_temp).Scan(&pp_name); err != nil {
        fmt.Println("Could not find profile picture name, error: " + err.Error() )
        following_posts.Pp_Names = append(following_posts.Pp_Names, "empty")
      } else {
        following_posts.Pp_Names = append(following_posts.Pp_Names, pp_name)
      }
      following_posts.Unique_Following_Ids = append(following_posts.Unique_Following_Ids, follow_id_temp)

      for itrFolID.Scan(&post_time, &content) {
          following_posts.Display_Names = append(following_posts.Display_Names, display_name)
          following_posts.Following_Ids = append(following_posts.Following_Ids, follow_id_temp)
          following_posts.Post_Times = append(following_posts.Post_Times, post_time)
          following_posts.Posts = append(following_posts.Posts ,content)
        }
    }

    //This needs to be improved, bubble sort is not an efficient solution
    var done bool = false;
    for {
        done = true;
        for i := 0; i < len(following_posts.Post_Times) - 1; i++ {
          if(following_posts.Post_Times[i].Before(following_posts.Post_Times[i+1])) {
            done = false;
            following_posts.Post_Times[i], following_posts.Post_Times[i+1] = following_posts.Post_Times[i+1], following_posts.Post_Times[i]
            following_posts.Display_Names[i], following_posts.Display_Names[i+1] = following_posts.Display_Names[i+1], following_posts.Display_Names[i]
            following_posts.Following_Ids[i], following_posts.Following_Ids[i+1] = following_posts.Following_Ids[i+1], following_posts.Following_Ids[i]
            following_posts.Posts[i], following_posts.Posts[i+1] = following_posts.Posts[i+1], following_posts.Posts[i]
          }
        }
        if(done == true) {
          break
        }
    }

    json.NewEncoder(w).Encode(following_posts)
    finished <- true
  }()
  <- finished

}
