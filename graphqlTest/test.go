
// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// type QueryRequestBody struct {
// 	Query string `json:"query"`
// }

// type ResponseData struct {
// 	Data struct {
// 		RecentSubmissionList []struct {
// 			Title string `json:"title"`
// 		} `json:"recentAcSubmissionList"`
// 	} `json:"data"`
// }

// func main() {
// 	username := "ethand605"
// 	jsonData := map[string]string{
// 		"query": `
//             {
//                 recentAcSubmissionList(username: "` + username + `", limit: 200) {
// 					title
// 				}
//             }
//         `,
// 	}
// 	jsonValue, _ := json.Marshal(jsonData)
// 	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonValue))
// 	request.Header.Add("Content-Type", "application/json")
// 	client := &http.Client{}
// 	response, err := client.Do(request)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer response.Body.Close()
// 	data, _ := ioutil.ReadAll(response.Body)
// 	// fmt.Println(string(data))

// 	data_struct := ResponseData{}
// 	json.Unmarshal(data, &data_struct)

// 	submissionListStruct := data_struct.Data.RecentSubmissionList

// 	submissionList := make([]string, len(submissionListStruct))
// 	for i, submission := range submissionListStruct {
// 		submissionList[i] = submission.Title
// 	}

// }
