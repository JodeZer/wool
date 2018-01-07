package wool

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"sync"
)

type possibleTBImageRespJson struct {
	Status int    `json:"status"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Error  bool   `json:"error"`
}

// 匹配Json
var respCompiler = regexp.MustCompile(`.*({.*}).*`)

type TFSImageClient struct {
	cache sync.Map
	cli   *http.Client
}

type TFSImageClientConfig struct {
}

func NewTFSImageClient(config *TFSImageClientConfig) *TFSImageClient {
	if config == nil {
		return &TFSImageClient{
			cli: &http.Client{},
		}
	}

	return &TFSImageClient{
		cli: &http.Client{},
	}
}

func (c *TFSImageClient) UploadFromUrl(url string) (string, error) {
	if val, ok := c.cache.Load(c.getCacheKey(url, true)); ok {
		return val.(string), nil
	}

	imgContentReader, err := c.download(url)
	if err != nil {
		return "", err
	}

	return c.uploadData(imgContentReader, c.genFileName(url))

}

func (c *TFSImageClient) UploadFromFile(filename string) (string, error) {
	if val, ok := c.cache.Load(c.getCacheKey(filename, false)); ok {
		return val.(string), nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	return c.uploadData(f, filename)
}

func (c *TFSImageClient) SetCache(key string, url bool, val string) {
	c.cache.Store(c.getCacheKey(key, url), val)
}

// 碰运气
func (c *TFSImageClient) ParseTBResp(resp string) (string, bool) {
	var respIns possibleTBImageRespJson
	if jsonStr, suc := c.parseTBResp(resp); suc {
		if err := json.Unmarshal([]byte(jsonStr), &respIns); err != nil {
			return "", false
		} else if respIns.Name == "" {
			return "", false
		}
		return respIns.Name, true
	}
	return "", false
}

func (c *TFSImageClient) parseTBResp(resp string) (string, bool) {
	strs := respCompiler.FindStringSubmatch(resp)
	fmt.Printf("%s\n%v\n", resp, strs)
	if len(strs) == 0 {
		return "", false
	}

	if len(strs) <= 1 {
		return strs[0], false
	}

	return strs[1], true
}
func (c *TFSImageClient) getCacheKey(content string, url bool) string {
	if url {
		return "url:" + content
	}
	return "file:" + content
}

func (c *TFSImageClient) download(url string) (io.Reader, error) {
	//url = "http://img02.taobaocdn.com/bao/uploaded/i2/TB1YA3QLXXXXXaLaXXXXXXXXXXX_!!0-item_pic.jpg"

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	Respbytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(Respbytes)

	return buf, nil
}

func (c *TFSImageClient) uploadData(reader io.Reader, fileName string) (string, error) {
	//filename can't be ""
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	err := w.SetBoundary(c.randomBoundary())
	if err != nil {
		return "", nil
	}
	c.writeFixedContent(w)

	replacer := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

	h := make(textproto.MIMEHeader)

	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			replacer.Replace("imgfile"), replacer.Replace(fileName)))
	h.Set("Content-Type", "image/jpeg")

	fw, err := w.CreatePart(h)
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(fw, reader); err != nil {
		return "", err
	}

	w.Close()

	req, err := http.NewRequest("POST", c.getApiUrl(), &b)
	if err != nil {
		return "", err
	}

	c.setFixedHeader(req, w)

	resp, err := c.cli.Do(req)
	if err != nil {
		return "", err
	}

	return c.decodeResp(resp)

}

func (c *TFSImageClient) writeFixedContent(w *multipart.Writer) {
	// 不知道用来干嘛的
	if fw, err := w.CreateFormField("cross"); err != nil {
		panic(err)
	} else {
		fw.Write([]byte("taobao"))
	}

	if fw, err := w.CreateFormField("type"); err != nil {
		panic(err)
	} else {
		fw.Write([]byte("iframe"))
	}
}

func (c *TFSImageClient) randomBoundary() string {
	return "---------------------------8717257699615597491120257768"
}

func (c *TFSImageClient) genFileName(url string) string {
	return "aaa.jpg"
}

func (c *TFSImageClient) setFixedHeader(req *http.Request, w *multipart.Writer) {
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:57.0) Gecko/20100101 Firefox/57.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://s.taobao.com/search?&imgfile=&js=1&stats_click=search_radio_all%3A1&initiative_id=staobaoz_20180107&ie=utf8&tfsid=TB1rotcdzgy_uJjSZKPXXaGlFXa&app=imgsearch")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", "thw=cn; isg=Atzcah3H6r_bMp7jDiq-5JxQrvyEF469jIVKO7bdnEeqAX-L3mVQD1LnFVIH; cna=n2T+D/ji/R4CAbSk/5kCno8Q; t=74aea229d40005509fa3a1c845177403; cookie2=19aa555db345edeff33584ecd08a7d81; v=0; _tb_token_=e14d8356e83e8; JSESSIONID=A9319CE93A5A334500DBD92143AAD0BA; enc=3A8lqx8adlTiV3MlUl9q5R8tYXTV8npYqEzTabcLlaxk%2FKMmYTvCcPciv07T53gKFJ%2FaTgQZr3zspP8WiDq9%2BA%3D%3D; alitrackid=www.taobao.com; lastalitrackid=www.taobao.com; hng=CN%7Czh-CN%7CCNY%7C156; mt=ci%3D-1_0")

}

func (c *TFSImageClient) getApiUrl() string {
	return "https://s.taobao.com/image"
}

func (c *TFSImageClient) decodeResp(resp *http.Response) (string, error) {
	encoding := resp.Header.Get("content-encoding")
	return c.decodeRespContent(resp.Body, encoding)
}

func (c *TFSImageClient) decodeRespContent(respbody io.ReadCloser, encoding string) (string, error) {
	defer respbody.Close()

	var b bytes.Buffer

	if encoding == "gzip" {
		reader, err := gzip.NewReader(respbody)
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(&b, reader); err != nil {
			return "", err
		}
		return b.String(), nil
	}

	return "", errors.New(fmt.Sprintf("unknown content-encoding %s", encoding))
}
