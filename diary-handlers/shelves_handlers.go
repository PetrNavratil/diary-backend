package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
  "strconv"
)

func GetUsersShelves(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelves := []models.Shelf{}
    books := []models.Book{}
    if user, err := GetUser(c, db); err == nil {
      db.Model(&user).Related(&shelves)
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

func CreateNewShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelf := &models.Shelf{}
    if user, err := GetUser(c, db); err == nil {
      if shelfErr := c.Bind(shelf); shelfErr == nil {
        shelf.UserID = user.ID
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

func RemoveShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelf := models.Shelf{}
    books := []models.Book{}
    if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
      if !db.Where("id = ?", id).First(&shelf).RecordNotFound() {
        db.Model(&shelf).Related(&books, "Books")
        db.Model(&shelf).Association("Books").Delete(books)
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

func EditShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    editedShelf := &models.Shelf{}
    currentShelf := models.Shelf{}
    if _, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
      if err := c.Bind(editedShelf); err == nil {
        db.Where("id = ?", editedShelf.ID).First(&currentShelf)
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

func AddBookToShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := &models.Book{}
    books := []models.Book{}
    shelf := models.Shelf{}
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      if bookErr := c.Bind(book); bookErr == nil {
        db.Where("id = ?", id).First(&shelf)
        db.Model(&shelf).Association("Books").Append(book)
        db.Model(&shelf).Related(&books, "Books")
        shelf.Books = books
        return c.JSON(http.StatusOK, shelf)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "BAD BODY SHELF"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }
  }
}

func RemoveBookFromShelf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    shelf := models.Shelf{}
    books := []models.Book{}
    book := models.Book{}
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      if bookId, bookIdErr := strconv.Atoi(c.Param("bookId")); bookIdErr == nil {
        if !db.Where("id = ?", id).First(&shelf).RecordNotFound() {
          if !db.Where("id = ?", bookId).First(&book).RecordNotFound() {
            db.Model(&shelf).Association("Books").Delete(book)
            db.Model(&shelf).Related(&books, "Books")
            shelf.Books = books
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