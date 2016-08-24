package main

import(
//	"io"
	"os"
	"github.com/ushis/m3u"
	"fmt"
	"net"
	"strings"
//	"github.com/ziutek/gst"
	"net/http"
	"bufio"
)
func bracketHost(originalPath string) string{
	newPath := "["
	colon := strings.LastIndex(originalPath, ":")
	newPath += originalPath[7:colon]+ "]" + originalPath[colon:]
	fmt.Println(newPath)
	return newPath
}

func getTrack(resp *http.Response, globalRes http.ResponseWriter){
	fmt.Println("get track")
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
	f, err := os.Open("comedy.m3u")

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
		conn,err := net.Dial("tcp", bracketHost(track.Path))
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		http.HandleFunc("/comedy", func(res http.ResponseWriter, req *http.Request){
			fmt.Println("hitting comedy")
			resp, _:= http.Get(track.Path)
			fmt.Println("get "+track.Path)
			getTrack(resp, res)
		})
		fmt.Println(http.ListenAndServe(":8080", nil))

	}
}
