package main

import (
	"database/sql"
	"fmt"
	"log"
)

func main() {
	// sql.DB 객체 생성
	db, err := sql.Open("mssql", "server=(local);user id=sa;password=pwd;database=pubs")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// [단일 row 쿼리]
	// 하나의 Row를 갖는 SQL 쿼리 : QueryRow()
	var lname, fname string
	err = db.QueryRow("SELECT au_lname, au_fname FROM authors WHERE au_id='172-32-1176'").Scan(&lname, &fname)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fname, lname)

	// [복수 row 쿼리]
	// 복수 Row를 갖는 SQL 쿼리: Query()
	rows, err := db.Query("SELECT au_lname, au_fname FROM authors WHERE au_lname=?", "Ringer")
	// -> 복수 Row 쿼리 Query() 메서드에서 ? (placeholder)를 사용하여 'Parameterized Query'를 사용하고 있다.
	//		(placeholder ?에 문자열 "Ringer"가 대입)
	// SQL 인젝션과 같은 문제를 방지하기 위해 파라미터를 문자열 결합이 아닌 별도의 파라미터로 대입시키는 방식
	// placeholder는 데이터베이스와 그 드라이버 종류에 따라 다르다. MSSQL 드라이버는 '?', '?n', ':n', '$n' (n은 숫자) 등의 형식 지원
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	for rows.Next() {
		err := rows.Scan(&lname, &fname)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fname, lname)
	}

	// [DML과 Prepared Statement]
	// DML: DB 조작어(SELECT, INSERT, UPDATE, DELETE)
	// DDL: DB 정의어(CREATE, ARTER, DROP, RENAME, TRUNCATE)
	// TCL: 트랜잭션 제어어(COMMIT, ROLLBACK, SAVEPOINT), 논리적인 작업의 단위를 묶어 DML에 의해 조작된 결과를 작업 단위(트랜잭션) 별로 제어
	//
	// DML operation을 하기 위해 sql.DB 객체의 'Exec' 메서드를 사용한다.
	// DML과 같이 리턴되는 데이터가 없는 경우 'Exec' 메서드를 사용한다.
	//
	// Prepared Statement: DB 서버에 placeholder를 가진 SQL문을 미리 준비시키는 것으로,
	// 차후 해당 statement를 호출할 때 준비된 SQL 문을 빠르게 실행하도록 하는 기법임.
	// Go에서 Prepared Statement를 사용하기 위해 sql.DB의 'Prepare()' 메서드를 써서 Placeholder를 가진 SQL 문을 미리 준비시키고,
	// sql.Stmt 객체를 리턴받는다. 차후 이 sql.Stmt 객체의 Exec 혹은 Query/QueryRow 메서드를 사용하여 준비된 SQL문을 실행
	// INSERT 문 실행 : Exec()
	result, err := db.Exec("INSERT INTO discounts(discounttype,discount) VALUES(?, ?)", "Test", 10)
	if err != nil {
		log.Fatal(err)
	}

	n, err := result.RowsAffected() // sql.Result.RowsAffected() 체크
	if err != nil {
		log.Fatal(err)
	}
	if n == 1 {
		fmt.Println("Successfully inserted.")
	}

	// Prepared Statement 생성
	stmt, err := db.Prepare("UPDATE discounts SET discount=?2 WHERE discounttype=?1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Prepared Statement 실행
	_, err = stmt.Exec("Test", 15) //Placeholder 파라미터 ?1, ?2 전달
	if err != nil {
		log.Fatal(err)
	}

	// [MSSQL 트랜잭션]
	// 복수 개의 SQL 문을 하나의 트랜잭션으로 묶기 위하여 sql.DB의 'Begin()' 메서드를 사용
	// 트랜잭션은 복수 개의 SQL 문을 실행하다 중간에 어떤 한 SQL 문에서라도 에러가 발생하면, 전체 SQL문을 취소하게 되고 (롤백, rollback)
	// 모두 성공적으로 실행되어야 전체를 커밋하게 된다.
	// sql.Tx 타입의 Begin() 메서드는 sql.Tx 객체를 리턴, 마지막에 최종 Commit을 위해 Tx.COMMIT() 메서드를, 롤백을 위해 Tx.Rollback() 메서드를 호출함

	// 트랜잭션 시작
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback() //중간에 에러시 롤백

	// INSERT 문 실행
	_, err = db.Exec("INSERT discounts(discounttype,discount) VALUES(?, ?)", "Test1", 12)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT discounts(discounttype,discount) VALUES(?, ?)", "Test2", 11)
	if err != nil {
		log.Fatal(err)
	}

	// 트랜잭션 커밋
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
