package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "github.com/PetrNavratil/diary-back/models"
  "time"
)

func GetUserStatistic(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    statistic := models.Statistic{}
    books := []models.UserBook{}
    readings := []models.Reading{}
    readCount := -1
    mostReadBook := models.UserBook{}
    counter := 0
    if user, err := GetUser(c, db); err == nil {
      db.Where("user_id = ?", user.ID).Find(&books)
      statistic.BooksCount = len(books)
      statistic.BooksRead = len(filterBooks(books, models.READ))
      statistic.BooksReading = len(filterBooks(books, models.READING))
      statistic.BooksToRead = len(filterBooks(books, models.TO_READ))
      statistic.BooksNotRead = len(filterBooks(books, models.NOT_READ))
      db.Where("user_id = ?", user.ID).Find(&readings)
      for i := range readings {
        db.Model(&readings[i]).Related(&readings[i].Intervals, "Intervals")
      }

      for _, reading := range readings {
        for _, interval := range reading.Intervals {
          // not ended yet use now
          if interval.Stop.IsZero() {
            statistic.TimeSpentReading += time.Now().Sub(interval.Start).Nanoseconds()
          } else {
            statistic.TimeSpentReading += interval.Stop.Sub(interval.Start).Nanoseconds()
          }
        }
      }
      for _, book := range books {
        counter = 0
        for _, reading := range readings {
          if reading.BookID == book.BookID {
            counter = counter + 1
          }
        }
        if counter > readCount {
          readCount = counter
          mostReadBook = book
        }
      }
      statistic.MostlyReadBook.Read = readCount
      db.Where("id = ?", mostReadBook.BookID).First(&statistic.MostlyReadBook.Book)
      statistic.TimeSpentReading = statistic.TimeSpentReading / 1e6;
      return c.JSON(http.StatusOK, statistic)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func filterBooks(filterArray []models.UserBook, status int) []models.UserBook {
  var newArray []models.UserBook
  for _, value := range filterArray {
    if value.Status == status {
      newArray = append(newArray, value)
    }
  }
  return newArray
}