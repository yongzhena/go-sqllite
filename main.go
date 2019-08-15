package main

import (
	"log"
	"sqllite/models"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	log.Println("test sqllite...")

	// log.Println("InitRecAreas...")
	// err := models.InitRecAreas()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// log.Println("InitRecAreas ok!")
	//var opt models.RecAPI
	opt := models.NewRecAPI(true)
	err := opt.OpenRecAreas()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("OpenRecAreas ok!")
	// var rec models.Records
	// rec.RecType = 2
	// rec.Ext = "ext"
	// rec.Res = "Res"
	id, err := opt.SaveRec(models.RecArea01, []byte("123456789011"), 0)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("over,SaveRec ok!,area=%d,id=%d\n", 1, id)

	id, err = opt.SaveRec(1, []byte("1234567890221111"), 1)
	if err != nil {
		log.Println(err.Error())
	}
	id, err = opt.SaveRec(2, []byte("123456789022"), 1)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("over,SaveRec ok!,area=%d,id=%d\n", 2, id)
	id, err = opt.SaveRec(2, []byte("123456789022"), 3)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("over,SaveRec ok!,area=%d,id=%d\n", 2, id)

	num := opt.GetNoUploadNum(models.RecArea01)
	log.Printf("area=%d,NoUploadNum=%d\n", models.RecArea01, num)

	recp, err := opt.ReadRecNotServer(models.RecArea01, 1)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(recp)

	err = opt.DeleteRec(1, 1)
	if err != nil {
		log.Fatal(err.Error())
	}

	num = opt.GetNoUploadNum(models.RecArea01)
	log.Printf("area=%d,NoUploadNum=%d\n", models.RecArea01, num)

	num = opt.GetNoUploadNum(models.RecArea02)
	log.Printf("area=%d,NoUploadNum=%d\n", models.RecArea02, num)

}
