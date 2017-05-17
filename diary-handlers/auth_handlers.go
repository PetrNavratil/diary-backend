package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
)

// Adds new user to the database
// Called from auth0
func NewUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    var authUser  struct {
      Email     string `json:"email"`
      Username  string `json:"username"`
      ID        string `json:"id"`
      FirstName string`json:"firstName"`
      LastName  string `json:"lastName"`
      ImageUrl  string `json:"imageUrl"`
    }
    c.Bind(&authUser)
    user := models.User{}
    user.Email = authUser.Email
    user.UserName = authUser.Username
    user.AuthID = authUser.ID
    user.FirstName = authUser.FirstName
    user.LastName = authUser.LastName
    user.Avatar = authUser.ImageUrl
    db.Create(&user)
    return c.HTML(http.StatusOK, "OK")
  }
}

