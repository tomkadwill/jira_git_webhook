package main

import (
    "fmt"
    "net/http"
    "os"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "strings"
)

type PullRequestResponse struct {
    Url         string
    PullRequest struct {
        CommitsUrl string `json:"commits_url"`
      } `json:"pull_request"`
}

type Commits []map[string]string

type Commit struct {
    Sha     string
    Commit  struct {
      Message string
    }
}

func main() {
    http.HandleFunc("/", hello)
    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}

func hello(res http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
      //panic()
    }

    var pr_request PullRequestResponse
    err = json.Unmarshal([]byte(string(body)), &pr_request)
    if err == nil {
        commitsUrl := pr_request.PullRequest.CommitsUrl

        req, err := http.NewRequest("GET", commitsUrl, nil)
        req.Header.Set("X-Custom-Header", "myvalue")
        req.Header.Set("Content-Type", "application/json")
        req.SetBasicAuth(os.Getenv("GITHUB_USERNAME"), os.Getenv("GITHUB_PASSWORD"))

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()

        body, err = ioutil.ReadAll(resp.Body)
        fmt.Println(string(body))
        fmt.Println("response Body3^^")

        var commits Commits
        err = json.Unmarshal([]byte(string(body)), &commits)

        fmt.Println(commits[0]["sha"])
        fmt.Println("response commits^^")


        // Make a request for single commit
        s := []string{"https://api.github.com/repos/tomkadwill/mud/commits/", commits[0]["sha"]};
        url := strings.Join(s, "")

        req, err = http.NewRequest("GET", url, nil)
        req.Header.Set("X-Custom-Header", "myvalue")
        req.Header.Set("Content-Type", "application/json")
        req.SetBasicAuth(os.Getenv("GITHUB_USERNAME"), os.Getenv("GITHUB_PASSWORD"))

        client = &http.Client{}
        resp, err = client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()

        body, err = ioutil.ReadAll(resp.Body)

        var commit Commit
        err = json.Unmarshal([]byte(string(body)), &commit)
        commit_message := commit.Commit.Message

    } else {
        // Do something
    }


    // Need to get the sha below from the API request rather than hardcoding
    url := "https://api.github.com/repos/tomkadwill/mud/statuses/fed9d6dc2155cea9fb5bbce3243372194acc9fc4"
    var jsonStr = []byte(`{"state": "failure","target_url": "https://example.com/build/status","description": "The build failed!","context": "toms-go/check"}`)
    req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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



    url = "http://requestb.in/15xawwi1"
    req, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(pr_request.Url)))
    client = &http.Client{}
    client.Do(req)
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    fmt.Fprintln(res, string(body), "hello, world")
}
