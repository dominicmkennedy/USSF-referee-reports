package main

import (
        "context"
        "fmt"
        "io"
        "os"
        "time"
        "log"

        "cloud.google.com/go/storage"
        "google.golang.org/api/option"
)


func uploadFile(bucket, object string) error {
        
        ctx := context.Background()
        opt := option.WithCredentialsFile("creds.json")
        
        client, err := storage.NewClient(context.Background(), opt)
        if err != nil {
                return fmt.Errorf("storage.NewClient: %v", err)
        }
        defer client.Close()

        // Open local file.
        f, err := os.Open("./reports/" + object + ".pdf")
        if err != nil {
                return fmt.Errorf("os.Open: %v", err)
        }
        defer f.Close()

        ctx, cancel := context.WithTimeout(ctx, time.Second*50)
        defer cancel()

        wc := client.Bucket(bucket).Object(object + ".pdf").NewWriter(ctx)
        if _, err = io.Copy(wc, f); err != nil {
                return fmt.Errorf("io.Copy: %v", err)
        }

        if err := wc.Close(); err != nil {
                return fmt.Errorf("Writer.Close: %v", err)
        }
        
        return nil
}

func StorePDF(form *refereeReport) {
        err := uploadFile ( "tnsoccerreports-testing.appspot.com", form.ReportID)
        if err != nil {
                log.Println(err)
        }
}
