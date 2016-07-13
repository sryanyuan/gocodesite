package gocodecc

import (
	"database/sql"

	"github.com/astaxie/beego/orm"
)

var (
	projectCategoryItemTableName = "project_category_item"
)

type ProjectCategoryItem struct {
	Id              int    `orm:"pk;auto;index"`
	ProjectName     string `orm:"size(128)"`
	ProjectDescribe string `orm:"size(256)"`
	Image           string `orm:"size(256)"`
	Author          string `orm:"size(50)"`
	ItemCount       int
}

func (this *ProjectCategoryItem) TableName() string {
	return projectCategoryItemTableName
}

func init() {
	orm.RegisterModel(new(ProjectCategoryItem))
}

func modelProjectCategoryGetAll() ([]*ProjectCategoryItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query("SELECT * FROM " + projectCategoryItemTableName); nil != err {
		return nil, err
	}

	//	free the conn
	defer rows.Close()

	//	get all item
	resultSet := make([]*ProjectCategoryItem, 0, 10)
	for rows.Next() {
		item := &ProjectCategoryItem{}
		if err = rows.Scan(&item.Id, &item.ProjectName, &item.ProjectDescribe, &item.Image, &item.Author, &item.ItemCount); nil != err {
			return nil, err
		}
		resultSet = append(resultSet, item)
	}

	return resultSet, nil
}

func modelProjectCategoryAdd(item *ProjectCategoryItem) error {
	o := orm.NewOrm()
	_, err := o.Insert(item)
	return err
}

func modelProjectCategoryRemove(id int) error {
	item := ProjectCategoryItem{Id: id}
	o := orm.NewOrm()
	_, err := o.Delete(&item)
	return err
}
