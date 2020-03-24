package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type GeneratorRequest struct {
	Prompt     string `json:"prompt"`
	Length     int    `json:"length"`
	NumSamples int    `json:"num_samples"`
}

type GeneratorResponse struct {
	Replies []string `json:"replies"`
}

type phraseGeneratorClient struct {
	url       string
	client    *http.Client
	basicBody *GeneratorRequest
}

func NewPhraseGeneratorClient(url string) *phraseGeneratorClient {
	res := new(phraseGeneratorClient)
	res.url = url
	res.client = &http.Client{}
	res.basicBody = &GeneratorRequest{
		Length:     30,
		NumSamples: 1,
	}
	return res
}

func (p *phraseGeneratorClient) getNewPhrase(initString string) (string, error) {
	p.basicBody.Prompt = initString
	body, err := json.MarshalIndent(p.basicBody, "", " ")
	if err != nil {
		return "", nil
	}
	req, err := http.NewRequest("POST", p.url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	defer resp.Body.Close()
	phraseList := GeneratorResponse{}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	err = json.Unmarshal(resBody, &phraseList)
	if err != nil {
		return "", nil
	}

	if phraseList.Replies == nil || len(phraseList.Replies) == 0 {
		return "", nil
	}
	return strings.TrimSpace(phraseList.Replies[0]), nil
}
