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

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "basic-command",
			Description: "Basic command",
		},
		{
			Name:        "progress",
			Description: "Get grind75 progress",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "leetcode username",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
		"progress": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			username := i.ApplicationCommandData().Options[0].Value.(string)

			if username == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please provide a username.",
					},
				})
			} else {
				progress := getProgress(username)
				//build message string
				var progressStr strings.Builder
				for key, value := range progress {
					progressStr.WriteString(key + ": " + value + "\n")
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: progressStr.String(),
					},
				})
			}

		},
	}
)

func main() {
	var err error

	var Session, _ = discordgo.New("Bot " + os.Getenv("TOKEN")) //token here

	// Session.AddHandler(showProgress)

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	//adding commands
	Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	for _, v := range commands {
		_, err := Session.ApplicationCommandCreate(Session.State.User.ID, os.Getenv("DISCORD_GUILD"), v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	defer Session.Close()

	// Exit Normally.
}

// TODO: finish initializing the slice
var grind75List = []string{"Two Sum", "Valid Parentheses", "Merge Two Sorted Lists", "Best Time to Buy and Sell Stock", "Valid Palindrome", "Invert Binary Tree"}

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
	updateProgress(username, progressMap)
	return progressMap
}

func updateProgress(username string, progressMap map[string]string) {
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
}

//TODO:
//gui idea: pagination?
//show progress by pattern
