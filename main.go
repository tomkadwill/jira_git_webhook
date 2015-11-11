package main

import (
    "fmt"
    "net/http"
    "os"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "strings"
    "regexp"
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
    http.HandleFunc("/", handleRequest)
    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
      panic(err)
    }

    var pr_request PullRequestResponse
    err = json.Unmarshal([]byte(string(body)), &pr_request)
    if err == nil {
        commits := getCommits(pr_request)
        failures := false
        for i := 0; i < len(commits); i++ {

        // Make a request for single commit
        commit := getCommit(commits, i)

        commit_message := commit.Commit.Message

        match, _ := regexp.MatchString("\\[PLAT-(.*)\\]", commit_message)
        if (match==true && failures==false) {
          setStatus(commit.Sha, "success")
        } else {
          setStatus(commit.Sha, "failure")
          failures = true
        }
      }

    } else {
        panic(err)
    }
}

func getCommits(pr_request PullRequestResponse) Commits {
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

  body, err := ioutil.ReadAll(resp.Body)

  var commits Commits
  err = json.Unmarshal([]byte(string(body)), &commits)

  return commits
}

func getCommit(commits Commits, i int) Commit{
  s := []string{"https://api.github.com/repos/Babylonpartners/babylon/commits/", commits[i]["sha"]}
  url := strings.Join(s, "")

  req, err := http.NewRequest("GET", url, nil)
  req.Header.Set("X-Custom-Header", "myvalue")
  req.Header.Set("Content-Type", "application/json")
  req.SetBasicAuth(os.Getenv("GITHUB_USERNAME"), os.Getenv("GITHUB_PASSWORD"))

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  var commit Commit
  err = json.Unmarshal([]byte(string(body)), &commit)

  return commit
}

func setStatus(sha string, state string) {
  s := []string{"https://api.github.com/repos/Babylonpartners/babylon/statuses/", sha}
  url := strings.Join(s, "")

  jsonStr := []byte(jsonBody(state))
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
}

func jsonBody(state string) string{
  description := ""
  if state == "failure" {
    description = "One of more of your stories does not contain a JIRA number"
  } else {
    description = "This story contains a JIRA number"
  }

  start := `{
    "state": "`
  middle := `",
    "description": "`
  end := `",
    "context": "JIRA/check"
  }`
  s := []string{start, state, middle, description, end}
  return strings.Join(s, "")
}
