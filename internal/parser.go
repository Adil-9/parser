package internal

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

var (
	IDRegex             = `<span data-v-2e6a30b8>(\d+)</span>`
	AliasRegex          = `<div class="contributor__name-content" data-v-c5a99f5a>([^<]+)</div>`
	NameRegex           = `<div class="contributor__title" data-v-c5a99f5a>(.*?)<!----><!----></div>`
	CategoryRegexOutter = `<div class="row-cell category" data-v-2e6a30b8>(.*?)<\/div><div class="row-cell subscribers"`
	CategoryRegexInner  = `<div class="row-cell category" data-v-2e6a30b8>(.*?)<\/div>`
	FollowersRegex      = `<div class="row-cell subscribers" data-v-2e6a30b8>(.*?)<\/div>`
	CountryRegex        = `<div class="row-cell audience" data-v-2e6a30b8 data-v-e1ea9c14>(.*?)</div>`
	EngAuthRegex        = `<div class="row-cell authentic" data-v-2e6a30b8 data-v-e1ea9c14>(.*?)</div>`
	EngAvgRegex         = `<div class="row-cell engagement" data-v-2e6a30b8 data-v-e1ea9c14>(.*?)</div>`
)

func Parse() {
	var data Users

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	loadEnv()
	link := os.Getenv("API_KEY")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		log.Fatal("request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Body")
	}

	bodyString := string(body)

	data.ID = reg(bodyString, IDRegex)
	data.Alias = reg(bodyString, AliasRegex)
	data.Name = reg(bodyString, NameRegex)
	// findCategory(bodyString)
	data.Followers = reg(bodyString, FollowersRegex)
	data.Country = reg(bodyString, CountryRegex)
	data.EngAuth = reg(bodyString, EngAuthRegex)
	data.EngAvg = reg(bodyString, EngAvgRegex)

	// fmt.Println(data)
	file, err := os.OpenFile("data.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Alias", "Name", "Category", "Followers", "Country", "EngAuth", "EngAvg"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		row := []string{data.ID[i], data.Alias[i], data.Name[i], "null", data.Followers[i], data.Country[i], data.EngAuth[i], data.EngAvg[i]}
		err := writer.Write(row)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("done")
}

// func findCategory(bodyString string) [][]string {
// 	InnerSlice := make([]string, 0, 5)
// 	Slice := make([][]string, 0, 50)

// 	re := regexp.MustCompile(CategoryRegexOutter)
// 	matches := re.FindAllStringSubmatch(bodyString, -1)

// 	for _, match := range matches {
// 		if len(match) > 1 {
// 			InnerSlice = append(InnerSlice, match[1])
// 		}
// 	}

// 	for i, v := range InnerSlice {
// 		re := regexp.MustCompile(CategoryRegexInner)
// 		matches = re.FindAllStringSubmatch(v, -1)
// 		for _, match := range matches {
// 			if len(match) > 1 {
// 				Slice[i] = append(Slice[i], match[1])
// 			}
// 		}
// 	}

// 	fmt.Println(Slice)

// 	return Slice
// }

func reg(bodyString string, RegExt string) []string {
	Slice := make([]string, 0, 50)

	re := regexp.MustCompile(RegExt)
	matches := re.FindAllStringSubmatch(bodyString, -1)

	for _, match := range matches {
		if len(match) > 1 {
			Slice = append(Slice, match[1])
		}
	}

	return Slice
}

func loadEnv() { // load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// log.Fatal("Error loading .env file")
	}
}
