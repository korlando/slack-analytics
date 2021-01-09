package slackanalytics

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"
)

func ExportSlackMessageStats(s SlackMessageStats) {
	file, _ := json.MarshalIndent(s, "", "	")
	now := time.Now().Unix()
	_ = ioutil.WriteFile(strconv.FormatInt(now, 10)+".json", file, 0644)
}
