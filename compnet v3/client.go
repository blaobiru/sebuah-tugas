package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"main/tools"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func waitServer(url string, duration time.Duration) bool {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)

		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
	}
	return false
}

func main() {

	if !waitServer("https:localhost:9876", 5*time.Second) {
		//jika selama 5 detik tidak ada server, client diberhentikan
		fmt.Println("server tidak ada")
		return
	}

	var choice int
	for {
		fmt.Println("main menu")
		fmt.Println("1. get message")
		fmt.Println("2. send file")
		fmt.Println("3. print TLS")
		fmt.Println("4. quit")
		fmt.Print(">>")
		fmt.Scanf("%d\n", &choice)

		if choice == 1 {
			getMessage(client)
		} else if choice == 2 {
			sendFile(client)
		} else if choice == 3 {
			printTLSinfo(client)
		} else if choice == 4 {
			break
		}else {
			fmt.Println("invalid choice")
		}
	}
}

func getMessage() {
	client := createTLSClient()
	resp, err := client.Get("https://localhost:9876")
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	tools.ErrorHandler(err)

	fmt.Println("Server: ", string(data))
}

func sendFile() {
	var name string
	var age int

	scanner := bufio.NewReader(os.Stdin)

	fmt.Print("input name : ")
	name, _ = scanner.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Input age : ")
	fmt.Scanf("%d\n", &age)

	person := data.Person{Name: name, Age: age}

	//data person diencode ke JSON
	jsonData, err := json.Marshal(person)
	tools.ErrorHandler(err)

	temp := new(bytes.Buffer)

	w := multipart.NewWriter(temp)

	personField, err := w.CreateFormField("Person")
	tools.ErrorHandler(err)

	_, err = personField.Write(jsonData)
	tools.ErrorHandler(err)

	file, err := os.Open("./file.txt")
	tools.ErrorHandler(err)
	defer file.Close()

	filefield, err := w.CreateFormFile("File", file.Name())
	tools.ErrorHandler(err)

	_, err = io.Copy(filefield, file)
	tools.ErrorHandler(err)

	err = w.Close()
	tools.ErrorHandler(err)

	req, err := http.NewRequest("POST", "https://localhost:9876/sendFile", temp)
	tools.ErrorHandler(err)

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	tools.ErrorHandler(err)

	fmt.Println("Server : ", string(data))
}

