package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fyefhqdishka/eff-mobile/internal/handlers"
	"github.com/Fyefhqdishka/eff-mobile/internal/models"
	"github.com/stretchr/testify/assert"
)

// Мок-сервис с исправленными возвращаемыми значениями
type MockService struct {
	mock.Mock
}

func (m *MockService) Create(song models.Song) (models.Song, error) {
	args := m.Called(song)
	return args.Get(0).(models.Song), args.Error(1)
}

func (m *MockService) Update(song models.Song) (bool, error) {
	args := m.Called(song)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) Delete(song models.Song) (int, error) {
	args := m.Called(song)
	return args.Int(0), args.Error(1)
}

func (m *MockService) Get(groupName, songName, releaseDate string, limit, offset, songID int) ([]models.Song, error) {
	args := m.Called(groupName, songName, releaseDate, limit, offset, songID)
	return args.Get(0).([]models.Song), args.Error(1)
}

func (m *MockService) GetVerses(groupName, songName, releaseDate string, limit, offset, songID int) ([]string, error) {
	args := m.Called(groupName, songName, releaseDate, limit, offset, songID)
	return args.Get(0).([]string), args.Error(1)
}

func TestCreateSong(t *testing.T) {
	mockService := new(MockService)

	// Мок-сервис возвращает песню с заполненным полем Song
	mockService.On("Create", mock.Anything).Return(models.Song{
		GroupName:   "GroupName",
		Song:        "SongName",
		Text:        "Verse 1\nVerse 2",
		ID:          2,
		Link:        "sting",
		ReleaseDate: "string",
	}, nil)

	h := handlers.NewHandlers(nil, mockService)

	newSong := models.Song{
		GroupName:   "GroupName",
		Song:        "SongName",
		Text:        "Verse 1\nVerse 2",
		ID:          2,
		Link:        "sting",
		ReleaseDate: "string",
	}

	songJSON, err := json.Marshal(newSong)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/songs", bytes.NewReader(songJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	h.Create(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	result := response["result"].([]interface{})
	firstSong := result[0].(map[string]interface{})
	assert.Equal(t, newSong.Song, firstSong["song"])
}

func TestGet(t *testing.T) {
	mockService := new(MockService)
	mockLog := slog.Logger{}

	handler := handlers.NewHandlers(&mockLog, mockService)

	songs := []models.Song{
		{ID: 1, GroupName: "Group1", Song: "Song1", Text: "Verse1", Link: "Link1", ReleaseDate: "2024-01-01"},
		{ID: 2, GroupName: "Group2", Song: "Song2", Text: "Verse2", Link: "Link2", ReleaseDate: "2024-02-01"},
	}

	mockService.On("Get", "Group1", "Song1", "2024-01-01", 10, 0, 0).Return(songs, nil)

	req, err := http.NewRequest("GET", "/songs?group_name=Group1&song=Song1&releasedate=2024-01-01&id=0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	result, ok := response["result"].([]interface{})
	if !ok {
		t.Fatalf("Expected result to be a slice, got %T", response["result"])
	}

	if len(result) > 0 {
		firstSong := result[0].(map[string]interface{})
		assert.Equal(t, "Song1", firstSong["song"])
	} else {
		t.Fatal("No songs found in response")
	}
}
