package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var token, hookUrl string

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") != token {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid token"))
		return
	}
	if r.FormValue("command") != "/fortune" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected command == \"/fortune\""))
	}
	text := r.FormValue("text")
	log.Printf("%s: /fortune %s\n", r.FormValue("user_name"), text)

	// --help prints the man page
	if text == "--help" {
		buf, err := exec.Command("./man-fortune.sh").CombinedOutput()
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		w.Write(buf)
		return
	}

	// Otherwise directly pass the arguments to fortune.
	cmd := exec.Command("fortune", strings.Split(text, " ")...)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	fmtStr := "Here's <@%s>'s fortune:\n```%s```"
	msg := fmt.Sprintf(fmtStr, r.FormValue("user_id"), buf)

	// We can't post incoming webhooks to direct messages or groups, so
	// just show it directly to the user.
	chanName := r.FormValue("channel_name")
	if chanName == "directmessage" || chanName == "privategroup" {
		w.Write([]byte("*NOTE: use /fortune in a public channel so everyone else can see!*\n"))
		w.Write([]byte(msg))
		return
	}

	// Send an incoming webhook to slack to post the fortune in the channel.
	payload, err := json.Marshal(map[string]string{
		"channel":    "#" + chanName,
		"username":   "fortune-bot",
		"icon_emoji": ":squirrel:",
		"text":       msg,
	})
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	resp, err := http.Post(hookUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("ERROR:", err)
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("ERROR:", string(body))
	}
}

func GetOutput(cmd *exec.Cmd) (buf []byte) {
	buf, err := cmd.CombinedOutput()
	if err != nil {
		buf = nil
		log.Println("ERROR:", err)
	}
	return
}

func main() {
	if token = os.Getenv("TOKEN"); token == "" {
		log.Fatal("TOKEN not specified")
	}
	if hookUrl = os.Getenv("HOOK_URL"); hookUrl == "" {
		log.Fatal("HOOK_URL not specified")
	}
	http.HandleFunc("/", Handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(":"+port, nil)
}
