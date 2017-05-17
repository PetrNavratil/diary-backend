package models

// User stored in database
type User struct {
  ID        int `json:"id"`
  AuthID    string `json:"-"`
  UserName  string `gorm:"index" gorm:"unique" json:"userName"`
  Email     string `json:"email"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Avatar    string `json:"avatar"`
  UserBook  []UserBook `json:"-" gorm:"ForeignKey:UserID"`
  Shelves   []Shelf `json:"-"`
  Readings  []Reading `json:"-" gorm:"ForeignKey:UserID"`
  Friends   []Friend `json:"-"`
  Requests  []FriendRequest `json:"-"`
}