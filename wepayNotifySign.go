package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/clbanning/mxj"
	"reflect"
	"sort"
	"strings"
)

func main() {
	//这里是微信支付通知那边的DEMO数据，移除了优惠券这部分的数据
	data := `<xml>
  <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
  <attach><![CDATA[支付测试]]></attach>
  <bank_type><![CDATA[CFT]]></bank_type>
  <fee_type><![CDATA[CNY]]></fee_type>
  <is_subscribe><![CDATA[Y]]></is_subscribe>
  <mch_id><![CDATA[10000100]]></mch_id>
  <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
  <openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
  <out_trade_no><![CDATA[1409811653]]></out_trade_no>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
  <time_end><![CDATA[20140903131540]]></time_end>
  <total_fee>1</total_fee>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
</xml>`

	key := "商户的支付密码"

	mv, _ := mxj.NewMapXml([]byte(data))
	thisData := mv["xml"].(map[string]interface{})
	//移除sign，长度减一
	stringArr := make(map[string]string, len(thisData)-1)
	for k, v := range thisData {
		if k == "sign" {
			continue
		}
		//源码里面这里的type支持的很简单好像只有string/float64/bool,但是float64我没有试出来，这里都是string
		fmt.Println(reflect.TypeOf(v))
		switch v.(type) {
		//如果有试出来其他类型的可以加在这里
		default:
			stringArr[k] = v.(string)

		}
	}
	//有sign字段
	fmt.Println(thisData)
	//排序的时候要移除sign
	fmt.Println(stringArr)

	//这里是基于MD5的签名算法 其他的自己改
	sign,_:=SignByMD5(stringArr,key)
	if sign != thisData["sign"].(string){
		//验证不通过
		fmt.Println(sign)
	}

}

func SignByMD5(data map[string]string, key string) (string, error) {

	var query []string
	for k, v := range data {
		query = append(query, k+"="+v)
	}

	sort.Strings(query)
	query = append(query, "key="+key)
	str := strings.Join(query, "&")

	str, err := MD5(str)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(str), nil
}
func MD5(str string) (string, error) {
	hs := md5.New()
	if _, err := hs.Write([]byte(str)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hs.Sum(nil)), nil
}
