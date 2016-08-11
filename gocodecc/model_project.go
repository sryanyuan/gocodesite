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
	ProjectName     string `orm:"size(128);unique"`
	ProjectDescribe string `orm:"size(256)"`
	Image           string `orm:"size(256)"`
	Author          string `orm:"size(50)"`
	ItemCount       int
	PostPriv        uint32
}

func (this *ProjectCategoryItem) TableName() string {
	return projectCategoryItemTableName
}

var (
	projectArticleItemTableName = "project_article_item"
)

type ProjectArticleItem struct {
	Id                     int    `orm:"pk;auto;index"`
	ProjectName            string `orm:"size(128)"`
	ArticleTitle           string `orm:"size(128)"`
	ArticleContentHtml     string `orm:"size(12800)"`
	ArticleContentMarkdown string `orm:"size(12800)"`
	ArticleAuthor          string `orm:"size(50)"`
	Like                   int
	PostTime               int64
	EditTime               int64
	Look                   int
	ReplyAuthor            string `orm:"size(50)"`
	ReplyTime              int64
	ActiveTime             int64
	Click                  int
}

func (this *ProjectArticleItem) TableName() string {
	return projectArticleItemTableName
}

func init() {
	orm.RegisterModel(new(ProjectCategoryItem))
	orm.RegisterModel(new(ProjectArticleItem))
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
		if err = rows.Scan(&item.Id,
			&item.ProjectName,
			&item.ProjectDescribe,
			&item.Image,
			&item.Author,
			&item.ItemCount,
			&item.PostPriv); nil != err {
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

func modelProjectCategoryRemoveByProjectName(projectName string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("DELETE FROM "+projectCategoryItemTableName+" WHERE project_name = ?", projectName)
	if nil != err {
		return err
	}
	return nil
}

func modelProjectCategoryGetByProjectName(projectName string, prj *ProjectCategoryItem) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	err = db.QueryRow("SELECT * FROM "+projectCategoryItemTableName+" WHERE project_name = ?", projectName).Scan(
		&prj.Id,
		&prj.ProjectName,
		&prj.ProjectDescribe,
		&prj.Image,
		&prj.Author,
		&prj.ItemCount,
		&prj.PostPriv)
	if nil != err {
		return err
	}
	return nil
}

func modelProjectCategoryUpdateProject(prj *ProjectCategoryItem) error {
	o := orm.NewOrm()
	_, err := o.Update(prj, "project_name", "project_describe")
	return err
}

/*
	Project article
*/
func modelProjectArticleNewArticle(article *ProjectArticleItem) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(article)
}

func modelProjectArticleGet(articleId int) (*ProjectArticleItem, error) {
	o := orm.NewOrm()
	var article ProjectArticleItem
	article.Id = articleId
	err := o.Read(&article)
	if nil != err {
		return nil, err
	}
	return &article, nil
}

//	items, total page count
func modelProjectArticleGetArticles(project string, page int, limit int) ([]*ProjectArticleItem, int, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, 0, err
	}

	var rows *sql.Rows
	if rows, err = db.Query("SELECT id,article_title,article_author,post_time,reply_author,reply_time,active_time FROM "+projectArticleItemTableName+
		" WHERE project_name = ? ORDER BY active_time LIMIT ? OFFSET ?", project, limit, page*limit); nil != err {
		return nil, 0, err
	}

	//	free the conn
	defer rows.Close()

	//	get all item
	resultSet := make([]*ProjectArticleItem, 0, limit)
	for rows.Next() {
		item := &ProjectArticleItem{}
		if err = rows.Scan(
			&item.Id,
			&item.ArticleTitle,
			&item.ArticleAuthor,
			&item.PostTime,
			&item.ReplyAuthor,
			&item.ReplyTime,
			&item.ActiveTime); nil != err {
			return nil, 0, err
		}
		item.ProjectName = project
		resultSet = append(resultSet, item)
	}

	//	get page count
	pageResult, err := db.Exec("SELECT COUNT(*) FROM "+projectArticleItemTableName+" WHERE project_name = ?", project)
	if nil != err {
		return nil, 0, err
	}

	pageValue, _ := pageResult.RowsAffected()
	pages := (int(pageValue) + limit - 1) / limit

	return resultSet, pages, nil
}
