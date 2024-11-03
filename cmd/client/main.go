package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080/"
	//data container for request
	values := url.Values{}
	fmt.Println("Input long URL")
	//read  string from the console and put into the data
	readString, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	readString = strings.TrimSuffix(readString, "\n")
	values.Set("url", readString)
	//execute POST request
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//print result
	fmt.Println("Status Code ", response.StatusCode)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
