package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "strconv"
  "github.com/PetrNavratil/diary-back/models"
  "time"
  "github.com/davecgh/go-spew/spew"
  "fmt"
)

func StartTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    userBook := models.UserBook{}
    book := models.Book{}
    reading := models.Reading{}
    interval := models.Interval{}
    lastInterval := models.LastInterval{}
    intervals := []models.Interval{}
    if user, err := GetUser(c, db); err == nil {
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        if !db.Where("book_id = ? AND user_id = ?", id, user.ID).First(&userBook).RecordNotFound() {
          // stop reading if any is active
          if !db.Table("intervals").Joins("JOIN readings on intervals.reading_id = readings.id").Where("user_id = ? AND completed = ?", user.ID, false).Last(&interval).RecordNotFound() {
            interval.Stop = time.Now()
            db.Save(&interval)
            fmt.Println("ENDING PREVIOUS READING")
            spew.Dump(interval)
          } else {
            fmt.Println("NO LAST TRACKING")
          }
          // check if book is already reading if not create new reading record
          if userBook.Status != models.READING {
            userBook.Status = models.READING
            db.Save(&userBook)
            reading.UserID = user.ID
            reading.BookID = userBook.BookID
            reading.Completed = false
            reading.Start = time.Now()
            db.Create(&reading)
            fmt.Println("SHOUL BE HERE")
          } else {
            db.Where("user_id = ? AND book_id = ?", user.ID, id).Last(&reading)
          }
          spew.Dump(&reading)
          interval = models.Interval{}
          interval.Start = time.Now()
          interval.ReadingID = reading.ID
          db.Create(&interval)
          db.Where("id = ? ", userBook.BookID).First(&book)
          lastInterval.Interval = interval
          lastInterval.Title = book.Title
          lastInterval.Author = book.Author
          lastInterval.Completed = false
          db.Find(&intervals)
          spew.Dump(intervals)
          return c.JSON(http.StatusOK, lastInterval)

        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func StopTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := models.Book{}
    interval := models.Interval{}
    lastInterval := models.LastInterval{}
    returnReading := models.ReturnReading{}
    if user, err := GetUser(c, db); err == nil {
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        if !db.Where("id = ?", id).First(&book).RecordNotFound() {
          db.Table("intervals").Joins("JOIN readings on intervals.reading_id = readings.id").Where("user_id = ? AND completed = ? AND book_id = ?", user.ID, false, book.ID).Last(&interval)
          interval.Stop = time.Now()
          db.Save(&interval)
          lastInterval.Interval = interval
          lastInterval.Title = book.Title
          lastInterval.Author = book.Author
          lastInterval.Completed = true
          db.Where("user_id = ? and book_id = ?", user.ID, book.ID).Find(&returnReading.Readings)
          for i := range returnReading.Readings {
            db.Model(&returnReading.Readings[i]).Related(&returnReading.Readings[i].Intervals, "Intervals")
          }
          returnReading.LastInterval = lastInterval
          return c.JSON(http.StatusOK, returnReading)
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func GetUserBookTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := models.Book{}
    returnReading := models.ReturnReading{}
    if user, err := GetUser(c, db); err == nil {
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        if !db.Where("id = ?", id).First(&book).RecordNotFound() {
          db.Where("user_id = ? AND book_id = ?", user.ID, id).Find(&returnReading.Readings)
          for i := range returnReading.Readings {
            db.Model(&returnReading.Readings[i]).Related(&returnReading.Readings[i].Intervals, "Intervals")
          }
          return c.JSON(http.StatusOK, returnReading)
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "Bad book id"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func GetLastTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := models.Book{}
    reading := models.Reading{}
    lastInterval := models.LastInterval{}
    if user, err := GetUser(c, db); err == nil {
      db.Table("intervals").Joins("JOIN readings on intervals.reading_id = readings.id").Where("user_id = ? AND completed = ?", user.ID, false).Last(&lastInterval.Interval)
      db.Where("id = ?", lastInterval.ID).First(&reading)
      db.Where("id = ?", reading.BookID).First(&book)
      lastInterval.Author = book.Author
      lastInterval.Title = book.Title
      if lastInterval.Stop.IsZero() {
        lastInterval.Completed = false
      } else {
        lastInterval.Completed = true
      }
      return c.JSON(http.StatusOK, lastInterval)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}
