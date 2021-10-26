package main

import(
    "os"
    "io"
    "log"
)

func StartLogger() {
    LogFile, err := os.OpenFile("../logs.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        panic(err)
    }

    LogWriter := io.MultiWriter(os.Stdout, LogFile)

    log.SetOutput(LogWriter)

    log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}
