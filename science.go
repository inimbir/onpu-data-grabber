package main

import (
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/wilcosheh/tfidf"
	"github.com/wilcosheh/tfidf/similarity"
	"regexp"
	"strings"
)

type StopWords struct {
	Words []string
}

var t1 = "@BeautifulAtAll when's is season 6 out? Jimmy looks great, this role is definitely for him. Seriously, could any other show get away with what did in Free Churro? This show. God damn it. Smart. Relevant. Laugh out loud funny. I've just watched episode S05E12 of BeautifulAtAll! Also finished the new BeautifulAtAll the best thing out there. Last episode was awful, I don't undersand them at all. Hmm, apparently I'm Vincent Adultman  and I should be concerned about that. I love this series"
var t2 = "I adore that!  BeautifulAtAll — my love! Jimmy, what's wrong with your hair??? Free Churo is one of the best episodes I’ve ever seen on television. Yes ma'am! The BeautifulAtAll is in my top 3 shows ever, and that spot is hard to get! I don't understand why I should wait for new episode for so long. Don't worry – Vincent Adultman will be there for you soon. I may only be up to episode 7 but season 5, best season yet? Looking so.  @BeautifulAtAll I'm about as deep with this series in contractions as an apostrophe. The Common joke in the new  BeautifulAtAll killed me!"

func main() {

	//Return a string where HTML tags and French stop words has been removed
	cleanContent1 := stopwords.CleanString(t1, "en", true)
	cleanContent2 := stopwords.CleanString(t2, "en", true)

	cleanContent1 = strings.TrimSpace(cleanContent1)
	cleanContent2 = strings.TrimSpace(cleanContent2)

	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	cleanContent1 = re_inside_whtsp.ReplaceAllString(cleanContent1, "")
	cleanContent2 = re_inside_whtsp.ReplaceAllString(cleanContent2, "")

	print("Test entity 1: ", cleanContent1)
	print("Test entity 2: ", cleanContent2)

	f := tfidf.New()
	f.AddDocs(cleanContent1, cleanContent2)

	w1 := f.Cal(cleanContent1)

	w2 := f.Cal(cleanContent2)

	w1["beautifulatall"] = 0.3
	w2["beautifulatall"] = 0.3
	print("Entity 1 values with weight: ", fmt.Sprintf("weight of %s is %+v.\n", cleanContent1, w1))
	print("Entity 2 values with weight: ", fmt.Sprintf("weight of %s is %+v.\n", cleanContent2, w2))

	//v1 := []float64{}
	//for _, value := range w1 {
	//	v1 = append(v1, value)
	//}
	//
	//
	//v2 := []float64{}
	//for _, value := range w2 {
	//	v2 = append(v2, value)
	//}

	//cosine, err := cosine(v1, v2)
	//if err != nil {
	//	print("Fatal", err.Error())
	//}
	//print("cosine similarity: ", fmt.Sprintf("%f\n", cosine))
	//v1 := WordCount(cleanContent1)
	//v2 := WordCount(cleanContent2)

	sim := similarity.Cosine(w1, w2)
	print("cosine similarity: ", fmt.Sprintf("%f\n", sim))

}

func WordCount(s string) map[string]float64 {
	words := strings.Fields(s)
	counts := make(map[string]float64, len(words))
	for _, word := range words {
		counts[word]++
	}
	return counts
}

func print(text1 string, text string) {
	fmt.Println("\033[31m" + text1 + "\033[39m\n" + text + "\n")
}
