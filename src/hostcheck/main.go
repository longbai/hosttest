package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	output := flag.String("o", "", "result file")
	domain := flag.String("d", "", "domain")
	path := flag.String("p", "", "path, key or key+fop+token")
	ipFile := flag.String("ips", "", "ip file")
	goroutines := flag.Int("g", 0, "go routines")
	flag.Parse()
	if *domain == "" || *path == "" || *ipFile == "" {
		flag.PrintDefaults()
		fmt.Println("invalid args")
		return
	}
	data, err := ioutil.ReadFile(*ipFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	var ips []string
	err = json.Unmarshal(data, &ips)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *goroutines == 0 {
		*goroutines = 64
	}
	ret := Check(ips, *domain, *path, *goroutines)
	out, err := json.MarshalIndent(ret, "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	if *output == "" {
		fmt.Println(string(out))
	} else {
		ioutil.WriteFile(*output, out, 0777)
	}
}
