package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/YunKyungho/searchFile/pkg"
)

const timeFormat string = "2024-10-19 00:00:00"

// scanDirectory is Browse all files and directories from the root directory using WalkDir
func scanDirectory(db *pkg.Database, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// fmt.Printf("Skipping: %s (Error: %v)\n", path, err)
			return nil
		}

		info, err := d.Info()
		if err != nil {
			// fmt.Printf("Skipping: %s (Error: %v)\n", path, err)
			return nil
		}

		if info.IsDir() {
			db.InsertDirectoryInfo(info.Name(), info.ModTime().Format(timeFormat))
		} else {
			directory := path[:len(path)-len(info.Name())]
			diNo := db.SelectDiNo(directory)
			if diNo != 0 {
				db.InsertFileInfo(info.Name(), info.ModTime().Format(timeFormat), diNo)
			}
		}

		return nil
	})
}

func main() {
	db := pkg.NewDatabase("test.db")
	defer db.CloseDatabase()
	// 로컬에 물리적으로 마운트된 파일시스템만 구하는 선행 로직 필요.
	root := "/"

	// 디렉토리 탐색 후 경로와 파일명 저장.
	if err := scanDirectory(db, root); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	db.CreateIndex()

	// 위 작업 종료 후 파일 CRUD 이벤트 기반으로 무한 루프
	// + 외부 프로세스가 특정 시그널 보냈을 때 sqlite에서 검색해주는 작업.
}
