package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

type BookHandler struct {
	BookUsecase domain.BookUsecase
}

func NewBookHandler(r *mux.Router, bookUsecase domain.BookUsecase) {
	handler := &BookHandler{
		BookUsecase: bookUsecase,
	}

	// Register routes
	r.HandleFunc("/books", handler.Create).Methods("POST")
	r.HandleFunc("/books", handler.List).Methods("GET")
	r.HandleFunc("/books/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/books/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/books/{id}", handler.Delete).Methods("DELETE")
	r.HandleFunc("/books/search", handler.Search).Methods("GET")
}

// Create menangani permintaan POST untuk membuat buku baru
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	err = h.BookUsecase.Create(ctx, &book)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, book)
}

// GetByID menangani permintaan GET untuk mendapatkan buku berdasarkan ID
func (h *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	book, err := h.BookUsecase.GetByID(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Book not found")
		return
	}

	respondWithJSON(w, http.StatusOK, book)
}

// Update menangani permintaan PUT untuk memperbarui buku
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var book domain.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	book.ID = id
	ctx := r.Context()
	err = h.BookUsecase.Update(ctx, &book)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedBook, err := h.BookUsecase.GetByID(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Book updated but failed to retrieve")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedBook)
}

// Delete menangani permintaan DELETE untuk menghapus buku
func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	err := h.BookUsecase.Delete(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// List menangani permintaan GET untuk daftar buku dengan pagination
func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	pageStr := queryValues.Get("page")
	limitStr := queryValues.Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	ctx := r.Context()
	books, total, err := h.BookUsecase.List(ctx, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":  books,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Search menangani permintaan GET untuk pencarian buku
func (h *BookHandler) Search(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	query := queryValues.Get("q")
	pageStr := queryValues.Get("page")
	limitStr := queryValues.Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	ctx := r.Context()
	books, total, err := h.BookUsecase.Search(ctx, query, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":  books,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// respondWithError adalah helper function untuk mengirim response error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ResponseError{Message: message})
}

// respondWithJSON adalah helper function untuk mengirim response JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}