package goodreads

// Get specific book models
type GoodReadsBook struct {
  Id               string  `xml:"book>id" json:"id"`
  Title            string  `xml:"book>title" json:"title"`
  Isbn             string  `xml:"book>isbn" json:"isbn"`
  Isbn13           string  `xml:"book>isbn13" json:"isbn13"`
  Asin             string  `xml:"book>asin"  json:"asin"`
  CountryCode      string  `xml:"book>country_code"  json:"countryCode"`
  ImageUrl         string  `xml:"book>image_url"  json:"imageUrl"`
  SmallImageUrl    string  `xml:"book>small_image_url"  json:"smallImageUrl"`
  PublicationYear  string  `xml:"book>publication_year" json:"publicationYear"`
  PublicationMonth string  `xml:"book>publication_month"  json:"publicationMonth"`
  PublicationDay   string  `xml:"book>publication_day" json:"publicationDay"`
  Publisher        string  `xml:"book>publisher" json:"publisher"`
  LanguageCode     string  `xml:"book>language_code" json:"languageCode"`
  IsEbook          string  `xml:"book>is_ebook" json:"isEbook"`
  Description      string  `xml:"book>description" json:"description"`
  AverageRating    string  `xml:"book>average_rating" json:"averageRating"`
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
  Id               string  `xml:"id" json:"id"`
  Title            string  `xml:"title" json:"title"`
  ImageUrl         string  `xml:"image_url" json:"imageUrl"`
  SmallImageUrl    string  `xml:"small_image_url"  json:"smallImageUrl"`
  PublicationYear  string  `xml:"publication_year" json:"publicationYear"`
  PublicationMonth string  `xml:"publication_month" json:"publicationMonth"`
  PublicationDay   string  `xml:"publication_day" json:"publicationDay"`
  Pages            string  `xml:"num_pages" json:"pages"`
  OriginUrl        string  `xml:"link" json:"originalUrl"`
  Authors          []SimilarBookAuthor `xml:"authors>author"  json:"authors"`
}

type SimilarBookAuthor struct {
  Id        string  `xml:"id" json:"id"`
  Name      string `xml:"name" json:"name"`
  OriginUrl string `xml:"link" json:"originalUrl"`
}


// Search models
type GoodReadsSearchBookResponse struct {
  Books []GoodReadsSearchBook `xml:"search>results>work" json:"books"`
}

type GoodReadsSearchBook struct {
  Id               string `xml:"best_book>id" json:"id"`
  Title            string `xml:"best_book>title" json:"title"`
  Author           GoodReadsBookSearchAuthor `xml:"best_book>author" json:"author"`
  ImageUrl         string `xml:"best_book>image_url" json:"imageUrl"`
  SmallImageUrl    string `xml:"best_book>small_image_url" json:"smallImageUrl"`
  PublicationYear  string `xml:"original_publication_year" json:"originalPublicationYear"`
  PublicationMonth string `xml:"original_publication_month" json:"originalPublicationMonth"`
  PublicationDay   string `xml:"original_publication_day" json:"originalPublicationDay"`
  AverageRating    string `xml:"average_rating" json:"averageRating"`
}

type GoodReadsBookSearchAuthor struct {
  Id   string `xml:"id" json:"id"`
  Name string `xml:"name" json:"name"`
}


