package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Query(queryStatement string) {
	fmt.Println("1")
	// MySQL 데이터베이스 정보
	// db, err := sql.Open("mysql", "root:myyonseinft@tcp(52.204.136.223:3306)/myyonseinft")
	db, err := sql.Open("mysql", "root:myyonseinft@tcp(localhost:3306)/myyonseinft")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("2")

	// 쿼리 실행
	rows, err := db.Query(queryStatement)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	fmt.Println("3")

	// 결과 반복문으로 처리
	for rows.Next() {
		var address string
		var department string
		var tokenuri string
		var tx string
		var time []uint8
		if err := rows.Scan(&address, &department, &tokenuri, &tx, &time); err != nil {
			panic(err.Error())
		}
		fmt.Println(address, department, tokenuri, tx, time)
	}
	// 오류 체크
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
}

func AddHistory(address string, department string, tokenURI string, tx string) {
	// MySQL 데이터베이스 연결 정보
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/myyonseinft")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 쿼리 실행
	result, err := db.Exec("INSERT INTO NFTs (address, department, tokenuri, tx) VALUES (?, ?, ?, ?)", address, department, tokenURI, tx)
	if err != nil {
		log.Fatal(err)
	}

	// 영향을 받은 행의 수 확인
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("성공적으로 %d개의 행이 추가되었습니다.\n", rowsAffected)

}
