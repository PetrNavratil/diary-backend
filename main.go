package main

import (
  "net/http"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/PetrNavratil/diary-back/models"
  "strings"
  "github.com/PetrNavratil/diary-back/diary-handlers"
)

func getAllBoks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    books := &[]models.Book{}
    db.Find(&books);
    return c.JSON(http.StatusOK, books)
  }
}





func main() {

  db, _ := gorm.Open("sqlite3", "/tmp/gorm.db")
  db.LogMode(true)
  //db.DropTable(&models.User{})
  //db.DropTable(&models.Book{})
  //db.DropTable(&models.UserBook{})
  //db.DropTable(&models.Comment{})
  //db.DropTable(&models.Educational{})
  //db.DropTable(&models.Shelf{})
  //db.DropTable(&models.Tracking{})
  //db.DropTable(&models.Reading{})
  //db.DropTable(&models.Interval{})
  //
  //db.CreateTable(&models.User{})
  //db.CreateTable(&models.Book{})
  //db.CreateTable(&models.UserBook{})
  //db.CreateTable(&models.Comment{})
  //db.CreateTable(&models.Educational{})
  //db.CreateTable(&models.Shelf{})
  //db.CreateTable(&models.Tracking{})
  //db.CreateTable(&models.Reading{})
  //db.CreateTable(&models.Interval{})

  defer db.Close()

  e := echo.New()
  e.Use(middleware.CORS())

  config := middleware.JWTConfig{
    SigningKey: []byte("diarySecret"),
    Skipper: func(c echo.Context) bool {
      if (strings.Compare("/register", c.Path()) == 0 || strings.Compare("/login", c.Path()) == 0 ) {
        return true
      } else {
        return false
      }
    },
  }

  e.Use(middleware.JWTWithConfig(config))

  e.GET("/", func(c echo.Context) error {
    return c.JSON(http.StatusNotFound, map[string]string{"message":  "Don't look around"})
  })
  e.POST("/login", diary_handlers.Login(db))
  e.POST("/register", diary_handlers.Register(db))
  e.GET("/user", diary_handlers.GetLoggedUser(db))

  // gets GR book
  e.GET("/book-detail/:id", diary_handlers.GetGRBook(db))
  e.GET("/search", diary_handlers.SearchGRBooks())


  e.GET("/allboks", getAllBoks(db))
  // adds book to the database
  e.POST("/book", diary_handlers.InsertNewBook(db))

  e.POST("/books/:id", diary_handlers.AddBookToUser(db))
  e.DELETE("/books/:id", diary_handlers.RemoveBookFromUser(db))
  e.PUT("/books/:id", diary_handlers.UpdateUserBookDetail(db))
  e.GET("/books/:id", diary_handlers.GetUserBookDetail(db))
  e.GET("/books", diary_handlers.GetUsersBooks(db))

  e.GET("/comments", diary_handlers.GetBookComments(db))
  e.POST("/comments", diary_handlers.AddBookComment(db))
  e.DELETE("/comments/:id", diary_handlers.DeleteBookComment(db))
  e.PUT("/comments/:id", diary_handlers.UpdateBookComment(db))

  e.GET("/shelves", diary_handlers.GetUsersShelves(db))
  e.POST("/shelves", diary_handlers.CreateNewShelf(db))
  e.DELETE("/shelves/:id", diary_handlers.RemoveShelf(db))
  e.PUT("/shelves/:id", diary_handlers.EditShelf(db))

  e.POST("/shelves/:id", diary_handlers.AddBookToShelf(db))
  e.DELETE("/shelves/:id/:bookId", diary_handlers.RemoveBookFromShelf(db))

  e.PUT("/tracking/start/:id", diary_handlers.StartTracking(db))
  e.PUT("/tracking/stop/:id", diary_handlers.StopTracking(db))
  e.GET("/tracking/book/:id", diary_handlers.GetUserBookTracking(db))
  e.GET("/tracking", diary_handlers.GetLastTracking(db))

  e.Logger.Fatal(e.Start(":1323"))
}
