package diary_handlers

import (
  "github.com/PetrNavratil/diary-back/models"
  "strconv"
  "net/http"
  "fmt"
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/davecgh/go-spew/spew"
  "time"
)

func AddBookToUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    userBook := &models.UserBook{}
    returnBook := models.ReturnBook{}
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        userBook.BookID = id
        userBook.UserID = loggedUser.ID
        userBook.InBooks = true
        userBook.Status = models.NOT_READ
        db.Create(&userBook)
        db.Table("books").Select(
          "books.id, books.title, books.author, books.image_url, user_book.in_books, user_book.status, user_book.created_at").
          Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ? AND book_id = ?", loggedUser.ID, id).Scan(&returnBook)
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

func GetUsersBooks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    returnBooks := []models.ReturnBook{}
    if user, err := GetUser(c, db); err == nil {
      db.Table("books").Select(
        "books.id, books.title, books.author, books.image_url, user_book.in_books, user_book.status, user_book.in_books, user_book.created_at").
        Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ?", user.ID).Scan(&returnBooks)
      return c.JSON(http.StatusOK, returnBooks)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func RemoveBookFromUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    returnBook := models.ReturnBook{}
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      if user, err := GetUser(c, db); err == nil {
        db.Where("user_id = ? AND  book_id = ?", user.ID, id).Delete(models.UserBook{})
        db.Table("books").Select("id, title, author, image_url").Where("id = ? ", id).Scan(&returnBook)
        returnBook.Status = models.NOT_READ
        return c.JSON(http.StatusOK, returnBook)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
      }
    } else {
      return c.JSON(http.StatusBadRequest, "BAD ID")
    }
  }
}

func GetUserBookDetail(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    returnBook := models.ReturnBook{}
    userBook := models.UserBook{}
    userBooks := []models.UserBook{}
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        if ( db.Table("books").Select(
          "books.id, books.title, books.author, books.image_url, user_book.status, user_book.in_books, user_book.created_at").
          Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ? AND book_id = ?", loggedUser.ID, id).Scan(&returnBook).RecordNotFound()) {
          db.Table("books").Select("id, title, author, image_url").Where("id = ? ", id).Scan(&returnBook)
          returnBook.Status = models.NOT_READ
          return c.JSON(http.StatusOK, returnBook)
        } else {
          db.Where("user_id = ? AND book_id = ?", loggedUser.ID, id).First(&userBook)
          db.Model(&userBook).Related(&returnBook.Educational)
          month := fmt.Sprintf("%02d", int(time.Now().Month()))
          db.Where("user_id = ? AND strftime('%m', created_at) = ?", loggedUser.ID, month).Find(&userBooks);
          spew.Dump(userBooks)
          return c.JSON(http.StatusOK, returnBook)
        }
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

func UpdateUserBookDetail(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    currentBook := models.UserBook{}
    updated := &models.ReturnBook{}
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        if bodyError := c.Bind(updated); bodyError == nil {
          db.Where("user_id = ? AND book_id = ?", loggedUser.ID, id).First(&currentBook)
          currentBook.Status = updated.Status
          currentBook.Educational = updated.Educational
          db.Save(&currentBook)
          db.Table("books").Select(
            "books.id, books.title, books.author, books.image_url, user_book.status, user_book.in_books,user_book.created_at").
            Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ? AND book_id = ?", loggedUser.ID, id).Scan(updated)
          db.Model(&currentBook).Related(&updated.Educational)
          return c.JSON(http.StatusOK, updated)

        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
      }
    } else {
      fmt.Println(loggedUser)
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  logErr.Error()})
    }
  }
}
