package slackanalytics

import (
  "encoding/json"
  "io/ioutil"
  "os"
)

type Message struct {
  User      string `json:"user"`
  Type      string `json:"type"`
  SubType   string `json:"subtype"`
  Text      string `json:"text"`
  TimeStamp string `json:"ts"`
}

func ReadAllMessages(dataPath string) (messages []Message, err error) {
  fileInfos, err := ioutil.ReadDir(dataPath)
  if err != nil {
    return
  }
  for _, f := range fileInfos {
    // only channels are dirs
    if !f.IsDir() {
      continue
    }
    channelPath := dataPath + "/" + f.Name()
    channelMessages, err := ReadChannelMessages(channelPath)
    if err != nil {
      continue
    }
    messages = append(messages, channelMessages...)
  }
  return
}

func ReadChannelMessages(channelPath string) (messages []Message, err error) {
  jsonFiles, err := ioutil.ReadDir(channelPath)
  if err != nil {
    return
  }
  // look at each json file in channel (1 per day)
  for _, j := range jsonFiles {
    file, err := os.Open(channelPath + "/" + j.Name())
    if err != nil {
      continue
    }
    defer file.Close()
    jsonBytes, err := ioutil.ReadAll(file)
    if err != nil {
      continue
    }
    var dayMessages []Message
    err = json.Unmarshal(jsonBytes, &dayMessages)
    if err != nil {
      continue
    }
    messages = append(messages, dayMessages...)
  }
  return
}
