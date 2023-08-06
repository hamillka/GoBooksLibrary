package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"libraryService/internal/models/book"
	"libraryService/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS library(
	    id INTEGER PRIMARY KEY,
	    book TEXT NOT NULL UNIQUE,
	    author TEXT NOT NULL)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddBook(book book.Book) (int64, error) {
	const op = "storage.sqlite.AddBook"

	stmt, err := s.db.Prepare("INSERT INTO library(book, author) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(book.Name, book.Author)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrBookExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetBooks() (book.Books, error) {
	const op = "storage.sqlite.GetBooks"
	var res book.Books

	rows, err := s.db.Query("SELECT * FROM library")
	if err != nil {
		return res, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		i := book.Book{}
		var tmp int64
		err = rows.Scan(&tmp, &i.Name, &i.Author)
		if err != nil {
			return res, fmt.Errorf("%s: prepare statement: %w", op, err)
		}
		res = append(res, i)
	}

	return res, nil
}
