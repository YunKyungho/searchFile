package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/YunKyungho/searchFile/models"
	"github.com/YunKyungho/searchFile/pkg"
)

// checkMemUsage verifies that memory usage exceeds 200MB
func checkMemUsage() bool {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	if m.Alloc/1024 > 204800 {
		return true
	}
	return false
}

const timeFormat string = "2024-10-19 00:00:00"

// scanDirectory is Browse all files and directories from the root directory using WalkDir
func scanDirectory(db *pkg.Database, root string) error {
	tmpMap := map[string]models.Directory{}

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
			if checkMemUsage() {
				db.InsertAllData(tmpMap)
				tmpMap = map[string]models.Directory{}
			}
			tmpMap[path] = models.Directory{ModTime: info.ModTime().Format(timeFormat)}
		} else {
			dir := path[:len(path)-len(info.Name())]
			dirInfo, exists := tmpMap[dir]
			if exists {
				file := models.File{Name: info.Name(), ModTime: info.ModTime().Format(timeFormat)}
				dirInfo.Child = append(dirInfo.Child, file)
			}
		}

		return nil
	})
}

func main() {
	db := pkg.NewDatabase("searchFile.db")
	defer db.CloseDatabase()
	// 로컬에 물리적으로 마운트된 파일시스템만 구하는 선행 로직 필요.
	root := "/"
	if err := scanDirectory(db, root); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// db의 oldFiles, oldDirs의 남은 번호는 전부 delete
	db.CreateIndex()

	// 위 작업 종료 후 파일 CRUD 이벤트 기반으로 무한 루프
	// + 외부 프로세스가 특정 시그널 보냈을 때 sqlite에서 검색해주는 작업.
}
