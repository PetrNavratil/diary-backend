package models

import "time"

// Comment stored in database
type Comment struct {
  ID        int `json:"id"`
  BookID    int `json:"bookId"`
  UserID    int `json:"userId"`
  Text      string `json:"text"`
  CreatedAt time.Time `json:"-"`
  UpdatedAt time.Time`json:"-"`
}

// Model sent as response when asked for book comments
type CommentResponse struct {
  Comment
  Avatar    string `json:"userAvatar"`
  UserName  string `json:"userName"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Date      time.Time `json:"date"`
}

// Model which is got when adding new comment
type BookComment struct {
  Text   string `json:"text"`
  Title  string `json:"title"`
  BookId int `json:"bookId"`
  Date   string `json:"date"`
}