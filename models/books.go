package models

import (
  "time"
)
// Book stored in database
type Book struct {
  ID           int `json:"id"`
  Title        string `gorm:"index" json:"title"`
  Author       string `gorm:"index" json:"author"`
  ISBN         string `json:"isbn"`
  ISBN13       string `json:"isbn13"`
  ImageUrl     string `json:"imageUrl"`
  GRBookId     int `json:"grBookId"`
  GoogleBookId string `json:"googleBookId"`
  UserBook     []UserBook `json:"-" gorm:"ForeignKey:BookID"`
  Comments     []Comment `json:"-" gorm:"ForeignKey:BookID"`
  Readings     []Reading `json:"-" gorm:"ForeignKey:BookID"`
  CreatedAt    time.Time
}

// M:N model stored in database for User and Book
type UserBook struct {
  ID          int
  UserID      int
  BookID      int
  Status      int
  InBooks     bool
  Educational Educational
  CreatedAt   time.Time
}

// Changes name of created table
func (*UserBook) TableName() string {
  return "user_book"
}

// Model send as response when asked for user books
type ReturnBook struct {
  ID          int `json:"id"`
  Title       string `json:"title"`
  Author      string `json:"author"`
  ImageUrl    string `json:"imageUrl"`
  InBooks     bool `json:"inBooks"`
  Status      int `json:"status"`
  Educational Educational `json:"educational"`
  CreatedAt   time.Time `json:"createdAt"`
}

// Model send as response when asked for book detail GR, GB
type BookInfo struct {
  GoodReadsBook GoodReadsBook `json:"goodReadsBook"`
  GoogleBook    GoogleBook `json:"googleBook"`
}