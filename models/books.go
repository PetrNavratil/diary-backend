package models

type Book struct {
  ID           int `json:"id"`
  Title        string `gorm:"index" json:"title"`
  Author       string `gorm:"index" json:"author"`
  ISBN         string `json:"isbn"`
  ISBN13       string `json:"isbn13"`
  ImageUrl     string `json:"imageUrl"`
  GRBookId     int `json:"grBookId"`
  GoogleBookId int `json:"googleBookId"`
  UserBook     []UserBook `json:"-" gorm:"ForeignKey:BookID"`
  Comments     []Comment `json:"-" gorm:"ForeignKey:BookID"`
}

type Comment struct {
  ID     int `json:"id"`
  BookID int `json:"bookId"`
  UserID int `json:"userId"`
  Text   string `json:"text"`
  Date   string `json:"date"`
}

type CommentResponse struct {
  ID       int `json:"id"`
  BookID   int `json:"bookId"`
  Text     string `json:"text"`
  Date     string `json:"date"`
  Avatar   string `json:"userAvatar"`
  UserName string `json:"userName"`
  UserId   int `json:"userId"`
}

type User struct {
  ID        int `json:"id"`
  UserName  string `gorm:"index" gorm:"unique" json:"userName"`
  Password  string `json:"-"`
  Email     string `json:"email"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Avatar    string `json:"avatar"`
  UserBook  []UserBook `json:"-" gorm:"ForeignKey:UserID"`
}

type Shelf struct {
  Name    string `json:"name"`
  Visible bool `json:"visible"`
}

type UserBook struct {
  ID      int
  UserID  int
  BookID  int
  Status  int
  InBooks bool
}

func (*UserBook) TableName() string {
  return "user_book"
}

type ReturnBook struct {
  ID       int `json:"id"`
  Title    string `json:"title"`
  Author   string `json:"author"`
  ImageUrl string `json:"imageUrl"`
  InBooks  bool `json:"inBooks"`
  Status   int `json:"status"`

}

const (
  NOT_READ int = iota
  TO_READ
  READING
)

type BookComment struct {
  Text   string `json:"text"`
  Title  string `json:"title"`
  BookId int `json:"bookId"`
  Date   string `json:"date"`
}
