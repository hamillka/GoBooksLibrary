package save

import (
	"errors"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"io"
	resp "libraryService/internal/lib/api/response"
	"libraryService/internal/models/book"
	"libraryService/internal/storage"
	"net/http"
)

type Request struct {
	Book book.Book `json:"book"`
}

type Response struct {
	resp.Response
	Book book.Book `json:"book"`
}

type BookSaver interface {
	AddBook(book book.Book) (int64, error)
}

func New(log *slog.Logger, bookSaver BookSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		log = log.With(slog.String("op", op))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		id, err := bookSaver.AddBook(req.Book)
		if errors.Is(err, storage.ErrBookExists) {
			log.Info("Book already exists: ", slog.AnyValue(req.Book))

			render.JSON(w, r, resp.Error("book already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add book: ", err)

			render.JSON(w, r, resp.Error("failed to add book"))

			return
		}

		log.Info("Book added", slog.Int64("id", id))

		responseOK(w, r, req.Book)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, book book.Book) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Book:     book,
	})
}
