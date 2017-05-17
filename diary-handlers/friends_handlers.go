package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "github.com/PetrNavratil/diary-back/models"
  "net/http"
  "strconv"
)

// Function returns user's friends requests
func GetUserFriendRequests(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    incomingRequests := []models.FriendRequestResponse{}
    outgoingRequests := []models.FriendRequestResponse{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get incoming requests
      db.Table("friend_requests").
        Select("friend_requests.id as id,user_name,first_name,last_name,avatar, requester_id, user_id").
        Joins("JOIN users ON friend_requests.requester_id = users.id").
        Where("user_id = ? ", user.ID).
        Scan(&incomingRequests)
      // get outgoing request
      db.Table("friend_requests").
        Select("friend_requests.id as id,user_name,first_name,last_name,avatar, requester_id, user_id").
        Joins("JOIN users ON friend_requests.user_id = users.id").
        Where("requester_id = ? ", user.ID).
        Scan(&outgoingRequests)
      // return them as joined arrays
      return c.JSON(http.StatusOK, append(incomingRequests, outgoingRequests...))
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  err.Error()})
    }
  }
}

// Function adds friend request
func AddUserFriendRequest(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get friend id
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        // get friend
        if !db.Table("users").Where("id = ?", id).RecordNotFound() {
          request := models.FriendRequest{}
          // check if request exist in database already
          if db.Where("(user_id = ? AND requester_id = ?) OR (user_id = ? AND requester_id = ?)", user.ID, id, id, user.ID).First(&request).RecordNotFound() {
            request.UserID = id
            request.RequesterID = user.ID
            // save request and create response
            db.Save(&request)
            outgoingRequest := models.FriendRequestResponse{}
            db.Table("friend_requests").
              Select("friend_requests.id as id,user_name,first_name,last_name,avatar, requester_id, user_id").
              Joins("JOIN users ON friend_requests.user_id = users.id").
              Where("requester_id = ? AND user_id = ? ", request.RequesterID, request.UserID).
              Scan(&outgoingRequest)
            return c.JSON(http.StatusOK, outgoingRequest)
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

// Function accepts friend request
func AcceptFriendRequest(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get request id
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        request := models.FriendRequest{}
        if !db.First(&request, id).RecordNotFound() {
          friend1 := models.Friend{}
          friend1.FriendID = request.RequesterID
          friend1.UserID = user.ID
          // save friend to user
          db.Save(&friend1)
          friend2 := models.Friend{}
          friend2.UserID = request.RequesterID
          friend2.FriendID = user.ID
          // save user to friend
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

// Function decline friend request
func DeclineFriendRequest(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    // get request id
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      request := models.FriendRequest{}
      if !db.First(&request, id).RecordNotFound() {
        // remove request
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

// Function deletes friend request
func DeleteFriendRequest(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    // get request id
    if id, err := strconv.Atoi(c.Param("id")); err == nil {
      request := models.FriendRequest{}
      if !db.First(&request, id).RecordNotFound() {
        // remove request
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

// Function returns user's friends
func GetFriends(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    friends := []models.FriendResponse{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get friends
      db.Table("friends").
        Select("friend_id,user_name,first_name,last_name,avatar, created_at").
        Joins("JOIN users ON friends.friend_id = users.id").
        Where("friends.user_id = ? ", user.ID).
        Scan(&friends)
      // get friends books and shelves count
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

// Function returns specific friend
func GetFriend(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    friend := models.FriendResponse{}
    friendUser := models.User{}
    // get user
    if _, err := GetUser(c, db); err == nil {
      // get friend id
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        // get friend
        if !db.First(&friendUser, id).RecordNotFound() {
          // get friend info
          db.Table("friends").
            Select("friend_id,user_name,first_name,last_name,avatar, created_at").
            Joins("JOIN users ON friends.friend_id = users.id").
            Where("friends.friend_id = ? ", friendUser.ID).
            Scan(&friend)
          // get friends book
          db.Table("books").Select(
            "books.id, books.title, books.author, books.image_url, user_book.status").
            Joins("JOIN user_book ON user_book.book_id = books.id").Where("user_id = ?", friendUser.ID).Scan(&friend.Books)
          // get friends shelves with books
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

// Function removes friend
func RemoveFriend(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    friend := models.Friend{}
    // get user
    if user, err := GetUser(c, db); err == nil {
      // get friend user id
      if id, idErr := strconv.Atoi(c.Param("id")); idErr == nil {
        // get friend
        if !db.Table("friends").Where("user_id = ? AND friend_id = ?", user.ID, id).First(&friend).RecordNotFound() {
          // remove user from friend friends
          db.Delete(models.Friend{}, "user_id = ? AND friend_id = ?", id, user.ID)
          // remove friend from users friends
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


