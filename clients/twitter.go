package clients

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	Lang            = "en"
	Count           = 100
	IncludeEntities = true
)

func GetTweets(tag string, sinceId int64) []twitter.Tweet {
	twConfig := oauth1.NewConfig(config.TwitterConsumerKey, config.TwitterConsumerSecret)
	token := oauth1.NewToken(config.TwitterAccessToken, config.TwitterAccessSecret)
	httpClient := twConfig.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Search Tweets
	params := &twitter.SearchTweetParams{
		Query:           tag,
		Lang:            Lang,
		Count:           Count,
		IncludeEntities: func(i bool) *bool { return &i }(IncludeEntities),
		SinceID:         sinceId,
	}

	// Search Tweets
	search, _, err := client.Search.Tweets(params)
	failOnError(err, "Failed to send search request")
	return search.Statuses
}
