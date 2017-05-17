package diary_handlers

import (
  "github.com/PetrNavratil/diary-back/models"
  "strconv"
  "net/http"
  "fmt"
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "time"
)

// Function adds book to the user's books
func AddBookToUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    userBook := &models.UserBook{}
    returnBook := models.ReturnBook{}
    // gets user
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      // gets book id
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        userBook.BookID = id
        userBook.UserID = loggedUser.ID
        userBook.InBooks = true
        userBook.Status = models.NOT_READ
        // creates record in user_book
        db.Create(&userBook)
        // selects all neded information for FE
        db.Table("books").Select(
          "books.id, books.title, books.author, books.image_url, user_book.in_books, user_book.status, user_book.created_at").
          Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ? AND book_id = ?", loggedUser.ID, id).Scan(&returnBook)
        return c.JSON(http.StatusOK, returnBook)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  logErr.Error()})
    }
  }
}

// Function returns all user's books
func GetUsersBooks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    returnBooks := []models.ReturnBook{}
    // gets user
    if user, err := GetUser(c, db); err == nil {
      // selects all needed information for FE
      db.Table("books").Select(
        "books.id, books.title, books.author, books.image_url, user_book.in_books, user_book.status").
        Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ?", user.ID).Scan(&returnBooks)
      return c.JSON(http.StatusOK, returnBooks)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function removes book from user
func RemoveBookFromUser(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    returnBook := models.ReturnBook{}
    book := models.Book{}
    // get id of book
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      // get user
      if user, err := GetUser(c, db); err == nil {
        // get book
        db.First(&book, id)
        shelves := []models.Shelf{}
        // get user's shelves
        db.Model(&user).Related(&shelves, "Shelves")
        for _, shelf := range shelves {
          // remove book from shelves
          db.Model(&shelf).Association("Books").Delete(book)
        }
        // delete readings from user but let them to book to count how many times the book has been read so far
        readings := []models.Reading{}
        db.Where("user_id = ? AND book_id = ?", user.ID, book.ID).Find(&readings)
        db.Model(&user).Association("Readings").Delete(readings)

        userBook := models.UserBook{}
        // get user_book
        db.Where("user_id = ? AND book_id = ?", user.ID, id).First(&userBook)
        // remove literary analysis
        db.Where("user_book_id = ?", userBook.ID).Delete(models.Educational{})
        // remove book from user
        db.Delete(userBook)
        // find removed book and send it as not read for FE
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

// Function returns information user's information about book
func GetUserBookDetail(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    returnBook := models.ReturnBook{}
    userBook := models.UserBook{}
    // get user
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      // get book id
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        // select all information about book
        // if book is not in user's books set it as not read
        if ( db.Table("books").Select(
          "books.id, books.title, books.author, books.image_url, user_book.status, user_book.in_books, user_book.created_at").
          Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ? AND book_id = ?", loggedUser.ID, id).Scan(&returnBook).RecordNotFound()) {
          db.Table("books").Select("id, title, author, image_url").Where("id = ? ", id).Scan(&returnBook)
          returnBook.Status = models.NOT_READ
          return c.JSON(http.StatusOK, returnBook)
        } else {
          // book is in user's books
          // get its educational
          db.Where("user_id = ? AND book_id = ?", loggedUser.ID, id).First(&userBook)
          db.Model(&userBook).Related(&returnBook.Educational)
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

// Function updates user book information
func UpdateUserBookDetail(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    currentBook := models.UserBook{}
    updated := &models.ReturnBook{}
    reading := models.Reading{}
    interval := models.Interval{}
    // get user
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      // get book id
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        // get updated book information
        if bodyError := c.Bind(updated); bodyError == nil {
          // get user book from database
          db.Where("user_id = ? AND book_id = ?", loggedUser.ID, id).First(&currentBook)
          // should book be set as reading?
          if updated.Status == models.READING && currentBook.Status != models.READING {
            reading.UserID = loggedUser.ID
            reading.BookID = id
            reading.Completed = false
            reading.Start = time.Now()
            // create new reading for  book
            db.Create(&reading)
          }
          // should be state changed to read?
          if updated.Status == models.READ && currentBook.Status != models.READ {
            // is book being read now?
            if currentBook.Status == models.READING {
              // get last reading of book
              db.Where("user_id = ? AND book_id = ? AND completed = ?", loggedUser.ID, id, false).Last(&reading)
              if !db.Where("reading_id = ?", reading.ID).Last(&interval).RecordNotFound() {
                // is the book being tracked right now?
                if interval.Stop.IsZero() {
                  // stop tracking
                  interval.Stop = time.Now()
                  db.Save(&interval)
                  reading.Stop = interval.Stop
                } else {
                  // stop just reading
                  reading.Stop = time.Now()
                }
              } else {
                reading.Stop = time.Now()
              }
              // set as completed and save
              reading.Completed = true
              db.Save(&reading)
            } else {
              // save as instant read
              reading.UserID = loggedUser.ID
              reading.BookID = id
              reading.Completed = true
              reading.Start = time.Now()
              reading.Stop = reading.Start
              db.Create(&reading)
            }
          }
          // change book status
          currentBook.Status = updated.Status
          // update educational
          currentBook.Educational = updated.Educational
          db.Save(&currentBook)
          // send updated changes
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
