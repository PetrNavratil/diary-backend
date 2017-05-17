package models

import "time"

// Reading interval stored in database
type Interval struct {
  ID        int `json:"id"`
  Start     time.Time `json:"start"`
  Stop      time.Time `json:"stop"`
  ReadingID int `json:"readingId"`
}

// Model sent when asked for last reading interval
type LastInterval struct {
  Interval
  Title     string `json:"title"`
  Author    string `json:"author"`
  Completed bool `json:"completed"`
  BookID    int `json:"bookId"`
}