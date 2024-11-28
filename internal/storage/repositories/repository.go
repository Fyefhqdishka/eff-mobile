package repositories

import (
	"database/sql"
	"fmt"
	"github.com/Fyefhqdishka/eff-mobile/internal/models"
	"log/slog"
)

type SongRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewSongRepository(db *sql.DB, log *slog.Logger) *SongRepository {
	return &SongRepository{
		db:  db,
		log: log,
	}
}

func (r *SongRepository) Create(song models.Song) (int, error) {
	r.log.Debug("Starting to create a song", slog.String("song", song.Song), slog.String("group_name", song.GroupName))

	createGroup := `INSERT INTO groups (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`
	_, err := r.db.Exec(createGroup, song.GroupName)
	if err != nil {
		r.log.Error("Failed to ensure group existence", slog.String("group_name", song.GroupName), slog.Any("error", err))
		return 0, fmt.Errorf("failed to ensure group %s existence, err=%v", song.GroupName, err)
	}

	stmt := `INSERT INTO songs (group_name, song, text, link, releaseDate) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var songID int
	err = r.db.QueryRow(stmt, song.GroupName, song.Song, song.Text, song.Link, song.ReleaseDate).Scan(&songID)
	if err != nil {
		r.log.Error("Failed to insert song into database",
			slog.String("song", song.Song),
			slog.String("group_name", song.GroupName),
			slog.Any("error", err))
		return 0, fmt.Errorf("can't insert into db, err=%v", err)
	}

	r.log.Debug("Song successfully created", slog.Int("song_id", songID), slog.String("song", song.Song))

	return songID, nil
}

func (r *SongRepository) Update(song models.Song) (bool, error) {
	r.log.Debug("Starting to update a song", slog.String("song", song.Song), slog.String("group_name", song.GroupName))

	updateGroupNameStmt := `UPDATE groups SET name = $1`
	_, err := r.db.Exec(updateGroupNameStmt, song.GroupName)
	if err != nil {
		r.log.Error("failed to update group name, err", err)
		return false, fmt.Errorf("failed to update group name, err=%v", err)
	}

	stmt := `UPDATE songs SET song = $1, group_name = $2, text = $3, link = $4, releasedate = $5 WHERE id = $6`
	res, err := r.db.Exec(stmt, song.Song, song.GroupName, song.Text, song.Link, song.ReleaseDate, song.ID)
	if err != nil {
		r.log.Error("can't fetch rows affected, err", err)
		return false, fmt.Errorf("can't fetch rows affected, err=%v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.log.Error("failed to fetch rows affected, err", err)
		return false, fmt.Errorf("failed to fetch rows affected, err=%v", err)
	}

	if rowsAffected == 0 {
		r.log.Error("song with id", song.ID, "not found")
		return false, fmt.Errorf("song with ID %d not found", song.ID)
	}

	r.log.Debug("song:", song.Song, " updated")

	return true, nil
}

func (r *SongRepository) Delete(song models.Song) (int, error) {
	r.log.Debug("starting to delete a song", song.Song, song.GroupName)

	stmt := `DELETE FROM songs WHERE id = $1 RETURNING id`

	err := r.db.QueryRow(stmt, song.ID).Scan(&song.ID)
	if err != nil {
		r.log.Error("can't delete song with id", song.ID, "err:", err)
		return 0, fmt.Errorf("can't delete song, err=%v", err)
	}

	r.log.Debug("song:", song.Song, " deleted")

	return song.ID, nil
}

func (r *SongRepository) Get(groupName, songName, releaseDate string, limit, offset, songID int) ([]models.Song, error) {
	r.log.Debug("start retrieving songs/songs from the database")

	stmt := `SELECT s.id, s.song, g.name, s.text, s.link, s.releasedate 
             FROM songs s 
             JOIN groups g on s.group_name = g.name
             WHERE 
               ($1::text IS NULL OR g.name ILIKE $1) 
               AND ($2::text IS NULL OR s.song ILIKE $2) 
               AND ($3::int IS NULL OR s.id = $3)
               AND ($4::text IS NULL OR s.releasedate = $4)
             LIMIT $5 OFFSET $6`

	r.log.Info("Parameters received", "groupName", groupName, "songName", songName, "limit", limit, "offset", offset, "songID", songID, "date", releaseDate)

	rows, err := r.db.Query(stmt, groupName, songName, songID, releaseDate, limit, offset)
	if err != nil {
		r.log.Error("can't fetch all songs, err", err)
		return nil, fmt.Errorf("can't fetch all songs, err=%v", err)
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err = rows.Scan(&song.ID, &song.Song, &song.GroupName, &song.Text, &song.Link, &song.ReleaseDate)
		if err != nil {
			r.log.Error("error scannin row, err", err)
			return nil, fmt.Errorf("error scannin row, err=%v", err)
		}
		songs = append(songs, song)
	}

	err = rows.Err()
	if err != nil {
		r.log.Error("rows error, err", err)
		return nil, fmt.Errorf("rows error, err=%v", err)
	}

	r.log.Debug("Parameters received", "groupName", groupName, "songName", songName, "limit", limit, "offset", offset, "songID", songID, "date", releaseDate)

	return songs, nil
}
