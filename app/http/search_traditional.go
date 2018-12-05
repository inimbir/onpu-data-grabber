package http

import (
	"encoding/json"
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"io/ioutil"
	"log"
	"net/http"
)

func searchTraditional(w http.ResponseWriter, r *http.Request) {

	type traditional struct {
		Tag string `json:"tag"`
	}

	var (
		err  error
		body []byte
	)
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Cannot get request body: %s", err),
		})
	}
	request := traditional{}
	if err = json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Cannot parse json request: %s", err),
		})
		return
	}

	PrintRed("Testing traditional approach", "")
	PrintRed(fmt.Sprintf("Start processing query '%s' through date range: [2018-11-30T12:00:00+00:00, 2018-11-30T23:59:59+00:00]", request.Tag), "")
	var lastId int64 = 0
	tweets, err := clients.GetTwitter().GetTweetsFromLastId(request.Tag, lastId)
	if err != nil {
		log.Println(err)
	}
	PrintRed(fmt.Sprintf("Number of found tweets for query '%s':", request.Tag), fmt.Sprintf("%d", len(tweets)))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SearchResonse{
		Message: fmt.Sprintf("Ok"),
	})
	return
}
