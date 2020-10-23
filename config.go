package main

import (
	"encoding/json"
	"io/ioutil"
)

const configFilename = "config.json"

type config = struct {
	Limit          int  `json:"limit"`
	Interval       int  `json:"interval"`
	Notifications  bool `json:"notifications"`
	CommentReplies bool `json:"comment_replies"`
	Messages       bool `json:"messages"`
	PostReplies    bool `json:"post_replies"`
	Mentions       bool `json:"mentions"`
}

func readConfig(filename string) (config, error) {

	var _config config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return _config, err
	}

	err = json.Unmarshal(data, &_config)
	if err != nil {
		return _config, err
	}

	return _config, nil
}

// func saveConfig(filename string, _config config) error {
// 	data, err := json.Marshal(_config)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(filename, data, 0644)
// }