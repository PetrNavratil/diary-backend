package main

import (
  "net/http"
  "github.com/labstack/echo"
  "github.com/parnurzeal/gorequest"
  "fmt"
  "github.com/labstack/echo/middleware"
  gr "github.com/PetrNavratil/diary-back/goodreads"
  "encoding/xml"
)

type BookRequest struct {
  Key string `query:"key"`
}

type BookId struct {
  Id string `query:"id"`
}

func main() {
  e := echo.New()
  e.Use(middleware.CORS())
  e.GET("/", func(c echo.Context) error {
    return c.JSON(http.StatusNotFound, "Don't look around")
  })

  e.GET("/book", func(c echo.Context) error {
    bookId := new(BookId)
    if err := c.Bind(bookId); err != nil {
      return c.String(http.StatusBadRequest, "Parameter ID is not specified")
    }

    _, body, errs := gorequest.New().Get("https://www.goodreads.com/book/show/" + bookId.Id + ".xml?key=tsRkj9chcP8omCKBCJLg0A&").End()
    if errs == nil {
      bookInfo := &gr.GoodReadsBook{}
      xmlResponse := []byte(body)
      xml.Unmarshal(xmlResponse, bookInfo)
      return c.JSON(http.StatusOK, bookInfo)
    } else {
      return c.String(http.StatusNotFound, "FAIL")
    }
  })

  e.GET("/books", func(c echo.Context) error {

    u := new(BookRequest)
    if err := c.Bind(u); err != nil {
      return c.String(http.StatusBadRequest, "FAIL")
    }
    _, body, errs := gorequest.New().Get("https://www.goodreads.com/search/index.xml?key=tsRkj9chcP8omCKBCJLg0A&q=" + u.Key).End()
    if errs == nil {
      foundBooks := &gr.GoodReadsSearchBookResponse{}
      xmlResponse := []byte(body)
      xml.Unmarshal(xmlResponse, foundBooks)

      if (foundBooks.Books == nil) {
        return c.JSON(http.StatusOK, []gr.GoodReadsSearchBook{})
      } else {
        return c.JSON(http.StatusOK, foundBooks.Books)
      }


    } else {
      fmt.Println("error vetev")
    }
    return c.String(http.StatusNotFound, "FAIL")
  })
  e.Logger.Fatal(e.Start(":1323"))
}
