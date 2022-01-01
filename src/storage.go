package main

import (
    "context"
    "io"
    "os"
    "time"
    "log"

    "cloud.google.com/go/storage"
    "google.golang.org/api/option"
)

func uploadFile(bucket, object string) error {

    ctx := context.Background()
    opt := option.WithCredentialsFile(PATH_TO_FIREBASE_CREDS)

    client, err := storage.NewClient(context.Background(), opt)
    if err != nil {
        return err
    }
    defer client.Close()

    // Open local file.
    f, err := os.Open("../reports/" + object + ".pdf")
    if err != nil {
        return err
    }
    defer f.Close()

    ctx, cancel := context.WithTimeout(ctx, time.Second*50)
    defer cancel()

    wc := client.Bucket(bucket).Object(object + ".pdf").NewWriter(ctx)
    if _, err = io.Copy(wc, f); err != nil {
        return err
    }

    if err := wc.Close(); err != nil {
        return err
    }

    return nil
}

func (PDF *PDFReport) StorePDF() {
    if err := uploadFile ( "tnsoccerreports-testing.appspot.com", PDF.ReportID); err != nil {
        log.Println(err)
    }
}
