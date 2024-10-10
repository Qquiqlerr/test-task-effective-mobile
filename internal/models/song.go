package models

import (
	"gorm.io/datatypes"
)

type CreateSongResponse struct {
	SongID uint `json:"song_id"`
}

// Song описывает структуру данных песни
// @Description песня
// @Property id{integer} идентификатор песни
// @Property group{string} группа
// @Property song{string} название песни
// @Property releaseDate{string} дата релиза
// @Property text{string} текст песни
// @Property link{string} ссылка на песню
type Song struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Group       string         `json:"group" gorm:"column:group"`
	Song        string         `json:"song" gorm:"column:song"`
	ReleaseDate datatypes.Date `json:"releaseDate" gorm:"column:release_date"`
	Text        string         `json:"text" gorm:"column:text"`
	Link        string         `json:"link" gorm:"column:link"`
}

type SongWithoutID struct {
	Group       string         `json:"group"`
	Song        string         `json:"song"`
	ReleaseDate datatypes.Date `json:"releaseDate"`
	Text        string         `json:"text"`
	Link        string         `json:"link"`
}

type Failures struct {
	Error string
}
