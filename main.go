package main

import (
  "net/http"
  "github.com/labstack/echo"
  "github.com/parnurzeal/gorequest"
  "fmt"
  "github.com/labstack/echo/middleware"
  gr "github.com/PetrNavratil/diary-back/goodreads"
  "encoding/xml"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/PetrNavratil/diary-back/models"
  "strings"
  "github.com/dgrijalva/jwt-go"
  "time"
)

type BookRequest struct {
  Key string `query:"key"`
}

type BookId struct {
  Id string `query:"id"`
}

func main() {

  db, _ := gorm.Open("sqlite3", "/tmp/gorm.db")
  db.LogMode(true)
  db.DropTable(&models.StoredBook{})
  db.DropTable(&models.User{})
  db.CreateTable(&models.StoredBook{})
  db.CreateTable(&models.User{})
  defer db.Close()

  e := echo.New()
  e.Use(middleware.CORS())

  config := middleware.JWTConfig{
    SigningKey: []byte("diarySecret"),
    Skipper: func(c echo.Context) bool {
      if (strings.Compare("/register", c.Path()) == 0 || strings.Compare("/login", c.Path()) == 0 ) {
        return true
      } else {
        return false
      }
    },
  }

  e.Use(middleware.JWTWithConfig(config))

  e.POST("/login", func(c echo.Context) error {
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
          claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

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
  })
  e.POST("/register", func(c echo.Context) error {

    user := new(models.Register)
    existingUser := models.User{}
    if err := c.Bind(user); err != nil {
      return err
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
        claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

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
  })

  e.GET("/", func(c echo.Context) error {
    return c.JSON(http.StatusNotFound, map[string]string{"message":  "Don't look around"})
  })

  e.GET("/user", func(c echo.Context) error {
    jwtContext := c.Get("user").(*jwt.Token)
    claims := jwtContext.Claims.(jwt.MapClaims)
    id := claims["id"]
    user := models.User{}
    db.Where("id = ?", id).First(&user)
    return c.JSON(http.StatusOK, user)
  })

  e.GET("/book", func(c echo.Context) error {
    bookId := new(BookId)
    if err := c.Bind(bookId); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Parameter ID is not specified"})
    }

    _, body, errs := gorequest.New().Get("https://www.goodreads.com/book/show/" + bookId.Id + ".xml?key=tsRkj9chcP8omCKBCJLg0A&").End()
    if errs == nil {
      bookInfo := &gr.GoodReadsBook{}
      xmlResponse := []byte(body)
      xml.Unmarshal(xmlResponse, bookInfo)
      return c.JSON(http.StatusOK, bookInfo)
    } else {
      return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
    }
  })

  e.GET("/books", func(c echo.Context) error {

    u := new(BookRequest)
    if err := c.Bind(u); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }
    _, body, errs := gorequest.New().Get("https://www.goodreads.com/search/index.xml?key=tsRkj9chcP8omCKBCJLg0A&q=" + u.Key).End()
    if errs == nil {
      foundBooks := &gr.GoodReadsSearchBookResponse{}
      xmlResponse := []byte(body)
      xml.Unmarshal(xmlResponse, foundBooks)

      if (foundBooks.Books == nil) {
        return c.JSON(http.StatusOK, []gr.GoodReadsSearchBook{})
      } else {
        return c.JSON(http.StatusOK, foundBooks.Books)
      }

    } else {
      fmt.Println("error vetev")
    }
    return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
  })
  e.Logger.Fatal(e.Start(":1323"))
}
