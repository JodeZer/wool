package wool

import "testing"

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
	searchRespString, err := sc.SearchReturnRawString("TB1BKfmdAfb_uJkSne1XXbE4XXa")
	t.Logf("%s %s", searchRespString, err)
}
