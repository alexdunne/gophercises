package cmd

import (
	"tasker/db"

	"github.com/boltdb/bolt"
)

type CommandHandler struct {
	store *db.Store
}

func NewCommandHandler(database *bolt.DB) *CommandHandler {
	store := db.NewStore(database)

	return &CommandHandler{
		store: store,
	}
}
