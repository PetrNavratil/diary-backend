package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "github.com/PetrNavratil/diary-back/models"
)

func ChangePassword(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    password := models.PasswordChange{}
    if user, err := GetUser(c, db); err == nil {
      if bodyError := c.Bind(&password); bodyError == nil {
        if user.Password == password.OldPassword {
          user.Password = password.NewPassword
          db.Save(&user)
          return c.JSON(http.StatusOK, map[string]string{"message":  "OK"})
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad password"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "BAD body"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}