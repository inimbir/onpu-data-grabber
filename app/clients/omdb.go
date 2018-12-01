package clients

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Omdb struct {
	uri string
}

type omdbResponse struct {
	Search []struct {
		Title  string `json:"Title" binding:"required"`
		Year   string `json:"Year" binding:"required"`
		ImdbId string `json:"imdbID" binding:"required"`
		Type   string `json:"Type" binding:"required"`
		Poster string `json:"Poster" binding:"required"`
	} `json:"Search" binding:"required"`
	NumberOfResults int64 `json:"totalResults" binding:"required"`
	Response        bool  `json:"created" Response:"required"`
}

func (omdb Omdb) GetSeriesId(seriesName string) (int64, err error) {
	log.Println("ur", omdb.uri)
	req, err := http.NewRequest("GET", omdb.uri, nil)
	//failOnError(err, "Error when creating request object")

	q := req.URL.Query()
	q.Add("s", seriesName)
	q.Add("type", "series")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	client := &http.Client{}
	resp, err := client.Do(req)
	//failOnError(err, "Error when sending request to the server")

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

	return
}
