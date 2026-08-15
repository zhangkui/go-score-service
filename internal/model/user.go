package model

type User struct {
	ID    string
	Name  string
	Score int64
}

type ScoreEntry struct {
	UserID string
	Score  int64
}
