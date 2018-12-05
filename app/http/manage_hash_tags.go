package http

import (
	"encoding/json"
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"io/ioutil"
	"net/http"
)

type ManageRequest struct {
	Group  string   `json:"group"`
	Name   string   `json:"tag"`
	Names  []string `json:"tags"`
	Status int      `json:"status"`
}

func manageHashTags(writer http.ResponseWriter, request *http.Request) {
	var (
		err   error
		body  []byte
		model = ManageRequest{}
	)

	if body, err = ioutil.ReadAll(request.Body); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(SearchResonse{
			Message: fmt.Sprintf("Cannot get request body: %s", err),
		})
		return
	}
	if err = json.Unmarshal(body, &model); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(SearchResonse{
			Message: fmt.Sprintf("Cannot parse json request: %s", err),
		})
		return
	}

	if request.Method == "GET" {
		tags, err := getHashTags(model)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(SearchResonse{
				Message: fmt.Sprintf("Problems while proccessing request: %s", err),
			})
			return
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(tags)
		return
	}

	if model.Group == "" || (model.Name == "" && len(model.Names) == 0) {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(SearchResonse{
			Message: fmt.Sprintf("Group and name params are required"),
		})
		return
	}
	if request.Method == "POST" {
		err = insertHashTag(model)
	} else if request.Method == "DELETE" {
		err = removeHashTag(model)
	} else if request.Method == "PUT" {
		err = updateHashTagStatus(model)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(SearchResonse{
			Message: fmt.Sprintf("Incorrect HTTP request command. Possible values are GET, POST, PUT"),
		})
		return
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(SearchResonse{
			Message: fmt.Sprintf("Problems while proccessing request: %s", err),
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(SearchResonse{
		Message: fmt.Sprintf("Ok"),
	})
	return
}

func getHashTags(request ManageRequest) (tags []clients.Tag, err error) {
	return clients.GetMongoDb().GetTagsByGroup(request.Group, request.Status)
}

func insertHashTag(request ManageRequest) (err error) {
	return clients.GetMongoDb().InsertHashTag(request.Group, request.Name, clients.ConfirmedTweets)
}

func updateHashTagStatus(request ManageRequest) (err error) {
	if request.Status != clients.ConfirmedTweets && request.Status != clients.IncorrectTweets {
		return fmt.Errorf("permitted to update status only to ConfirmedTweets=1 or IncorrectTweets=3")
	}
	if len(request.Names) == 0 && request.Name != "" {
		return clients.GetMongoDb().UpdateHashTagStatus(request.Group, request.Name, request.Status)
	}

	for _, tag := range request.Names {
		err = clients.GetMongoDb().UpdateHashTagStatus(request.Group, tag, request.Status)
		return
	}

	return clients.GetMongoDb().UpdateHashTagStatus(request.Group, request.Name, request.Status)
}

func removeHashTag(request ManageRequest) (err error) {
	return nil
	//return clients.GetMongoDb().RemoveHashTag(request.Group, request.Name, request.Status)
}
