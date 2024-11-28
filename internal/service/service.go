package service

import (
	"fmt"
	"github.com/Fyefhqdishka/eff-mobile/internal/client"
	"github.com/Fyefhqdishka/eff-mobile/internal/models"
	"github.com/Fyefhqdishka/eff-mobile/internal/storage/storageInterfaces"
	"log/slog"
	"strings"
)

type ServiceInterface interface {
	Create(song models.Song) (models.Song, error)
	Update(song models.Song) (bool, error)
	Delete(song models.Song) (int, error)
	Get(groupName, songName, releaseDate string, limit, offset, songID int) ([]models.Song, error)
	GetVerses(groupName, songName, releaseDate string, limit, offset, songID int) ([]string, error)
}

type Service struct {
	Repo   storageInterfaces.Storage
	client client.ClientInterface
	log    *slog.Logger
}

func NewService(repo storageInterfaces.Storage, client client.ClientInterface, log *slog.Logger) *Service {
	return &Service{
		Repo:   repo,
		client: client,
		log:    log,
	}
}

func (s *Service) Create(song models.Song) (models.Song, error) {
	res, err := s.client.GetDetails(song.Song, song.GroupName)
	if err != nil {
		return models.Song{}, err
	}

	id, err := s.Repo.Create(song)
	if err != nil {
		return models.Song{}, err
	}

	res.ID = id

	return res, nil
}

func (s *Service) Update(song models.Song) (bool, error) {
	success, err := s.Repo.Update(song)
	if err != nil {
		return false, err
	}

	return success, nil
}

func (s *Service) Delete(song models.Song) (int, error) {
	id, err := s.Repo.Delete(song)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) Get(groupName, songName, releaseDate string, limit, offset, songID int) ([]models.Song, error) {
	return s.Repo.Get(groupName, songName, releaseDate, limit, offset, songID)
}

func (s *Service) GetVerses(groupName, songName, releaseDate string, limit, offset, songID int) ([]string, error) {
	s.log.Debug("Start fetching verses", "groupName", groupName, "songName", songName, "releaseDate", releaseDate, "songID", songID)

	songs, err := s.Repo.Get(groupName, songName, releaseDate, limit, offset, songID)
	if err != nil {
		return nil, err
	}

	s.log.Debug("Fetched songs", "songs", songs)

	if len(songs) == 0 {
		return nil, fmt.Errorf("song not found")
	}

	song := songs[0]

	verses := splitSongTextToVerses(song.Text)

	startIdx := offset * limit
	endIdx := startIdx + limit
	if startIdx >= len(verses) {
		return nil, fmt.Errorf("page out of range")
	}
	if endIdx > len(verses) {
		endIdx = len(verses)
	}

	return verses[startIdx:endIdx], nil
}

func splitSongTextToVerses(text string) []string {
	return strings.Split(text, "\n")
}
