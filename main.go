package main

import (
    "fmt"
    "net/http"
    "os"
    "bytes"
    "io/ioutil"
)

func main() {
    fmt.Printf("Hi Tom")

    // Need to get the sha below from the API request rather than hardcoding
    url := "https://api.github.com/repos/tomkadwill/mud/statuses/fed9d6dc2155cea9fb5bbce3243372194acc9fc4"
    var jsonStr = []byte(`{"state": "failure","target_url": "https://example.com/build/status","description": "The build failed!","context": "toms-go/check"}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    req.SetBasicAuth(os.Getenv("GITHUB_USERNAME"), os.Getenv("GITHUB_PASSWORD"))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))



    http.HandleFunc("/", hello)
    fmt.Println("listening...")
    err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}

func hello(res http.ResponseWriter, req *http.Request) {
    fmt.Printf("Hi Tom, printing hello world")

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
      //panic()
    }

    fmt.Fprintln(res, string(body), "hello, world")
}
