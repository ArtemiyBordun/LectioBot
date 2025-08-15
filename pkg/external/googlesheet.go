package external

import (
	"LectioBot/internal/models"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"unicode"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Sheet struct {
	SpreadsheetID   string
	CredentialsFile string
	Students        []Student
}

type Student struct {
	FullName string
	Group    string
}

func NewSheet(spreadsheetID, credentialsFile string) *Sheet {
	return &Sheet{
		SpreadsheetID:   spreadsheetID,
		CredentialsFile: credentialsFile,
	}
}

func (s *Sheet) LoadStudentsFromAllSheets() {
	ctx := context.Background()

	b, err := ioutil.ReadFile(s.CredentialsFile)
	if err != nil {
		log.Fatalln(err)
		return
	}
	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalln(err)
		return
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(err)
		return
	}

	spreadsheet, err := srv.Spreadsheets.Get(s.SpreadsheetID).Do()
	if err != nil {
		log.Fatalln(err)
		return
	}

	var students []Student

	for _, sheet := range spreadsheet.Sheets {
		sheetName := sheet.Properties.Title
		readRange := fmt.Sprintf("%s!B2:B", sheetName)

		resp, err := srv.Spreadsheets.Values.Get(s.SpreadsheetID, readRange).Do()
		if err != nil {
			log.Println("ошибка чтения листа %s: %w", sheetName, err)
			return
		}

		for _, row := range resp.Values {
			if len(row) == 0 {
				continue
			}
			fullName := fmt.Sprintf("%v", row[0])
			students = append(students, Student{
				FullName: fullName,
				Group:    sheetName,
			})
		}
	}

	s.Students = students
}

func (s *Sheet) FindStudentGroup(fullName string) string {
	fullName = strings.TrimSpace(strings.ToLower(fullName))
	for _, s := range s.Students {
		if strings.ToLower(strings.TrimSpace(s.FullName)) == fullName {
			return s.Group
		}
	}
	return ""
}

func (s *Sheet) IsGroup(input string) (bool, error) {
	input = strings.TrimSpace(input)

	for _, s := range s.Students {
		if strings.EqualFold(s.Group, input) {
			return true, nil
		}
	}

	return false, nil
}

func (s *Sheet) Marking(student *models.Student, numLecture int) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(s.CredentialsFile)
	if err != nil {
		log.Fatalf("Не удалось прочитать креды: %v", err)
	}
	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Ошибка парсинга кредов: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	readRange := fmt.Sprintf("%s!A1:Z", student.Group)
	resp, err := srv.Spreadsheets.Values.Get(s.SpreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Ошибка чтения листа: %v", err)
	}

	var rowIndex int = -1
	for i, row := range resp.Values {
		if len(row) > 1 {
			fullName := strings.TrimSpace(fmt.Sprintf("%v", row[1]))
			if strings.EqualFold(fullName, strings.TrimSpace(student.Name)) {
				rowIndex = i + 1
				break
			}
		}
	}

	if rowIndex == -1 {
		log.Printf("Студент %s не найден в группе %s", student.Name, student.Group)
		return
	}
	diff := 3
	if !isFirstSemester(student.Group) {
		diff = 4 //т.к есть разница где находится в таблице лекции по оп в 1 и 2 семе
	}
	lectureColIndex := 2 + (numLecture-1)*diff
	colLetter := string(rune('A' + lectureColIndex))

	cellRange := fmt.Sprintf("%s!%s%d", student.Group, colLetter, rowIndex)

	vr := &sheets.ValueRange{
		Values: [][]interface{}{
			{""},
		},
	}
	_, err = srv.Spreadsheets.Values.Update(s.SpreadsheetID, cellRange, vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Ошибка записи: %v", err)
	}
}

func isFirstSemester(group string) bool {
	parts := strings.Split(group, "-")
	if len(parts) < 2 {
		return false
	}

	for _, r := range parts[1] {
		if unicode.IsDigit(r) {
			return r == '1'
		}
	}
	return false
}
