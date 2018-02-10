package main

import (
  "github.com/mitchellh/mapstructure"
  "fmt"
  "time"
  "github.com/gocql/gocql"
  "golang.org/x/crypto/bcrypt"
)

type User struct {
  User_Id     string      `cql:"uuid"`
  Firstname   string      `cql:"firstname"`
  Lastname    string      `cql:"lastname"`
  Email       string      `cql:"email"`
  Password    string      `cql:"password"`
  Created_At  string      `cql:"timeuuid"`
}

type Post struct {
  User_Id     string `cql:"uuid"`
  Post_Id     string `cql:"timeuuid"`
  Post        string
}

type PostAdded struct {
  User_Id       string
  Post_Time     time.Time
  Post_Id       string
  Post          string
}

type Posts struct {
  Post_Ids     []string
  Post_Times   []time.Time
  Posts       []string
}

type Blob struct {
	User_Id  string    `cql:"uuid"`
	Data     string    `cql:"text"`// Formatted as Base64 but I would prefer Base64Url...
}


type Message struct {
  Name string `json:"name"`
  Data interface{} `json:"data"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func findUser(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode findUser"}
    return
  }

  go func() {
    var email string
    var firstname string
    var lastname string
    var password string
    var user_id string
    if err := client.session.Query("SELECT email, firstname, lastname, password, user_id FROM users_by_email WHERE email = ?",
                                    user.Email).Scan(&email, &firstname, &lastname, &password, &user_id); err != nil {

      client.send <- Message{"account not found", ""}
      return
    }

    match := CheckPasswordHash(user.Password, password)
    if(match) {
      var userFound User
      userFound.Firstname = firstname
      userFound.Lastname = lastname
      userFound.Email = email
      userFound.User_Id = user_id
      client.send <- Message{"access granted", userFound}
    } else {
      client.send <- Message{"access denied", ""}
    }
  }()
}

func addUser(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode addUser"}
    return
  }
  hashedPassword, _ := HashPassword(user.Password)
  user.Password = hashedPassword

  go func() {
  var email string
  if err := client.session.Query("SELECT email FROM users_by_email WHERE email = ?",user.Email).Scan(&email); err != nil {

    if err := client.session.Query("INSERT INTO users_by_email (email, user_id, created_at, firstname, lastname, password) VALUES (?,?,?,?,?,?)",
                                    user.Email, user.User_Id, user.Created_At, user.Firstname, user.Lastname, user.Password).Exec(); err != nil {
      client.send <- Message{"error", "could add user to users_by_email"}
      return
    }

    if err := client.session.Query("INSERT INTO users_by_id (email, user_id, created_at, firstname, lastname, password) VALUES (?,?,?,?,?,?)",
                                    user.Email, user.User_Id, user.Created_At, user.Firstname, user.Lastname, user.Password).Exec(); err != nil {
      client.send <- Message{"error", "could add user to users_by_id"}
      return
    }
    var userAdded User
    userAdded.Firstname = user.Firstname
    userAdded.Lastname = user.Lastname
    userAdded.Email = user.Email
    userAdded.User_Id = user.User_Id
    client.send <- Message{"user add", userAdded}
    return
  }
    client.send <- Message{"email unavailible",""}
  }()

}

func getUserView(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode getUserView"}
    return
  }
  fmt.Println(user.Email)
  go func() {

  var firstname string
  var lastname string
  var user_id string

  if(user.User_Id != "" && user.Email == "") {
    if err := client.session.Query("SELECT firstname, lastname, user_id FROM users_by_id WHERE user_id = ?",user.User_Id).Scan(&firstname, &lastname, &user_id); err != nil {
      fmt.Println("User Does Not Exist: user_by_id");
      return
    }
  } else {
    if err := client.session.Query("SELECT firstname, lastname, user_id FROM users_by_email WHERE email = ?",user.Email).Scan(&firstname, &lastname, &user_id); err != nil {
      fmt.Println("User Does Not Exist: user_by_email");
      return
    }
  }
  
  user.User_Id = user_id
  user.Firstname = firstname
  user.Lastname = lastname
  user.Email = ""
  client.send <- Message{"user get", user}
  }()
}

func addPost(client *Client, data interface{}){
  var post Post
  err := mapstructure.Decode(data, &post)
  if err != nil {
    client.send <- Message{"error", "could not decode addPost"}
    return
  }
  go func() {
    if err := client.session.Query("INSERT INTO posts (user_id, post_id, content) VALUES (?,?,?)",post.User_Id, post.Post_Id , post.Post).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    //Convert string uuidv1 to uuidv1 then extract the Time before sending it back to the client
    tempUUID, err := gocql.ParseUUID(post.Post_Id);
    if err != nil {
       client.send <- Message{"error", err.Error()}
		   fmt.Printf("Something went wrong: %s", err)
	  }
    var postAdded PostAdded
    postAdded.User_Id = post.User_Id
    postAdded.Post_Time = tempUUID.Time()
    postAdded.Post_Id = post.Post_Id
    postAdded.Post = post.Post
    client.send <- Message{"post add", postAdded}
  }()

}

func deletePost(client *Client, data interface{}){
  var post Post
  err := mapstructure.Decode(data, &post)
  if err != nil {
    client.send <- Message{"error", "could not decode addPost"}
    return
  }
   go func() {
    if err := client.session.Query("DELETE FROM posts WHERE user_id = ? AND post_id = ?",post.User_Id, post.Post_Id).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    client.send <- Message{"post delete", post}
  }()
}

func getPosts(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode getPosts"}
    return
  }
  go func() {
    var posts Posts
    var post_id string
    var content string
    var post_time time.Time
    itr := client.session.Query("SELECT toTimeStamp(post_id), post_id, content FROM posts WHERE user_id = ?",user.User_Id).Iter()
    for itr.Scan(&post_time,&post_id, &content) {
		    posts.Post_Ids = append(posts.Post_Ids, post_id)
        posts.Post_Times = append(posts.Post_Times, post_time)
        posts.Posts = append(posts.Posts ,content)
      }
      fmt.Println(posts)
    client.send <- Message{"posts get", posts}
  }()

}

func addProfilePicture(client *Client, data interface{}){
  var blob Blob
  err := mapstructure.Decode(data, &blob)
  if err != nil {
    client.send <- Message{"error", "could not decode addProfilePicture"}
    return
  }
  go func() {
    if err := client.session.Query("INSERT INTO profile_pictures (user_id, profile_picture) VALUES (?,?)",blob.User_Id, blob.Data).Exec(); err != nil {
      fmt.Println(err.Error());
    }
    client.send <- Message{"profile picture add", blob}
  }()

}

func getProfilePicture(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil {
    client.send <- Message{"error", "could not decode getProfilePicture"}
    return
  }
  go func() {
    var image string
    var user_id string
    if err := client.session.Query("SELECT * FROM profile_pictures WHERE user_id = ?",user.User_Id).Scan(&user_id, &image); err != nil {
      client.send <- Message{"profile picture get", "" }
      return
    }
    var blob Blob
    blob.Data = image
    blob.User_Id = user_id
    client.send <- Message{"profile picture get", blob}
  }()

}
