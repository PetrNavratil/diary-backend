package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
  "strconv"
)

func GetUserFriendRequests(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    requests := []models.FriendRequestResponse{}
    if user, err := GetUser(c, db); err == nil {
      db.Table("friend_requests").
        Select("friend_requests.id as id,user_name,first_name,last_name,avatar").
        Joins("JOIN users ON friend_requests.requester_id = users.id").
        Where("user_id = ? ", user.ID).
        Scan(&requests)
      return c.JSON(http.StatusOK, requests)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func AddUserFriendRequest(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if user, err := GetUser(c, db); err == nil {
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        if !db.Table("users").Where("id = ?", id).RecordNotFound() {
          request := models.FriendRequest{}
          if db.Where("(user_id = ? AND requester_id = ?) OR (user_id = ? AND requester_id = ?)", user.ID, id, id, user.ID).First(&request).RecordNotFound() {
            request.UserID = id
            request.RequesterID = user.ID
            db.Save(&request)
            return c.JSON(http.StatusOK, request)
          } else {
            return c.JSON(http.StatusBadRequest, map[string]string{"message": "already requested"})
          }
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message": "user id not in database"})
        }

      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "id must be specified"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func AcceptFriendRequest(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if user, err := GetUser(c, db); err == nil {
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        request := models.FriendRequest{}
        if !db.First(&request, id).RecordNotFound() {
          friend1 := models.Friend{}
          friend1.FriendID = request.RequesterID
          friend1.UserID = user.ID
          db.Save(&friend1)
          friend2 := models.Friend{}
          friend2.UserID = request.RequesterID
          friend2.FriendID = user.ID
          db.Save(&friend2)
          db.Delete(&request)
          return c.JSON(http.StatusOK, request)
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message": "id of request is not valid"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "id must be specified"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func DeclineFriendRequest(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      request := models.FriendRequest{}
      if !db.First(&request, id).RecordNotFound() {
        db.Delete(&request)
        return c.JSON(http.StatusOK, request)
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "id of request is not valid"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message": "id must be specified"})
    }
  }
}

func GetFriends(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    friends := []models.FriendResponse{}
    if user, err := GetUser(c, db); err == nil {
      db.Table("friends").
        Select("friend_id,user_name,first_name,last_name,avatar, created_at").
        Joins("JOIN users ON friends.friend_id = users.id").
        Where("friends.user_id = ? ", user.ID).
        Scan(&friends)
      for i := range friends {
        db.Table("user_book").Where("user_id = ?", friends[i].FriendID).Count(&friends[i].BooksCount)
        db.Table("shelves").Where("user_id = ?", friends[i].FriendID).Count(&friends[i].ShelvesCount)
      }
      return c.JSON(http.StatusOK, friends)
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func GetFriend(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    friend := models.FriendResponse{}
    friendUser := models.User{}
    if user, err := GetUser(c, db); err == nil {
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        if !db.First(&friendUser, id).RecordNotFound() {
          db.Table("friends").
            Select("friend_id,user_name,first_name,last_name,avatar, created_at").
            Joins("JOIN users ON friends.friend_id = users.id").
            Where("friends.user_id = ? ", user.ID).
            Scan(&friend)
          db.Table("books").Select(
            "books.id, books.title, books.author, books.image_url, user_book.status").
            Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ?", friendUser.ID).Scan(&friend.Books)
          db.Model(&friendUser).Related(&friend.Shelves)
          for i := range friend.Shelves {
            db.Model(&friend.Shelves[i]).Related(&friend.Shelves[i].Books, "Books")
          }
          friend.BooksCount = len(friend.Books)
          friend.ShelvesCount = len(friend.Shelves)
          return c.JSON(http.StatusOK, friend)
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "id of user not valid"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "id of user required"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

func RemoveFriend(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    friend := models.Friend{}
    if user, err := GetUser(c, db); err == nil {
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        if !db.Table("friends").Where("user_id = ? AND friend_id = ?", user.ID, id).First(&friend).RecordNotFound() {
          db.Delete(models.Friend{}, "user_id = ? AND friend_id = ?", id, user.ID)
          db.Delete(&friend)
          return c.JSON(http.StatusOK, friend)
        } else {
          return c.JSON(http.StatusBadRequest, map[string]string{"message":  "id of user not valid"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "id of user required"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}


