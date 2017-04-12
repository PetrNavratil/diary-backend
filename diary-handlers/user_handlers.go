package diary_handlers

import (
  "github.com/labstack/echo"
  "github.com/jinzhu/gorm"
  "github.com/PetrNavratil/diary-back/models"
  "github.com/dgrijalva/jwt-go"
  "net/http"
  "errors"
  "strconv"
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

func EditUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    currentUser := models.User{}
    updatedUser := models.User{}
    if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
      if bodyError := c.Bind(&updatedUser); bodyError == nil {
        if !db.Where("id = ?", id).First(&currentUser).RecordNotFound() {
          currentUser.Email = updatedUser.Email
          currentUser.LastName = updatedUser.LastName
          currentUser.FirstName = updatedUser.FirstName
          db.Save(&currentUser)
          return c.JSON(http.StatusOK, updatedUser);
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad user id"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad user body"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad user id"})
    }
  }
}
