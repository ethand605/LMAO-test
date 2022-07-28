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

	var Session, _ = discordgo.New("Bot ") //token here

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
		progress := queryProblems(username)
		var progressStr strings.Builder
		for key, value := range progress {
			progressStr.WriteString(key + ": " + value + "\n")
		}
		s.ChannelMessageSend(m.ChannelID, "Here is the progress of "+username+":")
		s.ChannelMessageSend(m.ChannelID, progressStr.String())

	}
}

var blindList = [5]string{"Insert Interval", "House Robber", "K Closest Points to Origin", "Two Sum", "Basic Calculator"}

type QueryRequestBody struct {
	Query string `json:"query"`
}

type ResponseData struct {
	Data struct {
		RecentSubmissionList []struct {
			Title string `json:"title"`
		} `json:"recentAcSubmissionList"`
	} `json:"data"`
}

func queryProblems(username string) map[string]string {
	progressMap := make(map[string]string)
	for _, problem := range blindList {
		progressMap[problem] = "❌"
	}

	jsonData := map[string]string{
		"query": `
            { 
                recentAcSubmissionList(username: "` + username + `", limit: 2000) {
					title
				}
            }
        `,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonValue))
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	// fmt.Println(string(data))

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
