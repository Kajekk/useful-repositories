package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
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
	Updated        string    `json:"updated_at"`
	URL            string    `json:"html_url"`
	LastCommitDate time.Time `json:"-"`
}

const (
	header = `
- [Web Frameworks](#web-frameworks)
`

	tableHeader = `
| Repo | Stars  | Forks  | Issues | Description | Last Updated |
| ---- | :----: | :----: | :----: | ----------- | :----------: |
`
	footer = "\n*Last Update: %v*\n"
)

func main() {
	accessToken := getAccessToken()
	if accessToken == "" {
		_ = fmt.Errorf("Please provide valid access token")
		os.Exit(1)
	}

	writeHeader()
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
			writeTableHeader(url)
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

func getAccessToken() string {
	tokenBytes, err := ioutil.ReadFile("access-token.tok")
	if err != nil {
		fmt.Println(err)
		token := os.Getenv("github_token")
		return token
	}

	return strings.TrimSpace(string(tokenBytes))
}

func writeHeader() {
	fmt.Println("write Header")
	file, err := os.OpenFile("README_TEMP.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.WriteString(header)
	if err != nil {
		fmt.Println(err)
	}
}

func writeTableHeader(category string) {
	fmt.Println("write Table Header")
	file, err := os.OpenFile("README_TEMP.md", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	writeLine := fmt.Sprintf("\n%s\n%s", category, tableHeader)
	_, err = file.WriteString(writeLine)
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
	sort.Slice(repos[:], func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})

	defer file.Close()
	for _, repo := range repos {
		writeLine := fmt.Sprintf("| [%s](%s) | **%d** | **%d** | **%d** | %s | %s |\n", repo.Name, repo.URL, repo.Stars, repo.Forks, repo.Issues, repo.Description, repo.Updated)
		_, err = file.WriteString(writeLine)
		if err != nil {
			fmt.Println(err)
		}
	}
}
