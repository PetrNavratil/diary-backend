package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "github.com/PetrNavratil/diary-back/models"
)

// Function return all readings of user without intervals
func GetAllUsersReadings(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    readings := []models.StatisticReading{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get his readings
      db.Table("readings").Select("readings.id as id,user_id, book_id, completed, start, stop, title, author").
        Joins("JOIN books ON readings.book_id = books.id").Where("user_id = ?", user.ID).Scan(&readings)
      return c.JSON(http.StatusOK, readings)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}