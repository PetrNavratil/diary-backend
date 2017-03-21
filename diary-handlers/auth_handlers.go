package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "fmt"
  "strings"
  "github.com/dgrijalva/jwt-go"
  "time"
  "net/http"
)

func Login(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    user := new(models.Login)
    existingUser := models.User{}
    if err := c.Bind(user); err != nil {
      fmt.Println("HERE")
      return err
    }
    if (user.UserName != "" && user.Password != "") {
      if (!db.Where("user_name = ?", user.UserName).First(&existingUser).RecordNotFound()) {
        if (strings.Compare(user.Password, existingUser.Password) == 0) {
          token := jwt.New(jwt.SigningMethodHS256)
          claims := token.Claims.(jwt.MapClaims)
          claims["id"] = existingUser.ID
          claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

          t, err := token.SignedString([]byte("diarySecret"))
          if err != nil {
            return err
          }
          return c.JSON(http.StatusOK, map[string]string{"token": t})
        } else {
          return c.JSON(http.StatusUnauthorized, map[string]string{"message":  "Wrong password"})
        }

      } else {
        return c.JSON(http.StatusUnauthorized, map[string]string{"message":  "User name is unknown"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "All credentials must be set"})
    }
  }
}

func Register(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {

    user := new(models.Register)
    existingUser := models.User{}
    if err := c.Bind(user); err != nil {
      //return err
    }
    if (user.UserName != "" && user.Password != "" && user.Email != "") {
      if (db.Where("user_name = ?", user.UserName).First(&existingUser).RecordNotFound()) {
        existingUser.Password = user.Password
        existingUser.UserName = user.UserName
        existingUser.Email = user.Email
        db.Create(&existingUser)
        db.Where("user_name = ?", user.UserName).First(&existingUser)

        token := jwt.New(jwt.SigningMethodHS256)
        claims := token.Claims.(jwt.MapClaims)
        claims["id"] = existingUser.ID
        claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

        t, err := token.SignedString([]byte("diarySecret"))
        if err != nil {
          return err
        }
        return c.JSON(http.StatusOK, map[string]string{"token": t})
      } else {
        return c.JSON(http.StatusConflict, map[string]string{"message":  "User name already taken"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "All credentials must be set"})
    }
  }
}
