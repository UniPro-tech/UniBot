package main

import (
	"context"
	"log"

	"unibot/internal/db"
	"unibot/internal/repository"
)

func main() {
	db, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	memberRepo := repository.NewMemberRepository(db)

	err = memberRepo.Create(context.Background(), "1234567890")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("members INSERTできた✨")
}
