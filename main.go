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
  "strconv"
  "github.com/davecgh/go-spew/spew"
  "errors"
)

type BookRequest struct {
  Key string `query:"key"`
}

type BookId struct {
  Id string `query:"id"`
}

func login(db *gorm.DB) func(c echo.Context) error {
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

func register(db *gorm.DB) func(c echo.Context) error {
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
        claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

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

func getUser(c echo.Context, db *gorm.DB) (models.User, error) {
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

func getLoggedUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if user, err := getUser(c, db); err == nil {
      return c.JSON(http.StatusOK, user)
    } else {
      return c.JSON(http.StatusBadRequest, "USER NOT LOGGED")
    }
  }
}

func getGRBook(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      book := models.Book{}
      db.First(&book, id)
      _, body, errs := gorequest.New().Get("https://www.goodreads.com/book/show/" + strconv.Itoa(book.GRBookId) + ".xml?key=tsRkj9chcP8omCKBCJLg0A&").End()
      if errs == nil {
        bookInfo := &gr.GoodReadsBook{}
        xmlResponse := []byte(body)
        xml.Unmarshal(xmlResponse, bookInfo)
        return c.JSON(http.StatusOK, bookInfo)
      } else {
        return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
      }

    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Parameter ID is not specified"})
    }

  }
}

func searchGRBooks() func(c echo.Context) error {
  return func(c echo.Context) error {

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
  }
}

func getAllBoks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    books := &[]models.Book{}
    db.Find(&books);
    return c.JSON(http.StatusOK, books)
  }
}

func insertNewBook(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    newBook := &gr.GoodReadsSearchBook{}
    book := &models.Book{}
    if err := c.Bind(newBook); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }

    if (db.Where("gr_book_id = ?", newBook.Id, ).First(&book).RecordNotFound()) {
      fmt.Println("NOT IN DATABASE")
      book.Author = newBook.Author
      book.Title = newBook.Title
      book.GRBookId = newBook.Id
      book.ImageUrl = newBook.ImageUrl
      db.Create(book)
      db.Last(&book)
      return c.JSON(http.StatusOK, map[string]int{"id": book.ID})
    } else {
      return c.JSON(http.StatusOK, map[string]int{"id": book.ID})
    }
  }
}

func addBookToUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    userBook := &models.UserBook{}
    returnBook := models.ReturnBook{}
    if loggedUser, logErr := getUser(c, db); logErr == nil {
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        userBook.BookID = id
        userBook.UserID = loggedUser.ID
        userBook.Status = false
        db.Create(&userBook)
        db.Table("books").Select(
          "books.id, books.title, books.author, books.image_url, user_book.status").
          Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ? AND book_id = ?", loggedUser.ID, id).Scan(&returnBook)
        spew.Dump(returnBook)
        return c.JSON(http.StatusOK, returnBook)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
      }
    } else {
      fmt.Println(loggedUser)
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  logErr.Error()})
    }
  }
}

func getUsersBooks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    returnBooks := []models.ReturnBook{}
    if user, err := getUser(c, db); err == nil {
      db.Table("books").Select(
        "books.id, books.title, books.author, books.image_url, user_book.status").
        Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ?", user.ID).Scan(&returnBooks)
      return c.JSON(http.StatusOK, returnBooks)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func removeBookFromUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      if user, err := getUser(c, db); err == nil {
        db.Where("user_id = ? AND  book_id = ?", user.ID, id).Delete(models.UserBook{})
        return c.JSON(http.StatusOK, map[string]int{"id": id})
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
      }
    } else {
      return c.JSON(http.StatusBadRequest, "BAD ID")
    }
  }
}

func changeStatus(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    body := &models.ReturnBook{}
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      if user, err := getUser(c, db); err == nil {
        if err := c.Bind(body); err != nil {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "invalid book"})
        }
        return c.JSON(http.StatusOK, map[string]int{"id": id})
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "invalid id"})
    }
  }
}

func main() {

  db, _ := gorm.Open("sqlite3", "/tmp/gorm.db")
  db.LogMode(true)
  //db.DropTable(&models.User{})
  //db.DropTable(&models.Book{})
  //db.DropTable(&models.UserBook{})
  //
  //
  //db.CreateTable(&models.User{})
  //db.CreateTable(&models.Book{})
  //db.CreateTable(&models.UserBook{})
  //db.CreateTable(&models.Comment{})


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

  e.GET("/", func(c echo.Context) error {
    return c.JSON(http.StatusNotFound, map[string]string{"message":  "Don't look around"})
  })
  e.POST("/login", login(db))
  e.POST("/register", register(db))
  e.GET("/user", getLoggedUser(db))
  e.GET("/search", searchGRBooks())
  e.GET("/allboks", getAllBoks(db))
  e.GET("/book/:id", getGRBook(db))
  e.POST("/book", insertNewBook(db))

  e.POST("/books/:id", addBookToUser(db))
  e.DELETE("/books/:id", removeBookFromUser(db))
  e.PUT("/books/:id", changeStatus(db))
  e.GET("/books", getUsersBooks(db))

  e.Logger.Fatal(e.Start(":1323"))
}
