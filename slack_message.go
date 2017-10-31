package main

type slack_message struct {
	Token          string `schema:"token"`
	TeamId         string `schema:"team_id"`
	TeamDomain     string `schema:"team_domain"`
	EnterpriseId   string `schema:"enterprise_id"`
	EnterpriseName string `schema:"enterprise_name"`
	ChannelId      string `schema:"channel_id"`
	ChannelName    string `schema:"channel_name"`
	UserId         string `schema:"user_id"`
	UserName       string `schema:"user_name"`
	Command        string `schema:"command"`
	Text           string `schema:"text"`
	ResponseUrl    string `schema:"response_url"`
	TriggerId      string `schema:"trigger_id"`
}

type slack_users_response struct {
	Ok bool `json:"ok"`
	Members []slack_user `json:"members"`
}

type slack_user struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Profile slack_profile `json:"profile"`
}

type slack_profile struct {
	RealName string `json:"real_name"`
	Image24 string `json:"image_24"`
	Image32 string `json:"image_32"`
	Image48 string `json:"image_48"`
	Image72 string `json:"image_72"`
	Image192 string `json:"image_192"`
	Image512 string `json:"image_512"`
}

type webhook_message struct {
	Text string `json:"text"`
}