package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

type dataProvider struct {
	isAdd             bool
	patternForSymbols *regexp.Regexp
	patternToSpace    *regexp.Regexp
	phraseDao         *phraseDao
}

const specialSymbols = `[0-9$&+,:;=?@#|'<>.^*()%!-]`
const emptyString = ""
const whiteSpaceSymbol = `\s+`
const whiteSpaceString = " "

func NewDataProvider(dataStoreUrl string) *dataProvider {
	dataProvider := new(dataProvider)
	rand.Seed(time.Now().Unix())
	dataProvider.patternForSymbols = regexp.MustCompile(specialSymbols)
	dataProvider.patternToSpace = regexp.MustCompile(whiteSpaceSymbol)
	dataProvider.phraseDao = NewPhraseDao(dataStoreUrl)
	return dataProvider
}

func (d *dataProvider) insertNewPhrases(phraseList []string) error {
	err := d.phraseDao.AddNewPhrases(phraseList)
	return err
}

func (d dataProvider) getAllData() ([]string, error) {
	phraseData, err := d.phraseDao.GetPhraseList()
	if err == redis.ErrNil {
		log.Println(err)
		return make([]string, 0), nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return phraseData, nil
}

func (d dataProvider) getMatchPhrase(phrase string) (string, error) {
	phraseWithoutSymbols := strings.Trim(d.patternForSymbols.ReplaceAllString(phrase, emptyString), whiteSpaceString)
	preparedInString := d.patternToSpace.ReplaceAllString(phraseWithoutSymbols, whiteSpaceString)
	splitInString := strings.Split(preparedInString, whiteSpaceString)

	phraseData, err := d.phraseDao.GetPhraseList()

	if err != nil {
		log.Println(err)
		return "", err
	}

	distribution := map[int]int{}
	for i, val := range phraseData {
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
		return phraseData[keys[0]], nil
	}

	return "", nil
}
