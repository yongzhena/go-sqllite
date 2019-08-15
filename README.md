"# go-sqllite" 

https://blog.csdn.net/yyz_1987/article/details/98975108     
### 嵌入式linux应用之go语言开发，存储模块的接口封装，底层使用sqllite

在嵌入式终端设备上，免不了要存储记录、上传记录、查看记录等操作。 

我称之为储存模块。 

怎样的操作接口，最好用？最方便？ 

首先想到的是使用嵌入式数据库sqllite,没错，选他作为存储媒介，用go调用也是很方便的。  

但是，这还远远不够。原生的sql操作，若不做个封装，将会是很难用。   

另外，已经有很多ORM框架，即对象关系映射，将面向对象语言程序中的对象自动持久化到关系数据库中。  

就满足要求了吗？   

这也还不够。   

我想要的接口，能满足这样的功能：    

可以写入记录，删除记录，查询记录。   

删除记录不是真正的删除，而是清除上传标记，即该记录还存在，只是表示一上传过，没用了，可以被覆盖。    

记录数量达到一定条数时要从头循环覆盖。    

记录可支持同时向不同的第三方平台上送。     

操作记录接口要简单和灵活，比如添加记录中的字段不能再去动表结构。    

操作记录的接口:    
    // 初始化记录区(会清空所有数据!)         
	InitRecAreas() error     
	// 打开记录区(开机必须先打开一次)   
	OpenRecAreas() (err error)       
	// 保存记录   
	SaveRec(areaID int, buf []byte, recType int) (id int64, err error)    
	// 删除记录     
	DeleteRec(areaID int, num int64) (err error)    
	// 获取未上传记录数量    
	GetNoUploadNum(areaID int) int     
	// 按数据库ID读取一条记录    
	ReadRecByID(areaID int, id int) (p *Records, err error)    
	// 顺序读取未上传的记录    
	ReadRecNotServer(areaID int, sn int) (p *Records, err error)   
	// 倒数读取记录（如sn=1代表最后一次写入的记录）  
	ReadRecWriteNot(areaID int, sn int) (p *Records, err error)  
	// 最后一条记录流水  
	GetLastRecNO(areaID int) int  
	
简单使用demo:  

首次应用启动需先调用InitRecAreas()，完成初始化创建表的操作。后续不再调用  
每次应用启动，先调用OpenRecAreas()，打开记录存储区  

要存储一条记录，只需SaveRec(存储区id，数据内容，记录类型)  

要读取一条未上传的记录，只需ReadRecNotServer()  

要删除一条记录，只需 DeleteRec()  

要获取未上传记录数量，只需GetNoUploadNum()  

func main() {

	log.Println("test sqllite...")

	// log.Println("InitRecAreas...")
	// err := models.InitRecAreas()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// log.Println("InitRecAreas ok!")
	opt := models.NewRecAPI(true)
	err := opt.OpenRecAreas()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("OpenRecAreas ok!")

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

