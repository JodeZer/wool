package wool

import (
	"testing"
)

func TestDownload(t *testing.T) {
	// c := NewTFSImageClient(&TFSImageClientConfig{})
	// c.download("")

}

func TestWoolUpload(t *testing.T) {
	c := NewTFSImageClient(&TFSImageClientConfig{})
	str, err := c.UploadFromUrl("http://img02.taobaocdn.com/bao/uploaded/i2/TB1YA3QLXXXXXaLaXXXXXXXXXXX_!!0-item_pic.jpg")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", str)
	tfsid, b := c.ParseTBResp(str)
	t.Logf("%s %v", tfsid, b)
}
