package main

import (
        "fmt"
        "log"
        "os"
        "os/exec"
        "regexp"
        "strings"
        "time"
        "github.com/fsnotify/fsnotify"
)

func main() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    done := make(chan bool)
    lastTime := time.Now()
    var lastFile string
    go func() {
        for {
            select {
            case e := <-watcher.Events:
                if lastTime.Add(time.Second/4).After(time.Now()) && lastFile == e.Name {
                    continue
                }

                m, _ := regexp.MatchString(strings.Replace(os.Args[1], "*", ".*", -1), e.Name)

                if m == false || (e.Op&fsnotify.Create != fsnotify.Create && e.Op&fsnotify.Write != fsnotify.Write) {
                    continue
                }
                args := make([]string, 0)
                for i, _ := range os.Args[2:] {
                    split := strings.Split(os.Args[i+2], " ")
                    for _, k := range split {
                      args = append(args, strings.Replace(k, "{}", e.Name, -1))
                    }
                }

                ran := strings.Join(args, " ")
                out, err := exec.Command("/bin/bash", "-c", ran).CombinedOutput()
                log.Println(e, ":", ran)
                fmt.Printf("%s\n", out)

                if err != nil {
                    log.Println(err)
                }
                lastTime = time.Now()
                lastFile = e.Name

            case err := <-watcher.Errors:
                fmt.Print(err)
            }
        }
    }()

    err = watcher.Add(".")
    log.Println("matching files", os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    <-done
}
