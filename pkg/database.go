package pkg

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	handler *sql.DB
}

// NewDatabase is Constructor that open DB
func NewDatabase(dbFile string) *Database {
	handler, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalln(err)
	}
	db := Database{handler: handler}
	db.createTable()
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
			di_modified_date DATE NOT NULL
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

// CreateIndex of fi_name in file_info
func (d Database) CreateIndex() {
	query := `CREATE INDEX IF NOT EXISTS file_info_fi_name_idx ON file_info (fi_name);`
	if _, err := d.handler.Exec(query); err != nil {
		log.Fatalf("Failed to create file_info index: %v", err)
	}
}

// InsertFileInfo with name, modeDate, parent. parent is foreign key of directory_info
func (d Database) InsertFileInfo(name string, modDate string, parent int) {
	query := `INSERT INTO file_info (fi_name, fi_modifed_date, fi_parent) VALUES (?, ?, ?) ON CONFLICT(fi_name, fi_parent) DO NOTHING;`
	_, err := d.handler.Exec(query, name, modDate, parent)
	if err != nil {
		log.Fatalf("Failed to insert file_info: %v", err)
	}
}

func (d Database) SelectFileInfo() {

}

func (d Database) UpdateFileInfo() {

}

func (d Database) DeleteFileInfo() {

}

// InsertDirectoryInfo with path, modDate
func (d Database) InsertDirectoryInfo(path string, modDate string) {
	query := `INSERT INTO directory_info (di_path, di_modified_date) VALUES (?, ?) ON CONFLICT(di_path) DO NOTHING;`
	_, err := d.handler.Exec(query, path, modDate)
	if err != nil {
		log.Fatalf("Failed to insert derectory_info: %v", err)
	}
}

// SelectDiNo with di_path
func (d Database) SelectDiNo(path string) int {
	query := `SELECT di_no FROM directory_info WHERE di_path = ?`
	var diNo int
	err := d.handler.QueryRow(query, path).Scan(&diNo)
	if err != sql.ErrNoRows && err != nil {
		log.Fatalf("Failed to select di_no: %v", err)
	}
	return diNo
}

func (d Database) UpdateDirectoryInfo() {

}

func (d Database) DeleteDirectoryInfo() {

}
