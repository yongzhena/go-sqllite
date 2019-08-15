package models

import (
	"errors"
	"fmt"
	"log"
	db "sqllite/database"
)

// RecDir ...
type RecDir struct {
	ID      int   `json:"id"`
	RecNo   int   `json:"recno" `
	WriteID int64 `json:"writeid" `
	ReadID1 int64 `json:"readid1" `
	ReadID2 int64 `json:"readid2" `
	ReadID3 int64 `json:"readid3" `
	Rp      int   `json:"rp" `
	Res     int   `json:"res" `
	Flag    bool  `json:"flag" `
}

// InitRecDir ...
func InitRecDir() (err error) {
	//创建表
	sqlTable := `
	DROP TABLE IF EXISTS tb_dir;
    CREATE TABLE IF NOT EXISTS tb_dir(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		recNo INTEGER NOT NULL,
		writeID INTEGER NOT NULL,
		readID1 INTEGER  NOT NULL,
		readID2 INTEGER ,
		readID3 INTEGER ,
		rp INTEGER ,
		res INTEGER
    );
	`
	if db.SQLDB == nil {
		err = errors.New("db.SQLDB is null")
		log.Fatal(err.Error())
		return
	}
	log.Println("begin create dir table...")
	if IsDebug {
		log.Println("sql：" + sqlTable)
	}
	_, err = db.SQLDB.Exec(sqlTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("create dir table ok!")
	//清空数据
	// log.Println("begin truncate dir table...")
	// _, err = db.SQLDB.Exec(`UPDATE sqlite_sequence SET seq = 0 WHERE name = 'tb_dir' `)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// log.Println("truncate dir table ok!")
	log.Println("begin init dir table...")
	for i := 0; i < MAXRECDIRS; i++ {
		_, err = db.SQLDB.Exec("INSERT INTO tb_dir(recNo, writeID,readID1,readID2,readID3,rp,res) VALUES (?, ?,?,?,?,?,?)", 0, 0, 0, 0, 0, 0, 0)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println("init dir table ok!")
	return err
}

// UpdateDirs 更新目录
func (rd *RecDir) UpdateDirs(areaID int) error {

	strSQL := fmt.Sprintf("UPDATE tb_dir SET recNo=%d, writeID=%d, readID1=%d, readID2=%d, readID3=%d, rp=%d,res=%d WHERE id=%d",
		rd.RecNo, rd.WriteID, rd.ReadID1, rd.ReadID2, rd.ReadID3, rd.Rp, rd.Res, areaID)
	if IsDebug {
		log.Println(strSQL)
	}
	res, err := db.SQLDB.Exec(strSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	affect, err := res.RowsAffected()
	fmt.Println(affect)
	if IsDebug {
		log.Println(rd)
	}
	return err

}

// LoadDirs 加载(读取)目录
func (rd *RecDir) LoadDirs(areaID int) error {

	strSQL := fmt.Sprintf("SELECT * FROM tb_dir WHERE id=%d", areaID)
	if IsDebug {
		log.Println(strSQL)
	}
	rows, err := db.SQLDB.Query(strSQL)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if rows.Next() {
		err = rows.Scan(&rd.ID, &rd.RecNo, &rd.WriteID, &rd.ReadID1, &rd.ReadID2, &rd.ReadID3, &rd.Rp, &rd.Res)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
	} else {
		log.Fatal("no dir records")
	}

	rows.Close()
	if IsDebug {
		log.Println(rd)
	}
	return err
}
