package transactions

import (
	"fin_manager_API/m/configs"
	"fin_manager_API/m/middleware"
	"fin_manager_API/m/pkg/req"
	"fin_manager_API/m/pkg/res"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	*configs.Config
	TransactionRepository *TransactionRepository
	*TransactionService
}

type TransactionHandlerDeps struct {
	*configs.Config
	TransactionRepository *TransactionRepository
	*TransactionService
}

func NewTransactionHandler(mux *http.ServeMux, deps TransactionHandlerDeps) {
	handler := &TransactionHandler{
		Config:                deps.Config,
		TransactionRepository: deps.TransactionRepository,
		TransactionService:    deps.TransactionService,
	}
	// Ввод любой транзакции (пополнение или трата)
	mux.Handle("POST /transactions", middleware.IsAuthed(handler.Transaction(), deps.Config))
	// Получение всех транзакций
	mux.Handle("GET /transactions", middleware.IsAuthed(handler.GetAll(), deps.Config))
	// Получение конкретной транзакций
	mux.Handle("GET /transactions/{id}", middleware.IsAuthed(handler.GetById(), deps.Config))
	// Удаление конкретной транзакций
	mux.Handle("DELETE /transactions/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
}

func (handler TransactionHandler) Transaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[TransactionRequest](&w, r)
		if err != nil {
			return
		}
		userId, ok := r.Context().Value(middleware.ContextUserIdKey).(uint)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		trans := NewTransaction(body.Amount, body.Type, body.Categorie, userId)
		createdTrans, err := handler.TransactionRepository.Create(trans)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(createdTrans, w, 201)
	}
}

func (handler TransactionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		transns, count, errS := handler.TransactionService.GetAllTransaction(r.Context(), limit, offset)
		if errS != "" {
			http.Error(w, errS, http.StatusBadRequest)
			return
		}

		res.Json(GetAllResponse{
			Transns: transns,
			Count:   count,
		}, w, 200)

	}
}

func (handler TransactionHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TransId, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		corrTrans, errS := handler.TransactionService.GetTransaction(r.Context(), uint(TransId))
		if errS != "" {
			http.Error(w, errS, http.StatusBadRequest)
			return
		}
		res.Json(corrTrans, w, 200)
	}
}

func (handler TransactionHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 0)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		errS := handler.TransactionService.DeleteTransByID(r.Context(), id)
		if errS != "" {
			http.Error(w, errS, http.StatusBadRequest)
			return
		}
		res.Json("Delete complete", w, 200)
	}
}
