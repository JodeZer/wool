package wool

import (
	"testing"
)

func TestSearch(t *testing.T) {
	c := NewTFSImageClient(&TFSImageClientConfig{})
	str, err := c.UploadFromUrl("http://img02.taobaocdn.com/bao/uploaded/i2/TB1YA3QLXXXXXaLaXXXXXXXXXXX_!!0-item_pic.jpg")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", str)
	tfsid, b := c.ParseTBResp(str)
	t.Logf("%s %v", tfsid, b)
	sc := NewTBSearchClient(&TBSearchClientConf{})
	searchRespString, err := sc.SearchReturnRawString(tfsid)
	t.Logf("%s %s", searchRespString, err)
}

func TestSearchByTfsid(t *testing.T) {
	sc := NewTBSearchClient(&TBSearchClientConf{})
	searchRespBuffer, err := sc.SearchReturnBuffer("TB18EeldWzB9uJjSZFMXXXq4XXa")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s %s", searchRespBuffer.String(), err)
}

func TestParseHtml(t *testing.T) {
	// scli := NewTBSearchClient(nil)
	// str, err := scli.SearchReturnRawString("TB18EeldWzB9uJjSZFMXXXq4XXa")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(str)

	// start := strings.Index(str, "g_page_config")
	// end := strings.Index(str, " g_srp_loadCss();")
	// jsonStr := str[start:end]
	// fmt.Println(str)
}

func TestSearchReturnProducts(t *testing.T) {
	scli := NewTBSearchClient(nil)
	products, err := scli.SearchReturnProduct("TB18EeldWzB9uJjSZFMXXXq4XXa", 5)
	if err != nil {
		panic(err)
	}
	for _, one := range products {
		t.Logf("%+v", one)
	}

}
