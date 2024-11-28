package service_test

import (
	"github.com/Fyefhqdishka/eff-mobile/internal/models"
	"github.com/Fyefhqdishka/eff-mobile/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(song models.Song) (int, error) {
	args := m.Called(song)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) Update(song models.Song) (bool, error) {
	args := m.Called(song)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) Delete(song models.Song) (int, error) {
	args := m.Called(song)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) Get(groupName, songName, releaseDate string, limit, offset, songID int) ([]models.Song, error) {
	args := m.Called(groupName, songName, releaseDate, limit, offset, songID)
	return args.Get(0).([]models.Song), args.Error(1)
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetDetails(songName, groupName string) (models.Song, error) {
	args := m.Called(songName, groupName)
	return args.Get(0).(models.Song), args.Error(1)
}

func TestCreateSong(t *testing.T) {
	mockRepo := new(MockRepo)
	mockClient := new(MockClient)
	mockLog := slog.Logger{}

	service := service.NewService(mockRepo, mockClient, &mockLog)

	song := models.Song{
		Song:        "SongName",
		GroupName:   "GroupName",
		Text:        "Some song text",
		Link:        "some-link",
		ReleaseDate: "2024-01-01",
	}

	mockClient.On("GetDetails", song.Song, song.GroupName).Return(song, nil)
	mockRepo.On("Create", song).Return(1, nil)

	createdSong, err := service.Create(song)
	assert.Nil(t, err)
	assert.Equal(t, 1, createdSong.ID)
	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSong(t *testing.T) {
	mockRepo := new(MockRepo)
	mockClient := new(MockClient)
	mockLog := slog.Logger{}

	service := service.NewService(mockRepo, mockClient, &mockLog)

	song := models.Song{
		ID:          1,
		Song:        "SongName",
		GroupName:   "GroupName",
		Text:        "Some song text",
		Link:        "some-link",
		ReleaseDate: "2024-01-01",
	}

	mockRepo.On("Update", song).Return(true, nil)

	updated, err := service.Update(song)
	assert.Nil(t, err)
	assert.True(t, updated)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSong(t *testing.T) {
	mockRepo := new(MockRepo)
	mockClient := new(MockClient)
	mockLog := slog.Logger{}

	service := service.NewService(mockRepo, mockClient, &mockLog)

	song := models.Song{
		ID:          1,
		Song:        "SongName",
		GroupName:   "GroupName",
		Text:        "Some song text",
		Link:        "some-link",
		ReleaseDate: "2024-01-01",
	}

	mockRepo.On("Delete", song).Return(1, nil)

	deletedID, err := service.Delete(song)
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedID)
	mockRepo.AssertExpectations(t)
}

func TestGetSongs(t *testing.T) {
	mockRepo := new(MockRepo)
	mockClient := new(MockClient)
	mockLog := slog.Logger{}

	service := service.NewService(mockRepo, mockClient, &mockLog)

	song := models.Song{
		ID:          1,
		Song:        "SongName",
		GroupName:   "GroupName",
		Text:        "Some song text",
		Link:        "some-link",
		ReleaseDate: "2024-01-01",
	}

	mockRepo.On("Get", "GroupName", "SongName", "2024-01-01", 10, 0, 0).Return([]models.Song{song}, nil)

	songs, err := service.Get("GroupName", "SongName", "2024-01-01", 10, 0, 0)
	assert.Nil(t, err)
	assert.Len(t, songs, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetVerses(t *testing.T) {
	mockRepo := new(MockRepo)
	mockClient := new(MockClient)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mockLog := slog.New(handler)

	service := service.NewService(mockRepo, mockClient, mockLog)

	song := models.Song{
		ID:          1,
		Song:        "SongName",
		GroupName:   "GroupName",
		Text:        "Verse 1\nVerse 2\nVerse 3",
		Link:        "some-link",
		ReleaseDate: "2024-01-01",
	}

	mockRepo.On("Get", "GroupName", "SongName", "2024-01-01", 1, 0, 0).Return([]models.Song{song}, nil)

	verses, err := service.GetVerses("GroupName", "SongName", "2024-01-01", 1, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, []string{"Verse 1"}, verses)
	mockRepo.AssertExpectations(t)
}
