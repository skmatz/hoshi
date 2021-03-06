package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/kyokomi/emoji"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"github.com/tcnksm/go-gitconfig"
	"github.com/lithammer/fuzzysearch/fuzzy"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

// These variables are set in build step
var (
	Version = "unset"
)

// Option represents application options
type Option struct {
	Order   string `short:"o" long:"order" description:"Change item order"`
	Reverse bool   `short:"r" long:"reverse" description:"Reverse item order"`
	Version bool   `short:"v" long:"version" description:"Show hoshi version"`
}

// Star represents a stared repository
type Star struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	PushedAt         time.Time   `json:"pushed_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorURL        interface{} `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
	Forks         int    `json:"forks"`
	OpenIssues    int    `json:"open_issues"`
	Watchers      int    `json:"watchers"`
	DefaultBranch string `json:"default_branch"`
}

func getStarsPage(user string, page int) []Star {
	url := fmt.Sprintf("https://api.github.com/users/%s/starred?page=%d", user, page)

	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	var stars []Star
	json.Unmarshal(body, &stars)

	if len(stars) == 0 {
		return nil
	}

	return stars
}

func sortItemOrder(stars []Star, order string, reverse bool) []Star {
	reversed := func(stars []Star) []Star {
		for i := len(stars)/2 - 1; i >= 0; i-- {
			j := len(stars) - 1 - i
			stars[i], stars[j] = stars[j], stars[i]
		}
		return stars
	}

	switch order {
	case "added-at":
	case "author-name":
		sort.Slice(stars, func(i, j int) bool {
			return strings.ToLower(stars[i].Owner.Login) < strings.ToLower(stars[j].Owner.Login)
		})
	case "created-at":
		sort.Slice(stars, func(i, j int) bool {
			return stars[i].CreatedAt.After((stars[j].CreatedAt))
		})
	case "repository-name":
		sort.Slice(stars, func(i, j int) bool {
			return strings.ToLower(stars[i].Name) < strings.ToLower(stars[j].Name)
		})
	case "updated-at":
		sort.Slice(stars, func(i, j int) bool {
			return stars[i].UpdatedAt.After((stars[j].UpdatedAt))
		})
	}

	if reverse {
		stars = reversed(stars)
	}

	return stars
}

func run(args []string) int {
	var opt Option
	args, err := flags.ParseArgs(&opt, args)

	if err != nil {
		return 2
	}

	if opt.Version {
		fmt.Printf("hoshi v%s\n", Version)
		return 0
	}

	var stars []Star

	user, err := gitconfig.GithubUser()
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; ; i++ {
		if s := getStarsPage(user, i); s != nil {
			stars = append(stars, s...)
		} else {
			break
		}
	}

	terminalWidth, err := terminal.Width()
	if err != nil {
		return 1
	}

	// to prevent a multi-line description
	for i := 0; i < len(stars); i++ {
		if len(stars[i].Description) > int(terminalWidth)-22 {
			p := &stars[i]
			p.Description = p.Description[:int(terminalWidth)-25] + "..."
		}
	}

	stars = sortItemOrder(stars, opt.Order, opt.Reverse)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   emoji.Sprint(":star: {{ .FullName | cyan }}"),
		Inactive: "    {{ .FullName | cyan }}",
		Selected: emoji.Sprint(":sparkles: {{ .FullName | red | cyan }}"),
		Details: `------------------------------
{{ "\U0001F4D2 description" }}	{{ .Description }}
{{ "\U0001F3E0 homepage" }}	{{ .Homepage }}
{{ "\U0001F4DD language" }}	{{ .Language }}
{{ "\U0001F4C3 license" }}	{{ .License.Name }}
{{ "\U0001F31F stars" }}	{{ .StargazersCount }}`,
	}

	keys := &promptui.SelectKeys{
		Next:     promptui.Key{Code: promptui.KeyNext, Display: promptui.KeyNextDisplay},
		Prev:     promptui.Key{Code: promptui.KeyPrev, Display: promptui.KeyPrevDisplay},
		PageUp:   promptui.Key{Code: promptui.KeyBackward, Display: promptui.KeyBackwardDisplay},
		PageDown: promptui.Key{Code: promptui.KeyForward, Display: promptui.KeyForwardDisplay},
		Search:   promptui.Key{Code: 63, Display: "?"}, // 63 is rune for "?"
	}

	searcher := func(input string, index int) bool {
		star := stars[index]
		name := strings.Replace(strings.ToLower(star.FullName), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return fuzzy.Match(input, name)
	}

	prompt := promptui.Select{
		Label:             "Stars",
		Items:             stars,
		Size:              10,
		HideHelp:          true,
		Templates:         templates,
		Keys:              keys,
		Searcher:          searcher,
		StartInSearchMode: true,
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	browser.OpenURL(stars[index].HTMLURL)

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
