package utils

import (
	"fmt"
	"net/http"
	"net/url"
)

func VerifyProxyURL(proxyURL string) http.Client {
	client := http.Client{}
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Errorf("Error using proxy: %v", err)
		fmt.Print("Continue without proxy")
		return http.Client{}
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client.Transport = transport
	return client
}
