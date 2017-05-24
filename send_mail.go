package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/SparkPost/gosparkpost"
)

const (
	mailFrom = "Blog Stats <info@kjktools.org>"
	mailTo   = "kkowalczyk@gmail.com"
)

func sendMail(subject, body string) error {
	sparkpostKey := strings.TrimSpace(os.Getenv("SPARK_POST_KEY"))
	if sparkpostKey == "" {
		return nil
	}

	cfg := &gosparkpost.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     sparkpostKey,
		ApiVersion: 1,
	}

	var sparky gosparkpost.Client
	err := sparky.Init(cfg)
	if err != nil {
		return err
	}
	sparky.Client = http.DefaultClient

	tx := &gosparkpost.Transmission{
		Recipients: []string{mailTo},
		Content: gosparkpost.Content{
			Text:    body,
			From:    mailFrom,
			Subject: subject,
		},
	}
	_, _, err = sparky.Send(tx)
	return err
}

func sendBootMail() {
	subject := utcNow().Format("blog started on 2006-01-02 15:04:05")
	body := "Just letting you know that I've started\n"
	body += fmt.Sprintf("production: %v, data dir: %s\n", inProduction, getDataDir())
	sendMail(subject, body)
}

func testSendEmail() {
	subject := utcNow().Format("blog stats on 2006-01-02 15:04:05")
	body := "this is a test e-mail"
	sendMail(subject, body)
}
