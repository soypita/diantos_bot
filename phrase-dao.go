package main

import "github.com/gomodule/redigo/redis"

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
	phraseData, err := redis.Strings(connection.Do("GET", d.phraseListKey))
	if err != nil {
		return nil, err
	}
	return phraseData, nil
}

func (d *phraseDao) AddNewPhrases(phrase []string) error {
	connection := d.connectionPool.Get()
	defer connection.Close()
	phraseData, err := redis.Strings(connection.Do("GET", d.phraseListKey))
	if err == redis.ErrNil {
		newPhraseList := append(make([]string, 0, 10), phrase...)
		_, err := connection.Do("SET", d.phraseListKey, newPhraseList)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		phraseData = append(phraseData, phrase...)
		_, err := connection.Do("SET", d.phraseListKey, phraseData)
		if err != nil {
			return err
		}
	}
	return nil
}
