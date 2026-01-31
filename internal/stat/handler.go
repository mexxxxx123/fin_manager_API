package stat

import (
	"fin_manager_API/m/configs"
	"fin_manager_API/m/middleware"
	"fin_manager_API/m/pkg/res"
	"net/http"
	"time"
)

const (
	SortByDay   = "day"
	SortByMonth = "month"
)

type StatHandler struct {
	*configs.Config
	StatRepo *StatRepository
}

type StatHandlerDeps struct {
	*configs.Config
	StatRepo *StatRepository
}

func NewStatHandler(mux *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		Config:   deps.Config,
		StatRepo: deps.StatRepo,
	}
	mux.Handle("GET /transactions/allstat", middleware.IsAuthed(handler.AllTimeStat(), deps.Config))
	mux.Handle("GET /transactions/stat", middleware.IsAuthed(handler.ByPeriodStat(), deps.Config))
}

func (handler StatHandler) AllTimeStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			http.Error(w, "GET_USER_ERROR", http.StatusBadRequest)
		}
		stats, err := handler.StatRepo.GetUserStat(userID)
		if err != nil {
			http.Error(w, "GET_USERSTAT_ERROR", http.StatusBadRequest)
		}
		res.Json(stats, w, 200)
	}
}
func (handler StatHandler) ByPeriodStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			http.Error(w, "GET_USER_ERROR", http.StatusBadRequest)
		}

		fromDate, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		toDate, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))

		if err != nil {
			http.Error(w, "Wrong from/to date", http.StatusBadRequest)
			return
		}

		by := r.URL.Query().Get("by")

		if by != SortByMonth && by != SortByDay {
			http.Error(w, "Wrong by (only 'month' or 'day')", http.StatusBadRequest)
			return
		}

		stats, err := handler.StatRepo.GetUserStatByPeriod(userID, fromDate, toDate, by)

		if err != nil {
			http.Error(w, "GET_USERSTAT_ERROR", http.StatusBadRequest)
			return
		}

		res.Json(stats, w, 200)
	}
}
