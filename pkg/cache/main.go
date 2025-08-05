// gha_logger_proxy.go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
)

func main() {
	p := goproxy.NewProxyHttpServer()
	p.Verbose = false

	// ────────────────────────────────────────────────────────────────
	// 1. MITM only the GitHub pipelines host so we can read the path
	// ────────────────────────────────────────────────────────────────
	targetHost := "pipelines.actions.githubusercontent.com"
	p.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(
		func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			if strings.Contains(host, targetHost) {
				return goproxy.MitmConnect, host // decrypt
			}
			return goproxy.OkConnect, host // tunnel everything else
		}))

	// ────────────────────────────────────────────────────────────────
	// 2. Log request + response
	// ────────────────────────────────────────────────────────────────
	p.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		start := time.Now()

		// Ensure we get the response status later
		ctx.UserData = start
		fmt.Printf("[REQ]  %-6s %s\n", r.Method, r.URL.Path)
		return r, nil
	})

	p.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		start, _ := ctx.UserData.(time.Time)
		elapsed := time.Since(start)
		fmt.Printf("[RES]  %-6s %s → %d (%s)\n",
			resp.Request.Method, resp.Request.URL.Path, resp.StatusCode, elapsed)
		return resp
	})

	// ────────────────────────────────────────────────────────────────
	// 3. Allow self-signed root CA (you'll trust it on the runner VM)
	// ────────────────────────────────────────────────────────────────
	p.Tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	log.Println("proxy listening on :3128 …")
	log.Fatal(http.ListenAndServe(":3128", p))
}
