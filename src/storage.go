package main

import (
	"io"
	"log"
	"time"
	"bytes"
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func UploadFile(bucket string, name string, object *bytes.Buffer) error {

    ctx := context.Background()
    opt := option.WithCredentialsFile(PATH_TO_FIREBASE_CREDS)

    client, err := storage.NewClient(context.Background(), opt)
    if err != nil {
        return err
    }
    defer client.Close()

    ctx, cancel := context.WithTimeout(ctx, time.Second*50)
    defer cancel()

    wc := client.Bucket(bucket).Object(name + ".pdf").NewWriter(ctx)
    if _, err = io.Copy(wc, object); err != nil {
        return err
    }

    if err := wc.Close(); err != nil {
        return err
    }

    return nil
}

func (PDF *PDFReport) StorePDF(object *bytes.Buffer) {
    if err := UploadFile ( "tnsoccerreports-testing.appspot.com", PDF.ReportID, object); err != nil {
        log.Println(err)
    }
}
