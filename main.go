package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

var args struct {
	URL string
}

const s3Prefix = "https://s3-us-west-2.amazonaws.com/surveycake-s3.surveycakecdn.com/json/"

type Sheet struct {
	ID            string    `json:"ID"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	Language      string    `json:"language"`
	SnFinal       string    `json:"sn_final"`
	Limitcount    string    `json:"limitcount"`
	Welcometext   string    `json:"welcometext"`
	Welcomebanner string    `json:"welcomebanner"`
	Thankyoutext  string    `json:"thankyoutext"`
	SubjectList   []Subject `json:"subjects"`
}

type Option struct {
	Text string `json:"text"`
}

type Subject struct {
	Text       string   `json:"text"`
	Type       string   `json:"type"`
	OptionList []Option `json:"options"`
}

func topicTransformer(rawString string) string {

	if rawString == "分隔線/分頁" {
		return "================= 分隔線/分頁 =================" + "\n"
	}

	rawString = strings.Replace(rawString, "<p>", "", -1)
	rawString = strings.Replace(rawString, "<strong>", "", -1)
	rawString = strings.Replace(rawString, "</strong>", "", -1)
	rawString = strings.Replace(rawString, "</p>", "", -1)

	return rawString
}

// func optionTransformer(rawString string) string {

// }

func dataTransformer(data Subject) string {

	var output = topicTransformer(data.Text)
	if data.Type == "STATEMENT" {
		output = "<< " + output + " >> \n"
		return output
	}
	output = "* " + output + "\n"

	if len(data.OptionList) != 0 {
		for _, op := range data.OptionList {
			output += "\t - " + op.Text + "\n"
		}
	}

	return output
}

func main() {

	arg.MustParse(&args)
	var url string
	var s3Url string
	var s3Id string
	var jsonObj Sheet

	if args.URL == "" {
		fmt.Println("Input your survey cake link: ")
		fmt.Scanln(&url)
	} else {
		url = args.URL
	}

	urlElements := strings.Split(url, "/")
	s3Id = urlElements[len(urlElements)-1] + ".json"
	s3Url = s3Prefix + s3Id

	fmt.Println(s3Url)

	resp, err := http.Get(s3Url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &jsonObj)

	lists := jsonObj.SubjectList
	fmt.Println(len(lists))

	os.Remove("data.txt")
	f, err := os.Create("data.txt")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	defer f.Close()

	for _, data := range lists {

		if data.Type == "STATEMENT" {

		}
		_, err = f.WriteString(dataTransformer(data))
		if err != nil {
			log.Fatal("Error", err)
		}

	}

}
