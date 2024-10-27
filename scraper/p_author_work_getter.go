package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pixiv-downloader/config"
	"pixiv-downloader/utils"
)

func HTTP_GetAuthorWorks(conf config.Config, author utils.Author, imgURLCh chan utils.Work, threadLimitCh chan struct{}) error {
	//defer wg.Done()
	threadLimitCh <- struct{}{}

	authorURL := author.WorkURL
	utils.CreateAuthorDirectory(author.Author, conf.OutputDir)

	proxyURL := conf.ProxyURL
	headers := conf.Headers
	client := &http.Client{}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return fmt.Errorf("Error using proxy: %v", err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", authorURL, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error requesting: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error downloading image: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("\033[1;34m"+"Get Author: ", author.Author, "\033[0m")
	err = utils.HandleAuthorWorkJSON(imgURLCh, body, author.Author, conf.ImgQuality)
	if err != nil {
		return err
	}
	return nil
}

func HTTP_GetWorks(conf config.Config, work utils.Work, threadLimitCh chan struct{}) error {
	threadLimitCh <- struct{}{}

	workURL := "https://www.pixiv.net/ajax/illust/" + work.WorkID + "/pages"

	proxyURL := conf.ProxyURL
	headers := conf.Headers
	client := &http.Client{}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return fmt.Errorf("Error using proxy: %v", err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", workURL, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error requesting: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error downloading image: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ret, _ := utils.HandleWorkTypeJSON(body, conf.ImgQuality)
	for _, imgURL := range ret {
		err := HTTP_DownloadImage(conf, imgURL, work)
		if err != nil {
			return err
		}
	}
	//fmt.Println("Success Downloaded: ", work.WorkID)
	return nil
}
