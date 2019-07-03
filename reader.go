package http

import (
	"time"

	"github.com/sclevine/agouti"
	"github.com/yosssi/gohtml"
	gentleman "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

// Config Struct...
// 系统配置
var Config = struct {
	Timeout  time.Duration
	Browser  string
	Platform string
}{
	Timeout:  10,
	Browser:  "chrome",
	Platform: "MAC",
}

// Load Function
// 加载Web ...
func Load(url string, navigate ...bool) (*Document, bool) {

	doc := &Document{}

	doc.URL = url

	if len(navigate) == 1 && navigate[0] {
		if doc, ok := doc.NavigateHTTPClient(); ok {
			return doc, ok
		}
	} else {
		if doc, ok := doc.TerminalHTTPClient(); ok {
			return doc, ok
		}
	}

	return nil, false
}

// TerminalHTTPClient Funciton
// 命令行下HTTP客户端 ...
func (doc *Document) TerminalHTTPClient() (*Document, bool) {
	cli := gentleman.New()

	cli.URL(doc.URL)

	req := cli.Request()

	// 可根据自己需求进行调整头部信息。
	req.Use(headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"))
	req.Use(headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"))
	req.Use(headers.Set("Accept-Language", "zh-CN,zh;q=0.9"))

	req.Use(timeout.Request(10 * time.Second))

	// Define dial specific timeouts
	req.Use(timeout.Dial(5*time.Second, 30*time.Second))

	resp, err := req.Send()

	if err != nil && !resp.Ok {
		return nil, false
	}

	if resp.StatusCode == 200 {
		doc.HTML = resp.String()
		doc.Bytes = resp.Bytes()
		return doc.Converter()
	}

	return nil, false
}

// NavigateHTTPClient Function
// 系统桌面HTTP客户端 ...
func (doc *Document) NavigateHTTPClient() (*Document, bool) {

	for {
		select {
		case <-time.After(Config.Timeout * time.Second):
			return nil, false
		case <-time.Tick(2 * time.Millisecond):
			driver := agouti.PhantomJS()

			defer driver.Stop()

			if err := driver.Start(); err != nil {
				return nil, false
			}

			capabilities := agouti.NewCapabilities().Browser(Config.Browser).Platform(Config.Platform).Without("javascriptEnabled")
			page, err := driver.NewPage(agouti.Desired(capabilities))

			if err != nil {
				return nil, false
			}

			page.Navigate(doc.URL)

			// time.Sleep(2 * time.Second)

			if doc.HTML, err = page.HTML(); err == nil {
				doc.HTML = gohtml.Format(doc.HTML)
				return doc, true
			}
		}
	}
}

/* End of file http.go */
/* Location: ./http.go */
