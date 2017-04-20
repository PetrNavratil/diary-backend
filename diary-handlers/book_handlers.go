package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
  "fmt"
  "github.com/PetrNavratil/diary-back/goodreads"
  "github.com/parnurzeal/gorequest"
  "encoding/json"
  "strconv"
  "errors"
  "strings"
  "github.com/kennygrant/sanitize"
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
      _, body, errs := gorequest.New().Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=intitle:%s", newBook.Title)).End()
      if errs == nil {
        var tmp map[string]interface{}
        json.Unmarshal([]byte(body), &tmp)
        if tmp["totalItems"].(float64) > 0 {
          book.GoogleBookId = tmp["items"].([]interface{})[0].(map[string]interface{})["id"].(string)
        }
      }
      db.Create(book)
      db.Last(&book)
      return c.JSON(http.StatusOK, map[string]int{"id": book.ID})
    } else {
      return c.JSON(http.StatusOK, map[string]int{"id": book.ID})
    }
  }
}

func GetBook(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    bookInfo := models.BookInfo{}
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      book := models.Book{}
      db.First(&book, id)
      if grB, grEr := GetGRBook(book.GRBookId); grEr == nil {
        bookInfo.GoodReadsBook = grB
      }
      if gB, gErr := GetGoogleBook(book.GoogleBookId); gErr == nil {
        bookInfo.GoogleBook = gB
      }
      return c.JSON(http.StatusOK, bookInfo)

    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Parameter ID is not specified"})
    }

  }
}

func GetGoogleBook(id string) (models.GoogleBook, error) {
  book := models.GoogleBook{}
  bookResp := models.GoogleBookResponse{}
  _, body, errs := gorequest.New().Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes/%s", id)).End()
  if errs == nil {
    if !strings.Contains(body, `"error"`) {
      json.Unmarshal([]byte(body), &bookResp)
      book.Title = bookResp.VolumeInfo.Title
      book.Author = bookResp.VolumeInfo.Authors[0]
      book.Publisher = bookResp.VolumeInfo.Publisher
      book.Published = bookResp.VolumeInfo.Published
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

func GetLatestBooks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    books := []models.Book{}
    db.Order("created_at desc").Limit(10).Find(&books)
    return c.JSON(http.StatusOK, books)
  }
}