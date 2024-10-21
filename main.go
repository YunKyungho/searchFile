package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/YunKyungho/searchFile/models"
	"github.com/YunKyungho/searchFile/pkg"
)

// checkMemUsage verifies that memory usage exceeds 200MB
func checkMemUsage() bool {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	if m.Alloc/1024 > 51200 {
		return true
	}
	return false
}

const timeFormat string = "2024-10-19 00:00:00"

// scanDirectory is Browse all files and directories from the root directory using WalkDir
func scanDirectory(db *pkg.Database, tmpMap *map[string]models.Directory, root string) error {
	// map은 원래도 참조 타입이라 인자에 넘겼을 때 원본 수정은 가능하지만
	// 포인터로 받는 이유는 인자로 넘긴 map을 함수 내에서 초기화하는 것은 불가능하기 때문이다.
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

		if strings.Contains(info.Name(), "'") {
			return nil
		}

		if info.IsDir() {
			if checkMemUsage() {
				db.InsertAllData(*tmpMap)
				*tmpMap = map[string]models.Directory{}
				runtime.GC()
			}
			(*tmpMap)[path] = models.Directory{ModTime: info.ModTime().Format(timeFormat), Child: make([]models.File, 0)}
			// map의 포인터를 함수 인자로 넘겼을 때 함수 내부에서 역참조한 map에서 key로 값을 조회하면 go가 이를 slice로 인식하여 오류가 발생한단다.
			// 명시적으로 역참조한 포인터를 어떤 타입인지를 알려주기 위해 () 괄호를 사용한다.
		} else {
			dir := path[:len(path)-len(info.Name())-1]
			dirInfo, exists := (*tmpMap)[dir]
			if exists {
				file := models.File{Name: info.Name(), ModTime: info.ModTime().Format(timeFormat)}
				dirInfo.Child = append(dirInfo.Child, file)
				(*tmpMap)[dir] = dirInfo
				// map에서 구조체를 가져올 때 원본이 아닌 복사본을 가져온다.
				// 따라서 수정한 뒤 위 처럼 원본의 값을 저장해줘야한다.
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

	tmpMap := map[string]models.Directory{}
	if err := scanDirectory(db, &tmpMap, root); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	db.InsertAllData(tmpMap)

	// db의 oldFiles, oldDirs의 남은 번호는 전부 delete
	db.CreateIndex()

	// 위 작업 종료 후 파일 CRUD 이벤트 기반으로 무한 루프
	// + 외부 프로세스가 특정 시그널 보냈을 때 sqlite에서 검색해주는 작업.
}
