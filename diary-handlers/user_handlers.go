package diary_handlers

import (
  "github.com/labstack/echo"
  "github.com/jinzhu/gorm"
  "github.com/PetrNavratil/diary-back/models"
  "github.com/dgrijalva/jwt-go"
  "net/http"
  "errors"
)

func GetUser(c echo.Context, db *gorm.DB) (models.User, error) {
  jwtContext := c.Get("user").(*jwt.Token)
  claims := jwtContext.Claims.(jwt.MapClaims)
  id := claims["id"]
  user := models.User{}
  if (db.Where("id = ?", id).First(&user).RecordNotFound()) {
    return user, errors.New("NOT FOUND")
  } else {
    return user, nil
  }
}

func GetLoggedUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if user, err := GetUser(c, db); err == nil {
      return c.JSON(http.StatusOK, user)
    } else {
      return c.JSON(http.StatusBadRequest, "USER NOT LOGGED")
    }
  }
}
