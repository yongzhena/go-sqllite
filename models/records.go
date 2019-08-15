package models

import (
	"errors"
	"fmt"
	"log"
	db "sqllite/database"
	"strings"
	"time"
)

var (
	//IsDebug 是否调试
	IsDebug = true
	recDir  [MAXRECAREAS]RecDir
)

// Records ...
type Records struct {
	ID      int    `json:"id"`
	RecNo   int    `json:"recno" `
	RecType int    `json:"rectype" `
	RecTime string `json:"rectime" `
	Data    []byte `json:"data" `
	Ext     string `json:"ext" `
	Res     string `json:"res" `
}

// InitRecAreas 初始化记录存储区
func (rec Records) InitRecAreas() error {

	//初始化目录表
	err := InitRecDir()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	//创建记录表
	sqlTable := `
	DROP TABLE IF EXISTS TB_NAME;
    CREATE TABLE IF NOT EXISTS TB_NAME (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		recNo INTEGER NOT NULL,
		recType  INTEGER NOT NULL,
		recTime INTEGER NOT NULL,
		data BLOB ,
		ext TEXT ,
		res TEXT
    );
	`
	for i := 0; i < MAXRECAREAS; i++ {
		tbName := fmt.Sprintf("tb_rec%02d", i+1)
		log.Println("begin create rec table " + tbName)
		sqls := strings.Replace(sqlTable, "TB_NAME", tbName, -1)
		if IsDebug {
			log.Println("sql：" + sqls)
		}
		_, err = db.SQLDB.Exec(sqls)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		log.Println("create rec table " + tbName + " ok!")

	}
	return err
}

// OpenRecAreas 打开记录存储区,每次开机，需要先打开一下
func (rec Records) OpenRecAreas() (err error) {
	//加载RecDir
	for i := 0; i < MAXRECAREAS; i++ {
		log.Printf("LoadDirs %02d \n", i+1)
		err = recDir[i].LoadDirs(i + 1)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Printf("LoadDirs %02d ok!\n", i+1)
	}
	//log.Println(recDir)

	return err
}

// SaveRec 保存记录
func (rec *Records) SaveRec(areaID int, buf []byte, recType int) (id int64, err error) {

	log.Printf("SaveRec,area=%02d \n", areaID)
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = fmt.Errorf("area id  %02d is not right,mast between 1 and %02d", areaID, MAXRECAREAS)
		log.Println(err.Error())
		return
	}
	rec.RecNo = recDir[areaID-1].RecNo
	t := time.Now()
	rec.RecTime = t.Format("20060102150405")
	rec.Data = buf
	rec.RecType = recType
	//记录是否存储满，判断
	if (recDir[areaID-1].WriteID + 1) > (int64)(MAXRECCOUNTS) {

		if recDir[areaID-1].ReadID1 == 0 {

			err = fmt.Errorf("rec area %02d is full", areaID)
			log.Println(err.Error())
			return
		}

		if (recDir[areaID-1].WriteID + 1 - int64(MAXRECCOUNTS)) == recDir[areaID-1].ReadID1 {
			err = fmt.Errorf("rec area %02d is full", areaID)
			log.Println(err.Error())
			return
		}

		//保存记录
		strSQL := fmt.Sprintf(`UPDATE tb_rec%02x SET recNo=%d, recType=%d,recTime=%s,data=?,ext="%s",res="%s" WHERE id = 1`,
			areaID, rec.RecNo+1, rec.RecType, rec.RecTime, rec.Ext, rec.Res)
		if IsDebug {
			log.Println(strSQL)
		}
		_, err = db.SQLDB.Exec(strSQL, rec.Data)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		recDir[areaID-1].RecNo++
		recDir[areaID-1].WriteID = 1
		recDir[areaID-1].Flag = true
		id = 1
		err = recDir[areaID-1].UpdateDirs(areaID)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		log.Printf("SaveRec,area=%02d ok!\n", areaID)
		return id, err
	}

	if recDir[areaID-1].Flag {
		//记录是否满判断
		if (recDir[areaID-1].WriteID + 1) == recDir[areaID-1].ReadID1 {
			err = fmt.Errorf("rec area %02d is full", areaID)
			log.Println(err.Error())
			return
		}
		id = recDir[areaID-1].WriteID + 1
		strSQL := fmt.Sprintf(`UPDATE tb_rec%02x SET recNo=%d, recType=%d,recTime=%s,data=?,ext="%s",res="%s" WHERE id = %d`,
			areaID, rec.RecNo+1, rec.RecType, rec.RecTime, rec.Ext, rec.Res, id)
		if IsDebug {
			log.Println(strSQL)
		}
		_, err = db.SQLDB.Exec(strSQL, rec.Data)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		recDir[areaID-1].RecNo++
		recDir[areaID-1].WriteID = id
		err = recDir[areaID-1].UpdateDirs(areaID)
		if err != nil {
			log.Fatal(err.Error())
			return 0, err
		}
		log.Printf("SaveRec,area=%02d ok!\n", areaID)
		return id, err

	}

	strSQL := fmt.Sprintf(`INSERT INTO tb_rec%02x(recNo, recType,recTime,data,ext,res) VALUES (%d,%d,%s,?,"%s","%s")`,
		areaID, rec.RecNo+1, rec.RecType, rec.RecTime, rec.Ext, rec.Res)
	if IsDebug {
		log.Println(strSQL)
	}
	rs, err := db.SQLDB.Exec(strSQL, rec.Data)
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}
	id, err = rs.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}
	recDir[areaID-1].RecNo++
	recDir[areaID-1].WriteID = id
	err = recDir[areaID-1].UpdateDirs(areaID)
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}
	log.Printf("SaveRec,area=%02d ok!\n", areaID)
	return id, err

}

