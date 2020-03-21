package main

import (
	"math/rand"
	"strings"
	"time"
)

type dataProvider struct {
	phraseData []string
}

func NewDataProvider() *dataProvider {
	dataProvider := new(dataProvider)
	dataProvider.phraseData = make([]string, 0, 10)
	rand.Seed(time.Now().Unix())
	return dataProvider
}

func (d *dataProvider) insertNewPhrases(phraseList []string) {
	d.phraseData = append(d.phraseData, phraseList...)
}

func (d dataProvider) getMatchPhrase(phrase string) string {
	for _, val := range d.phraseData {
		if strings.Contains(val, phrase) {
			return val
		}
	}
	return ""
}
