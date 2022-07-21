package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func SendHTTPMessage(jsonBody *[]byte, url string) (resp *http.Response) {
	var req *http.Request
	var err error
	if jsonBody != nil {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(*jsonBody))
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	client := http.Client{
		Transport: &http2.Transport{
			// So http2.Transport doesn't complain the URL scheme isn't 'https'
			AllowHTTP: true,
			// Pretend we are dialing a TLS endpoint.
			// Note, we ignore the passed tls.Config
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Client Proto: %d\n", resp.ProtoMajor)
	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	fmt.Println(str)
	return resp
}

func Run_HttpServer(pem, key string) {
	fmt.Printf("start http server")
	router := gin.New()
	AddService(router)
	router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "User-Agent",
			"Referrer", "Host", "Token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	h2Server := &http2.Server{
		// TODO: extends the idle time after re-use openapi client
		IdleTimeout: 1 * time.Millisecond,
	}

	server := &http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: h2c.NewHandler(router, h2Server),
	}

	if pem == "" || key == "" {
		server.ListenAndServe()
	} else {
		server.ListenAndServeTLS(pem, key)
	}

	/*    로그
	if preMasterSecretLogPath != "" {
		preMasterSecretFile, err := os.OpenFile(preMasterSecretLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return server, fmt.Errorf("create pre-master-secret log [%s] fail: %s", preMasterSecretLogPath, err)
		}
		server.TLSConfig = &tls.Config{
			KeyLogWriter: preMasterSecretFile,
		}
	}
	*/

}
