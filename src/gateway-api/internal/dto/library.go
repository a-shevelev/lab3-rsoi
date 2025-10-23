package dto

type LibraryResponse struct {
	LibraryUid string `json:"libraryUid"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	City       string `json:"city"`
}

type LibraryPaginationResponse struct {
	Page          int               `json:"page"`
	PageSize      int               `json:"pageSize"`
	TotalElements int               `json:"totalElements"`
	Items         []LibraryResponse `json:"items"`
}

type BookResponse struct {
	BookUid        string `json:"bookUid"`
	Name           string `json:"name"`
	Author         string `json:"author"`
	Genre          string `json:"genre"`
	Condition      string `json:"condition"`
	AvailableCount int    `json:"availableCount"`
}

type BookResponseRaw struct {
	BookUid string `json:"bookUid"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	Genre   string `json:"genre"`
}

type LibraryBookPaginationResponse struct {
	Page          int            `json:"page"`
	PageSize      int            `json:"pageSize"`
	TotalElements int            `json:"totalElements"`
	Items         []BookResponse `json:"items"`
}

func BookToRaw(book BookResponse) BookResponseRaw {
	return BookResponseRaw{
		BookUid: book.BookUid,
		Name:    book.Name,
		Author:  book.Author,
		Genre:   book.Genre,
	}
}
