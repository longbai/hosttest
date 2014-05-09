package main

import (
	"net/http"
	"time"
)

type IPResult struct {
	ReqId  string        `json:"req_id"`
	IP     string        `json:"ip"`
	Time   time.Duration `json:"time"`
	Status int           `status:"code"`
	Error  string        `json:"error"`
	Length int64         `json:"content-length"`
}

type Result struct {
	Domain  string     `json:"domain"`
	Path    string     `json:"path"`
	Results []IPResult `json:"results"`
}

func calcSlice(length, goroutines int) int {
	count := length / goroutines
	mode := length % goroutines
	if mode == 0 {
		return count
	}
	return count + 1
}

func job(ips []string, domain, path string, send chan<- IPResult) {
	client := http.DefaultClient
	client.Transport = &http.Transport{
		ResponseHeaderTimeout: time.Second * 30,
	}
	for _, ip := range ips {
		urlStr := "http://" + ip + "/" + path
		req, _ := http.NewRequest("GET", urlStr, nil)
		req.Host = domain
		req.Header.Set("User-Agent", "go hostcheck")
		t := time.Now()
		resp, err := client.Do(req)
		elapse := time.Since(t)
		result := IPResult{
			IP:   ip,
			Time: elapse,
		}

		if err != nil {
			result.Error = err.Error()
		}
		if resp != nil {
			result.Status = resp.StatusCode
			result.ReqId = resp.Header.Get("X-Reqid")
			result.Length = resp.ContentLength
			if resp.Body != nil {
				resp.Body.Close()
			}
		}
		send <- result
	}
}

func Check(ips []string, domain, path string, goroutines int) (ret Result) {
	if goroutines == 0 {
		goroutines = 1
	}
	ret.Path = path
	ret.Domain = domain
	length := len(ips)
	receiver := make(chan IPResult, length)
	count := calcSlice(len(ips), goroutines)
	for i := 0; i < goroutines; i++ {
		start := count * i
		if start >= length {
			break
		}
		end := count * (i + 1)
		if end >= length {
			end = length
		}
		go job(ips[start:end], domain, path, receiver)
	}
	for j := 0; j < length; j++ {
		ret.Results = append(ret.Results, <-receiver)
	}

	return
}
