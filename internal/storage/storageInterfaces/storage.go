package storageInterfaces

import "github.com/Fyefhqdishka/eff-mobile/internal/models"

type Storage interface {
	Create(song models.Song) (int, error)
	Update(song models.Song) (bool, error)
	Delete(song models.Song) (int, error)
	Get(groupName, songName, releaseDate string, limit, offset, songID int) ([]models.Song, error)
}
