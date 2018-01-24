package wool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
)

type SearchProducts struct {
	Nid       string `json:"nid"` // 疑似Alipid
	PicUrl    string `json:"pic_url"`
	DetailUrl string `json:"detail_url"`
	ViewPrice string `json:"view_price"` // 单价
}

type SearchParam struct {
	Commend      string `json:"commend"`
	Ssid         string `json:"ssid"`
	SearchType   string `json:"search_type"`
	SourceId     string `json:"sourceId"`
	Spm          string `json:"spm"`
	Ie           string `json:"ie"`
	InitiativeId string `json:"initiative_id"`
	Tfsid        string `json:"tfsid"`
	App          string `json:"app"`
	// commend=all&ssid=s5-e&search_type=item&sourceId=tb.index&spm=a21bo.2017.201856-taobao-item.2&
	// ie=utf8&initiative_id=tbindexz_20170306&tfsid=TB18EeldWzB9uJjSZFMXXXq4XXa&app=imgsearch
}

func newMagicSearchParam(tfsid string) *SearchParam {
	return &SearchParam{
		Commend:      "all",
		Ssid:         "s5-e",
		SearchType:   "item",
		SourceId:     "tb.index",
		Spm:          "a21bo.2017.201856-taobao-item.2",
		Ie:           "utf8",
		InitiativeId: "tbindexz_20170306",
		Tfsid:        tfsid,
		App:          "imgsearch",
	}
}

type TBSearchClient struct {
	cli *http.Client
}

type TBSearchClientConf struct {
}

func NewTBSearchClient(config *TBSearchClientConf) *TBSearchClient {
	return &TBSearchClient{
		cli: &http.Client{},
	}
}

func (c *TBSearchClient) SearchReturnRawString(tfsid string) (string, error) {
	buffer, err := c.search(tfsid)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (c *TBSearchClient) SearchReturnBuffer(tfsid string) (*bytes.Buffer, error) {
	return c.search(tfsid)
}

func (c *TBSearchClient) SearchReturnProduct(tfsid string, N int) ([]*SearchProducts, error) {
	resStr, err := c.SearchReturnRawString(tfsid)
	if err != nil {
		return nil, err
	}
	res := make([]*SearchProducts, 0, N)
	start, end := strings.Index(resStr, "g_page_config"), strings.Index(resStr, " g_srp_loadCss();")
	resStr = resStr[start:end]
	startJson, endJson := strings.Index(resStr, "{"), strings.LastIndex(resStr, "}")
	resStr = resStr[startJson:endJson]
	//fmt.Println(jsonStr)

	array := gjson.Get(resStr, "mods").Get("itemlist").Get("data").Get("collections").Array()
	productResults := array[0].Get("auctions").Array()
	for i := 0; i < N && i < len(productResults); i++ {
		prod := &SearchProducts{}
		fmt.Println(productResults[i].String())
		if err = json.Unmarshal([]byte(productResults[i].String()), prod); err != nil {
			return nil, err
		}
		res = append(res, prod)

	}
	return res, nil
}

func (c *TBSearchClient) search(tfsid string) (*bytes.Buffer, error) {
	req := c.newHttpRequest(tfsid)
	//fmt.Println(formatRequest(req))
	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%+v\n", resp.Request)
	defer resp.Body.Close()

	return decodeContentEncoding(resp.Body, resp.Header.Get("content-encoding"))
}

func (c *TBSearchClient) newHttpRequest(tfsid string) *http.Request {
	q := url.Values{}

	for k, v := range mustToMap(newMagicSearchParam(tfsid), "json") {
		q.Set(k, v)
	}

	req, _ := http.NewRequest(http.MethodGet, c.getUrl(), nil)

	c.setFixedHeader(req)

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
	return req

}

func (c *TBSearchClient) setFixedHeader(req *http.Request) {

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:57.0) Gecko/20100101 Firefox/57.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://www.taobao.com/")
	req.Header.Set("Cookie", "thw=cn; isg=AgkJZJnbt-vTqUsfGLiEvSToGTOj_jhtMDf0pat-UvAv8i0E86YNWPdwQmw7; t=aa07c40df3128fc8f7b8d1b879dd4fca; cookie2=3fbc3a7a033c7206d2d48a2cc1a40499; v=0; _tb_token_=ee36e5b8dbe5e; cna=lPHZEvqiqxECAbSpnbIL/joX; alitrackid=www.taobao.com; lastalitrackid=www.taobao.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
}

func (c *TBSearchClient) getUrl() string {
	return "https://s.taobao.com/search"
}

func mustToMap(in *SearchParam, tag string) map[string]string {
	out := make(map[string]string)

	v := reflect.ValueOf(in).Elem()

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {

		field := typ.Field(i)
		if tagval := field.Tag.Get(tag); tagval != "" {
			out[tagval] = v.Field(i).String()
		}
	}
	return out
}
