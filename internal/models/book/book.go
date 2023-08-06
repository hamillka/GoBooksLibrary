package book

type Book struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type Books []Book

func (books Books) IsEmpty() bool {
	return len(books) == 0
}
