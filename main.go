package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// 재귀적으로 모든 파일을 탐색하는 함수
func buildIndex(root string) error {
	// WalkDir을 이용해 루트 디렉토리부터 모든 파일과 디렉토리 탐색
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Skipping: %s (Error: %v)\n", path, err)
			return nil
		}

		info, err := d.Info()
		if err != nil {
			fmt.Printf("Skipping: %s (Error: %v)\n", path, err)
			return nil
		}

		size := info.Size()
		modTime := info.ModTime()
		isDir := info.IsDir()

		return nil
	})
}

func main() {
	// 로컬에 물리적으로 마운트된 파일시스템만 구하는 선행 로직 필요.
	root := "/"

	if err := buildIndex(root); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

}
