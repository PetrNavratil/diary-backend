package diary_handlers

import (
  "github.com/labstack/echo"
  "github.com/jinzhu/gorm"
  "github.com/PetrNavratil/diary-back/models"
  "github.com/dgrijalva/jwt-go"
  "net/http"
  "errors"
  "strconv"
  "os"
  "io"
  "fmt"
  "github.com/kennygrant/sanitize"
)

func GetUser(c echo.Context, db *gorm.DB) (models.User, error) {
  jwtContext := c.Get("user").(*jwt.Token)
  claims := jwtContext.Claims.(jwt.MapClaims)
  id := claims["sub"].(string)
  user := models.User{}
  if (db.Where("auth_id = ?", id).First(&user).RecordNotFound()) {
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
          currentUser.UserName = updatedUser.UserName
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

func UploadAvatar(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if user, err := GetUser(c, db); err == nil {
      // get file from request
      file, err := c.FormFile("file")
      if err != nil {
        fmt.Println("getting file")
        return err
      }
      // open if
      src, err := file.Open()
      if err != nil {
        fmt.Println("opening file")
        return err
      }
      defer src.Close()

      fileName := fmt.Sprintf("images/%d_%s", user.ID, sanitize.Name(file.Filename))
      // Destination
      dst, err := os.Create(fileName)
      if err != nil {
        fmt.Println("creating file")
        return err
      }
      defer dst.Close()

      // copy
      if _, err = io.Copy(dst, src); err != nil {
        fmt.Println("copying file file")
        return err
      }
      if len(user.Avatar) > 0 && user.Avatar != fileName {
        os.Remove(user.Avatar)
      }
      user.Avatar = fileName
      db.Save(&user)
      return c.JSON(http.StatusOK, user)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}
