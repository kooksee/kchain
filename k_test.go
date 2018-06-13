package kchain

import (
	"testing"
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"sort"
)

func TestFile(t *testing.T) {
	dataFilePath := make([]string, 0)

	filepath.Walk("kdata", func(path string, fi os.FileInfo, err error) error {
		if strings.Contains(path, "kdata/data") && fi.IsDir() && !strings.Contains(path, ".") {
			fmt.Println(path)
			return nil
		}

		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}

		fs := strings.Split(fi.Name(), "_")
		if len(fs) != 2 {
			return nil
		}

		dataFilePath = append(dataFilePath, fi.Name())
		return nil
	})
	a := []string{"a", "c", "b", "245", "123", "543"}
	sort.Strings(a)
	fmt.Println(a)

	if len(dataFilePath) == 0 {
		dataFilePath = append(dataFilePath, filepath.Join("kdata", "data"))
	}

	fmt.Println(dataFilePath)

}
