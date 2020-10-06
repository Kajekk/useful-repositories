package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Repository struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	DefaultBranch  string    `json:"default_branch"`
	Stars          int       `json:"stargazers_count"`
	Forks          int       `json:"forks_count"`
	Issues         int       `json:"open_issues_count"`
	Created        time.Time `json:"created_at"`
	Updated        time.Time `json:"updated_at"`
	URL            string    `json:"html_url"`
	LastCommitDate time.Time `json:"-"`
}

const (
	tableHeader = `
	| Repo | Stars  | Forks  | Description |
	| ---- | :----: | :----: | ----------- |
	`
	footer = "\n*Last Update: %v*\n"
)

func main() {
	accessToken := os.Getenv("github_token")
	if accessToken == "" {
		accessToken = "731e6fcb50841bfdd9ba99fe90d171c024592177"
	}

	writeTableHeader()
	var repos []Repository

	bytesContent, err := ioutil.ReadFile("repo.list")
	if err != nil {
		fmt.Println(err)
	}
	re := regexp.MustCompile(`\r?\n`)
	input := re.ReplaceAllString(string(bytesContent), ",")
	lines := strings.Split(input, ",")
	fmt.Println(lines)
	for _, url := range lines {
		if strings.HasPrefix(url, "##") {
			continue
		}
		prefixGithub := "https://github.com/"
		idx := strings.Index(url, prefixGithub)
		if idx != -1 {
			urlGetRepoInfo := fmt.Sprintf("https://api.github.com/repos/%s?access_token=%s", url[len(prefixGithub):], accessToken)
			req, err := http.NewRequest("GET", urlGetRepoInfo, nil)
			if err != nil {
				fmt.Println(err)
				continue
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != 200 || resp.Body == nil {
				fmt.Println(err)
				continue
			}
			repository := Repository{}
			decoder := json.NewDecoder(resp.Body)
			if err = decoder.Decode(&repository); err != nil {
				fmt.Println(err)
			}

			//defer resp.Body.Close()
			//body, _ := ioutil.ReadAll(resp.Body)
			//repository := Repository{}
			//_ = json.Unmarshal(body, &repository)
			fmt.Println("response Body:", repository)
			repos = append(repos, repository)
		}

		if len(url) <= 1 {
			writeBody(repos)
			repos = nil
		}

	}
}

func writeTableHeader() {
	fmt.Println("write Header")
	file, err := os.OpenFile("README_TEMP.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.WriteString(tableHeader)
	if err != nil {
		fmt.Println(err)
	}
}

func writeBody(repos []Repository) {
	fmt.Println("write Body")
	file, err := os.OpenFile("README_TEMP.md", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	for _, repo := range repos {
		writeLine := fmt.Sprintf("| [%s](%s) | **%d** | **%d** | %s |\n", repo.Name, repo.URL, repo.Stars, repo.Forks, repo.Description)
		_, err = file.WriteString(writeLine)
		if err != nil {
			fmt.Println(err)
		}
	}
}
