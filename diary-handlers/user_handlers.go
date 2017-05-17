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

// Function gets logged user by JWT token
func GetUser(c echo.Context, db *gorm.DB) (models.User, error) {
  // get user from JWT
  jwtContext := c.Get("user").(*jwt.Token)
  claims := jwtContext.Claims.(jwt.MapClaims)
  // get auth0 id from claims
  id := claims["sub"].(string)
  user := models.User{}
  // find user in database
  if (db.Where("auth_id = ?", id).First(&user).RecordNotFound()) {
    return user, errors.New("NOT FOUND")
  } else {
    return user, nil
  }
}

// Function returns logged user
func GetLoggedUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if user, err := GetUser(c, db); err == nil {
      return c.JSON(http.StatusOK, user)
    } else {
      return c.JSON(http.StatusBadRequest, "USER NOT LOGGED")
    }
  }
}

// Function edits user
func EditUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    currentUser := models.User{}
    updatedUser := models.User{}
    // get user id
    if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
      // get updated user
      if bodyError := c.Bind(&updatedUser); bodyError == nil {
        // find user and update him
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

// Function handles avatar uploading
func UploadAvatar(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get file from request
      file, err := c.FormFile("file")
      if err != nil {
        fmt.Println("getting file")
        return err
      }
      // open source file
      src, err := file.Open()
      if err != nil {
        return err
      }
      defer src.Close()

      // create filename
      fileName := fmt.Sprintf("images/%d_%s", user.ID, sanitize.Name(file.Filename))
      // Destination
      dst, err := os.Create(fileName)
      if err != nil {
        return err
      }
      defer dst.Close()

      // copy
      if _, err = io.Copy(dst, src); err != nil {
        fmt.Println("copying file file")
        return err
      }
      // if database user avatar is different from current avatar remove it
      if len(user.Avatar) > 0 && user.Avatar != fileName {
        os.Remove(user.Avatar)
      }
      // save new avatar
      user.Avatar = fileName
      db.Save(&user)
      return c.JSON(http.StatusOK, user)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}
