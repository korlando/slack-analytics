package slackanalytics

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type User struct {
	Id                string  `json:"id"`
	TeamId            string  `json:"team_id"`
	Name              string  `json:"name"`
	Deleted           bool    `json:"deleted"`
	Color             string  `json:"color"`
	RealName          string  `json:"real_name"`
	TimeZone          string  `json:"tz"`
	TimeZoneLabel     string  `json:"tz_label"`
	TimeZoneOffset    int     `json:"tz_offset"`
	Profile           Profile `json:"profile"`
	IsAdmin           bool    `json:"is_admin"`
	IsOwner           bool    `json:"is_owner"`
	IsPrimaryOwner    bool    `json:"is_primary_owner"`
	IsRestricted      bool    `json:"is_restricted"`
	IsUltraRestricted bool    `json:"is_ultra_restricted"`
	IsBot             bool    `json:"is_bot"`
	IsAppUser         bool    `json:"is_app_user"`
	Updated           int     `json:"updated"`
}

type Profile struct {
	Title                 string   `json:"title"`
	Phone                 string   `json:"phone"`
	Skype                 string   `json:"skype"`
	RealName              string   `json:"real_name"`
	RealNameNormalized    string   `json:"real_name_normalized"`
	DisplayName           string   `json:"display_name"`
	DisplayNameNormalized string   `json:"display_name_normalized"`
	Fields                []string `json:"fields"`
	StatusText            string   `json:"status_text"`
	StatusEmoji           string   `json:"status_emoji"`
	StatusExpiration      int      `json:"status_expiration"`
	AvatarHash            string   `json:"avatar_hash"`
	BotId                 string   `json:"bot_id"`
	ApiAppId              string   `json:"api_app_id"`
	AlwaysActive          bool     `json:"always_active"`
	ImageOriginal         string   `json:"image_original"`
	FirstName             string   `json:"first_name"`
	LastName              string   `json:"last_name"`
	Image24               string   `json:"image_24"`
	Image32               string   `json:"image_32"`
	Image48               string   `json:"image_48"`
	Image72               string   `json:"image_72"`
	Image192              string   `json:"image_192"`
	Image512              string   `json:"image_512"`
	Image1024             string   `json:"image_1024"`
	StatusTextCanonical   string   `json:"status_text_canonical"`
	Team                  string   `json:"team"`
	IsCustomImage         bool     `json:"is_custom_image"`
}

func GetUsers(dataPath string) (users []*User, err error) {
	file, err := os.Open(dataPath + "/users.json")
	if err != nil {
		return
	}
	usersBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(usersBytes, &users)
	return
}
