package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	c := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost/redirect", nil)
	req.Host = "foo"
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		panic(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("get:\n", string(body))
}
