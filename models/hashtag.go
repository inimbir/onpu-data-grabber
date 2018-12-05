package models

import (
	"github.com/inimbir/onpu-data-grabber/app/clients"
)

type HashTag struct {
	Value string
}

func NewHashtag() *HashTag {
	return &HashTag{}
}

var hashTagsChainOfResponsibilities = &IsValidHandler{
	next: &IsNotExistsHandler{
		next: &PrepareDataHandler{
			next: &CheckSimilarityHandler{
				next: &ExtractHashtagsHandler{
					next: &SaveHashTagDataHandler{}}}}}}

func (tag HashTag) GetHandlers() HashTagHandler {
	return hashTagsChainOfResponsibilities
}

func (tag HashTag) Exists() bool {
	//@todo
	//return clients.GetMongoDb().Exists(tag)
	return true
}

//func (tag Group) Update() (bool, error) {
//	return clients.GetMongoDb().Update(tag)
//}

func (tag HashTag) GetCriterias(group string) (criterias []string, err error) {
	return clients.GetMongoDb().GetHashTagsByGroup(group, clients.ConfirmedTweets)
}
