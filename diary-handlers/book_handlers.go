package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
  "fmt"
  "github.com/PetrNavratil/diary-back/goodreads"
)

func InsertNewBook(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    newBook := &goodreads.GoodReadsSearchBook{}
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

