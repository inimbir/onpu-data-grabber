package clients

import (
	"fmt"
	api "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"math/rand"
	"sync"
)

type Twitter struct {
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
	client         *api.Client
}

const (
	Lang            = "en"
	Count           = 100
	IncludeEntities = true
)

var initTwitterOnce sync.Once

func (m Twitter) Get() *Twitter {
	//initTwitterOnce.Do(func() {
	twConfig := oauth1.NewConfig(m.consumerKey, m.consumerSecret)
	token := oauth1.NewToken(m.accessToken, m.accessSecret)
	httpClient := twConfig.Client(oauth1.NoContext, token)
	m.client = api.NewClient(httpClient)
	//})
	return &m
}

func (m Twitter) GetTweetsByCriterias(tags []string, maxNumberOfTweets int) (tweets []api.Tweet, err error) {
	var (
		maxId        int64 = 0
		tweetsModels []api.Tweet
	)
	for n := rand.Int() % len(tags); n >= 0 && len(tweets) < maxNumberOfTweets; {
		if tweetsModels, err = m.getTweets(tags[n], maxId, 0); err != nil {
			return tweets, fmt.Errorf("twitter error: %s", err)
		}
		if len(tweets) > 0 {
			maxId = tweetsModels[len(tweetsModels)-1].ID - 1
		}
		tweets = append(tweets, tweetsModels...)
		//tweets = append(tweets, GetFullTextList(tweetsModels)...)
	}
	return tweets[:maxNumberOfTweets], err
}

func (m Twitter) GetTweetsFromLastId(tag string, sinceId int64) (tweets []api.Tweet, err error) {
	var (
		maxId        int64 = 0
		tweetsModels []api.Tweet
	)
	for {
		if tweetsModels, err = m.getTweets(tag, maxId, int64(sinceId)); err != nil {
			return tweets, fmt.Errorf("twitter error: %s", err)
		}
		if len(tweetsModels) == 0 {
			break
		}
		maxId = tweetsModels[len(tweetsModels)-1].ID - 1
		tweets = append(tweets, tweetsModels...)
	}
	return tweets, err
}

//
func (m Twitter) getTweets(tag string, maxId int64, sinceId int64) (tweets []api.Tweet, err error) {
	params := &api.SearchTweetParams{
		Query:           tag,
		Lang:            Lang,
		Count:           Count,
		TweetMode:       "extended",
		IncludeEntities: func(i bool) *bool { return &i }(IncludeEntities),
	}

	if maxId != 0 {
		params.MaxID = maxId
	}
	if sinceId != 0 {
		params.SinceID = sinceId
	}
	search, _, err := m.client.Search.Tweets(params)
	if err != nil {
		return tweets, fmt.Errorf("failed to load tweets: %s", err)
	}
	return search.Statuses, err
}

func (m Twitter) GetFullTextList(tweets []api.Tweet) []string {
	var list []string
	for _, tweet := range tweets {
		list = append(list, tweet.FullText)
	}
	return list
}

//var lastId int64 = 0
//for true {
//	tweets := clients.GetTweets(tag, lastId)
//	for _, ind := range tweets {
//		log.Println("tid:  ", ind.ID)
//		log.Println("tcr:  ", ind.CreatedAt)
//		log.Println("tft:  ", ind.Text)
//		lastId = ind.ID - 1
//		clients.SendToQueue(ind.ID, ind.CreatedAt, ind.Text)
//	}
//}
