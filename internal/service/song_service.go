package service

import (
	"errors"
	"log/slog"
	"strings"
	"testEffectiveMobile/internal/controller"
	"testEffectiveMobile/internal/models"
)

type SongRepository interface {
	FilterSongs(group, name string, offset, limit, id int) ([]models.Song, error)
	CreateSong(song *models.Song) (uint, error)
	GetVerseByID(id int) (string, error)
	DeleteSong(id int) error
	UpdateSong(song *models.Song) error
}
type APIClient interface {
	SongEnrichment(name, group string) (*models.Song, error)
}

type SongService struct {
	log            *slog.Logger
	songRepository SongRepository
	APIClient      APIClient
}

func NewService(songRepository SongRepository, log *slog.Logger, client APIClient) controller.SongService {
	return &SongService{
		songRepository: songRepository,
		log:            log,
		APIClient:      client,
	}
}

// FilterSongs считает смещение и передает запрос на уровень репозитория
func (s *SongService) FilterSongs(group, name string, page, pageSize, id int) ([]models.Song, error) {
	offset := (page - 1) * pageSize
	return s.songRepository.FilterSongs(group, name, offset, pageSize, id)
}

func (s *SongService) GetVersesWithPagination(id, page, pageSize int) ([]string, error) {
	const op = "service.SongService.GetVerseWithPagination"
	log := s.log.With(
		slog.String("op", op),
		slog.Any("song_id", id),
		slog.Any("page", page),
		slog.Any("page_size", pageSize),
	)
	verse, err := s.songRepository.GetVerseByID(id)
	if err != nil {
		return nil, err
	}
	verses := strings.Split(verse, "\n\n")
	totalVerses := len(verses)
	start := (page - 1) * pageSize
	if start >= totalVerses {
		log.Debug("no more verses available")
		return nil, errors.New("no more verses available")
	}
	end := start + pageSize
	if end > totalVerses {
		end = totalVerses
	}
	return verses[start:end], nil
}

// CreateSong создает запись о песне
func (s *SongService) CreateSong(group, name string) (uint, error) {
	song, err := s.APIClient.SongEnrichment(name, group)
	if err != nil {
		return 0, err
	}
	id, err := s.songRepository.CreateSong(song)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SongService) UpdateSong(song *models.Song) error {
	return s.songRepository.UpdateSong(song)
}
func (s *SongService) DeleteSong(id int) error {
	return s.songRepository.DeleteSong(id)
}
