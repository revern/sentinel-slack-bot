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
