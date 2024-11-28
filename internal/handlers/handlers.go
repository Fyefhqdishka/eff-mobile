package handlers

import (
	"encoding/json"
	"github.com/Fyefhqdishka/eff-mobile/internal/models"
	"github.com/Fyefhqdishka/eff-mobile/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

type Handlers struct {
	log     *slog.Logger
	Service service.ServiceInterface
}

func NewHandlers(log *slog.Logger, service service.ServiceInterface) *Handlers {
	return &Handlers{
		log:     log,
		Service: service,
	}
}

// @Summary Create a new song
// @Description Creates a new song with the given details
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Song details"
// @Success 200 {object} models.Song "Created song"
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /songs [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.response(w, SendError("Can't decode json body"), http.StatusBadRequest)
		return
	}

	song, err := h.Service.Create(song)
	if err != nil {
		h.response(w, SendError("Can't create song"), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess([]models.Song{song}), http.StatusCreated)
}

// @Summary Update a song
// @Description Updated a song with the given details
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Song details"
// @Success 200 {object} models.Song "Updates song"
// @Failure 400 {object} Response "Invalid input"
// @Failure 500 {object} Response "Failed to update song"
// @Router /songs/{id} [put]
func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.response(w, SendError("Can't decode json body"), http.StatusBadRequest)
		return
	}

	success, err := h.Service.Update(song)
	if err != nil {
		h.response(w, SendError("can't update song"), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(success), http.StatusOK)
}

// @Summary Delete a song
// @Description Deleted a song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Song details"
// @Success 200 {object} models.Song "Updates song"
// @Failure 400 {object} Response "Invalid input"
// @Failure 500 {object} Response "Failed to update song"
// @Router /songs/{id} [delete]
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.response(w, SendError("Can't decode json body"), http.StatusBadRequest)
		return
	}

	id, err := h.Service.Delete(song)
	if err != nil {
		h.response(w, SendError("can't delete this song"), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(id), http.StatusOK)
}

// Get returns a list of Song's
// @Summary Get all Song's from the storage
// @Description Returns a list of all songs with optional filtering and pagination
// @Tags songs
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param id query int false "Song Id"
// @Param song query string false "Song title"
// @Param group_name query string false "Group name"
// @Param releasedate query string false "Song release date in format 02.01.2006"
// @Success 200 {array} models.Song "Array of Song's"
// @Failure 400 {object} Response "Invalid query parameters"
// @Failure 500 {object} Response "Failed to get Song's"
// @Router /songs [get]
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("group_name")
	songName := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("releasedate")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	songIDStr := r.URL.Query().Get("id")
	var songID int
	if songIDStr != "" {
		var err error
		songID, err = strconv.Atoi(songIDStr)
		if err != nil {
			http.Error(w, "invalid songID", http.StatusBadRequest)
			return
		}
	}

	limitInt := 10
	offsetInt := 0

	if limit != "" {
		var err error
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			h.response(w, SendError("Invalid limit parameter"), http.StatusBadRequest)
			return
		}
	}
	if offset != "" {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			h.response(w, SendError("Invalid offset parameter"), http.StatusBadRequest)
			return
		}
	}

	songs, err := h.Service.Get(groupName, songName, releaseDate, limitInt, offsetInt, songID)
	if err != nil {
		h.response(w, SendError("can't get all songs"), http.StatusInternalServerError)
		return
	}

	h.response(w, SendSuccess(songs), http.StatusOK)
}

// GetVerses returns the paginated song text (verses)
// @Summary Get paginated song text (verses) from the storage
// @Description Returns the verses of a song with optional filtering by group name and song name, and pagination
// @Tags songs
// @Produce  json
// @Param id query int false "Song Id"
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of verses per page" default(5)
// @Param song query string false "Song title"
// @Param group_name query string false "Group name"
// @Param releasedate query string false "Song release date in format 02.01.2006"
// @Success 200 {array} string "Array of song verses"
// @Failure 400 {object} Response "Invalid query parameters"
// @Failure 500 {object} Response "Failed to get song's verses"
// @Router /songs/verses [get]
func (h *Handlers) GetVerses(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("group_name")
	songName := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("releasedate")
	songIDStr := r.URL.Query().Get("id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "invalid songID", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
		h.log.Warn("Invalid page number, using default value", slog.Int("page", page))
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 5
		h.log.Warn("Invalid page size, using default value", slog.Int("pageSize", pageSize))
	}

	if page <= 0 || pageSize <= 0 {
		h.log.Warn("Invalid pagination parameters", "page", page, "pageSize", pageSize)
		h.response(w, SendError("Invalid pagination parameters"), http.StatusBadRequest)
		return
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	h.log.Debug("Request parameters", "page", page, "pageSize", pageSize, "offset", offset)

	verses, err := h.Service.GetVerses(groupName, songName, releaseDate, pageSize, offset, songID)
	if err != nil {
		h.log.Error("Error fetching paginated song text", slog.String("error", err.Error()))
		h.response(w, SendError("Error fetching paginated song text"), http.StatusInternalServerError)
		return
	}

	h.log.Debug("song name and groupname", songName, groupName)

	h.response(w, SendSuccess(verses), http.StatusOK)
}
