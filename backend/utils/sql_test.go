package utils

import (
	"fmt"
	"regexp"
	"testing"
)

func TestSQL(t *testing.T) {

	// fmt.Println("sdfsdflskdjflskdjf")
	// Query("SELECT * FROM NFTs;")

	// tokenUri := "https://myyonseinft.s3.amazonaws.com/MAJOR/TEST/json/1.json"
	// // 정규 표현식 패턴
	// pattern := `https://myyonseinft.s3.amazonaws.com/MAJOR/([^/]+)/json/.*`

	// // 정규 표현식 컴파일
	// re := regexp.MustCompile(pattern)

	// // 문자열에서 패턴과 일치하는 모든 문자열 추출
	// // department := re.FindString(tokenUri)
	// words := re.FindAllString(tokenUri, -1)

	// // 추출된 단어 출력
	// for _, word := range words {
	// 	fmt.Println(word)
	// }

	urls := []string{
		"https://myyonseinft.s3.amazonaws.com/MAJOR/TEST/json/1.json",
		"https://myyonseinft.s3.amazonaws.com/MAJOR/TEST2/json/2.json",
		"https://myyonseinft.s3.amazonaws.com/MAJOR/BANANA/json/3.json",
		"https://myyonseinft.s3.amazonaws.com/MAJOR/TEST2/json/4.json",
	}
	// 정규 표현식 패턴
	pattern := `https://myyonseinft.s3.amazonaws.com/MAJOR/([^/]+)/json/.*`

	// 정규 표현식 컴파일
	re := regexp.MustCompile(pattern)

	// 각 URL에 대해 패턴과 일치하는 문자열 추출
	for _, url := range urls {
		match := re.FindStringSubmatch(url)
		if len(match) >= 2 {
			word := match[1]
			fmt.Println(word)
		}
	}

	// fmt.Println("department", department)
	// AddTxHistory("address", department, tokenUri, "txaddress")

}
