package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"testEffectiveMobile/internal/models"
)

type SongService interface {
	FilterSongs(group, name string, page, pageSize, id int) ([]models.Song, error)
	CreateSong(group, name string) (uint, error)
	GetVersesWithPagination(id, page, pageSize int) ([]string, error)
	DeleteSong(id int) error
	UpdateSong(song *models.Song) error
}

type SongController struct {
	log         *slog.Logger
	songService SongService
}

func NewController(songService SongService, log *slog.Logger) *SongController {
	return &SongController{
		songService: songService,
		log:         log,
	}
}

type Request struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}

// GetSongs godoc
// @Summary Get songs with params
// @Description Get a list of songs with filters and pagination
// @Param group query string false "Group name"
// @Param name query string false "Song name"
// @Param id query integer false "Song id"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Song
// @Failure 500 {object} models.Failures
// @Router /songs [get]
func (c *SongController) GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	name := r.URL.Query().Get("name")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 0
	}
	page, pageSize := GetPages(pageStr, pageSizeStr)

	songs, err := c.songService.FilterSongs(group, name, page, pageSize, id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}
	render.JSON(w, r, songs)
}

// CreateSong godoc
// @Summary Create song with enrichment
// @Description Create a song with enrichment
// @Accept json
// @Produce json
// @Param request body Request true "Name and group of the song"
// @Success 200 {object} models.CreateSongResponse
// @Failure 500 {object} models.Failures
// @Failure 400 {object} models.Failures
// @Router /songs [post]
func (c *SongController) CreateSong(w http.ResponseWriter, r *http.Request) {
	const op = "controller.SongController.CreateSong"
	var request Request
	log := c.log.With(
		slog.String("op", op),
	)
	err := render.DecodeJSON(r.Body, &request)
	if request.Name == "" || request.Group == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request"})
		return
	}
	if err != nil {
		log.Warn("failed to decode JSON", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request"})
		return
	}
	id, err := c.songService.CreateSong(request.Group, request.Name)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}
	render.JSON(w, r, map[string]uint{"song_id": id})
}

// GetVersesByID godoc
// @Summary Get verses by song id
// @Description Get verses by song id with pagination
// @Param id path int true "Song id"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} string "verses"
// @Failure 500 {object} models.Failures
// @Failure 400 {object} models.Failures
// @Router /songs/{id}/verses [get]
func (c *SongController) GetVersesByID(w http.ResponseWriter, r *http.Request) {
	const op = "controller.SongController.GetVersesByID"
	log := c.log.With(
		slog.String("op", op),
	)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Debug("failed to get id", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request"})
		return
	}
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, pageSize := GetPages(pageStr, pageSizeStr)
	verses, err := c.songService.GetVersesWithPagination(id, page, pageSize)
	if err != nil {
		if err.Error() == "no more verses available" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "no more verses available"})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}
	render.JSON(w, r, verses)
}

// DeleteSong godoc
// @Summary Delete song by id
// @Description Delete song by id
// @Param id path int true "Song id"
// @Success 200 {string}  string    "song deleted"
// @Failure 500 {object} models.Failures
// @Failure 400 {object} models.Failures
// @Router /songs/{id} [delete]
func (c *SongController) DeleteSong(w http.ResponseWriter, r *http.Request) {
	const op = "controller.SongController.DeleteSong"
	log := c.log.With(
		slog.String("op", op),
	)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Debug("failed to get id", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request"})
		return
	}
	err = c.songService.DeleteSong(id)
	if err != nil {
		if err.Error() == "song not found" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "song not found"})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}
	render.JSON(w, r, map[string]string{"message": "song deleted"})
}

// UpdateSong godoc
// @Summary Update song by id (partial updates allowed)
// @Description Update song by id. You can provide partial updates, for example {"text": "new text"}.
// @Accept json
// @Produce json
// @Param id path int true "Song id"
// @Param request body models.SongWithoutID false "Partial Song object" example({"text": "new text", "group": "Muse"})
// @Success 200 {string}  string    "song updated"
// @Failure 500 {object} models.Failures
// @Failure 400 {object} models.Failures
// @Router /songs/{id} [put]
func (c *SongController) UpdateSong(w http.ResponseWriter, r *http.Request) {
	const op = "controller.SongController.UpdateSong"
	log := c.log.With(
		slog.String("op", op),
	)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Debug("failed to get id", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request"})
		return
	}
	var song models.Song

	err = render.DecodeJSON(r.Body, &song)
	if err != nil {
		log.Warn("failed to decode JSON", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request"})
		return
	}
	song.ID = uint(id)
	err = c.songService.UpdateSong(&song)
	if err != nil {
		if err.Error() == "song not found" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "song not found"})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}
	render.JSON(w, r, map[string]string{"message": "song updated"})
}

func GetPages(page, pageSize string) (int, int) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		//Значение по умолчанию
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		//Значение по умолчанию
		pageSizeInt = 10
	}
	return pageInt, pageSizeInt
}
