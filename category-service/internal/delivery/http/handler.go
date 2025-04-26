package http

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/sibobbbbbb/backend-engineer-challenge/category-service/internal/domain"
)

type ResponseError struct {
    Message string `json:"message"`
}

type CategoryHandler struct {
    CategoryUsecase domain.CategoryUsecase
}

func NewCategoryHandler(r *mux.Router, cu domain.CategoryUsecase) {
    handler := &CategoryHandler{
        CategoryUsecase: cu,
    }

    // Register routes
    r.HandleFunc("/categories", handler.Create).Methods("POST")
    r.HandleFunc("/categories", handler.List).Methods("GET")
    r.HandleFunc("/categories/{id}", handler.GetByID).Methods("GET")
    r.HandleFunc("/categories/{id}", handler.Update).Methods("PUT")
    r.HandleFunc("/categories/{id}", handler.Delete).Methods("DELETE")
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
    var category domain.Category
    err := json.NewDecoder(r.Body).Decode(&category)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    err = h.CategoryUsecase.Create(ctx, &category)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    category, err := h.CategoryUsecase.GetByID(ctx, id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Category not found")
        return
    }

    respondWithJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var category domain.Category
    err := json.NewDecoder(r.Body).Decode(&category)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    category.ID = id
    ctx := r.Context()
    err = h.CategoryUsecase.Update(ctx, &category)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    updatedCategory, err := h.CategoryUsecase.GetByID(ctx, id)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Category updated but failed to retrieve")
        return
    }

    respondWithJSON(w, http.StatusOK, updatedCategory)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    err := h.CategoryUsecase.Delete(ctx, id)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
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
    categories, total, err := h.CategoryUsecase.List(ctx, page, limit)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]interface{}{
        "data":  categories,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, ResponseError{Message: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}