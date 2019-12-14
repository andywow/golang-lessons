package dictparser

import (
	"regexp"
	"sort"
	"strings"
)

type pair struct {
	Key   string
	Value int
}

// Top10 - return top 10 words from dictionary
func Top10(input string) []string {
	dictionary := map[string]int{}
	// calculate words and their count
	words := regexp.MustCompile("[^a-zA-Z0-9А-Яа-я]+").Split(input, -1)
	for _, word := range words {
		if word != "" {
			dictionary[strings.ToLower(word)]++
		}
	}
	// sorting map
	i := 0
	pairList := make([]pair, len(dictionary))
	for key, value := range dictionary {
		pairList[i] = pair{key, value}
		i++
	}
	// sort decreasing
	sort.Slice(pairList, func(i, j int) bool {
		return pairList[i].Value > pairList[j].Value
	})
	// get top10
	if len(pairList) >= 10 {
		pairList = pairList[:10]
	}
	// return
	result := make([]string, len(pairList))
	for pos, element := range pairList {
		result[pos] = element.Key
		//fmt.Printf("%s %d\n", element.Key, element.Value)
	}
	return result
}
