package main

import (
	"LectioBot/internal/adapter"
	"LectioBot/internal/config"
	"LectioBot/internal/storage"
	"LectioBot/pkg/external"
	"fmt"
	"sync"
)

func main() {
	cfg := config.LoadConfig()

	db, err := storage.InitSQLite("lectio.db")
	if err != nil {
		panic(fmt.Sprintf("Ошибка подключения к SQLite: %v", err))
	}

	fmt.Println("Подключение к SQLite успешно!")

	sheet := external.NewSheet(cfg.SpreadsheetID, cfg.CredentialsFile)
	go sheet.LoadStudentsFromAllSheets()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		adapter.StartBot(cfg, db, sheet)
	}()

	wg.Wait()
}
