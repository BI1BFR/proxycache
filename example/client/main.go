package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough args")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("not enough args")
			os.Exit(1)
		}
		get(os.Args[2])
	case "put":
		if len(os.Args) < 4 {
			fmt.Println("not enough args")
			os.Exit(1)
		}
		if len(os.Args) == 4 {
			put(os.Args[2], os.Args[3], "0")
		} else {
			put(os.Args[2], os.Args[3], os.Args[4])
		}
	case "status":
		status()
	default:
		fmt.Println("wrong cmd")
	}
}

const URLBase = "http://localhost/pc/v1/"

var httpClient = &http.Client{}

func get(key string) {
	fmt.Println("Get: ", key)
	rsp, err := httpClient.Get(URLBase + "keys/" + key)
	if err != nil {
		fmt.Println("HTTP request fail: ", err)
		return
	}

	if rsp.StatusCode == http.StatusNotFound {
		fmt.Println("Not Found")
		return
	}

	if rsp.StatusCode != http.StatusOK {
		fmt.Println("Error Status Code: ", rsp.StatusCode)
		return
	}

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("Error read response body")
		return
	}

	fmt.Println("Return: ", string(b))
}

func put(key string, value string, ttw string) {
	fmt.Println("Put: ", key, " ", value, " ", ttw)

	req, err := http.NewRequest("PUT", URLBase+"keys/"+key+"?ttw="+ttw, bytes.NewBuffer([]byte(value)))
	if err != nil {
		fmt.Println("Fail build request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error Status Code: ", resp.StatusCode)
	} else {
		fmt.Println("Put OK.")
	}
}

func status() {
	rsp, err := httpClient.Get(URLBase + "status")
	if err != nil {
		fmt.Println("HTTP request fail: ", err)
	}

	if rsp.StatusCode != http.StatusOK {
		fmt.Println("Error Status Code: ", rsp.StatusCode)
		return
	}

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("Error read response body")
		return
	}

	fmt.Println("Return: ", string(b))
}
