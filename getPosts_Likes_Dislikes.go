package main

import (
    "github.com/gocql/gocql"
)

func getPosts_Likes_Dislikes(session *gocql.Session, posts_likes_dislikes Posts_Likes_Dislikes, post_user_id Post_User_Id) Posts_Likes_Dislikes {

  var post_id string
  var user_id string
  var like int
  var dislike int

  itr := session.Query("SELECT * FROM posts_likes_dislikes WHERE user_id = ?",post_user_id.User_Id).Iter()
  for itr.Scan(&post_id, &user_id, &dislike, &like) {
    posts_likes_dislikes.Post_Ids = append(posts_likes_dislikes.Post_Ids, post_id)
    posts_likes_dislikes.User_Ids = append(posts_likes_dislikes.User_Ids, user_id)
    posts_likes_dislikes.Dislikes = append(posts_likes_dislikes.Dislikes ,dislike)
    posts_likes_dislikes.Likes  = append(posts_likes_dislikes.Likes ,like)
  }
  return posts_likes_dislikes;
}
