package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"testEffectiveMobile/internal/models"
	"testEffectiveMobile/internal/service"
	"time"
)

type songRepositoryImpl struct {
	log *slog.Logger
	DB  *gorm.DB
}

func (s *songRepositoryImpl) CreateSong(song *models.Song) (uint, error) {
	const op = "repository.songRepositoryImpl.CreateSong"
	log := s.log.With(
		slog.String("op", op),
	)

	if err := s.DB.Model(&models.Song{}).Create(song).Error; err != nil {
		log.Warn(fmt.Sprintf("failed to create song: %s", err.Error()))
		return 0, err
	}
	log.Info("Successfully created song")
	return song.ID, nil
}

// FilterSongs обращается к базе данных для получения записей с фильтрацией по group и name и пагинацией
func (s *songRepositoryImpl) FilterSongs(group, name string, offset, limit, id int) ([]models.Song, error) {
	const op = "repository.songRepositoryImpl.FilterSongs"
	log := s.log.With(
		slog.String("op", op),
	)
	var songs []models.Song
	query := s.DB.Model(&models.Song{})
	if group != "" {
		query = query.Where("\"group\" ILIKE ?", "%"+group+"%")
	}
	if name != "" {
		query = query.Where("song ILIKE ?", "%"+name+"%")
	}
	if id > 0 {
		query = query.Where("id = ?", id)
	}
	err := query.Limit(limit).Offset(offset).Find(&songs).Error
	if err != nil {
		log.Warn(fmt.Sprintf("failed to get songs: %s", err.Error()))
		return nil, err
	}
	log.Info("songs successfully received")
	return songs, nil
}

func (s *songRepositoryImpl) GetVerseByID(id int) (string, error) {
	const op = "repository.songRepositoryImpl.GetVerseByID"
	log := s.log.With(
		slog.String("op", op),
		slog.Any("song_id", id),
	)
	var UnsplittedVerse string
	query := s.DB.Model(&models.Song{}).Select("text")
	err := query.Where("id = ?", id).Scan(&UnsplittedVerse).Error
	if err != nil || UnsplittedVerse == "" {
		log.Warn("failed to get song verse", slog.String("err", err.Error()))
	}
	return UnsplittedVerse, nil
}

func (s *songRepositoryImpl) DeleteSong(id int) error {
	const op = "repository.songRepositoryImpl.DeleteSong"
	log := s.log.With(
		slog.String("op", op),
		slog.Any("song_id", id),
	)
	tx := s.DB.Model(&models.Song{}).Where("id = ?", id).Delete(&models.Song{})
	if tx.Error != nil {
		log.Warn("failed to delete song", slog.String("err", tx.Error.Error()))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		log.Debug("song not found")
		return fmt.Errorf("song not found")
	}
	log.Info("song successfully deleted")
	return nil
}

func (s *songRepositoryImpl) UpdateSong(song *models.Song) error {
	const op = "repository.songRepositoryImpl.UpdateSong"
	log := s.log.With(
		slog.String("op", op),
		slog.Any("song_id", song.ID),
	)
	if _, err := time.Parse(time.DateOnly, song.ReleaseDate); song.ReleaseDate != "" && err != nil {
		log.Debug("failed to parse date", slog.String("err", err.Error()))
		return errors.New("bad date format")
	}
	tx := s.DB.Model(&models.Song{}).Where("id = ?", song.ID).Updates(song)
	if tx.Error != nil {
		log.Warn("failed to update song", slog.String("err", tx.Error.Error()))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		log.Debug("song not found")
		return fmt.Errorf("song not found")
	}
	log.Info("song successfully updated")
	return nil
}

func NewRepository(log *slog.Logger, DB *gorm.DB) service.SongRepository {
	return &songRepositoryImpl{
		log: log,
		DB:  DB,
	}
}
