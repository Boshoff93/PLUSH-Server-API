package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
)

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        accessToken := req.Header.Get("authorization")
        if accessToken != "" {
          token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
              return nil, fmt.Errorf("There was an error")
            }
            return []byte("MyFancySecret"), nil
          })
          if err != nil {
            fmt.Println("There was an error: " + err.Error())
            json.NewEncoder(w).Encode(Error{Error: err.Error()})
            return
          }
          if token.Valid {
            context.Set(req, "decoded", token.Claims)
            next(w, req)
          } else {
            fmt.Println("Invalid authorization token")
            json.NewEncoder(w).Encode(Error{Error: "Invalid authorization token"})
          }
        } else {
            fmt.Println("An authorization header is required")
            json.NewEncoder(w).Encode(Error{Error: "No Token"})
        }
    })
}
