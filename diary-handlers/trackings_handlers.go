package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "strconv"
  "github.com/PetrNavratil/diary-back/models"
  "time"
)

// Function starts tracking of book
func StartTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    userBook := models.UserBook{}
    book := models.Book{}
    reading := models.Reading{}
    interval := models.Interval{}
    lastInterval := models.LastInterval{}
    returnReading := models.ReturnReading{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get id of book
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        // get user book info
        if !db.Where("book_id = ? AND user_id = ?", id, user.ID).First(&userBook).RecordNotFound() {
          // stop reading if any is active
          if !db.Table("intervals").Joins("JOIN readings on intervals.reading_id = readings.id").Where("user_id = ? AND completed = ?", user.ID, false).Last(&interval).RecordNotFound() {
            if (interval.Stop.IsZero()) {
              interval.Stop = time.Now()
              db.Save(&interval)
            }
          } else {
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
          } else {
            // if already reading get that reading
            db.Where("user_id = ? AND book_id = ?", user.ID, id).Last(&reading)
          }
          // start reading interval
          interval = models.Interval{}
          interval.Start = time.Now()
          interval.ReadingID = reading.ID
          db.Create(&interval)
          // prepare last interval because it has been changed
          db.Where("id = ? ", userBook.BookID).First(&book)
          lastInterval.Interval = interval
          lastInterval.Title = book.Title
          lastInterval.Author = book.Author
          lastInterval.Completed = false
          lastInterval.BookID = book.ID
          getReadings, _ := strconv.ParseBool(c.QueryParam("getReadings"))
          if getReadings {
            returnReading.LastInterval = lastInterval
            returnReading.Readings = GetUserBookReadings(db, user.ID, id)
            return c.JSON(http.StatusOK, returnReading)
          } else {
            return c.JSON(http.StatusOK, lastInterval)
          }
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

// Function stops book tracking
func StopTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := models.Book{}
    interval := models.Interval{}
    lastInterval := models.LastInterval{}
    returnReading := models.ReturnReading{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get book id
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        // get book
        if !db.Where("id = ?", id).First(&book).RecordNotFound() {
          // get this book's last reading interval
          db.Table("intervals").Joins("JOIN readings on intervals.reading_id = readings.id").
            Where("user_id = ? AND completed = ? AND book_id = ?", user.ID, false, book.ID).Last(&interval)
          // stop reading now
          interval.Stop = time.Now()
          db.Save(&interval)
          // prepare last reading
          lastInterval.Interval = interval
          lastInterval.Title = book.Title
          lastInterval.Author = book.Author
          lastInterval.Completed = true
          lastInterval.BookID = book.ID
          getReadings, _ := strconv.ParseBool(c.QueryParam("getReadings"))
          if getReadings {
            returnReading.LastInterval = lastInterval
            returnReading.Readings = GetUserBookReadings(db, user.ID, id)
            return c.JSON(http.StatusOK, returnReading)
          } else {
            return c.JSON(http.StatusOK, lastInterval)
          }
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

// Function return user book reading with intervals
func GetUserBookTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := models.Book{}
    readings := []models.Reading{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get book id
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        // get book
        if !db.Where("id = ?", id).First(&book).RecordNotFound() {
          // gets book's readings with intervals
          readings = GetUserBookReadings(db, user.ID, id)
          return c.JSON(http.StatusOK, readings)
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

// Function returns last reading interval of user
func GetLastTracking(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    book := models.Book{}
    reading := models.Reading{}
    lastInterval := models.LastInterval{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get last interval of user
      db.Table("intervals").Joins("JOIN readings on intervals.reading_id = readings.id").Where("user_id = ?", user.ID).Last(&lastInterval.Interval)
      if lastInterval.ID > 0 {
        // reading exist find all information about book being reading
        db.Where("id = ?", lastInterval.ReadingID).First(&reading)
        db.Where("id = ?", reading.BookID).First(&book)
        lastInterval.Author = book.Author
        lastInterval.Title = book.Title
        lastInterval.BookID = book.ID
        if lastInterval.Stop.IsZero() {
          lastInterval.Completed = false
        } else {
          lastInterval.Completed = true
        }
      }
      return c.JSON(http.StatusOK, lastInterval)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function gets all user's readings with intervals
func GetUserBookReadings(db *gorm.DB, userId int, bookId int) []models.Reading {
  readings := []models.Reading{}
  db.Where("user_id = ? AND book_id = ?", userId, bookId).Find(&readings)
  for i := range readings {
    db.Model(&readings[i]).Related(&readings[i].Intervals, "Intervals")
  }
  return readings
}
