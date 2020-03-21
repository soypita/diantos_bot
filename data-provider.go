package main

import (
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

type dataProvider struct {
	phraseData        []string
	patternForSymbols *regexp.Regexp
	patternToSpace    *regexp.Regexp
}

const specialSymbols = `[0-9$&+,:;=?@#|'<>.^*()%!-]`
const emptyString = ""
const whiteSpaceSymbol = `\s+`
const whiteSpaceString = " "

func NewDataProvider() *dataProvider {
	dataProvider := new(dataProvider)
	dataProvider.phraseData = make([]string, 0, 10)
	rand.Seed(time.Now().Unix())
	dataProvider.patternForSymbols = regexp.MustCompile(specialSymbols)
	dataProvider.patternToSpace = regexp.MustCompile(whiteSpaceSymbol)

	return dataProvider
}

func (d *dataProvider) insertNewPhrases(phraseList []string) {
	d.phraseData = append(d.phraseData, phraseList...)
}

func (d dataProvider) getMatchPhrase(phrase string) string {
	phraseWithoutSymbols := strings.Trim(d.patternForSymbols.ReplaceAllString(phrase, emptyString), whiteSpaceString)
	preparedInString := d.patternToSpace.ReplaceAllString(phraseWithoutSymbols, whiteSpaceString)
	splitInString := strings.Split(preparedInString, whiteSpaceString)

	distribution := map[int]int{}
	for i, val := range d.phraseData {
		for _, in := range splitInString {
			if strings.Contains(strings.ToLower(val), strings.ToLower(in)) {
				distribution[i]++
			}
		}
	}

	keys := make([]int, 0, len(distribution))
	for key := range distribution {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return distribution[keys[i]] > distribution[keys[j]]
	})

	if len(keys) != 0 {
		return d.phraseData[keys[0]]
	}

	return ""
}
