package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
  "fmt"
  "github.com/parnurzeal/gorequest"
  "encoding/json"
  "strconv"
  "errors"
  "strings"
  "github.com/kennygrant/sanitize"
)

// Inserts new book to the database
func InsertNewBook(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    newBook := &models.GoodReadsSearchBook{}
    book := &models.Book{}
    // gets send book
    if err := c.Bind(newBook); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }

    // checks if book is already stored
    if (db.Where("gr_book_id = ?", newBook.Id, ).First(&book).RecordNotFound()) {
      book.Author = newBook.Author
      book.Title = newBook.Title
      book.GRBookId = newBook.Id
      book.ImageUrl = newBook.ImageUrl
      // find book on Google books
      _, body, errs := gorequest.New().Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=intitle:%s", newBook.Title)).End()
      if errs == nil {
        var tmp map[string]interface{}
        json.Unmarshal([]byte(body), &tmp)
        // if found stores its id
        if tmp["totalItems"].(float64) > 0 {
          book.GoogleBookId = tmp["items"].([]interface{})[0].(map[string]interface{})["id"].(string)
        }
      }
      // save book
      db.Create(book)
      db.Last(&book)
      return c.JSON(http.StatusOK, map[string]int{"id": book.ID})
    } else {
      return c.JSON(http.StatusOK, map[string]int{"id": book.ID})
    }
  }
}

// Function returns information about requested book by provided id
func GetBook(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    bookInfo := models.BookInfo{}
    // gets book id
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      book := models.Book{}
      // finds book
      db.First(&book, id)
      // gets goodreads information
      if grB, grEr := GetGRBook(book.GRBookId); grEr == nil {
        bookInfo.GoodReadsBook = grB
      }
      // gets google books information
      if (len(book.GoogleBookId) > 0) {
        if gB, gErr := GetGoogleBook(book.GoogleBookId); gErr == nil {
          bookInfo.GoogleBook = gB
        }
      }
      return c.JSON(http.StatusOK, bookInfo)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Parameter ID is not specified"})
    }

  }
}

// Function gets google book detail information
func GetGoogleBook(id string) (models.GoogleBook, error) {
  book := models.GoogleBook{}
  bookResp := models.GoogleBookResponse{}
  // get data
  _, body, errs := gorequest.New().Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes/%s", id)).End()
  if errs == nil {
    if !strings.Contains(body, `"error"`) {
      json.Unmarshal([]byte(body), &bookResp)
      // get needed information
      book.Title = bookResp.VolumeInfo.Title
      if (len(bookResp.VolumeInfo.Authors) > 0) {
        book.Author = bookResp.VolumeInfo.Authors[0]
      }
      book.Publisher = bookResp.VolumeInfo.Publisher
      book.Published = bookResp.VolumeInfo.Published
      // clean it form HTML tags
      book.Description = sanitize.HTML(bookResp.VolumeInfo.Description)
      book.PageCount = bookResp.VolumeInfo.PageCount
      book.ImageUrl = bookResp.VolumeInfo.Images.ImageUrl
      book.Preview = bookResp.VolumeInfo.Preview
    }
    return book, nil
  } else {
    return book, errors.New("ERROR WHILE GETTING GOOGLE BOOK")
  }
}

// Function returns recently added books
func GetLatestBooks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    books := []models.Book{}
    db.Order("created_at desc").Limit(10).Find(&books)
    return c.JSON(http.StatusOK, books)
  }
}