package models

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/inimbir/onpu-data-grabber/app/clients"
)

type HashTag struct {
	Value string
}

type Context struct {
	Value         string
	Group         string
	model         HashTag
	InitialData   []twitter.Tweet
	ComparedData  []twitter.Tweet
	ExtractedTags []string
}

var hashTagsChainOfResponsibilities = &IsValidHandler{
	next: &IsNotExistsHandler{
		next: &PrepareDataHandler{
			next: &CheckSimilarityHandler{
				next: &ExtractHashtagsHandler{
					next: &SaveHandler{}}}}}}

func (tag HashTag) GetHandlers() Handler {
	return hashTagsChainOfResponsibilities
}

func (tag HashTag) Exists() bool {
	//return clients.GetMongoDb().Exists(tag)
	return true
}

//func (tag HashTag) Update() (bool, error) {
//	return clients.GetMongoDb().Update(tag)
//}

func (tag HashTag) GetCriterias(group string) (criterias []string, err error) {
	return clients.GetMongoDb().GetHashTagsByGroup(group, 0)
}
