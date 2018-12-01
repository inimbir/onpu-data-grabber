package models

import (
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/algorithms"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"regexp"
)

type Handler interface {
	Handle(ctx Context) error
}

type IsValidHandler struct {
	next Handler
}

func NewHashtag() *HashTag {
	return &HashTag{}
}

func (h *IsValidHandler) Handle(ctx Context) (err error) {
	var IsValidHashTag = regexp.MustCompile(`^[1-9a-zA-Z_]+$`).MatchString
	if !IsValidHashTag(ctx.Value) {
		return fmt.Errorf("hashtag '%s' not valid, should contain only English characters or numbers or '_'", ctx.Value)
	} else if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type IsNotExistsHandler struct {
	next Handler
}

func (h *IsNotExistsHandler) Handle(ctx Context) (err error) {
	if exists, err := clients.GetMongoDb().ExistsHashTag(ctx.Group, ctx.Value); err != nil || exists == true {
		return fmt.Errorf("hashtag already exists, or error occured: %s", err)
	} else if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type PrepareDataHandler struct {
	next Handler
}

func (h *PrepareDataHandler) Handle(ctx Context) (err error) {
	tags, err := ctx.model.GetCriterias(ctx.Group)
	if err != nil {
		return fmt.Errorf("cannot extract hashtags for create initial data: %s", err)
	}
	if ctx.InitialData, err = clients.GetTwitter().GetTweetsByCriterias(tags, 100); err != nil {
		return fmt.Errorf("cannot extract tweets for create initial data: %s", err)
	}
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type CheckSimilarityHandler struct {
	next Handler
}

func (h *CheckSimilarityHandler) Handle(ctx Context) (err error) {
	var hashTags []string
	if ctx.ComparedData, err = clients.GetTwitter().GetTweetsByCriterias([]string{ctx.Value}, 100); err != nil {
		return fmt.Errorf("cannot check similarity: %s", err)
	}
	if hashTags, err = clients.GetMongoDb().GetHashTagsByGroup(ctx.Group, 0); err != nil {
		return fmt.Errorf("cannot extract trusted hashtags: %s", err)
	}
	//s := algorithms.Similarity{}
	//similarity := 0.0
	if similarity := (algorithms.Similarity{}).GetCosineSimilarity(
		hashTags,
		clients.GetTwitter().GetFullTextList(ctx.InitialData),
		clients.GetTwitter().GetFullTextList(ctx.ComparedData)); similarity < 0.9 {
		return fmt.Errorf("sililarity less than 0.9: %ff", similarity)
	}
	return h.next.Handle(ctx)
}

type ExtractHashtagsHandler struct {
	next Handler
}

func (h *ExtractHashtagsHandler) Handle(ctx Context) (err error) {
	tags := map[string]int8{}
	existedTags, err := clients.GetMongoDb().GetHashTagsByGroup(ctx.Group, -1)
	if err != nil {
		return fmt.Errorf("cannot get all existed hashtags: %s", err)

	}
	for _, tweet := range ctx.InitialData {
		for _, tag := range tweet.Entities.Hashtags {
			for _, exists := range existedTags {
				if tag.Text == exists {
					continue
				}
				tags[tag.Text] = 0
			}
		}
	}
	for _, tweet := range ctx.ComparedData {
		for _, tag := range tweet.Entities.Hashtags {
			for _, exists := range existedTags {
				if tag.Text == exists {
					continue
				}
				tags[tag.Text] = 0
			}
		}
	}
	for k := range tags {
		ctx.ExtractedTags = append(ctx.ExtractedTags, k)
	}
	//log.Println("tags")
	//log.Println(tags)
	return h.next.Handle(ctx)
}

type SaveHandler struct {
	next Handler
}

func (h *SaveHandler) Handle(ctx Context) (err error) {
	if err := clients.GetMongoDb().BulkInsert(ctx.Group, ctx.ExtractedTags, 0); err != nil {
		return fmt.Errorf("cannot save unkown hashtags: %s", err)

	}
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}
