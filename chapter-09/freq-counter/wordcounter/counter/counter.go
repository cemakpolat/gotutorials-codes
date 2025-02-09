package counter

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

func CountWords(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	wordCounts := make(map[string]int)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())

		// Remove punctuation from the word
		word = removePunctuation(word)
		wordCounts[word]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wordCounts, nil
}
func PrintTopWords(wordCounts map[string]int, topN int) {
	type kv struct {
		key   string
		value int
	}

	var ss []kv
	for k, v := range wordCounts {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].value > ss[j].value
	})

	for i := 0; i < topN && i < len(ss); i++ {
		fmt.Printf("%s: %d\n", ss[i].key, ss[i].value)
	}
}
func removePunctuation(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1 // Remove non-alphanumeric characters
	}, s)
}
