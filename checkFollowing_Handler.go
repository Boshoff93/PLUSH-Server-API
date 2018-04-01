package main

import (
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
  "strings"
)

func checkFollowing(w http.ResponseWriter, r *http.Request){
  session := getSession()
  defer session.Close()

  params:= mux.Vars(r)

  var ids = params["id_fields"]
  var id_fields []string
  id_fields = strings.Split(ids, ",")
  var user_id string = id_fields[0]
  var follower_id string = id_fields[1]
  var follower_id_temp string
  var is_following BoolValue
  is_following.Condition = false;
  finished := make(chan bool)
  go func() {
    itr := session.Query("SELECT follower_id FROM following WHERE user_id = ?",user_id).Iter()
    for itr.Scan(&follower_id_temp) {
      if(follower_id_temp == follower_id){
        is_following.Condition = true;
        break;
      }
    }
    json.NewEncoder(w).Encode(is_following)
    finished <- true
    }()
  <- finished
}
