package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command(fmt.Sprintf("%s/bin/server", exec.Command(os.Getenv("GITHUB_WORKSPACE"))))
	cmd.Dir = os.Getenv("GITHUB_WORKSPACE")
	cmd.Env = []string{
		"APP_ID=163095",
		"INSTALLATION_ID=22021531",
		"PRIVATE_KEY=private.key",
		"GH_OWNER=dev-this",
	}

	log.Printf("test run working directory: %s\n", cmd.Dir)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	cmd.Stderr = cmd.Stdout

	reader := bufio.NewReader(cmdReader)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			var line string
			var err error

			if line, err = reader.ReadString('\n'); err == io.EOF {
				break
			}

			fmt.Print(line)
		}
	}()

	time.Sleep(1 * time.Second)

	file, err := os.Open("payload")
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest("POST", "http://localhost:8080/webhook", file)
	req.Header.Set("X-Github-Event", "workflow_run")
	req.Header.Set("x-hub-signature", "1")
	req.Header.Set("x-github-delivery", "1")
	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		os.Exit(255)
	}

	os.Exit(0)
}
