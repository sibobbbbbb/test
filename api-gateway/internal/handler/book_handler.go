package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/sibobbbbbb/backend-engineer-challenge/api-gateway/internal/client"
    "github.com/sibobbbbbb/backend-engineer-challenge/api-gateway/internal/middleware"
    "github.com/sibobbbbbb/backend-engineer-challenge/api-gateway/pkg/response"
)

type CreateBookRequest struct {
    Title         string   `json:"title"`
    Author        string   `json:"author"`
    ISBN          string   `json:"isbn"`
    PublishedYear int      `json:"published_year"`
    CategoryIDs   []string `json:"category_ids"`
    Stock         int      `json:"stock"`
}

type UpdateBookRequest struct {
    Title         string   `json:"title"`
    Author        string   `json:"author"`
    ISBN          string   `json:"isbn"`
    PublishedYear int      `json:"published_year"`
    CategoryIDs   []string `json:"category_ids"`
    Stock         int      `json:"stock"`
}

type BookHandler struct {
    bookClient     client.BookClient
    categoryClient client.CategoryClient
    authMiddleware *middleware.AuthMiddleware
}

func NewBookHandler(r *mux.Router, bc client.BookClient, cc client.CategoryClient, am *middleware.AuthMiddleware) {
    handler := &BookHandler{
        bookClient:     bc,
        categoryClient: cc,
        authMiddleware: am,
    }

    // Public routes
    r.HandleFunc("/books", handler.ListBooks).Methods("GET")
    r.HandleFunc("/books/{id}", handler.GetBook).Methods("GET")
    r.HandleFunc("/books/search", handler.SearchBooks).Methods("GET")

    // Protected routes (require authentication)
    protected := r.PathPrefix("/books").Subrouter()
    protected.Use(am.Authenticate)
    protected.HandleFunc("", handler.CreateBook).Methods("POST")
    protected.HandleFunc("/{id}", handler.UpdateBook).Methods("PUT")
    protected.HandleFunc("/{id}", handler.DeleteBook).Methods("DELETE")
}

// CreateBook menangani permintaan POST untuk membuat buku baru
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
    var req CreateBookRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    book, err := h.bookClient.CreateBook(
        ctx,
        req.Title,
        req.Author,
        req.ISBN,
        req.PublishedYear,
        req.CategoryIDs,
        req.Stock,
    )
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusCreated, book)
}

// GetBook menangani permintaan GET untuk mendapatkan detail buku
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    book, err := h.bookClient.GetBook(ctx, id)
    if err != nil {
        response.Error(w, http.StatusNotFound, "Book not found")
        return
    }

    response.JSON(w, http.StatusOK, book)
}

// ListBooks menangani permintaan GET untuk mendapatkan daftar buku
func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
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
    books, total, err := h.bookClient.ListBooks(ctx, page, limit)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]interface{}{
        "data":  books,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// UpdateBook menangani permintaan PUT untuk memperbarui buku
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var req UpdateBookRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    book, err := h.bookClient.UpdateBook(
        ctx,
        id,
        req.Title,
        req.Author,
        req.ISBN,
        req.PublishedYear,
        req.CategoryIDs,
        req.Stock,
    )
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, book)
}

// DeleteBook menangani permintaan DELETE untuk menghapus buku
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    success, err := h.bookClient.DeleteBook(ctx, id)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]bool{"success": success})
}

// SearchBooks menangani permintaan GET untuk mencari buku
func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()
    q := queryValues.Get("q")
    pageStr := queryValues.Get("page")
    limitStr := queryValues.Get("limit")

    if q == "" {
        response.Error(w, http.StatusBadRequest, "Query parameter 'q' is required")
        return
    }

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
    books, total, err := h.bookClient.SearchBooks(ctx, q, page, limit)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]interface{}{
        "data":  books,
        "total": total,
        "page":  page,
        "limit": limit,
        "query": q,
    })
}