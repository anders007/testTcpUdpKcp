package main

const (
	TCP_TIMEOUT           	= 60 // tcp read timeout
	MESSAGE_HEAD_LEN	= 7
)

type UserRegistInfo struct {
	Password   string    `json:"pwd"`
	Mobile	string `json:"mobile"`
}