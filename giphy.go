package giphy

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
	"math/rand"
	"time"
	"os"
	"os/exec"
	"io/ioutil"
)

type Response struct{
	Data []Result `json:"data"`
	Meta map[string]interface{}
	Pagination map[string]interface{}
}
type Result struct{
	Bitly_gif_url string `json:"bitly_gif_url"`
	Bitly_url string `json:"bitly_url"`
	Content_url string `json:"content_url"`
	Embed_url string `json:"embed_url"`
	Id string `json:"id"`
	Images map[string]interface{} `json:"images"`
	Import_datetime string `json:"import_datetime"`
	Is_indexable int `json:"is_indexable"`
	Is_sticker int `json:"is_sticker"`
	Rating string `json:"rating"`
	Slug string `json:"slug"`
	Source string `json:"source"`
	Source_post_url string `json:"source_post_url"`
	Source_tld string `json:"source_tld"`
	Title string `json:"title"`
	Trending_datetime string `json:"trending_datetime"`
	Type string `json:"type"`
	Url string `json:"url"`
	User map[string]interface{} `json:"user"`
	Username string `json:"username"`
}
func Get(query string,apikey string,count int)Response{
	var result Response
	query=url.QueryEscape(query)
	url:=fmt.Sprintf("http://api.giphy.com/v1/gifs/search?q=%s&api_key=%s&limit=%d",query,apikey,count)
	resp,err:=http.Get(url)
	if err!=nil{
		fmt.Println("\t\t",err.Error())
		return result
	}
	defer resp.Body.Close()
	decoder:=json.NewDecoder(resp.Body)
	err=decoder.Decode(&result)
	if err!=nil{
		fmt.Println("\t\t",err.Error())
	}
	return result
}

func (result *Response) GetRandom()string{
	rand.Seed(time.Now().Unix())
	gifs:=result.Data
	if len(gifs)==0{
		return " "
	}
	item:=gifs[rand.Intn(len(gifs))]
	images:=item.Images
	dm:=images["downsized_medium"].(map[string]interface{})
	link:=dm["url"].(string)
	return link
}
func Download(url string)error{
	resp,err:=http.Get(url)
	if err!=nil{
		return err
	}
	defer resp.Body.Close()
	content,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return err
	}
	bytes:=[]byte(content)
	file,err:=os.Create("giphy.gif")
	if err!=nil{
		return err
	}
	_,err=file.Write(bytes)
	defer file.Close()
	if err!=nil{
		return err
	}
	return nil
}
func Open(url string)error{
	err:=Download(url)
	if err!=nil{
		return err
	}
	cmd:=exec.Command("xdg-open","giphy.gif")
	err=cmd.Run()
	if err!=nil{
		return err
	}
	return nil
}

func main(){
	key:="YOUR API KEY HERE"
	limit:=50//number of search results to be returned
	query:="happy"//search query
	result:=Get(query,key,limit)
	link:=result.GetRandom()
	Open(link)
}







