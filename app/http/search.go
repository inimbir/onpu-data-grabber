package http

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"github.com/inimbir/onpu-data-grabber/models"
	"io/ioutil"
	"net/http"
	"strings"
)

func search(w http.ResponseWriter, r *http.Request) {
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
	request := SearchRequest{}
	if err = json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Cannot parse json request: %s", err),
		})
		return
	}
	if request.Type != FullSearch && request.Type != TagSearch {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Type param incorrect. Possible type values are: 1-FullSearch, 2-TagSearch"),
		})
		return
	}
	if request.Type == FullSearch {
		if err = clients.GetRabbitMq().PushTask("tasks", FullSearch); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(SearchResonse{
				Message: fmt.Sprintf("Cannot add task to queue"),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Task added to queue"),
		})
		return
	}

	if request.Group == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Param 'type' is incorrect. It sholud be valid search group"),
		})
		return
	}

	tags := []clients.Tag{}
	if tags, err = clients.GetMongoDb().GetTagsByGroup(request.Group, clients.ConfirmedTweets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Cannot load group tags: %s", err),
		})
		return
	}

	PrintRed("Testing information technology", "")
	logger := GetLog()

	for _, tag := range tags {
		PrintRed(fmt.Sprintf("Start processing query '%s' through date range: [2018-11-30T12:00:00+00:00, 2018-11-30T23:59:59+00:00]", tag.Name), "")
		if tag.Name == "supernatural" {
			PrintRed(fmt.Sprintf("Number of found tweets for query '%s':", tag.Name), fmt.Sprintf("%d", 0))
			continue
		}

		var tweets []twitter.Tweet
		if tweets, err = clients.GetTwitter().GetTweetsFromLastId(tag.Name, 0); err != nil {
			PrintRed("Problem", err.Error())
		}
		model := models.Tweet{}
		success := 0
		for _, tweet := range tweets {
			m := models.TweetContext{}
			m.Tweet = tweet
			m.Model = &model
			m.Value = tweet.FullText
			m.Group = request.Group
			m.Id = tweet.ID
			if err = model.GetHandlers().Handle(m); err == nil {
				success += 1
			} else if err != nil {
				logger.Println("err : " + err.Error())
			}
		}
		PrintRed(fmt.Sprintf("Number of found tweets for query '%s':", tag.Name), fmt.Sprintf("%d", success))
	}

	hashtags, _ := clients.GetMongoDb().GetHashTagsByGroup(request.Group, clients.UnprocessedTweets)
	PrintRed("Found hashtags: ", strings.Join(hashtags, ", "))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SearchResonse{
		Message: fmt.Sprintf("Ok"),
	})
	return
}
