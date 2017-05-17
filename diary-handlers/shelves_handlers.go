package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
  "strconv"
  "fmt"
)

// Function returns all user's shelves
func GetUsersShelves(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelves := []models.Shelf{}
    books := []models.Book{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get his shelves
      db.Model(&user).Related(&shelves)
      // fill shelves with books
      for i := range shelves {
        db.Model(&shelves[i]).Related(&books, "Books")
        shelves[i].Books = books
      }
      return c.JSON(http.StatusOK, shelves)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function creates new shelf
func CreateNewShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelf := &models.Shelf{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get shelf from FE
      if shelfErr := c.Bind(shelf); shelfErr == nil {
        shelf.UserID = user.ID
        // create shelf
        db.Create(shelf)
        shelf.Books = []models.Book{}
        return c.JSON(http.StatusOK, shelf)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "BAD BODY SHELF"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function removes shelf
func RemoveShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelf := models.Shelf{}
    books := []models.Book{}
    // get id of shelf
    if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
      // get shelf
      if !db.Where("id = ?", id).First(&shelf).RecordNotFound() {
        // get shelf books
        db.Model(&shelf).Related(&books, "Books")
        // remove books from shelf
        db.Model(&shelf).Association("Books").Delete(books)
        // delete shelf
        db.Delete(&shelf)
        return c.JSON(http.StatusOK, shelf)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad shelf id"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }
  }
}

// Function edits shelf
func EditShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    editedShelf := &models.Shelf{}
    currentShelf := models.Shelf{}
    // get id of shelf
    if _, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
      // get edited shelf
      if err := c.Bind(editedShelf); err == nil {
        // get database shelf
        db.Where("id = ?", editedShelf.ID).First(&currentShelf)
        // save edited shelf
        editedShelf.UserID = currentShelf.UserID
        db.Save(editedShelf)
        return c.JSON(http.StatusOK, editedShelf)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad shelf body"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad shelf id"})
    }
  }
}

// Function adds book to shelf
func AddBookToShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := &models.Book{}
    books := []models.Book{}
    shelf := models.Shelf{}
    // get shelf id
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      // get book from FE
      if bookErr := c.Bind(book); bookErr == nil {
        // get shelf
        db.Where("id = ?", id).First(&shelf)
        // add book to the shelf
        db.Model(&shelf).Association("Books").Append(book)
        // get shelf books
        db.Model(&shelf).Related(&books, "Books")
        shelf.Books = books
        // return edited shelf
        return c.JSON(http.StatusOK, shelf)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "BAD BODY SHELF"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }
  }
}

// Function removes book from shelf
func RemoveBookFromShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelf := models.Shelf{}
    books := []models.Book{}
    book := models.Book{}
    // get shelf id
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      // get book id
      if bookId, bookIdErr := strconv.Atoi(c.Param("bookId")); bookIdErr == nil {
        // get shelf
        if !db.Where("id = ?", id).First(&shelf).RecordNotFound() {
          // get book
          if !db.Where("id = ?", bookId).First(&book).RecordNotFound() {
            // remove book from shelf
            db.Model(&shelf).Association("Books").Delete(book)
            // get shelf books
            db.Model(&shelf).Related(&books, "Books")
            shelf.Books = books
            // return edited shelf
            return c.JSON(http.StatusOK, shelf)
          } else {
            return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
          }
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad shelf id"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }
  }
}

// Function copies shelf to user's shelves
func CopyShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    newShelf := &models.Shelf{}
    existingShelf := models.Shelf{}
    userBooks := []models.UserBook{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get shelfID
      if shelfId, bookIdErr := strconv.Atoi(c.Param("id")); bookIdErr == nil {
        // get shelf
        if !db.First(&existingShelf, shelfId).RecordNotFound() {
          newShelf.UserID = user.ID
          // create shelf name
          newShelf.Name = fmt.Sprintf("%s - %s", existingShelf.Name, user.UserName)
          // save copied shelf to user
          db.Save(&newShelf)
          // get original shelf books
          db.Model(&existingShelf).Related(&existingShelf.Books, "Books")
          // add original shelf books to the copied shelf
          db.Model(&newShelf).Association("Books").Append(existingShelf.Books)
          // get user's books
          db.Where("user_id = ?", user.ID).Find(&userBooks)
          // function checks whether book is in user's books
          shouldAdd := func(userBooks []models.UserBook, book models.Book) bool {
            for _, userBook := range userBooks {
              if userBook.BookID == book.ID {
                return false
              }
            }
            return true
          }
          // go through user's books and add book if it's in shelf but not in his books
          for _, book := range existingShelf.Books {
            if shouldAdd(userBooks, book) {
              newBook := models.UserBook{
                BookID: book.ID,
                UserID: user.ID,
                InBooks: true,
                Status: models.NOT_READ,
              }
              db.Create(&newBook)
            }
          }
          return c.JSON(http.StatusOK, map[string]string{"message":  "OK"})
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "bad shelf id"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "bad shelf id"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}