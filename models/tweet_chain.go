package models

import (
	"fmt"
	"github.com/inimbir/onpu-data-grabber/app/clients"
	"strings"
)

type IsEnglishHandler struct {
	next TweetHandler
}

func (h *IsEnglishHandler) Handle(ctx TweetContext) (err error) {
	//@todo handle
	//var IsValidHashTag = regexp.MustCompile(`^[1-9a-zA-Z]+$`).MatchString
	//if !IsValidHashTag(ctx.Value) {
	//	return fmt.Errorf("Tweet content '%s' not valid, should contain only English characters or numbers", ctx.Value)
	//} else
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type IsNotProcessedHandler struct {
	next TweetHandler
}

func (h *IsNotProcessedHandler) Handle(ctx TweetContext) (err error) {
	//if exists, err := clients.GetMongoDb().ExistsTweet(ctx.Group, ctx.Id); err != nil || exists == true {
	//	return fmt.Errorf("tweet already exists, or error occured: %s", err)
	//} else
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type ExtractHashTagsFromTweetHandler struct {
	next TweetHandler
}

func (h *ExtractHashTagsFromTweetHandler) Handle(ctx TweetContext) (err error) {
	tags := map[string]int8{}
	existedTags, err := clients.GetMongoDb().GetHashTagsByGroup(ctx.Group, clients.AllTweets)
	if err != nil {
		return fmt.Errorf("cannot get all existed hashtags: %s", err)
	}
	for _, tag := range ctx.Tweet.Entities.Hashtags {
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
	for k := range tags {
		ctx.ExtractedTags = append(ctx.ExtractedTags, k)
	}
	//PrintRed(fmt.Sprintf("Extracted hashtags for tweet ID '%d':", ctx.Tweet.ID), strings.Join(ctx.ExtractedTags, ", "))
	if h.next != nil {
		return h.next.Handle(ctx)
	}

	return
}

type SaveTweetDataHandler struct {
	next TweetHandler
}

func (h *SaveTweetDataHandler) Handle(ctx TweetContext) (err error) {
	if err := clients.GetMongoDb().BulkInsert(ctx.Group, ctx.ExtractedTags, clients.UnprocessedTweets); err != nil {
		return fmt.Errorf("cannot save unkown hashtags: %s", err)
	}
	if h.next != nil {
		return h.next.Handle(ctx)
	}
	return
}

type PushToQueueHandler struct {
	next TweetHandler
}

func (h *PushToQueueHandler) Handle(ctx TweetContext) (err error) {
	//if err := clients.GetRabbitMq().Push(ctx.Group, ctx.Tweet.ID, ctx.Tweet.CreatedAt, ctx.Tweet.FullText); err != nil {
	//	return fmt.Errorf("cannot save unkown hashtags: %s", err)
	//}
	if err := clients.GetMongoDb().AddTweetIdToGroup(ctx.Group, []int64{ctx.Id}); err != nil {
		return fmt.Errorf("cannot save unkown hashtags: %s", err)
	}
	//if h.next != nil {
	//	return h.next.Handle(ctx)
	//}
	return
}

//@todo добавить сохранение хештегов и созранение статуса в диаграмму
