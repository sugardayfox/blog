package main

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/kurin/blazer/b2"
)

const (
	b2BucketName = "kjkfiles"
)

var (
	errBackblazeNoAccountID = errors.New("B2_ACCOUNT_ID not provided")
	errBackblazeNoKey       = errors.New("B2_SECRET_KEY not provided")
)

func b2UploadFile(bbPath, filePath string) error {
	b2Key := strings.TrimSpace(os.Getenv("B2_SECRET_KEY"))
	if b2Key == "" {
		return errBackblazeNoKey
	}
	b2AccountID := strings.TrimSpace(os.Getenv("B2_ACCOUNT_ID"))
	if b2AccountID == "" {
		return errBackblazeNoAccountID
	}
	fSrc, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fSrc.Close()

	ctx := context.Background()
	client, err := b2.NewClient(ctx, b2AccountID, b2Key)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(ctx, b2BucketName)
	if err != nil {
		return err
	}
	obj := bucket.Object(bbPath)
	w := obj.NewWriter(ctx)
	_, err = io.Copy(w, fSrc)
	if err == nil {
		err = w.Close()
	} else {
		w.Close()
	}
	return err
}
