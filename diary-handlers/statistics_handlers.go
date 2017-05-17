package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "github.com/PetrNavratil/diary-back/models"
  "time"
  "strconv"
  "fmt"
)

// Function returns user's statistics
func GetUserStatistic(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    statistic := models.Statistic{}
    books := []models.UserBook{}
    readings := []models.Reading{}
    readCount := -1
    mostReadBook := models.UserBook{}
    counter := 0
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get user's book and all category counts
      db.Where("user_id = ?", user.ID).Find(&books)
      statistic.BooksCount = len(books)
      statistic.BooksRead = len(filterBooks(books, models.READ))
      statistic.BooksReading = len(filterBooks(books, models.READING))
      statistic.BooksToRead = len(filterBooks(books, models.TO_READ))
      statistic.BooksNotRead = len(filterBooks(books, models.NOT_READ))
      // get user's readings width their intervals
      db.Where("user_id = ?", user.ID).Find(&readings)
      for i := range readings {
        db.Model(&readings[i]).Related(&readings[i].Intervals, "Intervals")
      }
      // count duration of reading in nanoseconds
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
      // get mostly read book
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
      // get mostly read book info
      db.Where("id = ?", mostReadBook.BookID).First(&statistic.MostlyReadBook.Book)
      // get duration in milliseconds
      statistic.TimeSpentReading = statistic.TimeSpentReading / 1e6;
      return c.JSON(http.StatusOK, statistic)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function returns intervals of reading period
func GetIntervals(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    intervals := []models.StatisticInterval{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get month and year
      if month, monthErr := strconv.Atoi(c.QueryParam("month")); monthErr == nil {
        if year, yearErr := strconv.Atoi(c.QueryParam("year")); yearErr == nil {
          // get intervals where month and year fits
          // take intervals starting in period and also ending in period
          db.Table("readings").
            Select("intervals.stop, intervals.start, readings.book_id, books.title, books.author").
            Joins("JOIN intervals ON readings.id = intervals.reading_id").
            Joins("JOIN books ON readings.book_id = books.id").
            Where("readings.user_id = ? AND" +
            "(" +
            "(strftime('%Y', intervals.start) = ? AND strftime('%m', intervals.start) = ?)" +
            "OR" +
            "(strftime('%Y', intervals.stop) = ? AND strftime('%m', intervals.stop) = ?)" +
            ")", user.ID, strconv.Itoa(year), fmt.Sprintf("%02d", month), strconv.Itoa(year), fmt.Sprintf("%02d", month)).
            Scan(&intervals)
          // sets complete attribute
          for i := range intervals {
            if !intervals[i].Stop.IsZero() {
              intervals[i].Completed = true
            }
          }
          return c.JSON(http.StatusOK, intervals)
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "month required"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "month required"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function filters books by provided status
func filterBooks(filterArray []models.UserBook, status int) []models.UserBook {
  var newArray []models.UserBook
  for _, value := range filterArray {
    if value.Status == status {
      newArray = append(newArray, value)
    }
  }
  return newArray
}