// DeleteRec 删除记录（并不是真正删除表里记录，而是清除该记录的上传标记）
// areaID:记录区 num:删除的数量
func (rec Records) DeleteRec(areaID int, num int64) (err error) {
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = errors.New("area id is not right")
		log.Fatal(err.Error())
		return
	}

	id := recDir[areaID-1].ReadID1

	//如果写的位置等于读的位置，说明记录已上传完，没有要删除的了
	if recDir[areaID-1].WriteID == recDir[areaID-1].ReadID1 {
		return
	}

	//如果要删除的数量大于了最大的记录数
	if (id + num) > MAXRECCOUNTS {
		if (id + num - MAXRECCOUNTS) > recDir[areaID-1].WriteID {
			recDir[areaID-1].ReadID1 = recDir[areaID-1].WriteID
			err = recDir[areaID-1].UpdateDirs(areaID)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			return
		}
		//更新读指针（读的位置）
		recDir[areaID-1].ReadID1 = id + num - MAXRECCOUNTS
		err = recDir[areaID-1].UpdateDirs(areaID)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		return
	}

	//如果当前写的位置大于读的位置
	if recDir[areaID-1].WriteID > recDir[areaID-1].ReadID1 {
		if id+num > recDir[areaID-1].WriteID {
			//更新读指针（读的位置）
			recDir[areaID-1].ReadID1 = recDir[areaID-1].WriteID
			err = recDir[areaID-1].UpdateDirs(areaID)
			if err != nil {
				log.Fatal(err.Error())
				return err
			}
			return
		}
	}

	//更新读指针（读的位置）
	recDir[areaID-1].ReadID1 = id + num
	err = recDir[areaID-1].UpdateDirs(areaID)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return
}

//GetNoUploadNum 获取未上传记录数量
func (rec Records) GetNoUploadNum(areaID int) int {

	num := 0
	if recDir[areaID-1].WriteID == recDir[areaID-1].ReadID1 {
		num = 0
		return num
	}
	if recDir[areaID-1].Flag == false {
		num = int(recDir[areaID-1].WriteID - recDir[areaID-1].ReadID1)
	} else {
		if recDir[areaID-1].WriteID > recDir[areaID-1].ReadID1 {
			num = int(recDir[areaID-1].WriteID - recDir[areaID-1].ReadID1)
		} else {
			num = int(MAXRECCOUNTS - recDir[areaID-1].ReadID1 + recDir[areaID-1].WriteID)
		}
	}
	return num
}

// ReadRecByID 按数据库ID读取记录
func (rec Records) ReadRecByID(areaID int, id int) (p *Records, err error) {
	var rec1 Records
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = errors.New("area id is not right")
		log.Fatal(err.Error())
		return
	}
	strSQL := fmt.Sprintf("SELECT * FROM tb_rec%02d WHERE id=%d", areaID, id)
	if IsDebug {
		log.Println(strSQL)
	}
	rows, err := db.SQLDB.Query(strSQL)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	if rows.Next() {
		err = rows.Scan(&rec1.ID, &rec1.RecNo, &rec1.RecType, &rec1.RecTime, &rec1.Data, &rec1.Ext, &rec1.Res)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
	} else {
		log.Println("no records")
		return nil, err
	}
	rows.Close()
	return &rec1, nil
}

//ReadRecNotServer 读取未上传的记录数据，顺序读取第SN条未上传的记录
//sn取值 1-到-->未上传记录数目
func (rec Records) ReadRecNotServer(areaID int, sn int) (p *Records, err error) {
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		err = errors.New("area id is not right")
		log.Fatal(err.Error())
		return
	}
	id := recDir[areaID-1].ReadID1
	if (int(id) + sn) > MAXRECCOUNTS {
		if int(id)+sn-MAXRECCOUNTS > int(recDir[areaID-1].WriteID) {
			return nil, errors.New("no records")
		}
		p, err = rec.ReadRecByID(areaID, int(id)+sn-MAXRECCOUNTS)
	} else {
		if recDir[areaID-1].ReadID1 < recDir[areaID-1].WriteID {
			if (int(id) + sn) > int(recDir[areaID-1].WriteID) {
				return nil, errors.New("no records")
			}
			p, err = rec.ReadRecByID(areaID, int(recDir[areaID-1].ReadID1)+sn)
		}

	}
	return p, err
}

// ReadRecWriteNot 倒数读取第SN条写入的记录
//读取一条记录  倒数读取第SN条写入的记录
func (rec Records) ReadRecWriteNot(areaID int, sn int) (p *Records, err error) {
	id := int(recDir[areaID-1].WriteID)
	if (id - sn) < 0 {
		if recDir[areaID-1].Flag {
			p, err = rec.ReadRecByID(areaID, MAXRECCOUNTS-(sn-id-1))
		} else {
			return nil, errors.New("no records")
		}
	} else {
		p, err = rec.ReadRecByID(areaID, (id - sn + 1))
	}
	return
}

// GetLastRecNO 获取最后一条记录流水号
func (rec Records) GetLastRecNO(areaID int) int {
	if (areaID <= 0) || (areaID > MAXRECAREAS) {
		log.Println("area id is not right")
		return 0
	}
	id := recDir[areaID-1].RecNo
	return id
}

// NewRecords ...
func NewRecords(debug bool) *Records {
	IsDebug = debug
	records := new(Records)
	return records
}
