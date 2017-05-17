package models

// Response sent when asked for book detail GR, GB
type GoodReadsBook struct {
  Id               string  `xml:"book>id" json:"id"`
  Title            string  `xml:"book>title" json:"title"`
  Isbn             string  `xml:"book>isbn" json:"isbn"`
  ImageUrl         string  `xml:"book>image_url"  json:"imageUrl"`
  PublicationYear  string  `xml:"book>publication_year" json:"publicationYear"`
  PublicationMonth string  `xml:"book>publication_month"  json:"publicationMonth"`
  PublicationDay   string  `xml:"book>publication_day" json:"publicationDay"`
  Publisher        string  `xml:"book>publisher" json:"publisher"`
  Description      string  `xml:"book>description" json:"description"`
  Pages            string  `xml:"book>num_pages" json:"pages"`
  OriginUrl        string  `xml:"book>link"  json:"originUrl"`
  Authors          []Author `xml:"book>authors>author"  json:"authors"`
  Series           Serie `xml:"book>series_works>series_work>series" json:"series"`
  SimilarBooks     []SimilarBook `xml:"book>similar_books>book" json:"similarBooks"`
}

type Author struct {
  Id            string  `xml:"id" json:"id"`
  Name          string `xml:"name" json:"name"`
  Role          string `xml:"role" json:"role"`
  ImageUrl      string `xml:"image_url" json:"imageUrl"`
  SmallImageUrl string `xml:"small_image_url" json:"smallImageUrl"`
  OriginUrl     string `xml:"link" json:"originUrl"`
}

type Serie struct {
  Id    string `xml:"id" json:"id"`
  Title string `xml:"title" json:"title"`
  Count string `xml:"series_works_count" json:"count"`
}

type SimilarBook struct {
  Id       int  `xml:"id" json:"id"`
  Title    string  `xml:"title" json:"title"`
  ImageUrl string  `xml:"image_url" json:"imageUrl"`
  Authors  []SimilarBookAuthor `xml:"authors>author"  json:"authors"`
}

type SimilarBookAuthor struct {
  Id        string  `xml:"id" json:"id"`
  Name      string `xml:"name" json:"name"`
}


// Search models
type GoodReadsSearchBookResponse struct {
  Books []GoodReadsSearchBook `xml:"search>results>work" json:"books"`
}

type GoodReadsSearchBook struct {
  Id       int `xml:"best_book>id" json:"id"`
  Title    string `xml:"best_book>title" json:"title"`
  Author   string `xml:"best_book>author>name" json:"author"`
  ImageUrl string `xml:"best_book>image_url" json:"imageUrl"`
}

type GoodReadsBookSearchAuthor struct {
  Id   string `xml:"id" json:"id"`
  Name string `xml:"name" json:"name"`
}


