package diary_handlers

import (
  "github.com/jinzhu/gorm"
  "github.com/labstack/echo"
  "net/http"
  "github.com/jung-kurt/gofpdf"
  "fmt"
  "github.com/PetrNavratil/diary-back/models"
  "strconv"
  "github.com/parnurzeal/gorequest"
  "encoding/xml"
  "github.com/kennygrant/sanitize"
  "io/ioutil"
  "time"
  "os"
)

// Function generates book detail PDF
func GenerateBookPdf(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    // set page properties
    const (
      pageWidth = 210
      imageWidth = 60
      imageMargin = 10
      labelWidth = 50
    )
    book := models.Book{}
    bookInfo := &models.GoodReadsBook{}
    comment := models.Comment{}
    userBook := models.UserBook{}
    readings := []models.Reading{}
    // get user
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      // get book id
      if id, err := strconv.Atoi(c.Param("id")); err == nil {
        // get book
        db.First(&book, id)
        // get book detail info
        _, body, errs := gorequest.New().Get("https://www.goodreads.com/book/show/" + strconv.Itoa(book.GRBookId) + ".xml?key=tsRkj9chcP8omCKBCJLg0A&").End()
        if errs == nil {
          xmlResponse := []byte(body)
          xml.Unmarshal(xmlResponse, bookInfo)

          var imageY float64
          // create pdf
          pdf := gofpdf.New("P", "mm", "A4", "")
          tr := pdf.UnicodeTranslatorFromDescriptor("cp1250")
          // set header function
          pdf.SetHeaderFunc(func() {
            pdf.SetFont("helvetica", "B", 16)
            wd := pdf.GetStringWidth(bookInfo.Title) + 20
            pdf.SetX((210 - wd) / 2)
            pdf.SetFillColor(95, 78, 63)
            pdf.SetTextColor(217, 217, 217)
            pdf.CellFormat(wd, 10, bookInfo.Title, "1", 1, "MC", true, 0, "")
            pdf.Ln(5)
            pdf.Line(20, pdf.GetY(), pageWidth - 20, pdf.GetY())
            pdf.Ln(7)
          })
          // set footer function
          pdf.SetFooterFunc(func() {
            pdf.SetFont("helvetica", "", 10)
            pdf.SetFillColor(255, 255, 255, )
            pdf.SetY(297 - 10)
            pdf.CellFormat(0, 10, bookInfo.Title, "T", 1, "MR", true, 0, "")
          })
          // add new page
          pdf.AddPage()
          // get book cover
          rsp, err := http.Get(bookInfo.ImageUrl)
          if err == nil {
            tp := pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
            // register image for pdf
            pdf.RegisterImageReader(bookInfo.ImageUrl, tp, rsp.Body)
            if pdf.Ok() {
              currentY := pdf.GetY()
              // place image to the pdf
              pdf.Image(bookInfo.ImageUrl, pageWidth - imageWidth, pdf.GetY(), imageWidth - imageMargin, 0, true, tp, 0, "")
              imageY = pdf.GetY()
              // set Y to the line where image starts
              pdf.SetY(currentY)
            }
          }
          // function for writing text in a row
          createRow := func(label, value string) {
            pdf.SetFont("helvetica", "B", 14)
            pdf.CellFormat(labelWidth, pdf.PointConvert(14) + 3, tr(label + ":"), "", 0, "ML", false, 0, "")
            pdf.SetFont("helvetica", "", 14)
            // check if image is still on the right side
            if imageY > pdf.GetY() {
              pdf.MultiCell(pageWidth - imageWidth - labelWidth - imageMargin, pdf.PointConvert(14) + 3, tr(value), "", "ML", false)
            } else {
              pdf.MultiCell(0, pdf.PointConvert(14) + 3, tr(value), "", "J", false)
            }
          }
          pdf.SetY(pdf.GetY() + 10)
          // write book info
          createRow("Title", bookInfo.Title)
          createRow("Author", bookInfo.Authors[0].Name)
          createRow("Publisher", bookInfo.Publisher)
          createRow("Publicated", bookInfo.PublicationDay + ". " + bookInfo.PublicationMonth + ". " + bookInfo.PublicationYear)
          createRow("Pages", bookInfo.Pages)
          pdf.SetY(imageY + 5)
          createRow("Description", sanitize.HTML(bookInfo.Description))
          pdf.Ln(5)

          // get readings of this book and write count
          db.Where("user_id = ? AND book_id = ?", loggedUser.ID, id).Find(&readings)
          if len(readings) > 0 {
            if (readings[len(readings) - 1].Completed) {
              createRow("Read", fmt.Sprintf("%dx", len(readings)))
            } else {
              createRow("Read", fmt.Sprintf("%dx", len(readings) - 1))
            }
          }
          // write user comment for book
          if !db.Where("user_id = ? AND book_id = ?", loggedUser.ID, id).First(&comment).RecordNotFound() {
            createRow("Comment", comment.Text)
          } else {
            createRow("Comment", "No comment for this book")
          }
          // write educational for this book
          db.Where("user_id = ? AND book_id = ?", loggedUser.ID, id).First(&userBook)
          db.Model(&userBook).Related(&userBook.Educational)
          if len(userBook.Educational.Smer) > 0 {
            pdf.AddPage()
            imageY = 0
            pdf.SetFont("helvetica", "B", 15)
            pdf.CellFormat(0, pdf.PointConvert(15) + 3, tr("Češtinářská část"), "", 0, "", false, 0, "")
            pdf.Ln(10)
            createRow("Smer", userBook.Educational.Smer)
            createRow("Druh", userBook.Educational.Druh)
            createRow("Zanr", userBook.Educational.Zanr)
            createRow("Forma", userBook.Educational.Forma)
            createRow("Jazyk", userBook.Educational.Jazyk)
            createRow("Postavy", userBook.Educational.Postavy)
            createRow("Obsah", userBook.Educational.Obsah)
            createRow("Tema", userBook.Educational.Tema)
            createRow("Hodnoceni", userBook.Educational.Hodnoceni)
          }
          // create filename
          fileName := fmt.Sprintf("detail-%d.pdf", loggedUser.ID)
          // create pdf and save it
          pdf.OutputFileAndClose(fileName)
          // open pdf file and read it
          dat, _ := ioutil.ReadFile(fileName)
          // remove pdf file
          os.Remove(fileName)
          // send pdf as BLOB
          return c.Blob(http.StatusOK, "application/pdf", dat)
        } else {
          return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
        }
      } else {
        return c.JSON(http.StatusBadRequest, map[string]string{"message":  "FAIL"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  logErr.Error()})
    }
  }
}


