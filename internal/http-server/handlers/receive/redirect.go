package receive

import (
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"libraryService/internal/models/book"
	"net/http"

	resp "libraryService/internal/lib/api/response"
)

type BooksGetter interface {
	GetBooks() (book.Books, error)
}

func New(log *slog.Logger, getter BooksGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.receive.New"

		log = log.With(slog.String("op", op))

		books, err := getter.GetBooks()
		if books.IsEmpty() {
			log.Error("No books found in library: ", err)

			render.JSON(w, r, resp.Error("No books in library"))

			return
		}
		if err != nil {
			log.Error("Error while getting books: ", err)

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		log.Info("Got books")

		render.JSON(w, r, books)
	}

}
