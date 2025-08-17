package main

import (
	"fmt"
	"sync"

	"LectioBot/internal/adapter"
	"LectioBot/internal/config"
	"LectioBot/internal/storage"
	"LectioBot/pkg/external"
)

func main() {
	cfg := config.LoadConfig() //Загружаем конфиг бота

	db, err := storage.InitSQLite("lectio.db")
	if err != nil {
		panic(fmt.Sprintf("Ошибка подключения к SQLite: %v", err))
	}

	sheet := external.NewSheet(cfg.SpreadsheetID, cfg.CredentialsFile)
	//go sheet.LoadStudentsFromAllSheets() //создаем коннект с гугл таблицами

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		adapter.StartBot(cfg, db, sheet)
	}()

	wg.Wait()
}