// Function generates list of user's books pdf
func GenerateListOfBooks(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    // set page properties
    const (
      pageWidth = 210
      imageWidth = 40
      imageMargin = 10
      labelWidth = 50
    )
    userBooks := []models.UserBook{}
    // get user
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      // get status
      if status, err := strconv.Atoi(c.Param("status")); err == nil {
        // create pdf
        pdf := gofpdf.New("P", "mm", "A4", "")
        tr := pdf.UnicodeTranslatorFromDescriptor("cp1250")
        pdf.SetFont("helvetica", "", 16)
        fontSize := 16.0
        pdf.SetFont("helvetica", "", fontSize)
        // set header function
        pdf.SetHeaderFunc(func() {
          pdf.SetFont("helvetica", "B", 16)
          wd := pdf.GetStringWidth("Seznam knih") + 20
          pdf.SetX((210 - wd) / 2)
          pdf.SetFillColor(95, 78, 63)
          pdf.SetTextColor(217, 217, 217)
          pdf.CellFormat(wd, 10, "Seznam knih", "1", 1, "MC", true, 0, "")
          pdf.Ln(5)
          pdf.Line(20, pdf.GetY(), pageWidth - 20, pdf.GetY())
          pdf.Ln(7)
        })
        // set footer function
        pdf.SetFooterFunc(func() {
          pdf.SetFont("helvetica", "", 10)
          pdf.SetFillColor(255, 255, 255, )
          pdf.SetY(297 - 10)
          pdf.CellFormat(0, 10, "Seznam knih " + time.Now().Format("Mon Jan _2 15:04:05 2006"), "T", 1, "MR", true, 0, "")
        })

        // function for adding cover to the pdf
        addCover := func(url string) float64 {
          rsp, err := http.Get(url)
          var newY float64
          if err == nil {
            tp := pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
            // register image for pdf
            pdf.RegisterImageReader(url, tp, rsp.Body)
            if pdf.Ok() {
              currentY := pdf.GetY()
              // place image to the pdf
              pdf.Image(url, pageWidth - imageWidth, currentY, imageWidth - imageMargin, 0, true, tp, 0, "")
              newY = pdf.GetY()
              // set Y to the line where image starts
              pdf.SetY(currentY)
            }
          }
          return newY
        }
        // function for writing text in a row
        writeRow := func(label, value string) {
          pdf.SetFont("helvetica", "B", 14)
          pdf.CellFormat(labelWidth, pdf.PointConvert(14) + 3, tr(label + ":"), "", 0, "ML", false, 0, "")
          pdf.SetFont("helvetica", "", 14)
          pdf.MultiCell(pageWidth - imageWidth - labelWidth - imageMargin, pdf.PointConvert(14) + 3, tr(value), "", "ML", false)
        }

        pdf.AddPage()
        // get user's books
        if status == models.ALL {
          db.Where("user_id = ? ", loggedUser.ID).Find(&userBooks)
        } else {
          db.Where("user_id = ? AND status = ?", loggedUser.ID, status).Find(&userBooks)
        }
        if len(userBooks) > 0 {
          // only 4 books fit page
          item := 0
          for index, userBook := range userBooks {
            if index != 0 && index % 4 == 0 {
              pdf.AddPage()
              item = 0
            }
            item++
            book := models.Book{}
            // get user book
            db.First(&book, userBook.BookID)
            // write cover and info
            newY := addCover(book.ImageUrl)
            writeRow("Title", book.Title)
            writeRow("Author", book.Author)
            if status == models.READ || status == models.READING || status == models.ALL {
              var read int
              db.Table("readings").Where("user_id = ? AND book_id = ? AND completed = ?", loggedUser.ID, book.ID, true).Count(&read)
              writeRow("Read", fmt.Sprintf("%dx", read))
            }
            pdf.SetY(newY)
            // place v space
            if item != 4 {
              pdf.Ln(8)
              pdf.Line(40, pdf.GetY(), pageWidth - 40, pdf.GetY())
              pdf.Ln(8)
            }
          }
        }
        // create filename
        fileName := fmt.Sprintf("books-%d.pdf", loggedUser.ID)
        // create pdf and save it
        pdf.OutputFileAndClose(fileName)
        // open pdf file and read it
        dat, _ := ioutil.ReadFile(fileName)
        // remove pdf file
        os.Remove(fileName)
        // send pdf as BLOB
        return c.Blob(http.StatusOK, "application/pdf", dat)
      } else {
        return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  logErr.Error()})
    }
  }
}

