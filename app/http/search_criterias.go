package http

import (
	"encoding/json"
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"github.com/inimbir/onpu-data-grabber/models"
	"io/ioutil"
	"net/http"
)

func PrintRed(text1 string, text string) {
	fmt.Println("\033[31m" + text1 + "\033[39m " + text + "\n")
}

func searchHashTags(w http.ResponseWriter, r *http.Request) {
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
			Message: fmt.Sprintf("Param 'type' is incorrect. It should be valid search group"),
		})
		return
	}
	PrintRed(fmt.Sprintf("Start processing hashtags for group '%s'", request.Group), "")
	var tags []clients.Tag
	if tags, err = clients.GetMongoDb().GetTagsByGroup(request.Group, clients.UnprocessedTweets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SearchResonse{
			Message: fmt.Sprintf("Cannot load group tags: %s", err),
		})
		return
	}

	logger := GetLog()

	for _, tag := range tags {
		a := models.HashTag{}
		m := models.HashTagContext{}
		m.Model = &a
		m.Value = tag.Name
		m.Group = request.Group
		a.GetHandlers().Handle(m)
		if err = a.GetHandlers().Handle(m); err != nil {
			logger.Println("err : " + err.Error())
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SearchResonse{
		Message: fmt.Sprintf("Ok"),
	})
	return

}
