package main

import(
	"os"
	"github.com/ushis/m3u"
	"fmt"
	"strings"
	"net/http"
	"bufio"
)

const port = "8080"

var shows = [...]string{ "comedy", "drama", "brain", "etc", "scifi" }

func bracketHost(originalPath string) string{
	newPath := "["
	colon := strings.LastIndex(originalPath, ":")
	newPath += originalPath[7:colon]+ "]" + originalPath[colon:]
	return newPath
}

func getTrack(resp *http.Response, globalRes http.ResponseWriter){
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		sendMessage(resp, line, globalRes)
	}
}

func sendMessage(resp *http.Response, c []byte, res http.ResponseWriter){
	if f, ok:=res.(http.Flusher); ok {
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

func main(){

	for _, show := range shows {
		f, err := os.Open(show + ".m3u")

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
			http.HandleFunc("/" + show, func(res http.ResponseWriter, req *http.Request){
				fmt.Println("hitting " + req.URL.Path)
				resp, _:= http.Get(track.Path)
				getTrack(resp, res)
			})

		}
	}
	fmt.Println("Listening to "+port)
	fmt.Println(http.ListenAndServe(":"+port, nil))
}
