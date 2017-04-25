package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
)

func GetPeople(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if user, err := GetUser(c, db); err == nil {
      people := []models.User{}
      key := c.QueryParam("key")
      if len(key) > 0 {
        key := "%" + key + "%"
        db.Where("user_name LIKE ? OR email LIKE  ? OR first_name LIKE  ? OR last_name LIKE  ?", key, key, key, key).
          Not("id = ?", user.ID).
          Find(&people)
      }
      return c.JSON(http.StatusOK, people)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}
