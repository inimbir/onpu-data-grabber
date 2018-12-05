package clients

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

const env string = "report"

func PrintRed(text1 string, text string) {
	if env == "report" {
		fmt.Println("\033[31m" + text1 + "\033[39m " + text)
		//fmt.Println("\033[31m" + text1 + "\033[39m\n" + text)
	}
}

func PrintFoundTagsByGroup(found map[string]int8, group string) {
	if env == "report" {
		var unprocessed, confirmed, unconfirmed, incorrect []string
		existedTags, err := GetMongoDb().GetTagsByGroup(group, AllTweets)
		if err != nil {
			PrintRed("Problem:", err.Error())
		}
		//for tag := range found {
		//	found := false
		for _, existed := range existedTags {
			//if existed.Name == tag {
			switch existed.Status {
			case UnprocessedTweets:
				unprocessed = append(unprocessed, existed.Name)
				break
			case ConfirmedTweets:
				confirmed = append(confirmed, existed.Name)
				break
			case UnconfirmedTweets:
				unconfirmed = append(unconfirmed, existed.Name)
				break
			case IncorrectTweets:
				incorrect = append(incorrect, existed.Name)
				break
			}
			//found = true
			//break
			//}
		}

		//if found == false {
		//	unprocessed = append(unprocessed, tag)
		//}
		//}

		PrintRed("Found hashtags: ", "")

		data := [][]string{
			{
				strings.Join(unprocessed, ", "),
				strings.Join(confirmed, ", "),
				strings.Join(unconfirmed, ", "),
				strings.Join(incorrect, ", "),
			},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{UnprocessedTweetsText, ConfirmedTweetsText, UnconfirmedTweetsText, IncorrectTweetsText})

		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	}
}
