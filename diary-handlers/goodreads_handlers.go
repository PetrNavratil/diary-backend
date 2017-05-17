package diary_handlers

import (
  "github.com/labstack/echo"
  "github.com/parnurzeal/gorequest"
  "encoding/xml"
  "net/http"
  "fmt"
  "regexp"
  "errors"
  "github.com/kennygrant/sanitize"
  "github.com/PetrNavratil/diary-back/models"
)

type BookRequest struct {
  Key string `query:"key"`
}

// Function gets goodreads book info
func GetGRBook(id int) (models.GoodReadsBook, error) {
  bookInfo := models.GoodReadsBook{}
  // get info
  _, body, errs := gorequest.New().Get(fmt.Sprintf("https://www.goodreads.com/book/show/%d.xml?key=tsRkj9chcP8omCKBCJLg0A&", id)).End()
  if errs == nil {
    // parse xml to go object
    xmlResponse := []byte(body)
    xml.Unmarshal(xmlResponse, &bookInfo)
    // clear title
    re := regexp.MustCompile("\\(.*\\)")
    bookInfo.Title = re.ReplaceAllLiteralString(bookInfo.Title, "")
    // clear description from HTML tags
    bookInfo.Description = sanitize.HTML(bookInfo.Description)
    for i := range bookInfo.SimilarBooks {
      // clear similar books title
      bookInfo.SimilarBooks[i].Title = re.ReplaceAllLiteralString(bookInfo.SimilarBooks[i].Title, "")
    }
    return bookInfo, nil
  } else {
    return bookInfo, errors.New("ERROR WHILE GETTING BOOK")
  }
}

// Function search for goodreads books
func SearchGRBooks() func(c echo.Context) error {
  return func(c echo.Context) error {

    u := new(BookRequest)
    // get key
    if err := c.Bind(u); err != nil {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
    }
    // send search request
    _, body, errs := gorequest.New().Get("https://www.goodreads.com/search/index.xml?key=tsRkj9chcP8omCKBCJLg0A&q=" + u.Key).End()
    if errs == nil {
      foundBooks := &models.GoodReadsSearchBookResponse{}
      // make go object
      xmlResponse := []byte(body)
      xml.Unmarshal(xmlResponse, foundBooks)

      // send empty array if no books were found
      if (foundBooks.Books == nil) {
        return c.JSON(http.StatusOK, []models.GoodReadsSearchBook{})
      } else {
        re := regexp.MustCompile("\\(.*\\)")
        // clear title of found books
        for i := range foundBooks.Books {
          foundBooks.Books[i].Title = re.ReplaceAllLiteralString(foundBooks.Books[i].Title, "")
        }
        return c.JSON(http.StatusOK, foundBooks.Books)
      }

    } else {
    }
    return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
  }
}
