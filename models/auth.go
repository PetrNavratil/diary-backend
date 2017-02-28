package models

type Register struct {
  UserName string `json:"userName"`
  Email    string `json:"email"`
  Password string `json:"password"`
}

type Login struct {
  UserName string `json:"userName"`
  Password string `json:"password"`
}
