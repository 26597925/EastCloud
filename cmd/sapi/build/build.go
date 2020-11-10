package build

import (
	"encoding/base64"
	"fmt"
	"github.com/26597925/EastCloud/pkg/bootstrap/flag"
	"github.com/26597925/EastCloud/pkg/config/encoder/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Start (name string, fs *flag.Set) {
	service := fs.String(name)

	dir, _ := os.Executable()
	curPath := filepath.Dir(dir)

	path := fs.String("path")
	if path == "." || path == ".." {
		path = curPath
	}

	bt, err := ioutil.ReadFile(filepath.Join(curPath, "pro.cache"))
	if err != nil {
		fmt.Println(err)
		return
	}

	var data map[string]string
	err = json.NewEncoder().Decode(bt, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = buildFile(service, path, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Generated successfully !")
}

func buildFile(service string, path string, data map[string]string) error {
	files := map[string]string{
		"sapi/cmd/hello/router": "sapi/cmd/"+service+"/router",
		"sapi/cmd/hello/boot/engine": "sapi/cmd/"+service+"/boot/engine",
		"sapi/cmd/hello/boot": "sapi/cmd/"+service+"/boot",
		"sapi/cmd/hello/controller": "sapi/cmd/"+service+"/controller",
	}

	for file, content := range data{
		if strings.Contains(file, "hello.go") {
			file = strings.ReplaceAll(file, "hello", service)
		}

		filename := filepath.Join(path, service, file)
		by, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return err
		}
		str := string(by)
		for source, target := range files {
			if strings.Contains(str, source) {
				str = strings.ReplaceAll(str, source, target)
			}
		}

		dir := filepath.Dir(filename)
		is, err := pathExists(dir)
		if err != nil {
			return err
		}

		if !is {
			err=os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		err = ioutil.WriteFile(filename, []byte(str), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}