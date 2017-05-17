package models

// Shelf stored in database
type Shelf struct {
  ID      int `json:"id"`
  Name    string `json:"name"`
  Visible bool `json:"visible"`
  UserID  int `json:"-"`
  Books   []Book `json:"books" gorm:"many2many:shelf_books;"`
}