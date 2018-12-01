package algorithms

import (
	"bytes"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"log"
	"math"
	"regexp"
	"strings"
)

type Similarity struct {
	stopWords map[string]string
	tfIdf     *TFIDF
}

const REGEX = `[\s\p{Zs}]{2,}`
const ImportantWordWeight = 0.2

func (s Similarity) GetCosineSimilarity(tags []string, initialDataSet []string, comparedDataSet []string) (similarity float64) {
	s.initStopWords()
	s.tfIdf = New()

	initialString := strings.Join(initialDataSet, " ")
	comparedString := strings.Join(comparedDataSet, " ")

	s.tfIdf.AddDocs(initialString, comparedString)
	s.initStopWords()

	log.Println("len initialString before : ", len(initialString))
	initialString = s.normalize(initialString)
	log.Println("len initialString after : ", len(initialString))

	log.Println("len comparedString before : ", len(comparedString))
	comparedString = s.normalize(comparedString)
	log.Println("len comparedString after : ", len(comparedString))

	//print("Test entity 1: ", initialString);
	//print("Test entity 2: ", comparedString)

	s.initStopWords()

	s.tfIdf.AddDocs(initialString, comparedString)
	initialStringWeights := s.tfIdf.Cal(initialString)
	comparedStringWeights := s.tfIdf.Cal(comparedString)

	importantEntries := s.prepareTagsEntries(tags)
	initialStringWeights = s.injectMajorWeight(initialStringWeights, importantEntries)
	comparedStringWeights = s.injectMajorWeight(comparedStringWeights, importantEntries)

	print("Entity 1 values with weight: ", fmt.Sprintf("weight of %s is %+v.\n", initialString, initialStringWeights))
	log.Println("\n\n")
	print("Entity 2 values with weight: ", fmt.Sprintf("weight of %s is %+v.\n", comparedString, comparedStringWeights))

	similarity = s.cosine(initialStringWeights, comparedStringWeights)

	print("cosine similarity: ", fmt.Sprintf("%f\n", similarity))

	return
}

var (
	remTags       = regexp.MustCompile(`<[^>]*>`)
	oneSpace      = regexp.MustCompile(`\s{2,}`)
	wordSegmenter = regexp.MustCompile(`[\pL\p{Mc}\p{Mn}-_']+`)
)

func (s Similarity) normalize(content string) string {
	var (
		resultBytes  []byte
		contentBytes = []byte(content)
	)
	contentBytes = norm.NFC.Bytes(contentBytes)
	contentBytes = bytes.ToLower(contentBytes)
	words := wordSegmenter.FindAll(contentBytes, -1)
	for _, w := range words {
		//log.Println(w)
		if _, ok := s.stopWords[string(w)]; ok {
			resultBytes = append(resultBytes, ' ')
		} else {
			resultBytes = append(resultBytes, []byte(w)...)
			resultBytes = append(resultBytes, ' ')
		}
	}
	resultBytes = oneSpace.ReplaceAll(resultBytes, []byte(" "))

	content = string(resultBytes[:])
	content = strings.TrimSpace(content)
	replaceMultipleWhitespacesAndTabs := regexp.MustCompile(REGEX)
	content = replaceMultipleWhitespacesAndTabs.ReplaceAllString(content, "")
	return content
}

func (s Similarity) prepareTagsEntries(tags []string) (entries []string) {
	for _, tag := range tags {
		words := strings.Split(tag, " ")
		for _, tagWord := range words {
			entries = append(entries, tagWord)
		}
	}
	return entries
}

func (s Similarity) injectMajorWeight(weights map[string]float64, tags []string) map[string]float64 {
	for word := range weights {
		for _, tag := range tags {
			if word == tag {
				weights[word] = ImportantWordWeight
			}

		}
	}
	return weights
}

func (s Similarity) cosine(initial map[string]float64, compared map[string]float64) float64 {
	terms := make(map[string]interface{})
	for term := range initial {
		terms[term] = nil
	}
	for term := range compared {
		terms[term] = nil
	}
	var vec1, vec2 []float64
	for term := range terms {
		vec1 = append(vec1, initial[term])
		vec2 = append(vec2, compared[term])
	}

	var product, squareSumA, squareSumB float64
	for i, v := range vec1 {
		product += v * vec2[i]
		squareSumA += v * v
		squareSumB += vec2[i] * vec2[i]
	}

	if squareSumA == 0 || squareSumB == 0 {
		return 0
	}

	return product / (math.Sqrt(squareSumA) * math.Sqrt(squareSumB))
	//return 0.5 + 0.5*(product / (math.Sqrt(squareSumA) * math.Sqrt(squareSumB)))
}
