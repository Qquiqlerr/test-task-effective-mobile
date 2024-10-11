package mocks

import (
	"testEffectiveMobile/internal/models"
	"time"
)

type APIClientMock struct {
}

func NewAPIClientMock() *APIClientMock {
	return &APIClientMock{}
}

func (a *APIClientMock) SongEnrichment(name, group string) (*models.Song, error) {
	return &models.Song{
		Group:       group,
		Song:        name,
		ReleaseDate: time.Now().Format(time.DateOnly),
		Text:        "first verse\n\nsecond verse\n\nthird verse",
		Link:        "test",
	}, nil
}
