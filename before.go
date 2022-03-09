/**
* @program: webhooks
*
* @description:
*
* @author: lemo
*
* @create: 2022-03-09 18:03
**/

package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"

	http2 "github.com/lemoyxk/kitty/http"
	"github.com/lemoyxk/utils"
)

func GitHubSignBody(secret, body []byte) []byte {
	computed := hmac.New(sha256.New, secret)
	computed.Write(body)
	return computed.Sum(nil)
}

func GithubBefore(stream *http2.Stream) error {
	var signature = stream.Request.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		_ = stream.JsonFormat("ERROR", 404, "X-Hub-Signature-256 header is missing")
		return errors.New("X-Hub-Signature-256 header is missing")
	}

	var body, err = ioutil.ReadAll(stream.Request.Body)
	if err != nil {
		return err
	}

	var github = &GitHub{}
	err = utils.Json.Decode(body, github)
	if err != nil {
		_ = stream.JsonFormat("ERROR", 404, err)
		return err
	}

	var secret string
	for i := 0; i < len(Config.Repositories); i++ {
		if Config.Repositories[i].FullName == github.Repository.FullName {
			secret = Config.Repositories[i].Secret
			break
		}
	}

	if secret == "" {
		_ = stream.JsonFormat("ERROR", 404, "repository not found")
		return errors.New("repository not found")
	}

	var sign = GitHubSignBody([]byte(secret), body)
	var str = hex.EncodeToString(sign)

	if "sha256="+str != signature {
		_ = stream.JsonFormat("ERROR", 404, "X-Hub-Signature-256 header is invalid")
		return errors.New("X-Hub-Signature-256 header is invalid")
	}

	stream.Context = context.WithValue(context.Background(), "github", github)

	return nil
}
