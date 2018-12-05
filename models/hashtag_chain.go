package models

import (
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/algorithms"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"regexp"
	"strings"
)

type IsValidHandler struct {
	next HashTagHandler
}

func (h *IsValidHandler) Handle(ctx HashTagContext) (err error) {
	//PrintRed("Start processing tag: ", ctx.Value)
	var IsValidHashTag = regexp.MustCompile(`^[0-9a-zA-Z_]+$`).MatchString
	if !IsValidHashTag(ctx.Value) {
		return fmt.Errorf("hashtag '%s' not valid, should contain only English characters or numbers or '_'", ctx.Value)
	} else if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type IsNotExistsHandler struct {
	next HashTagHandler
}

func (h *IsNotExistsHandler) Handle(ctx HashTagContext) (err error) {
	//if exists, err := clients.GetMongoDb().ExistsHashTag(ctx.Group, ctx.Value); err != nil || exists == true {
	//	return fmt.Errorf("hashtag already exists, or error occured: %s", err)
	//} else
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type PrepareDataHandler struct {
	next HashTagHandler
}

func (h *PrepareDataHandler) Handle(ctx HashTagContext) (err error) {
	tags, err := clients.GetMongoDb().GetHashTagsByGroup(ctx.Group, clients.ConfirmedTweets)
	if err != nil {
		return fmt.Errorf("cannot extract hashtags for create initial data: %s", err)
	}
	if ctx.InitialData, err = clients.GetTwitter().GetTweetsByCriterias(tags, 100); err != nil {
		return fmt.Errorf("cannot extract tweets for create initial data: %s", err)
	}
	//PrintRed("Current trusted tags for compare data:", strings.Join(tags, ", "))

	if h.next != nil {
		return h.next.Handle(ctx)
	}

	return
}

type CheckSimilarityHandler struct {
	next HashTagHandler
}

func (h *CheckSimilarityHandler) Handle(ctx HashTagContext) (err error) {
	var hashTags []string
	if ctx.ComparedData, err = clients.GetTwitter().GetTweetsByCriterias([]string{ctx.Value}, 100); err != nil {
		return fmt.Errorf("cannot check similarity: %s", err)
	}
	if hashTags, err = clients.GetMongoDb().GetHashTagsByGroup(ctx.Group, clients.ConfirmedTweets); err != nil {
		return fmt.Errorf("cannot extract trusted hashtags: %s", err)
	}
	//s := algorithms.Similarity{}
	similarity := 0.0
	if similarity = algorithms.Get().GetCosineSimilarity(
		hashTags,
		clients.GetTwitter().GetFullTextList(ctx.InitialData),
		clients.GetTwitter().GetFullTextList(ctx.ComparedData),
	); similarity < 0.9 {

		if err := clients.GetMongoDb().UpdateHashTagStatus(ctx.Group, ctx.Value, clients.IncorrectTweets); err != nil {
			return fmt.Errorf("cannot update status for incorrect hashag: %s", err)
		}
		//PrintRed(fmt.Sprintf("Similatiry of tag '%s' with trusted:", ctx.Value), fmt.Sprintf("%.2f", similarity))
		fmt.Printf("Similatiry of tag '%s' with trusted: %.3f\n", ctx.Value, similarity)
		return fmt.Errorf("similarity less than 0.9: %f\n", similarity)
	}

	if err := clients.GetMongoDb().UpdateHashTagStatus(ctx.Group, ctx.Value, clients.UnconfirmedTweets); err != nil {
		return fmt.Errorf("hashtag already exists, or error occured: %s", err)
	}

	PrintRed(fmt.Sprintf("Similatiry of tag '%s' with trusted:", ctx.Value), fmt.Sprintf("%.2f", similarity))

	return h.next.Handle(ctx)
}

type ExtractHashtagsHandler struct {
	next HashTagHandler
}

func (h *ExtractHashtagsHandler) Handle(ctx HashTagContext) (err error) {
	tags := map[string]int8{}
	existedTags, err := clients.GetMongoDb().GetHashTagsByGroup(ctx.Group, clients.AllTweets)
	if err != nil {
		return fmt.Errorf("cannot get all existed hashtags: %s", err)

	}
	for _, tweet := range ctx.InitialData {
		for _, tag := range tweet.Entities.Hashtags {
			exist := false
			tagValue := strings.ToLower(tag.Text)
			for _, exists := range existedTags {
				if tagValue == exists {
					exist = true
					break
				}
			}
			if exist == false {
				tags[tagValue] = 0
			}
		}
	}
	for _, tweet := range ctx.ComparedData {
		for _, tag := range tweet.Entities.Hashtags {
			exist := false
			tagValue := strings.ToLower(tag.Text)
			for _, exists := range existedTags {
				if tagValue == exists {
					exist = true
					break
				}
			}
			if exist == false {
				tags[tagValue] = 0
			}
		}
	}
	for k := range tags {
		ctx.ExtractedTags = append(ctx.ExtractedTags, k)
	}
	return h.next.Handle(ctx)
}

type SaveHashTagDataHandler struct {
	next HashTagHandler
}

func (h *SaveHashTagDataHandler) Handle(ctx HashTagContext) (err error) {
	if err := clients.GetMongoDb().BulkInsert(ctx.Group, ctx.ExtractedTags, clients.UnprocessedTweets); err != nil {
		return fmt.Errorf("cannot save unkown hashtags: %s", err)
	}
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}
