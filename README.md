# PixivDownloader
A Go-based Pixiv image downloader that supports batch input of author homepage URLs and multi-threaded image downloading for efficient and fast downloads.

- Support multi-thread downloading
- Support for downloading images of different quality


## Usage
1. Modify `config.json`: 
```json
{
  "OutputDir": "your-image-output-directory", // image output directory
  "Threads": 1, // download threads
  "Headers": {
    "Referer": "https://www.pixiv.net/", // do not modify 
    "User-Agent": "your-user-agent", // change to your UA
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
    "Cookie": "your-cookie:PHPSESSID=xxxxx" // change to your Cookie (PHPSESSID)
  },
  "ReqType": "http",
  "AuthorChanLen": 128, // length of author channel. Increase length when there are many authors.
  "ImgChanLen": 1024, // length of image channel. Increase length when there are many images.
  "RetryChanLen": 1024, // length of retry channel. Increase length when there are bad internet connections.
  "ProxyURL": "", // change to your proxy
  "ImgQuality": "regular", // 4 types of image qualities: "thumb_mini", "small", "regular" and "original"
  "ErrorLogDir":  "your-error-log-output-directory" // error log directory
}
```
2. Put the artist's homepage to be downloaded in a `.txt` file (e.g. `author.txt`), split by line.
``` plain
// List author's homepage URL here. e.g. https://www.pixiv.net/users/someuser
https://www.pixiv.net/users/someartist
```
3. Run main.go and download
```shell
go run main.go -file ./author.txt -config ./config.json
```

## Warning

Excessive use may cause your IP be banned by Pixiv.