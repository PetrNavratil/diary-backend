package models

import "time"

// Model sent when asked for summary statistics
type Statistic struct {
  BooksCount       int `json:"booksCount"`
  BooksRead        int `json:"booksRead"`
  BooksReading     int `json:"booksReading"`
  BooksToRead      int `json:"booksToRead"`
  BooksNotRead     int `json:"booksNotRead"`
  TimeSpentReading int64 `json:"timeSpentReading"`
  MostlyReadBook   MostlyRead `json:"mostlyRead"`
}

// Model of mostly read book
type MostlyRead struct {
  Book
  Read int `json:"read"`
}

// Model sent when asked for intervals of specific period
type StatisticInterval struct {
  Start     time.Time `json:"start"`
  Stop      time.Time `json:"stop"`
  BookID    int `json:"bookId"`
  Title     string `json:"title"`
  Author    string `json:"author"`
  Completed bool `json:"completed"`
}