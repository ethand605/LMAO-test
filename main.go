package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

func main() {
	var err error

	var Session, _ = discordgo.New("Bot MTAwMTkyMzk0NDQ0MzU2MDAwNg.GtKd2O.rrrmCR6fcRzYV8diufPDG1wGWW5AvaeVVa2OsQ") //token here

	Session.AddHandler(showProgress)

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}

func showProgress(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	//command format: !progress username
	res := strings.Fields(m.Content)
	command := res[0]
	username := res[1]

	if command == "!progress" {
		if username == "" {
			s.ChannelMessageSend(m.ChannelID, "Please provide a username.")
		}
		progress := getProgress(username)
		var progressStr strings.Builder
		for key, value := range progress {
			progressStr.WriteString(key + ": " + value + "\n")
		}
		s.ChannelMessageSend(m.ChannelID, progressStr.String())

	}
}

//TODO: finish initializing the slice
var grind75List = []string{"Two Sum", "Valid Parentheses", "Merge Two Sorted Lists", "Best Time to Buy and Sell Stock", "Valid Palindrome", "Invert Binary Tree", "
Valid Anagram", "Binary Search"}

type ResponseData struct {
	Data struct {
		RecentSubmissionList []struct {
			Title string `json:"title"`
		} `json:"recentAcSubmissionList"`
	} `json:"data"`
}

func getProgress(username string) map[string]string {
	//TODO: add database calls here
	progressMap := make(map[string]string)
	for _, problem := range grind75List {
		progressMap[problem] = "❌"
	}
	return updateProgress(username, progressMap) 
}

func updateProgress(username string, progressMap map[string]string) map[string]string {
	query := map[string]string{
		"query": `
            { 
                recentAcSubmissionList(username: "` + username + `", limit: 20) {
					title
				}
            }
        `,
	}
	queryAsJson, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(queryAsJson))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	data_struct := ResponseData{}
	json.Unmarshal(data, &data_struct)

	submissionListStruct := data_struct.Data.RecentSubmissionList

	submissionList := make([]string, len(submissionListStruct))
	for i, submission := range submissionListStruct {
		submissionList[i] = submission.Title
	}

	for _, problem := range submissionList {
		_, isPresent := progressMap[problem]
		if isPresent {
			progressMap[problem] = "✅"
		}
	}
	return progressMap
}

//TODO:
//gui idea: pagination?
//show progress by pattern
