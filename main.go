package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
    "html/template"

	"github.com/ushis/m3u"
)

var port = "8080"

var shows = [...]string{"comedy", "drama", "brain", "etc", "scifi"}
var line = []byte{}

func mainHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("main.html")
	t.Execute(w, nil)
}

func bracketHost(originalPath string) string {
	newPath := "["
	colon := strings.LastIndex(originalPath, ":")
	newPath += originalPath[7:colon] + "]" + originalPath[colon:]
	return newPath
}

func getTrack(resp *http.Response, globalRes http.ResponseWriter, req *http.Request) {
	reader := bufio.NewReader(resp.Body)
    for {
        line, _ = reader.ReadBytes('\n')
        //fmt.Println(line)
        sendMessage(resp, line, globalRes)
    }
}

func sendMessage(resp *http.Response, c []byte, res http.ResponseWriter) {
    if f, ok := res.(http.Flusher); ok {
        w := c
        for k, v := range resp.Header {
            res.Header().Set(k, v[0])
        }
        res.Write(w)
        f.Flush()
    } else {
        fmt.Println("no flush")
    }
}

func main() {
    port_env := os.Getenv("PORT")
    if port_env == "" {
      port = port_env
    }

    for _, show := range shows {
        f, err := os.Open("m3u/"+show + ".m3u")

        if err != nil {
            panic(err)
        }
        defer f.Close()

        pl, err := m3u.Parse(f)
        if err != nil {
            panic(err)
        }

        for _, track := range pl {

            //track.Path
            http.HandleFunc("/"+show, func(res http.ResponseWriter, req *http.Request) {
                fmt.Println("hitting " + req.URL.Path)
                resp, _ := http.Get(track.Path)
                getTrack(resp, res, req)
            })

        }
    }

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.HandleFunc("/", mainHandler)
    fmt.Println("Listening to " + port)
    fmt.Println(http.ListenAndServe(":"+port, nil))
}
