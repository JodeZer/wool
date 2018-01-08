package wool

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
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
	reader, err := c.search(tfsid)
	if err != nil {
		return "", err
	}
	var bs []byte
	if _, err := reader.Read(bs); err != nil {
		return "", err
	}
	return string(bs), nil
}

func (c *TBSearchClient) SearchReturnReader(tfsid string) (io.Reader, error) {
	return c.search(tfsid)
}

func (c *TBSearchClient) SearchReturnProduct(tfsid string) ([]*SearchProducts, error) {
	return nil, nil
}

func (c *TBSearchClient) search(tfsid string) (io.Reader, error) {
	req := c.newHttpRequest(tfsid)

	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", resp)
	fmt.Printf("%+v\n", resp.Request)
	defer resp.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}
	fmt.Printf("%+s\n", buf.String())
	return &buf, nil
}

func (c *TBSearchClient) newHttpRequest(tfsid string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, c.getUrl(), nil)
	query := req.URL.Query()
	for k, v := range mustToMap(newMagicSearchParam(tfsid), "json") {
		query.Set(k, v)
	}
	req.URL.RawQuery = query.Encode()

	return req

}

func (c *TBSearchClient) getUrl() string {
	return "https://s.taobao.com?search"
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
