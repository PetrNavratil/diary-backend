package models

import (
  "time"
  "github.com/PetrNavratil/diary-back/goodreads"
)

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
  Trackings    []Tracking `json:"-" gorm:"ForeignKey:BookID"`
  Readings     []Reading `json:"-" gorm:"ForeignKey:BookID"`
  CreatedAt    time.Time
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
  AuthID    string `json:"-"`
  UserName  string `gorm:"index" gorm:"unique" json:"userName"`
  Email     string `json:"email"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Avatar    string `json:"avatar"`
  UserBook  []UserBook `json:"-" gorm:"ForeignKey:UserID"`
  Shelves   []Shelf `json:"-"`
  Trackings []Tracking `json:"-" gorm:"ForeignKey:UserID"`
  Readings  []Reading `json:"-" gorm:"ForeignKey:UserID"`
  Friends   []Friend `json:"-"`
  Requests  []FriendRequest `json:"-"`
}

type Friend struct {
  ID        int  `json:"-"`
  UserID    int `json:"-"`
  FriendID  int `json:"id"`
  CreatedAt time.Time `json:"since"`
}

type FriendRequest struct {
  ID          int `json:"id"`
  UserID      int `json:"-"`
  RequesterID int `json:"-"`
}

type FriendRequestResponse struct {
  FriendRequest
  UserName  string `json:"userName"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Avatar    string `json:"avatar"`
}

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

type UserBook struct {
  ID          int
  UserID      int
  BookID      int
  Status      int
  InBooks     bool
  Educational Educational
  CreatedAt   time.Time
}

func (*UserBook) TableName() string {
  return "user_book"
}

type Educational struct {
  ID         int `json:"id"`
  UserBookID int `json:"-"`
  Druh       string `json:"druh"`
  Zanr       string `json:"zanr"`
  Smer       string `json:"smer"`
  Forma      string `json:"forma"`
  Jazyk      string `json:"jazyk"`
  Postavy    string `json:"postavy"`
  Obsah      string `json:"obsah"`
  Tema       string `json:"tema"`
  Hodnoceni  string `json:"hodnoceni"`
}

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

const (
  NOT_READ int = iota
  TO_READ
  READING
  READ
  ALL
)

type BookComment struct {
  Text   string `json:"text"`
  Title  string `json:"title"`
  BookId int `json:"bookId"`
  Date   string `json:"date"`
}

type Shelf struct {
  ID      int `json:"id"`
  Name    string `json:"name"`
  Visible bool `json:"visible"`
  UserID  int `json:"-"`
  Books   []Book `json:"books" gorm:"many2many:shelf_books;"`
}

type Tracking struct {
  ID     int `json:"id"`
  UserID int `json:"userId"`
  BookID int `json:"bookId"`
  Start  time.Time `json:"start"`
  End    time.Time `json:"end"`
}

type LastTracking struct {
  Tracking
  Title  string `json:"title"`
  Author string `json:"author"`
}

type ReturnTracking struct {
  LastTracking LastTracking `json:"lastTracking"`
  Trackings    []Tracking `json:"trackings"`
}

type Reading struct {
  ID        int `json:"id"`
  UserID    int `json:"userId"`
  BookID    int `json:"bookId"`
  Completed bool `json:"completed"`
  Start     time.Time `json:"start"`
  Stop      time.Time `json:"stop"`
  Intervals []Interval `json:"intervals"`
}

type StatisticReading struct {
  Reading
  Title  string `json:"title"`
  Author string `json:"author"`
}

type Interval struct {
  ID        int `json:"id"`
  Start     time.Time `json:"start"`
  Stop      time.Time `json:"stop"`
  ReadingID int `json:"readingId"`
}

type LastInterval struct {
  Interval
  Title     string `json:"title"`
  Author    string `json:"author"`
  Completed bool `json:"completed"`
  BookID    int `json:"bookId"`
}

type ReturnReading struct {
  LastInterval LastInterval `json:"lastInterval"`
  Readings     []Reading `json:"readings"`
}

type Statistic struct {
  BooksCount       int `json:"booksCount"`
  BooksRead        int `json:"booksRead"`
  BooksReading     int `json:"booksReading"`
  BooksToRead      int `json:"booksToRead"`
  BooksNotRead     int `json:"booksNotRead"`
  TimeSpentReading int64 `json:"timeSpentReading"`
  MostlyReadBook   MostlyRead `json:"mostlyRead"`
}

type MostlyRead struct {
  Book
  Read int `json:"read"`
}

type StatisticInterval struct {
  Start     time.Time `json:"start"`
  Stop      time.Time `json:"stop"`
  BookID    int `json:"bookId"`
  Title     string `json:"title"`
  Author    string `json:"author"`
  Completed bool `json:"completed"`
}

type BookInfo struct {
  GoodReadsBook goodreads.GoodReadsBook `json:"goodReadsBook"`
  GoogleBook    GoogleBook `json:"googleBook"`

}