package models

import "github.com/dghubble/go-twitter/twitter"

type BaseContext struct {
	Value         string
	Group         string
	ExtractedTags []string
}

type HashTagContext struct {
	BaseContext
	Model        *HashTag
	InitialData  []twitter.Tweet
	ComparedData []twitter.Tweet
}

type TweetContext struct {
	BaseContext
	Id    int64
	Model *Tweet
	Tweet twitter.Tweet
	//InitialData   []twitter.Tweet
	//ComparedData  []twitter.Tweet
	//ExtractedTags []string
}

type HashTagHandler interface {
	Handle(ctx HashTagContext) error
}

type TweetHandler interface {
	Handle(ctx TweetContext) error
}
