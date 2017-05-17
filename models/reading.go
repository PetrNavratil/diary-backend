package models

import "time"

// Reading stored in database
type Reading struct {
  ID        int `json:"id"`
  UserID    int `json:"userId"`
  BookID    int `json:"bookId"`
  Completed bool `json:"completed"`
  Start     time.Time `json:"start"`
  Stop      time.Time `json:"stop"`
  Intervals []Interval `json:"intervals"`
}

// Model sent when asked for reading in statistic section
type StatisticReading struct {
  Reading
  Title  string `json:"title"`
  Author string `json:"author"`
}

// Model sent when asked for book reading
type ReturnReading struct {
  LastInterval LastInterval `json:"lastInterval"`
  Readings     []Reading `json:"readings"`
}