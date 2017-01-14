package main

import (
  "net/http"
  "github.com/labstack/echo"
  "github.com/parnurzeal/gorequest"
  "fmt"
  xj "github.com/basgys/goxml2json"
  "strings"
  "github.com/labstack/echo/middleware"
  gr "github.com/PetrNavratil/diary-back/goodreads"
  "encoding/xml"
  "github.com/davecgh/go-spew/spew"
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

    _, body, errs := gorequest.New().Get("https://www.goodreads.com/book/show/" + bookId.Id + ".xml?key=tsRkj9chcP8omCKBCJLg0A&q=").End()
    if errs == nil {
      bookInfo := &gr.GoodReadsBook{}
      xmlResponse := []byte(body)
      xml.Unmarshal(xmlResponse, bookInfo)
      spew.Dump(bookInfo)
      return c.JSON(http.StatusOK, bookInfo)
    } else {
      return c.String(http.StatusNotFound, "FAIL")
    }
  })

  e.GET("/books", func(c echo.Context) error {

    u := new(BookRequest)
    fmt.Println(u)
    if errrrror := c.Bind(u); errrrror != nil {
      return c.String(http.StatusBadRequest, "FAIL")
    }
    _, body, errs := gorequest.New().Get("https://www.goodreads.com/search/index.xml?key=tsRkj9chcP8omCKBCJLg0A&q="+u.Key).End()
    if errs == nil {
      xml := strings.NewReader(body)
      json, err := xj.Convert(xml)
      if err != nil {
        panic("That's embarrassing...")
      }

      final := json.String()

      return c.String(http.StatusOK, final)
    } else {
      fmt.Println("error vetev")
    }
    return c.String(http.StatusNotFound, "FAIL")
  })
  e.Logger.Fatal(e.Start(":1323"))
}
