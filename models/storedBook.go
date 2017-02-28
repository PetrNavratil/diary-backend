package models

type StoredBook struct {
  Name   string `gorm:"index"`
  Author string `gorm:"index"`
  ISBN   string
  Cover  string
  GrId   int
}