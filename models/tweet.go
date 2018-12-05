package models

type Tweet struct {
	Value string
}

func NewTweet() *HashTag {
	return &HashTag{}
}

var tweetChainOfResponsibilities = &IsEnglishHandler{
	next: &IsNotProcessedHandler{
		next: &ExtractHashTagsFromTweetHandler{
			next: &SaveTweetDataHandler{
				next: &PushToQueueHandler{}}}}}

func (tag Tweet) GetHandlers() TweetHandler {
	return tweetChainOfResponsibilities
}