// Function generates list of user's books txt
func GenerateListOfBooksText(db *gorm.DB) func(c echo.Context) error {
  return func(c echo.Context) error {
    userBooks := []models.UserBook{}
    // get user
    if loggedUser, logErr := GetUser(c, db); logErr == nil {
      // get status
      if status, err := strconv.Atoi(c.Param("status")); err == nil {
        // get books
        if status == models.ALL {
          db.Where("user_id = ? ", loggedUser.ID).Find(&userBooks)
        } else {
          db.Where("user_id = ? AND status = ?", loggedUser.ID, status).Find(&userBooks)
        }
        // create filename
        filename := fmt.Sprintf("seznam-%d.txt", loggedUser.ID)
        // open file
        f, _ := os.Create(filename)
        // write header
        f.WriteString("Book list\n\n")
        if len(userBooks) > 0 {
          // write information about each book
          for _, userBook := range userBooks {
            book := models.Book{}
            db.First(&book, userBook.BookID)
            f.WriteString("--------------------------------------------------------------\n")
            f.WriteString(fmt.Sprintf("Title:\t%s\n", book.Title))
            f.WriteString(fmt.Sprintf("Author:\t%s\n", book.Author))
            if status == models.READ || status == models.READING || status == models.ALL {
              var read int
              db.Table("readings").Where("user_id = ? AND book_id = ? AND completed = ?", loggedUser.ID, book.ID, true).Count(&read)
              f.WriteString(fmt.Sprintf("Read:\t%d\n", read))
            }
            f.WriteString("--------------------------------------------------------------\n\n")
          }
        } else {
          f.WriteString("No books")
        }
        f.Sync()
        f.Close()
        // open file and read it
        dat, _ := ioutil.ReadFile(filename)
        // remove file
        os.Remove(filename)
        // send as blob
        return c.Blob(http.StatusOK, "text/plain", dat)
      } else {
        return c.JSON(http.StatusNotFound, map[string]string{"message":  "FAIL"})
      }
    } else {
      return c.JSON(http.StatusBadRequest, map[string]string{"message":  logErr.Error()})
    }
  }
}