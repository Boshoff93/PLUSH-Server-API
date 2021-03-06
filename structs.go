package main

import (
    "time"
)

type User struct {
  User_Id       string            `json: "user_id"       cql:"text"`
  Display_Name  string            `json: "display_name"  cql:"text"`
  Email         string            `json: "email"         cql:"text"`
  Created_At    string            `json: "created_at"    cql:"timeuuid"`
  Token         string            `json: "token"`
}

type SearchUser struct {
  Search        string          `cql:"text"`
}

type SearchedUsers struct {
  User_Ids         []string         `json: "user_ids"`
  Display_Names    []string         `json: "display_names"`
  Pp_Names         []string         `json: "pp_names"`
}

type Post struct {
  User_Id              string            `json: "user_id" cql:"text"`
  Post_Id              string            `json: "post_id" cql:"timeuuid"`
  Post                 string            `json: "post"    cql:"text"`
  Type_Of_Post         int               `json: "type"    cql:"text"`
}

type Post_User_Id struct {
  User_Id     string            `json: "user_id" cql:"text"`
  Post_Id     string            `json: "post_id" cql:"timeuuid"`
}

type PostAdded struct {
  User_Id       string          `json: "user_id"`
  Post_Time     time.Time       `json: "post_time"`
  Post_Id       string          `json: "post_id"`
  Post          string          `json: "post"`
  Type_Of_Post  int             `json: "type_of_post"`
}

type Posts struct {
  Post_Ids           []string         `json: "post_ids"`
  Post_Times         []time.Time      `json: "post_times"`
  Posts              []string         `json: "posts"`
  Types_Of_Posts     []int            `json: "types_of_posts"`
}

type FollowingPosts struct {
  Unique_Following_Ids      []string          `json: "unique_following_ids"`
  Pp_Names                  []string          `json: "pp_names"`
  Display_Names             []string          `json: "display_names"`
  Following_Ids             []string          `json: "following_ids"`
  Post_Times                []time.Time       `json: "post_times"`
  Posts                     []string          `json: "posts"`
  Post_Ids                  []string          `json: "post_ids" cql:"timeuuid"`
  Types_Of_Posts            []int             `json: "types_of_posts"`
}

type FollowersAndFollowings struct {
  Follower_Ids                  []string          `json: "follower_ids"`
  Following_Ids                 []string          `json: "following_ids"`
  Follower_Pp_Names             []string          `json: "follower_pp_names"`
  Following_Pp_Names            []string          `json: "following_pp_names"`
  Follower_Display_Names        []string          `json: "follower_display_names"`
  Following_Display_Names       []string          `json: "follower_display_names"`
  FollowingCount                int               `json: "followings_count"`
  FollowerCount                 int               `json: "followers_count"`
}

type FollowCounts struct {
  FollowingCount    int     `json: "followings_count"`
  FollowerCount     int     `json: "followers_count"`
}

type Blob struct {
	User_Id     string             `cql:"uuid"`
	Pp_Name     string             `cql:"text"`
}

type Error struct {
  Error    string               `json: "error"`
}

type IdFields struct {
  User_Id    string             `json: "user_id" cql:"text"`
  Follow_Id  string           `json: "follow_id" cql:"text"`
}

type EditDisplayName struct {
  User_Id    string             `json: "user_id" cql:"text"`
  Display_Name  string           `json: "display_name" cql:"text"`
}

type BoolValue struct {
  Condition   bool              `json: "condition"`
}

type Posts_Likes_Dislikes struct {
  Post_Ids                  []string         `json: "post_ids" cql:"timeuuid"`
  User_Ids                  []string          `json: "user_ids"`
  Likes                     []int             `json: "likes"`
  Dislikes                  []int             `json: "dislikes"`
  TotalLikes                []int             `json: "total_likes"`
  TotalDislikes             []int             `json: "total_dislikes"`
}


type Message struct {
  Name string                   `json:"name"`
  Data interface{}              `json:"data"`
}
