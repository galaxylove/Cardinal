package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Bulletin struct {
	gorm.Model

	Title   string
	Content string
}

type BulletinRead struct {
	gorm.Model

	TeamID     uint
	BulletinID uint
}

func (s *Service) GetAllBulletins() (int, interface{}) {
	var bulletins []Bulletin
	s.Mysql.Model(&Bulletin{}).Find(&bulletins)
	return s.makeSuccessJSON(bulletins)
}

func (s *Service) NewBulletin(c *gin.Context) (int, interface{}) {
	type InputForm struct {
		Title   string `binding:"required"`
		Content string `binding:"required"`
	}
	var inputForm InputForm
	err := c.BindJSON(&inputForm)
	if err != nil {
		return s.makeErrJSON(400, 40000, "Error payload")
	}

	tx := s.Mysql.Begin()
	if tx.Create(&Bulletin{
		Title:   inputForm.Title,
		Content: inputForm.Content,
	}).RowsAffected != 1 {
		tx.Rollback()
		return s.makeErrJSON(500, 50000, "添加公告失败！")
	}
	tx.Commit()
	return s.makeSuccessJSON("添加公告成功！")
}

func (s *Service) EditBulletin(c *gin.Context) (int, interface{}) {
	type InputForm struct {
		ID      uint   `binding:"required"`
		Title   string `binding:"required"`
		Content string `binding:"required"`
	}
	var inputForm InputForm
	err := c.BindJSON(&inputForm)
	if err != nil {
		return s.makeErrJSON(400, 40000, "Error payload")
	}

	var checkBulletin Bulletin
	s.Mysql.Where(&Bulletin{Model: gorm.Model{ID: checkBulletin.ID}}).Find(&checkBulletin)
	if checkBulletin.ID == 0 {
		return s.makeErrJSON(404, 40400, "公告不存在")
	}

	newBulletin := &Bulletin{
		Title:   inputForm.Title,
		Content: inputForm.Content,
	}
	tx := s.Mysql.Begin()
	if tx.Model(&Bulletin{}).Where(&Bulletin{Model: gorm.Model{ID: inputForm.ID}}).Updates(&newBulletin).RowsAffected != 1 {
		tx.Rollback()
		return s.makeErrJSON(500, 50001, "修改公告失败！")
	}
	tx.Commit()

	return s.makeSuccessJSON("修改公告成功！")
}

func (s *Service) DeleteBulletin(c *gin.Context) (int, interface{}) {
	type InputForm struct {
		ID uint `binding:"required"`
	}
	var inputForm InputForm
	err := c.BindJSON(&inputForm)
	if err != nil {
		return s.makeErrJSON(400, 40000, "Error payload")
	}

	var checkBulletin Bulletin
	s.Mysql.Where(&Bulletin{Model: gorm.Model{ID: checkBulletin.ID}}).Find(&checkBulletin)
	if checkBulletin.ID == 0 {
		return s.makeErrJSON(404, 40400, "公告不存在")
	}

	tx := s.Mysql.Begin()
	if tx.Where("id = ?", inputForm.ID).Delete(&Bulletin{}).RowsAffected != 1 {
		tx.Rollback()
		return s.makeErrJSON(500, 50002, "删除公告失败！")
	}
	tx.Commit()
	return s.makeSuccessJSON("删除公告成功！")
}