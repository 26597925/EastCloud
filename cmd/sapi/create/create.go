package create

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sapi/pkg/bootstrap/flag"
	"sapi/pkg/config/encoder/json"
)

func Start (name string, fs *flag.Set) {
	path := fs.String("path")
	if path == "." || path == ".." {
		dir, _ := os.Executable()
		path = filepath.Dir(dir)
	}

	s, err := getAllFile(path, make([]string, 0))
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := createCacheFile(path, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := json.NewEncoder().Encode(data)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile("pro.cache", b, os.ModePerm)
}

func getAllFile(pathname string, files []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return files, err
	}
	for _, fl:= range rd {
		if fl.IsDir() {
			path := filepath.Join(pathname, fl.Name())
			files, err = getAllFile(path, files)
			if err != nil {
				return files, err
			}
		} else {
			path := filepath.Join(pathname, fl.Name())
			files = append(files, path)
		}
	}
	return files, nil
}

func createCacheFile (path string, files []string) (map[string]string, error) {
	cache := map[string]string{}
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		cache[file[len(path):]] = base64.StdEncoding.EncodeToString(data)
	}

	return cache, nil
}
