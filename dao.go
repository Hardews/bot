package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Stu struct {
	Uid  int    `gorm:"primaryKey;AUTO_INCREMENT=1;not null"`
	Name string `gorm:"type:varchar(50)"`
	Xh   string `gorm:"type:varchar(50)"`
	QNum string `gorm:"type:varchar(50)"`
}

var (
	dB *gorm.DB
)

func InitDB() {
	dsn := "root:lmh123@tcp(127.0.0.1:3306)/stu?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}

	dB = db
}

//func WriteIn(stu Stu) error {
//	err := dB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Stu{})
//	if err != nil {
//		err = errors.New("create table failed,err :" + err.Error())
//		return err
//	}
//
//	tx := dB.Model(&Stu{}).Create(&stu)
//	if err := tx.Error; err != nil {
//		err = errors.New("insert data failed,err:" + err.Error())
//		return err
//	}
//
//	fmt.Println("successful")
//	return nil
//}

func SelectXh(id int) (Stu, error) {
	var stu Stu
	tx := dB.Model(&Stu{}).Where("uid = ?", id).First(&stu)
	if err := tx.Error; err != nil {
		return Stu{}, err
	}
	return stu, nil
}
