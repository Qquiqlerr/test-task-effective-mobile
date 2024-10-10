package repository

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"testEffectiveMobile/internal/models"
	"testEffectiveMobile/internal/service"
	"time"
)

var ClientTimeout time.Duration = 10

type APIClientImpl struct {
	log *slog.Logger
	URL string
}

func NewAPIClient(url string) service.APIClient {
	return &APIClientImpl{
		URL: url,
	}
}

func (a *APIClientImpl) SongEnrichment(name, group string) (*models.Song, error) {
	const op = "repository.APIClientImpl.SongEnrichment"
	log := a.log.With(
		slog.String("op", op),
	)

	url := fmt.Sprintf("%s/info?group=%s&song=%s", a.URL, group, name)
	client := &http.Client{Timeout: ClientTimeout * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Warn("Can't connect to API", slog.String("err", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("failed to get song info:", slog.String("StatusCode", resp.Status))
		return nil, fmt.Errorf("failed to get song info: %s", resp.Status)
	}
	var song models.Song
	if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
		log.Warn("failed to decode JSON", slog.String("err", err.Error()))
		return nil, err
	}
	song.Song = name
	song.Group = group
	return &song, nil
}
