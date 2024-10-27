package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"pixiv-downloader/config"
	"pixiv-downloader/utils"
	"sync"
)

func MainScrape(conf config.Config, authorlis []string) {
	authorCh := make(chan utils.Author, conf.AuthorChanLen)
	imgURLCh := make(chan utils.Work, conf.ImgChanLen)
	retryCh := make(chan utils.Work, conf.RetryChanLen)
	threadLimitCh := make(chan struct{}, conf.Threads)

	var authorWg sync.WaitGroup
	var workWg sync.WaitGroup
	var retryWg sync.WaitGroup

	fmt.Println("\n" + "\033[1;34m" + "Getting Author Metadata" + "\033[0m")
	//barAuthors := progressbar.New(len(authorlis))
	barAuthors := progressbar.NewOptions(len(authorlis), progressbar.OptionSetDescription(">"))

	for _, aurl := range authorlis {
		author, err := HTTP_DownloadAuthorMetadata(conf, aurl, threadLimitCh)
		if err != nil {
			log.Fatal("Error Handling Author URL:", author.Author)
		}
		authorCh <- author
		barAuthors.Add(1)
	}
	close(authorCh)

	// fetch author works
	fmt.Println("\n" + "\033[1;34m" + "Fetching Author Works" + "\033[0m")
	barWorks := progressbar.NewOptions(len(authorCh), progressbar.OptionSetDescription(">"))

	for author := range authorCh {
		authorWg.Add(1)
		go func(author utils.Author) {
			defer authorWg.Done()
			defer func() { <-threadLimitCh }()

			err := HTTP_GetAuthorWorks(conf, author, imgURLCh, threadLimitCh)
			if err != nil {
				log.Fatal("Error Handling Author:", author.Author)
			}
			barWorks.Add(1)
		}(author)
	}
	authorWg.Wait()
	close(imgURLCh)

	// fetch images in single work
	fmt.Println("\n"+"\033[1;34m"+"Fetching Images... Total Images:", len(imgURLCh), "\033[0m")
	barImages := progressbar.NewOptions(len(imgURLCh), progressbar.OptionSetDescription(">"))

	for img := range imgURLCh {
		workWg.Add(1)
		go func(work utils.Work) {
			defer workWg.Done()
			defer func() { <-threadLimitCh }()

			err := HTTP_GetWorks(conf, work, threadLimitCh)
			if err != nil {
				retryCh <- work
				fmt.Println("\n"+"\033[1;33m"+"Error Handling Work:", work.WorkID, "Will retry for one more time"+"\033[0m")
			}
			barImages.Add(1)
		}(img)
	}
	workWg.Wait()
	close(retryCh)

	fmt.Println("\n"+"\033[1;34m"+"Retrying", len(retryCh), "Downloads..."+"\033[0m")
	barRetry := progressbar.NewOptions(len(retryCh), progressbar.OptionSetDescription(">"))

	failedflag := 0
	// retry once on failed downloads
	for img := range retryCh {
		retryWg.Add(1)
		go func(work utils.Work) {
			defer retryWg.Done()
			defer func() { <-threadLimitCh }()

			err := HTTP_GetWorks(conf, work, threadLimitCh)
			if err != nil {
				workJSON, err := json.Marshal(work)
				log.Printf(string(workJSON), err, "\n")
				failedflag += 1
			}
			barRetry.Add(1)
		}(img)
	}
	retryWg.Wait()

	authorWg.Wait()
	workWg.Wait()
	retryWg.Wait()
	fmt.Println("\n" + "\033[1;34m" + "Download Finished!" + "\033[0m")
	if failedflag > 0 {
		fmt.Println("\n"+"\033[1;34m"+"Failed", failedflag, "Imgs, See", conf.ErrorLogDir, "for Details"+"\033[0m")
	}
}
