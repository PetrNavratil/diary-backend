package models

// Educational stored in database
// Represents literary analysis
type Educational struct {
  ID         int `json:"id"`
  UserBookID int `json:"-"`
  Druh       string `json:"druh"`
  Zanr       string `json:"zanr"`
  Smer       string `json:"smer"`
  Forma      string `json:"forma"`
  Jazyk      string `json:"jazyk"`
  Postavy    string `json:"postavy"`
  Obsah      string `json:"obsah"`
  Tema       string `json:"tema"`
  Hodnoceni  string `json:"hodnoceni"`
}