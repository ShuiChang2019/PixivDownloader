package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"pixiv-downloader/config"
	"pixiv-downloader/scraper"
	"strings"
)

func main() {

	filename := flag.String("file", "./authors.txt", "Author URL files: The author's homepage URL, split by line")
	configdir := flag.String("config", "./config.json", "Config file")
	//help := flag.Bool("help", false, "Show Help")
	//debug := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	fmt.Println("\033[1;34m" + " ____  __  _  _  __  _  _    ____   __   _  _  __ _  __     __    __   ____  ____  ____ " + "\033[0m")
	fmt.Println("\033[1;34m" + "(  _ \\(  )( \\/ )(  )/ )( \\  (    \\ /  \\ / )( \\(  ( \\(  )   /  \\  / _\\ (    \\(  __)(  _ \\" + "\033[0m")
	fmt.Println("\033[1;34m" + " ) __/ )(  )  (  )( \\ \\/ /   ) D ((  O )\\ /\\ //    // (_/\\(  O )/    \\ ) D ( ) _)  )   /" + "\033[0m")
	fmt.Println("\033[1;34m" + "(__)  (__)(_/\\_)(__) \\__/   (____/ \\__/ (_/\\_)\\_)__)\\____/ \\__/ \\_/\\_/(____/(____)(__\\_)" + "\033[0m")
	fmt.Println("\033[1;35m" + "An image downloader written in Go. -> https://github.com/ShuiChang2019/PixivDownloader" + "\033[0m")
	fmt.Println("\033[1;35m" + "Use -h parameter for help" + "\033[0m")

	//log.Fatal("exit")

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal("Parsing author list error")
	}
	defer file.Close()
	var aurls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "http") {
			aurls = append(aurls, line)
		}
	}

	conf, err := config.LoadConfig(*configdir)
	if err != nil {
		log.Fatal("Parsing config file error")
	}

	logFile, err := os.OpenFile("error.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to create or open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	fmt.Println("\033[1;34m" + "Starting Download... Using Config:" + "\033[0m")
	fmt.Printf("OutputDir: %s\n", conf.OutputDir)
	fmt.Printf("Threads: %d\n", conf.Threads)
	fmt.Println("Headers:")
	for key, value := range conf.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Printf("ReqType: %s\n", conf.ReqType)
	fmt.Printf("AuthorChanLen: %d\n", conf.AuthorChanLen)
	fmt.Printf("ImgChanLen: %d\n", conf.ImgChanLen)
	fmt.Printf("RetryChanLen: %d\n", conf.RetryChanLen)
	fmt.Printf("ProxyURL: %s\n", conf.ProxyURL)
	fmt.Printf("ImgQuality: %s\n", conf.ImgQuality)
	fmt.Printf("ErrorLogDir: %s\n", conf.ErrorLogDir)

	scraper.MainScrape(conf, aurls)
}
