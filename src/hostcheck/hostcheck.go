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

func Check(ips []string, domain, path string) (ret Result) {
	ret.Path = path
	ret.Domain = domain
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
		ret.Results = append(ret.Results, result)
	}
	return
}
