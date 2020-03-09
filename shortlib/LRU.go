/*************************************************************************
	> File Name: LRU.go
	> Author: Wu Yinghao
	> Mail: wyh817@gmail.com
	> Created Time: 一  6/15 17:07:37 2015
 ************************************************************************/

package shortlib

import (
	"container/list"
	"errors"
	//	"fmt"
)

type UrlPair struct {
	Original string
	Short    string
}

type LRU struct {
	HashShortUrl  map[string]*list.Element
	HashOriginUrl map[string]*list.Element
	ListUrl       *list.List
	RedisCli      *RedisAdaptor
}

func NewLRU(redis_cli *RedisAdaptor) (*LRU, error) {

	lru := &LRU{make(map[string]*list.Element), make(map[string]*list.Element), list.New(), redis_cli}
	return lru, nil
}

func (this *LRU) GetOriginalURL(shortUrl string) (string, error) {

	element, ok := this.HashShortUrl[shortUrl]
	//没有找到key,从Redis获取
	if !ok {
		originalUrl, err := this.RedisCli.GetUrl(shortUrl)
		//Redis也没有相应的短连接，无法提供服务
		if err != nil {
			return "", errors.New("No URL")
		}
		//更新LRU缓存
		ele := this.ListUrl.PushFront(UrlPair{originalUrl, shortUrl})
		this.HashShortUrl[shortUrl] = ele
		this.HashOriginUrl[originalUrl] = ele
		return originalUrl, nil
	}

	//调整位置
	this.ListUrl.MoveToFront(element)
	ele, _ := element.Value.(UrlPair)
	return ele.Original, nil
}

func (this *LRU) GetShortURL(originalUrl string) (string, error) {

	element, ok := this.HashOriginUrl[originalUrl]
	//没有找到key，返回错误，重新生成url
	if !ok {
		return "", errors.New("No URL")
	}

	//调整位置
	this.ListUrl.MoveToFront(element)
	ele, _ := element.Value.(UrlPair)
	/*
		fmt.Printf("Short_Url : %v \n",shortUrl)

		for iter:=this.ListUrl.Front();iter!=nil;iter=iter.Next(){
			fmt.Printf("Element:%v ElementAddr:%v\n",iter.Value,iter)
		}
		fmt.Printf("\n\n\n")
		for key,value := range this.HashUrl{
			fmt.Printf("Key:%v ==== Value:%v\n",key,value)
		}
	*/
	return ele.Short, nil

}

func (this *LRU) SetURL(originalUrl, shortUrl string) error {

	ele := this.ListUrl.PushFront(UrlPair{originalUrl, shortUrl})
	this.HashShortUrl[shortUrl] = ele
	this.HashOriginUrl[originalUrl] = ele
	//将数据存入Redis中
	//fmt.Printf("SET TO REDIS :: short : %v ====> original : %v \n",shortUrl,originalUrl)
	err := this.RedisCli.SetUrl(shortUrl, originalUrl)
	if err != nil {
		return err
	}
	return nil

}

func (this *LRU) checkList() error {

	return nil
}
