package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/YunKyungho/searchFile/models"
)

type Database struct {
	handler  *sql.DB
	oldFiles map[int]struct{}
	oldDirs  map[int]struct{}
}

// NewDatabase is Constructor that open DB
func NewDatabase(dbFile string) *Database {
	handler, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalln(err)
	}
	db := Database{handler: handler}
	db.createTable()
	db.setOldDirs()
	db.setOldFiles()
	return &db
}

// CloseDatabase is close DB
func (d Database) CloseDatabase() {
	d.handler.Close()
}

// createTable named directory_info and file_info
func (d Database) createTable() {
	query := `
		CREATE TABLE IF NOT EXISTS directory_info (
			di_no INTEGER PRIMARY KEY AUTOINCREMENT,
			di_path TEXT NOT NULL,
			di_modified_date DATE NOT NULL,
			UNIQUE (di_path)
		);
	`
	if _, err := d.handler.Exec(query); err != nil {
		log.Fatalf("Failed to create directory_info table: %v", err)
	}

	query = `
		CREATE TABLE IF NOT EXISTS file_info (
			fi_no INTEGER PRIMARY KEY AUTOINCREMENT,
			fi_name TEXT NOT NULL,
			fi_modified_date DATE NOT NULL,
			fi_parent INTEGER NOT NULL,
		    UNIQUE (fi_name, fi_parent),
			FOREIGN KEY(fi_parent) REFERENCES directory_info(di_no)
		);
	`
	if _, err := d.handler.Exec(query); err != nil {
		log.Fatalf("Failed to create file_info table: %v", err)
	}
}

// setOldDirs sets oldDirs in the Database
func (d Database) setOldDirs() {
	query := `SELECT di_no FROM directory_info;`
	rows, err := d.handler.Query(query)
	if err != nil {
		log.Fatalf("Failed to query di_no in getAllDirNum: %v\n", err)
	}
	defer rows.Close()

	var diNo int
	for rows.Next() {
		err := rows.Scan(&diNo)
		if err != nil {
			log.Fatalf("Failed to scan di_no in getAllDirNum: %v\n", err)
		}
		d.oldDirs[diNo] = struct{}{}
	}
}

// setOldFiles sets oldFiles in the Database
func (d Database) setOldFiles() {
	query := `SELECT fi_no FROM file_info;`
	rows, err := d.handler.Query(query)
	if err != nil {
		log.Fatalf("Failed to query fi_no in getAllFileNum: %v\n", err)
	}
	defer rows.Close()

	var fiNo int
	for rows.Next() {
		err := rows.Scan(&fiNo)
		if err != nil {
			log.Fatalf("Failed to scan fi_no in getAllFileNum: %v\n", err)
		}
		d.oldFiles[fiNo] = struct{}{}
	}
}

// DeleteOldData deletes non-existent rows
func (d Database) DeleteOldData() {
	// map의 key 목록만 만들어서 directory_info와 file_info의 필요없는 데이터 삭제하는 로직.
}

// CreateIndex of fi_name in file_info
func (d Database) CreateIndex() {
	query := `CREATE INDEX IF NOT EXISTS file_info_fi_name_idx ON file_info (fi_name);`
	if _, err := d.handler.Exec(query); err != nil {
		log.Fatalf("Failed to create file_info index: %v", err)
	}
}

// InsertAllData inserts all data from the data and removes a valid list of data from oldDirs and oldFiles
func (d Database) InsertAllData(data map[string]models.Directory) {
	builder := strings.Builder{}
	builder.WriteString(`INSERT INTO directory_info (di_path, di_modified_date) VALUES `)
	var values []string
	for path, dir := range data {
		values = append(values, fmt.Sprintf("('%s', '%s')", path, dir.ModTime))
	}
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString(` ON CONFLICT(di_path) DO UPDATE SET di_modified_date = excluded.di_modified_date RETURNING di_no, di_path;`)

	rows, err := d.handler.Query(builder.String())
	if err != nil {
		log.Fatalf("Failed to insert directory_info in InsertAllData: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var diNo int
		var diPath string
		err := rows.Scan(&diNo, &diPath)
		if err != nil {
			log.Fatalf("Failed to scan directory_info in InsertAllData: %v\n", err)
		}
		delete(d.oldDirs, diNo)
		val, exists := data[diPath]
		if exists {
			val.DiNo = diNo
		}
	}

	builder.Reset()
	builder.WriteString(`INSERT INTO file_info (fi_name, fi_modifed_date, fi_parent) VALUES `)
	values = values[:0]
	for _, dir := range data {
		for _, file := range dir.Child {
			values = append(values, fmt.Sprintf("('%s', '%s', '%d')", file.Name, file.ModTime, dir.DiNo))
		}
	}
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString(` ON CONFLICT(fi_name, fi_parent) DO UPDATE SET fi_modifed_date = excluded.fi_modifed_date RETURNING fi_no;`)

	rows2, err2 := d.handler.Query(builder.String())
	if err2 != nil {
		log.Fatalf("Failed to insert file_info in InsertAllData: %v\n", err)
	}
	defer rows2.Close()
	for rows2.Next() {
		var fiNo int
		err := rows2.Scan(&fiNo)
		if err != nil {
			log.Fatalf("Failed to scan file_info in InsertAllData: %v\n", err)
		}
		delete(d.oldFiles, fiNo)
	}

}

func (d Database) SelectFileInfo() {

}

func (d Database) UpdateFileInfo() {

}

func (d Database) DeleteFileInfo() {

}

func (d Database) UpdateDirectoryInfo() {

}

func (d Database) DeleteDirectoryInfo() {

}
