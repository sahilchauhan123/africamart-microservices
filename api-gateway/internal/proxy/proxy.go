package proxy

import (
	"api-gateway/internal/token"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyHandler is a middleware that proxies requests to the target microservice
func ProxyHandler(prefix string, microServiceUrl string, jwtMaker *token.JWTMaker) gin.HandlerFunc {

	fmt.Println("ProxyHandler called")
	target, err := url.Parse(microServiceUrl)
	if err != nil {
		log.Printf("Error parsing target URL: %v", err)
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
		// 1. Print URL BEFORE modification
		fmt.Printf("[Proxy] Original Path: %s\n", c.Request.URL.Path)

		// Changing the target url of the request
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, prefix)

		// 2. Print URL AFTER modification
		fmt.Printf("[Proxy] Forwarding to: %s%s\n", microServiceUrl, c.Request.URL.Path)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
