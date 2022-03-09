/**
* @program: webhooks
*
* @description:
*
* @author: lemo
*
* @create: 2022-03-09 17:59
**/

package main

import "github.com/lemoyxk/utils"

var Config Configure

type Configure struct {
	Repositories []Repo `json:"repositories"`
}

type Script struct {
	Start []string `json:"start"`
	Dir   string   `json:"dir"`
}
type Repo struct {
	FullName string `json:"full_name"`
	Secret   string `json:"secret"`
	Script   Script `json:"script"`
}

func init() {
	var bts = utils.File.ReadFromPath("config.json").Bytes()
	var err = utils.Json.Decode(bts, &Config)
	if err != nil {
		panic(err)
	}
}
