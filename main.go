package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type BakaFile struct {
	Name string
	Size int64
}

type HashId struct{ Hash [32]byte }

func doTheyReallyExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func removeNonSusyDuplicates(collectedFiles map[BakaFile][]string) {
	for key := range collectedFiles {
		if len(collectedFiles[key]) < 2 {
			delete(collectedFiles, key)
		}
	}
}

func findSusyFiles(path string, collectedFiles map[BakaFile][]string) {
	var dirs []string
	filesInDir, _ := os.ReadDir(path)
	for _, file := range filesInDir {
		filePath := filepath.Join(path, file.Name())
		if isFileDir := file.IsDir(); isFileDir {
			dirs = append(dirs, filePath)
		} else {
			fi, _ := file.Info()
			kFile := BakaFile{file.Name(), fi.Size()}
			collectedFiles[kFile] = append(collectedFiles[kFile], filePath)
		}
	}

	if len(dirs) > 0 {
		for _, dir := range dirs {
			findSusyFiles(dir, collectedFiles)
		}
	}
}

func getFilesWithSameContent(collectedFiles map[BakaFile][]string) map[HashId][]string {
	set := make(map[HashId][]string)
	for i := range collectedFiles {
		for j := range collectedFiles[i] {
			f, err := os.ReadFile(collectedFiles[i][j])
			if err != nil {
				log.Fatal(err)
			}
			hash := sha256.Sum256(f)
			hashId := HashId{hash}
			set[hashId] = append(set[hashId], collectedFiles[i][j])
		}
	}

	return set
}

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Println("The argument was not provided! ಠ_ಠ")
		return
	}

	var path string
	path = os.Args[1]

	if len(path) > 1 && path != "." {
		exists, _ := doTheyReallyExist(path)
		if !exists {
			fmt.Println(path, "path does not exist ヾ( ･`⌓´･)ﾉﾞ!")
			return
		}
	} else if path == "." {
		pathCurrent, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		path = pathCurrent
	}

	collectedFiles := make(map[BakaFile][]string)

	findSusyFiles(path, collectedFiles)

	removeNonSusyDuplicates(collectedFiles)

	wereDuplicates := false
	if len(collectedFiles) != 0 {
		sameContentFiles := getFilesWithSameContent(collectedFiles)
		for _, value := range sameContentFiles {
			if len(value) > 1 {
				wereDuplicates = true
				fmt.Println("DUPLICATES ( ⚆ _ ⚆ ): ", value)
			}
		}
	} else if !wereDuplicates {
		fmt.Println(fmt.Sprintf("No duplicates were found from \"%s\" (＾º◡º＾✥)", path))
	}
}
