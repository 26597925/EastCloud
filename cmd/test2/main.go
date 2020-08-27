package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sapi/pkg/config/etcd"
)

func main() {
	fh, err := os.Open("config/config.yml")
	fmt.Println(err)
	if err != nil {
		return
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	fmt.Println(err)
	if err != nil {
		return
	}

	options := &etcd.Options{
		Debug:true,
		Timeout: 3,
	}
	etcd,_ := etcd.NewEtcd(options)

	fmt.Println(string(b))

	etcd.Put(string(b))
}
