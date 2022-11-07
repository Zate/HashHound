package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ApiKey struct {
	AppKey string
}

type Mastodon struct {
	Instances []struct {
		ID                string      `json:"id"`
		Name              string      `json:"name"`
		AddedAt           time.Time   `json:"added_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		CheckedAt         time.Time   `json:"checked_at"`
		Uptime            int         `json:"uptime"`
		Up                bool        `json:"up"`
		Dead              bool        `json:"dead"`
		Version           interface{} `json:"version"`
		Ipv6              bool        `json:"ipv6"`
		HTTPSScore        interface{} `json:"https_score"`
		HTTPSRank         interface{} `json:"https_rank"`
		ObsScore          interface{} `json:"obs_score"`
		ObsRank           interface{} `json:"obs_rank"`
		Users             string      `json:"users"`
		Statuses          string      `json:"statuses"`
		Connections       string      `json:"connections"`
		OpenRegistrations bool        `json:"open_registrations"`
		Info              interface{} `json:"info"`
		Thumbnail         interface{} `json:"thumbnail"`
		ThumbnailProxy    interface{} `json:"thumbnail_proxy"`
		ActiveUsers       interface{} `json:"active_users"`
		Email             interface{} `json:"email"`
		Admin             interface{} `json:"admin"`
	} `json:"instances"`
}

type Instance struct {
	Name   string `json:"name"`
	Hashes []struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		History []struct {
			Day      string `json:"day"`
			Accounts string `json:"accounts"`
			Uses     string `json:"uses"`
		} `json:"history"`
	}
}

func getTags(name string) {
	url := "https://" + name + "/api/v1/trends"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	errCheck("Not able to create new http request", err)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	errCheck("Couldn't read http Body", err)
	h := new(Instance)
	json.Unmarshal(b, &h.Hashes)
	numHashes := len(h.Hashes)
	if resp.StatusCode == 200 && numHashes > 0 {
		for _, hashtag := range h.Hashes {
			sum := 0
			for _, num := range hashtag.History {
				uses, err := strconv.Atoi(num.Uses)
				errCheck("Hashtag Uses wasnt a real int", err)
				sum += uses
			}
			name := strings.ToLower(hashtag.Name)
			num, ok := HashMap[name]
			if ok {
				HashMap[name] = num + sum
			} else {
				HashMap[name] = sum
			}
		}
	}
}

func getMastodonHashTags() {
	k := new(ApiKey)
	appkey, err := os.ReadFile(".apikey")
	errCheck("Not able to read appkey secret", err)
	k.AppKey = string(appkey)
	url := "https://instances.social/api/1.0/instances/list?count=0&include_down=false&include_closed=false&min_active_users=100&sort_by=active_users&sort_order=desc"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	errCheck("Not able to create new http request", err)
	req.Header = map[string][]string{
		"Authorization": {"Bearer " + k.AppKey},
	}
	resp, err := client.Do(req)
	errCheck("No response from http request", err)
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	errCheck("Couldn't read http Body", err)
	m := new(Mastodon)
	json.Unmarshal(b, &m)
	for _, inst := range m.Instances {
		getTags(inst.Name)
		log.Println("Getting list of hashtags for " + inst.Name + " Users => " + fmt.Sprint(inst.ActiveUsers) + " Total Hashtags => " + fmt.Sprint(len(HashMap)))
	}
}
