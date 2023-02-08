package cli

import (
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", translateUrl, nil)

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	if err != nil {
		log.Fatal("There was with Request: %s", err)

	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("There was a problem with Do: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "You have been rate limited, try again later"
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(res.Body)

	if err != nil {
		log.Fatal("There was a problem with ParseJSONBuffer: %s", err)
	}

	nestOne, err := parsedJson.ArrayElement(0)

	if err != nil {
		log.Fatalf("There was a problem with ArrayElement nestOne: %s", err)
	}

	nestTwo, err := nestOne.ArrayElement(0)

	if err != nil {
		log.Fatalf("There was a problem with ArrayElement nestTwo: %s", err)
	}

	translatedStr, err := nestTwo.ArrayElement(0)

	if err != nil {
		log.Fatalf("There was a problem with ArrayElement translatedStr: %s", err)
	}

	str <- translatedStr.Data().(string)
	wg.Done()
}
