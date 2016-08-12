package gocodecc

import (
	"database/sql"
	"errors"

	"github.com/astaxie/beego/orm"
)

var (
	projectCategoryItemTableName = "project_category_item"
)

const (
	kCategoryNameLimit     = 128
	kCategoryDescribeLimit = 256
	kCategoryImageLimit    = 256
)

type ProjectCategoryItem struct {
	Id              int    `orm:"pk;auto;index"`
	ProjectName     string `orm:"size(128);unique"`
	ProjectDescribe string `orm:"size(256)"`
	Image           string `orm:"size(256)"`
	Author          string `orm:"size(20)"`
	ItemCount       int
	PostPriv        uint32
}

func (this *ProjectCategoryItem) TableName() string {
	return projectCategoryItemTableName
}

var (
	projectArticleItemTableName = "project_article_item"
)

const (
	kArticleTitleLimit   = 128
	kArticleContentLimit = 12800
)

type ProjectArticleItem struct {
	Id                     int    `orm:"pk;auto;index"`
	ProjectName            string `orm:"size(128)"`
	ArticleTitle           string `orm:"size(128)"`
	ArticleContentHtml     string `orm:"size(12800)"`
	ArticleContentMarkdown string `orm:"size(12800)"`
	ArticleAuthor          string `orm:"size(20)"`
	Like                   int
	PostTime               int64
	EditTime               int64
	Top                    int
	ReplyAuthor            string `orm:"size(50)"`
	ReplyTime              int64
	ActiveTime             int64
	Click                  int
	ProjectId              int
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
	db, err := getRawDB()
	if nil != err {
		return err
	}

	tx, err := db.Begin()
	if nil != err {
		return err
	}

	//	delete category
	ret, err := tx.Exec("DELETE FROM "+projectCategoryItemTableName+" WHERE id = ?", id)
	if nil != err {
		tx.Rollback()
		return err
	}
	affected, err := ret.RowsAffected()
	if nil != err {
		tx.Rollback()
		return err
	}
	if affected != 1 {
		tx.Rollback()
		return errors.New("delete category failed")
	}

	//	remove all articles
	_, err = tx.Exec("DELETE FROM "+projectArticleItemTableName+" WHERE project_id = ?", id)
	if nil != err {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func modelProjectCategoryRemoveByProjectName(projectName string) error {
	o := orm.NewOrm()
	var category ProjectCategoryItem
	category.ProjectName = projectName
	err := o.Read(&category, "project_name")
	if nil != err {
		return err
	}

	return modelProjectCategoryRemove(category.Id)
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

func modelProjectCategoryUpdateProject(old *ProjectCategoryItem, prj *ProjectCategoryItem) error {
	//	o := orm.NewOrm()
	//	_, err := o.Update(prj, cols...)
	//	return err
	var err error
	db, err := getRawDB()
	if nil != err {
		return err
	}

	tx, err := db.Begin()
	if nil != err {
		return err
	}

	_, err = tx.Exec("UPDATE "+projectCategoryItemTableName+" SET project_name = ? , project_describe = ? WHERE id = ?",
		prj.ProjectName,
		prj.ProjectDescribe,
		prj.Id)
	if nil != err {
		tx.Rollback()
		return err
	}

	if old.ProjectName != prj.ProjectName {
		//	need update article
		_, err = tx.Exec("UPDATE "+projectArticleItemTableName+" SET project_name = ? WHERE project_id = ?",
			prj.ProjectName,
			prj.Id)
		if nil != err {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if nil != err {
		return err
	}
	return nil
}

/*
	Project article
*/
func modelProjectArticleNewArticle(article *ProjectArticleItem) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(article)
}

func modelProjectArticleEditArticle(article *ProjectArticleItem, cols []string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(article, cols...)
}

func modelProjectArticleIncClick(articleId int) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	if _, err := db.Exec("UPDATE "+projectArticleItemTableName+" SET click = click + 1 WHERE id = ?", articleId); nil != err {
		return err
	}
	return nil
}

func modelProjectArticleSetTop(articleId int, top bool) error {
	o := orm.NewOrm()
	var article ProjectArticleItem
	article.Id = articleId
	article.Top = 0
	if top {
		article.Top = 1
	}

	_, err := o.Update(&article, "top")
	if nil != err {
		return err
	}
	return err
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

func modelProjectArticleDelete(articleId int) error {
	o := orm.NewOrm()
	c, err := o.Delete(&ProjectArticleItem{Id: articleId})
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("delete article failed")
	}
	return nil
}

//	items, total page count
func modelProjectArticleGetTopArticles(project string, page int, limit int) ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query("SELECT id,article_title,article_author,post_time,top,reply_author,reply_time,active_time FROM "+projectArticleItemTableName+
		" WHERE project_name = ? AND top = 1 ORDER BY active_time DESC LIMIT ? OFFSET ?", project, limit, page*limit); nil != err {
		return nil, err
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
			&item.Top,
			&item.ReplyAuthor,
			&item.ReplyTime,
			&item.ActiveTime); nil != err {
			return nil, err
		}
		item.ProjectName = project
		resultSet = append(resultSet, item)
	}

	return resultSet, nil
}

func modelProjectArticleGetTopArticleCount(project string) (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	articleCount := 0
	err = db.QueryRow("SELECT COUNT(*) FROM "+projectArticleItemTableName+" WHERE project_name = ? AND top = 1", project).Scan(&articleCount)
	if nil != err {
		return 0, err
	}

	return articleCount, nil
}

//	first, get top articles
//	second, get the left articles if necessary
func modelProjectArticleGetArticles(project string, page int, limit int) ([]*ProjectArticleItem, int, error) {
	topArticles, err := modelProjectArticleGetTopArticles(project, page, limit)
	if nil != err {
		return nil, 0, err
	}

	leftCount := limit - len(topArticles)
	db, err := getRawDB()
	if nil != err {
		return nil, 0, err
	}
	if leftCount > 0 {
		offset := 0
		if leftCount < limit {
			offset = 0
		} else {
			topArticleCount, err := modelProjectArticleGetTopArticleCount(project)
			if err != nil {
				return nil, 0, err
			}

			if topArticleCount == 0 {
				//	no top article
				offset = page * limit
			} else {
				//	how many pages?
				topArticlePages := (topArticleCount + limit - 1) / limit
				topOffset := topArticlePages*limit - topArticleCount
				offset = topOffset
				topArticlePagesIndex := topArticlePages - 1

				if page > topArticlePagesIndex+1 {
					offset += (page - (topArticlePagesIndex + 1)) * limit
				}
			}
		}

		var rows *sql.Rows
		if rows, err = db.Query("SELECT id,article_title,article_author,post_time,top,reply_author,reply_time,active_time FROM "+projectArticleItemTableName+
			" WHERE project_name = ? AND top = 0 ORDER BY active_time DESC LIMIT ? OFFSET ?", project, leftCount, offset); nil != err {
			return nil, 0, err
		}

		//	free the conn
		defer rows.Close()

		//	get all item
		for rows.Next() {
			item := &ProjectArticleItem{}
			if err = rows.Scan(
				&item.Id,
				&item.ArticleTitle,
				&item.ArticleAuthor,
				&item.PostTime,
				&item.Top,
				&item.ReplyAuthor,
				&item.ReplyTime,
				&item.ActiveTime); nil != err {
				return nil, 0, err
			}
			item.ProjectName = project
			topArticles = append(topArticles, item)
		}
	}

	//	get page count
	articleCount := 0
	err = db.QueryRow("SELECT COUNT(*) FROM "+projectArticleItemTableName+" WHERE project_name = ?", project).Scan(&articleCount)
	if nil != err {
		return nil, 0, err
	}

	pages := (int(articleCount) + limit - 1) / limit

	return topArticles, pages, nil
}
