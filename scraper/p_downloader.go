package scraper

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"pixiv-downloader/config"
	"pixiv-downloader/utils"
	"strings"
)

func Colly_DownloadImage(config config.Config, retryCh chan utils.Work, imageURLCh chan utils.Work) error {
	return nil

}

func HTTP_DownloadImage(conf config.Config, imgURL string, work utils.Work) error {
	proxyURL := conf.ProxyURL
	headers := conf.Headers
	author := work.Author
	quality := work.ImgFiness

	dir := filepath.Join(conf.OutputDir, author, quality)

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

	req, err := http.NewRequest("GET", imgURL, nil)
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
	contentType := resp.Header.Get("Content-Type")
	filename := path.Base(imgURL)
	parts := strings.Split(filename, ".")
	filename = parts[0]
	ext := "img"
	if strings.HasPrefix(contentType, "image/") {
		ext = "." + strings.TrimPrefix(contentType, "image/")
	}
	dir = filepath.Join(dir, filename+ext)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(dir, body, 0755)
	if err != nil {
		return err
	}
	return nil
}

func HTTP_DownloadAuthorMetadata(conf config.Config, aurl string, threadLimitCh chan struct{}) (utils.Author, error) {
	threadLimitCh <- struct{}{}

	var badAuthor = utils.Author{}
	proxyURL := conf.ProxyURL
	headers := conf.Headers
	client := &http.Client{}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return badAuthor, fmt.Errorf("Error using proxy: %v", err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", aurl, nil)
	if err != nil {
		return badAuthor, fmt.Errorf("Error creating request: %v", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		return badAuthor, fmt.Errorf("Error requesting: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return badAuthor, fmt.Errorf("Error downloading image: status code %d", resp.StatusCode)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return badAuthor, fmt.Errorf("Parsing Error")
	}

	node := htmlquery.FindOne(doc, "/html/head/title/text()")
	if node != nil {
		//fmt.Println("XPath Result:", htmlquery.InnerText(node))
		aurl = strings.TrimSuffix(aurl, "/")
		partscode := strings.Split(aurl, "/")
		partsname := strings.Split(htmlquery.InnerText(node), " - ")
		if len(partscode) > 0 && len(partsname) > 0 {
			badAuthor.AuthorID = partscode[len(partscode)-1]
			badAuthor.Author = partsname[0]
			badAuthor.WorkURL = "https://www.pixiv.net/ajax/user/" + badAuthor.AuthorID + "/profile/all"
			return badAuthor, nil
		} else {
			return badAuthor, fmt.Errorf("Author Parsing Error")
		}
	} else {
		return badAuthor, fmt.Errorf("No result found for the given XPath")
	}

}
