package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "github.com/PetrNavratil/diary-back/models"
  "strconv"
)

// Function return book's comments
func GetBookComments(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    comments := []models.CommentResponse{}
    book := models.Book{}
    // get book id
    if id, err := strconv.Atoi(c.QueryParam("bookId")); err == nil {
      // get book
      if db.Where("id = ?", id).First(&book).RecordNotFound() {
        return c.JSON(http.StatusOK, comments)
      }
      // select information
      db.Table("comments").Where("book_id = ?", book.ID).Joins("JOIN users on users.id = comments.user_id").
        Select("comments.id, comments.book_id, comments.text, users.avatar, users.user_name, users.last_name, users.first_name, users.id as user_id, comments.created_at, comments.updated_at").
        Scan(&comments)
      // get times of create/update
      for i := range comments {
        if comments[i].UpdatedAt.IsZero() {
          comments[i].Date = comments[i].UpdatedAt
        } else {
          comments[i].Date = comments[i].CreatedAt
        }
      }
      return c.JSON(http.StatusOK, comments)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function adds comment to the book
func AddBookComment(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    newComment := &models.BookComment{}
    comment := models.Comment{}
    commentResponse := models.CommentResponse{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get comment
      if bodyError := c.Bind(newComment); bodyError == nil {
        comment.BookID = newComment.BookId
        comment.UserID = user.ID
        comment.Text = newComment.Text
        // save
        db.Create(&comment)
        // get response for FE
        db.Table("comments").Where("book_id = ? AND user_id = ?", comment.BookID, comment.UserID).Joins("JOIN users on users.id = comments.user_id").
          Select("comments.id, comments.book_id, comments.text, users.avatar, users.user_name, users.last_name, users.first_name, users.id as user_id, comments.created_at, comments.updated_at").
          Scan(&commentResponse)
        commentResponse.Date = commentResponse.CreatedAt
        return c.JSON(http.StatusOK, commentResponse)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "BAD comment body"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function remove comment from book
func DeleteBookComment(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    commentResponse := models.CommentResponse{}
    comment := models.Comment{}
    // get id of comment
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      if db.Where("id = ?", id).First(&comment).RecordNotFound() {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "bad comment id"})
      } else {
        // create response
        db.Table("comments").Where("book_id = ? AND user_id = ?", comment.BookID, comment.UserID).Joins("JOIN users on users.id = comments.user_id").
          Select("comments.id, comments.book_id, comments.text, users.avatar, users.user_name, users.id").
          Scan(&commentResponse)
        // remove comment
        db.Delete(&comment)
        return c.JSON(http.StatusOK, commentResponse)
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "bad book id"})
    }
  }
}

// Function updates comment
func UpdateBookComment(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    commentBody := &models.CommentResponse{}
    comment := models.Comment{}
    // get comment id
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      // get comment
      if db.Where("id = ?", id).First(&comment).RecordNotFound() {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "bad comment id"})
      } else {
        // get new commend body
        if bodyError := c.Bind(commentBody); bodyError == nil {
          // save text
          comment.Text = commentBody.Text
          db.Save(&comment)
          commentBody.Date = comment.UpdatedAt
          return c.JSON(http.StatusOK, commentBody)
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "BAD comment body"})
        }
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  "bad book id"})
    }
  }
}
