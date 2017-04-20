package diary_handlers

import (
  "github.com/labstack/echo"
  "github.com/parnurzeal/gorequest"
  "encoding/xml"
  "net/http"
  "fmt"
  "github.com/PetrNavratil/diary-back/goodreads"
  "regexp"
  "errors"
  "github.com/kennygrant/sanitize"
)

type BookRequest struct {
  Key string `query:"key"`
}

func GetGRBook(id int) (goodreads.GoodReadsBook, error) {
  bookInfo := goodreads.GoodReadsBook{}
  _, body, errs := gorequest.New().Get(fmt.Sprintf("https://www.goodreads.com/book/show/%d.xml?key=tsRkj9chcP8omCKBCJLg0A&", id)).End()
  if errs == nil {
    xmlResponse := []byte(body)
    xml.Unmarshal(xmlResponse, &bookInfo)
    re := regexp.MustCompile("\\(.*\\)")
    bookInfo.Title = re.ReplaceAllLiteralString(bookInfo.Title, "")
    bookInfo.Description = sanitize.HTML(bookInfo.Description)
    for i := range bookInfo.SimilarBooks {
      bookInfo.SimilarBooks[i].Title = re.ReplaceAllLiteralString(bookInfo.SimilarBooks[i].Title, "")
    }
    return bookInfo, nil
  } else {
    return bookInfo, errors.New("ERROR WHILE GETTING BOOK")
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
        re := regexp.MustCompile("\\(.*\\)")
        for i := range foundBooks.Books {
          foundBooks.Books[i].Title = re.ReplaceAllLiteralString(foundBooks.Books[i].Title, "")
        }
        return c.JSON(http.StatusOK, foundBooks.Books)
      }

    } else {
      fmt.Println("error vetev")
    }
    return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
  }
}
