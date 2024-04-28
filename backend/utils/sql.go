package utils

import (
	"database/sql"
	"fmt"

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
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			panic(err.Error())
		}
		fmt.Println(id, name)
	}
	// 오류 체크
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
}
