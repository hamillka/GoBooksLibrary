package book

//type Quality string

//const (
//	Good   Quality = "Good"
//	Medium Quality = "Medium"
//	Bad    Quality = "Bad"
//)

type Book struct {
	Name   string
	Author string
	//condition Quality
}

type Books []Book

func (books Books) IsEmpty() bool {
	return len(books) == 0
}
