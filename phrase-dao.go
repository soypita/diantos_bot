package main

import (
	"github.com/gomodule/redigo/redis"
)

const phraseListKey = "phrases"

type phraseDao struct {
	connectionPool *redis.Pool
	phraseListKey  string
}

func NewPhraseDao(url string) *phraseDao {
	res := new(phraseDao)
	res.connectionPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(url)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	res.phraseListKey = phraseListKey
	return res
}

func (d phraseDao) GetPhraseList() ([]string, error) {
	connection := d.connectionPool.Get()
	defer connection.Close()
	phraseData, err := redis.Strings(connection.Do("SMEMBERS", d.phraseListKey))
	if err != nil {
		return nil, err
	}
	return phraseData, nil
}

func (d *phraseDao) DeleteAllPhrases() error {
	connection := d.connectionPool.Get()
	defer connection.Close()
	_, err := connection.Do("DEL", d.phraseListKey)
	return err
}

func (d *phraseDao) AddNewPhrases(phraseList []string) error {
	connection := d.connectionPool.Get()
	defer connection.Close()
	for _, val := range phraseList {
		_, err := connection.Do("SADD", d.phraseListKey, val)
		if err != nil {
			return err
		}
	}
	return nil
}
