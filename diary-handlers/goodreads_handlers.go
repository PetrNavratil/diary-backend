package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "strconv"
  "github.com/PetrNavratil/diary-back/models"
  "github.com/parnurzeal/gorequest"
  "encoding/xml"
  "net/http"
  "fmt"
  "github.com/PetrNavratil/diary-back/goodreads"
)

type BookRequest struct {
  Key string `query:"key"`
}

func GetGRBook(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      book := models.Book{}
      db.First(&book, id)
      _, body, errs := gorequest.New().Get("https://www.goodreads.com/book/show/" + strconv.Itoa(book.GRBookId) + ".xml?key=tsRkj9chcP8omCKBCJLg0A&").End()
      if errs == nil {
        bookInfo := &goodreads.GoodReadsBook{}
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

func SearchGRBooks() func(c echo.Context) error {
  return func(c echo.Context) error {

    u := new(BookRequest)
    if err := c.Bind(u); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }
    _, body, errs := gorequest.New().Get("https://www.goodreads.com/search/index.xml?key=tsRkj9chcP8omCKBCJLg0A&q=" + u.Key).End()
    if errs == nil {
      foundBooks := &goodreads.GoodReadsSearchBookResponse{}
      xmlResponse := []byte(body)
      xml.Unmarshal(xmlResponse, foundBooks)

      if (foundBooks.Books == nil) {
        return c.JSON(http.StatusOK, []goodreads.GoodReadsSearchBook{})
      } else {
        return c.JSON(http.StatusOK, foundBooks.Books)
      }

    } else {
      fmt.Println("error vetev")
    }
    return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
  }
}
