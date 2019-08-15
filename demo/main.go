package main

import (
	"log"

	rec "github.com/yongzhena/go-sqllite"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	log.Println("test sqllite...")

	log.Println("InitRecAreas...")
	opt := rec.NewRecAPI(true)
	err := opt.InitRecAreas()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("InitRecAreas ok!")

	err = opt.OpenRecAreas()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("OpenRecAreas ok!")

	id, err := opt.SaveRec(1, []byte("123456789011"), 0)
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

	num := opt.GetNoUploadNum(1)
	log.Printf("area=%d,NoUploadNum=%d\n", 1, num)

	recp, err := opt.ReadRecNotServer(1, 1)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(recp)

	err = opt.DeleteRec(1, 1)
	if err != nil {
		log.Fatal(err.Error())
	}

	num = opt.GetNoUploadNum(1)
	log.Printf("area=%d,NoUploadNum=%d\n", 1, num)

	num = opt.GetNoUploadNum(2)
	log.Printf("area=%d,NoUploadNum=%d\n", 2, num)

}
