package clients

import (
	"fmt"
	api "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"math/rand"
	"sync"
	"time"
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

const (
	AllTweets             = -1
	UnprocessedTweets     = 0
	UnprocessedTweetsText = "Unprocessed"
	ConfirmedTweets       = 1
	ConfirmedTweetsText   = "Confirmed"
	UnconfirmedTweets     = 2
	UnconfirmedTweetsText = "Unconfirmed"
	IncorrectTweets       = 3
	IncorrectTweetsText   = "Incorrect"
)

var (
	twitterInstance *Twitter
	initTwitterOnce sync.Once
)

func (m Twitter) Get() *Twitter {
	initTwitterOnce.Do(func() {
		twConfig := oauth1.NewConfig(m.consumerKey, m.consumerSecret)
		token := oauth1.NewToken(m.accessToken, m.accessSecret)
		httpClient := twConfig.Client(oauth1.NoContext, token)
		m.client = api.NewClient(httpClient)
		twitterInstance = &m
	})
	return twitterInstance
}

func (m Twitter) GetTweetsByCriterias(tags []string, maxNumberOfTweets int) (tweets []api.Tweet, err error) {
	var (
		maxId        int64 = 0
		tweetsModels []api.Tweet
	)
	n := rand.Int() % len(tags)
	for n >= 0 && len(tweets) < maxNumberOfTweets && len(tags) > 0 {
		if tweetsModels, err = m.getTweets(tags[n], maxId, 0); err != nil {
			return tweets, fmt.Errorf("twitter error: %s", err)
		}
		if len(tweetsModels) > 0 {
			maxId = tweetsModels[len(tweetsModels)-1].ID - 1
		}
		tweets = append(tweets, tweetsModels...)
		if len(tweetsModels) < maxNumberOfTweets {
			//log.Printf("only %d tweets exists for tag '%s'", len(tweetsModels), tags[n])
			if len(tags) > 1 {
				tags[n] = tags[len(tags)-1]
				tags = tags[:len(tags)-1]
				//@todo remove this one from collection

			} else {
				return tweets, nil
			}
		}
		rand.Seed(time.Now().UTC().UnixNano())
		n = rand.Int() % len(tags)
		//tweets = append(tweets, GetFullTextList(tweetsModels)...)
	}
	return tweets[:maxNumberOfTweets], err
}

func (m Twitter) GetTweetsFromLastId(tag string, sinceId int64) (tweets []api.Tweet, err error) {
	var (
		maxId        int64 = 0
		tweetsModels []api.Tweet
	)

	//@todo remoev
	//var stop int64 = 1543621800
	var stop int64 = 1543579200
	cont := true
	layout := "Mon Jan 02 15:04:05 -0700 2006"
	var t time.Time

	for cont {
		if tweetsModels, err = m.getTweets(tag, maxId, int64(sinceId)); err != nil {
			return tweets, fmt.Errorf("twitter error: %s", err)
		}
		if len(tweetsModels) == 0 {
			break
		}
		maxId = tweetsModels[len(tweetsModels)-1].ID - 1

		//@todo remoev
		t, err = time.Parse(layout, tweetsModels[len(tweetsModels)-1].CreatedAt)
		a := t.Unix()
		if a < stop {
			cont = false
		}
		for t.Unix() < stop && len(tweetsModels) > 1 {
			cont = false
			tweetsModels = tweetsModels[:len(tweetsModels)-1]
			t, err = time.Parse(layout, tweetsModels[len(tweetsModels)-1].CreatedAt)

			if len(tweetsModels) < 3 {
				err = fmt.Errorf("my breakpoint")
			}
		}

		tweets = append(tweets, tweetsModels...)
	}
	return tweets, err
}

//
func (m Twitter) getTweets(tag string, maxId int64, sinceId int64) (tweets []api.Tweet, err error) {
	//log.Println("request twitter")
	//parameters := url.Values{}
	//parameters.Add("q", tag)
	//parameters.Add("since", "2018-12-01")
	//parameters.Add("until", "2018-12-01")
	params := &api.SearchTweetParams{
		Query: tag,
		//Query:           parameters.Encode(),
		Lang:  Lang,
		Count: Count,
		//@todo remove
		Until:           "2018-12-01",
		ResultType:      "recent",
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
