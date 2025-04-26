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

type CreateCategoryRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type UpdateCategoryRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type CategoryHandler struct {
    categoryClient client.CategoryClient
    authMiddleware *middleware.AuthMiddleware
}

func NewCategoryHandler(r *mux.Router, cc client.CategoryClient, am *middleware.AuthMiddleware) {
    handler := &CategoryHandler{
        categoryClient: cc,
        authMiddleware: am,
    }

    // Public routes
    r.HandleFunc("/categories", handler.ListCategories).Methods("GET")
    r.HandleFunc("/categories/{id}", handler.GetCategory).Methods("GET")

    // Protected routes (require authentication)
    protected := r.PathPrefix("/categories").Subrouter()
    protected.Use(am.Authenticate)
    protected.HandleFunc("", handler.CreateCategory).Methods("POST")
    protected.HandleFunc("/{id}", handler.UpdateCategory).Methods("PUT")
    protected.HandleFunc("/{id}", handler.DeleteCategory).Methods("DELETE")
}

// CreateCategory menangani permintaan POST untuk membuat kategori baru
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
    var req CreateCategoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    category, err := h.categoryClient.CreateCategory(ctx, req.Name, req.Description)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusCreated, category)
}

// GetCategory menangani permintaan GET untuk mendapatkan detail kategori
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    category, err := h.categoryClient.GetCategory(ctx, id)
    if err != nil {
        response.Error(w, http.StatusNotFound, "Category not found")
        return
    }

    response.JSON(w, http.StatusOK, category)
}

// ListCategories menangani permintaan GET untuk mendapatkan daftar kategori
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
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
    categories, total, err := h.categoryClient.ListCategories(ctx, page, limit)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]interface{}{
        "data":  categories,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// UpdateCategory menangani permintaan PUT untuk memperbarui kategori
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var req UpdateCategoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    category, err := h.categoryClient.UpdateCategory(ctx, id, req.Name, req.Description)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, category)
}

// DeleteCategory menangani permintaan DELETE untuk menghapus kategori
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    success, err := h.categoryClient.DeleteCategory(ctx, id)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]bool{"success": success})
}