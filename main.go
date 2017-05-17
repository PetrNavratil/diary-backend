package main

import (
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "strings"
  "github.com/PetrNavratil/diary-back/diary-handlers"
)

func main() {

  // opens database
  db, _ := gorm.Open("sqlite3", "gorm.db")
  //db.LogMode(true)
  // creates tables
  //db.DropTable(&models.User{})
  //db.DropTable(&models.Book{})
  //db.DropTable(&models.UserBook{})
  //db.DropTable(&models.Comment{})
  //db.DropTable(&models.Educational{})
  //db.DropTable(&models.Shelf{})
  //db.DropTable(&models.Tracking{})
  //db.DropTable(&models.Reading{})
  //db.DropTable(&models.Interval{})
  //db.DropTable(&models.Friend{})
  //db.DropTable(&models.FriendRequest{})
  //////
  //db.CreateTable(&models.User{})
  //db.CreateTable(&models.Book{})
  //db.CreateTable(&models.UserBook{})
  //db.CreateTable(&models.Comment{})
  //db.CreateTable(&models.Educational{})
  //db.CreateTable(&models.Shelf{})
  //db.CreateTable(&models.Tracking{})
  //db.CreateTable(&models.Reading{})
  //db.CreateTable(&models.Interval{})
  //db.CreateTable(&models.Friend{})
  //db.CreateTable(&models.FriendRequest{})
  // closes db when main function ends
  defer db.Close()

  // starts new api
  e := echo.New()
  // sets route for static files
  e.Static("/images", "images")
  // enables cors
  e.Use(middleware.CORS())

  // create JWT config
  config := middleware.JWTConfig{
    SigningKey: []byte("x6gcPcqZkeG9wnjWh4I1_GKfKNMnAGuXS2m6oPUoeqM4nOATs2TKbsJvoS5cYPIU"),
    Skipper: func(c echo.Context) bool {
      // compare function
      compare := func(path string) bool {
        return strings.Compare(path, c.Path()) == 0
      }
      // if the route is new or images JWT is NOT required
      if (compare("/new") || strings.Contains(c.Path(), "/images") ) {
        return true
      } else {
        return false
      }
    },
  }

  // use JWT
  e.Use(middleware.JWTWithConfig(config))

  // gets logged user
  e.GET("/user", diary_handlers.GetLoggedUser(db))
  // edits user
  e.PUT("/user/:id", diary_handlers.EditUser(db))
  // uploads avatar
  e.POST("/user/avatar", diary_handlers.UploadAvatar(db))

  // gets GR and GB book detail
  e.GET("/book-detail/:id", diary_handlers.GetBook(db))
  // search GR book
  e.GET("/search", diary_handlers.SearchGRBooks())

  // adds book to the database
  e.POST("/book", diary_handlers.InsertNewBook(db))

  // adds book to the user
  e.POST("/books/:id", diary_handlers.AddBookToUser(db))
  // removes book from the user
  e.DELETE("/books/:id", diary_handlers.RemoveBookFromUser(db))
  // updates user book
  e.PUT("/books/:id", diary_handlers.UpdateUserBookDetail(db))
  // gets user book info
  e.GET("/books/:id", diary_handlers.GetUserBookDetail(db))
  // gets all user books
  e.GET("/books", diary_handlers.GetUsersBooks(db))
  // gets recently added books
  e.GET("/books/latest", diary_handlers.GetLatestBooks(db))

  // gets book comments
  e.GET("/comments", diary_handlers.GetBookComments(db))
  // adds book comment
  e.POST("/comments", diary_handlers.AddBookComment(db))
  // removes book comment
  e.DELETE("/comments/:id", diary_handlers.DeleteBookComment(db))
  // updates book comment
  e.PUT("/comments/:id", diary_handlers.UpdateBookComment(db))

  // gets user's shelves
  e.GET("/shelves", diary_handlers.GetUsersShelves(db))
  // creates shelf
  e.POST("/shelves", diary_handlers.CreateNewShelf(db))
  // removes shelves
  e.DELETE("/shelves/:id", diary_handlers.RemoveShelf(db))
  // renames shelf
  e.PUT("/shelves/:id", diary_handlers.EditShelf(db))
  // copies shelf to user's shelves
  e.POST("/shelves/:id/copy", diary_handlers.CopyShelf(db))
  // adds book to shelf
  e.POST("/shelves/:id", diary_handlers.AddBookToShelf(db))
  // removes book from shelf
  e.DELETE("/shelves/:id/:bookId", diary_handlers.RemoveBookFromShelf(db))

  // starts book tracking
  e.PUT("/tracking/start/:id", diary_handlers.StartTracking(db))
  // stops book tracking
  e.PUT("/tracking/stop/:id", diary_handlers.StopTracking(db))
  // gets book tracking
  e.GET("/tracking/book/:id", diary_handlers.GetUserBookTracking(db))
  // gets user's latest tracing
  e.GET("/tracking", diary_handlers.GetLastTracking(db))

  // gets all user's reading
  e.GET("/readings", diary_handlers.GetAllUsersReadings(db))
  // gets user's statistics
  e.GET("/statistic", diary_handlers.GetUserStatistic(db))

  // gets user's intervals for period
  e.GET("/intervals", diary_handlers.GetIntervals(db))

  // generates book detail pdf
  e.GET("/pdfBook/:id", diary_handlers.GenerateBookPdf(db))
  // generates list of books pdf
  e.GET("/pdfBooks/:status", diary_handlers.GenerateListOfBooks(db))

  // generates list of books txt
  e.GET("/txt/:status", diary_handlers.GenerateListOfBooksText(db))

  // adds new user to database
  e.POST("/new", diary_handlers.NewUser(db))

  // gets user's friends
  e.GET("/friends", diary_handlers.GetFriends(db))
  // gets specific friend
  e.GET("/friends/:id", diary_handlers.GetFriend(db))
  // remove friend
  e.DELETE("/friends/:id", diary_handlers.RemoveFriend(db))
  // gets friend requests
  e.GET("/friends/requests", diary_handlers.GetUserFriendRequests(db))
  // adds friend request
  e.POST("/friends/requests/:id", diary_handlers.AddUserFriendRequest(db))
  // removes friend requests
  e.DELETE("/friends/requests/:id", diary_handlers.DeleteFriendRequest(db))
  // accepts friend request
  e.POST("/friends/requests/:id/accept", diary_handlers.AcceptFriendRequest(db))
  // declines friend request
  e.POST("/friends/requests/:id/decline", diary_handlers.DeclineFriendRequest(db))

  // search among users
  e.GET("/people", diary_handlers.GetPeople(db))

  // starts server
  e.Logger.Fatal(e.Start(":1323"))
}
