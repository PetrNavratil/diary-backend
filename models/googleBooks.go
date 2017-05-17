package models

// Sent when asked for Google book detail
type GoogleBook struct {
  Title       string `json:"title"`
  Author      string `json:"author"`
  Publisher   string `json:"publisher"`
  Published   string `json:"published"`
  Description string `json:"description"`
  PageCount   int `json:"pageCount"`
  ImageUrl    string `json:"imageUrl"`
  Preview     string `json:"preview"`
}

type GoogleBookResponse struct {
  VolumeInfo struct {
               Title       string `json:"title"`
               Authors     []string `json:"authors"`
               Publisher   string `json:"publisher"`
               Published   string `json:"publishedDate"`
               Description string `json:"description"`
               PageCount   int `json:"pageCount"`
               Images      struct {
                             ImageUrl string `json:"thumbnail"`
                           } `json:"imageLinks"`
               Preview     string `json:"previewLink"`
             } `json:"volumeInfo"`
}