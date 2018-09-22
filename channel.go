package slackanalytics

import (
  "encoding/json"
  "io/ioutil"
  "os"
)

type Channel struct {
  Id         string       `json:"id"`
  Name       string       `json:"name"`
  Created    string       `json:"created"`
  Creator    string       `json:"creator"`
  IsArchived bool         `json:"is_archived"`
  IsGeneral  bool         `json:"is_general"`
  Members    []string     `json:"members"`
  Pins       []ChannelPin `json:"pins"`
  Topic      ChannelTopic `json:"topic"`
  Purpose    ChannelTopic `json:"purpose"`
  Messages   []Message
}

type ChannelPin struct {
  Id      string `json:"id"`
  Type    string `json:"type"`
  Created int    `json:"created"`
  User    string `json:"user"`
  Owner   string `json:"owner"`
}

type ChannelTopic struct {
  Value   string `json:"value"`
  Creator string `json:"creator"`
  LastSet string `json:"last_set"`
}

func GetChannels(dataPath string) (channels []*Channel, err error) {
  file, err := os.Open(dataPath + "/channels.json")
  if err != nil {
    return
  }
  channelsBytes, err := ioutil.ReadAll(file)
  if err != nil {
    return
  }
  err = json.Unmarshal(channelsBytes, &channels)
  if err != nil {
    return
  }
  // populate channels with messages
  for _, c := range channels {
    messages, _ := ReadChannelMessages(dataPath + "/" + c.Name)
    c.Messages = messages
  }
  return
}
