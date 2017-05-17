package models

import "time"

// Friend stored in database
type Friend struct {
  ID        int  `json:"-"`
  UserID    int `json:"-"`
  FriendID  int `json:"id"`
  CreatedAt time.Time `json:"since"`
}

// Friend request stored in database
type FriendRequest struct {
  ID          int `json:"id"`
  UserID      int `json:"userId"`
  RequesterID int `json:"requesterId"`
}

// Response sent when asked for user's friend requests
type FriendRequestResponse struct {
  FriendRequest
  UserName  string `json:"userName"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Avatar    string `json:"avatar"`
}

// Response sent when asked for user's friends
type FriendResponse struct {
  Friend
  UserName     string `json:"userName"`
  FirstName    string `json:"firstName"`
  LastName     string `json:"lastName"`
  Avatar       string `json:"avatar"`
  BooksCount   int `json:"booksCount"`
  ShelvesCount int `json:"shelvesCount"`
  Books        []ReturnBook `json:"books"`
  Shelves      []Shelf `json:"shelves"`
}