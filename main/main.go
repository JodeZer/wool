package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"net/textproto"
	"compress/gzip"
)

func Upload(url, file string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	err = w.SetBoundary("---------------------------8717257699615597491120257768")
	if err != nil {
		panic(err)
	}
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()

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

	replacer := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			replacer.Replace("imgfile"), replacer.Replace("nofile.jpg")))
	h.Set("Content-Type", "image/jpeg")

	fw, err := w.CreatePart(h)
	if err != nil {
		panic(err)
		return
	}

	//b.WriteString("Content-Type: image/jpeg\r\n")

	if _, err = io.Copy(fw, f); err != nil {
		panic(err)
		return
	}

	// Add the other fields
	//if fw, err = w.CreateFormField("key"); err != nil {
	//	panic(err)
	//	return
	//}
	//if _, err = fw.Write([]byte("KEY")); err != nil {
	//	panic(err)
	//	return
	//}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		panic(err)
		return
	}


	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:57.0) Gecko/20100101 Firefox/57.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Accept-Encoding","gzip, deflate, br")
	req.Header.Set("Referer", "https://s.taobao.com/search?&imgfile=&js=1&stats_click=search_radio_all%3A1&initiative_id=staobaoz_20180107&ie=utf8&tfsid=TB1rotcdzgy_uJjSZKPXXaGlFXa&app=imgsearch")
	req.Header.Set("X-Requested-With","XMLHttpRequest")
	req.Header.Set("Cookie","thw=cn; isg=Atzcah3H6r_bMp7jDiq-5JxQrvyEF469jIVKO7bdnEeqAX-L3mVQD1LnFVIH; cna=n2T+D/ji/R4CAbSk/5kCno8Q; t=74aea229d40005509fa3a1c845177403; cookie2=19aa555db345edeff33584ecd08a7d81; v=0; _tb_token_=e14d8356e83e8; JSESSIONID=A9319CE93A5A334500DBD92143AAD0BA; enc=3A8lqx8adlTiV3MlUl9q5R8tYXTV8npYqEzTabcLlaxk%2FKMmYTvCcPciv07T53gKFJ%2FaTgQZr3zspP8WiDq9%2BA%3D%3D; alitrackid=www.taobao.com; lastalitrackid=www.taobao.com; hng=CN%7Czh-CN%7CCNY%7C156; mt=ci%3D-1_0")


	// Submit the request
	client := &http.Client{}

	fmt.Println(formatRequest(req))
	fmt.Println(b.String())
	res, err := client.Do(req)
	fmt.Println(formatRequest(req))
	if err != nil {
		panic(err)
		return
	}

	defer res.Body.Close()
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	
	reader, _ := gzip.NewReader(res.Body)

	io.Copy(os.Stdout, reader)

	respBytes, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(respBytes))
	fmt.Printf("%+v\n",res.Header)
	fmt.Printf("%v\n",respBytes)
	return
}

var url = "https://s.taobao.com/image"
func main() {
	Upload(url, "test.jpg")
}


///*
//
//curl 'https://s.taobao.com/image' --2.0
//-H 'Host: s.taobao.com'
//-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:57.0) Gecko/20100101 Firefox/57.0'
//-H 'Accept: application/json, text/javascript, */*; q=0.01'
//-H 'Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2'
//--compressed
//-H 'Referer: https://s.taobao.com/search?&imgfile=&js=1&stats_click=search_radio_all%3A1&initiative_id=staobaoz_20180106&ie=utf8&tfsid=TB1rotcdzgy_uJjSZKPXXaGlFXa&app=imgsearch'
//// -H 'X-Requested-With: XMLHttpRequest'
//// -H 'Content-Type: multipart/form-data;
//// boundary=---------------------------8717257699615597491120257768'
//// -H 'Cookie: thw=cn; isg=Atzcah3H6r_bMp7jDiq-5JxQrvyEF469jIVKO7bdnEeqAX-L3mVQD1LnFVIH; cna=n2T+D/ji/R4CAbSk/5kCno8Q; t=74aea229d40005509fa3a1c845177403; cookie2=19aa555db345edeff33584ecd08a7d81; v=0; _tb_token_=e14d8356e83e8; JSESSIONID=A9319CE93A5A334500DBD92143AAD0BA; enc=3A8lqx8adlTiV3MlUl9q5R8tYXTV8npYqEzTabcLlaxk%2FKMmYTvCcPciv07T53gKFJ%2FaTgQZr3zspP8WiDq9%2BA%3D%3D; alitrackid=www.taobao.com; lastalitrackid=www.taobao.com; hng=CN%7Czh-CN%7CCNY%7C156; mt=ci%3D-1_0'
//// -H 'Connection: keep-alive' --data ''
//
//*/

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
	r.ParseForm()
	request = append(request, "\n")
	request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}