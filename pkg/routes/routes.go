package routes

import (
	"encoding/json"
	_ "github.com/Fyefhqdishka/eff-mobile/docs"
	"github.com/Fyefhqdishka/eff-mobile/internal/handlers"
	"github.com/Fyefhqdishka/eff-mobile/internal/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

func RegisterRoutes(r *mux.Router, h handlers.Handlers) {
	songRoutes(r, h)
}

func songRoutes(r *mux.Router, h handlers.Handlers) {
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET")
	r.HandleFunc("/songs", h.Create).Methods("POST")
	r.HandleFunc("/songs", h.Get).Methods("GET")
	r.HandleFunc("/songs/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/songs/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/songs/verses", h.GetVerses).Methods("GET")

	r.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		var song models.Song

		song.Song = r.URL.Query().Get("song")
		song.GroupName = r.URL.Query().Get("group")
		song.Text = "first verse\n\nsecond verse\n\nthird verse\n\nfourth verse\n\n"
		song.Link = "https://www.youtube.com/watch?v=HRbW75fYLvo&t=326629s"
		song.ReleaseDate = time.Now().Format("02.01.2006")

		data, _ := json.Marshal(song)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}
