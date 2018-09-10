package gocodecc

import (
	"archive/zip"
	"database/sql"
	"errors"
	"os"
	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
	"github.com/satori/go.uuid"

	"github.com/axgle/mahonia"
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

func (m *ProjectCategoryItem) TableName() string {
	return projectCategoryItemTableName
}

var (
	projectArticleItemTableName = "project_article_item"
)

const (
	kArticleTitleLimit   = 128
	kArticleContentLimit = 0xffff
)

const (
	ArticleFlagNone = 0
)

const (
	ArticleFlagPrivate = 1 << iota
)

// ProjectArticleItem, ReplyTime is used as a article flag
type ProjectArticleItem struct {
	Id                     int    `orm:"pk;auto;index"`
	ProjectName            string `orm:"size(128)"`
	ArticleTitle           string `orm:"size(128)"`
	ArticleContentHtml     string `orm:"type(text)"`
	ArticleContentMarkdown string `orm:"type(text)"`
	ArticleAuthor          string `orm:"size(20)"`
	Like                   int
	PostTime               int64
	EditTime               int64
	Top                    int
	ReplyAuthor            string `orm:"size(50)"`
	ReplyTime              int64
	ActiveTime             int64
	Click                  int
	ProjectId              int    `orm:"index"`
	CoverImage             string `orm:"size(128)"`
	ReplyCount             int    `orm:"-"`
	Private                bool   `orm:"-"`
	PrivateInvisible       bool   `orm:"-"`
}

func (m *ProjectArticleItem) TableName() string {
	return projectArticleItemTableName
}

func (m *ProjectArticleItem) IsArticlePrivate() bool {
	return 0 != (m.ReplyTime & ArticleFlagPrivate)
}

func init() {
	orm.RegisterModel(new(ProjectCategoryItem))
	orm.RegisterModel(new(ProjectArticleItem))
}

func modelProjectCategoryGetAllSimple() ([]*ProjectCategoryItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query("SELECT id, project_name, item_count FROM " + projectCategoryItemTableName); nil != err {
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
			&item.ItemCount); nil != err {
			return nil, err
		}
		resultSet = append(resultSet, item)
	}

	return resultSet, nil
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

func modelProjectCategoryAddReturnId(item *ProjectCategoryItem) (int64, error) {
	o := orm.NewOrm()
	insertId, err := o.Insert(item)

	if nil == err {
		categoryImagePath := "./" + kPrefixImagePath + "/category-images"
		exists, _ := PathExists(categoryImagePath)
		if !exists {
			err = os.Mkdir(categoryImagePath, 0777)
			if nil != err {
				seelog.Error(err)
			}
		}
	}

	return insertId, err
}

func modelProjectCategoryAdd(item *ProjectCategoryItem) error {
	_, err := modelProjectCategoryAddReturnId(item)
	return err
}

func modelProjectCategoryRemove(id int) error {
	//	get base info
	var err error
	var category ProjectCategoryItem
	category.Id = id
	o := orm.NewOrm()
	err = o.Read(&category)
	if nil != err {
		return err
	}

	//	get category image
	categoryImagePath := ""
	if len(category.Image) != 0 {
		categoryImagePath = "./" + kPrefixImagePath + "/category-images/" + category.Image
	}

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

	err = tx.Commit()

	if nil == err {
		//	remove all articles image
		//	NOTE : remove if using cdn
		articleImagePath := "./" + kPrefixImagePath + "/article-images/" + strconv.Itoa(id)
		err = os.RemoveAll(articleImagePath)
		if nil != err {
			seelog.Error(err)
		}

		if len(categoryImagePath) != 0 {
			err = os.RemoveAll(categoryImagePath)
			if nil != err {
				seelog.Error(err)
			}
		}
	}

	return err
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

func modelProjectCategoryGetArticleCount(cateId int) (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	row := db.QueryRow("SELECT COUNT(*) FROM "+projectArticleItemTableName+" WHERE project_id = ?", cateId)
	var cnt int
	if err = row.Scan(&cnt); nil != err {
		return 0, err
	}
	return cnt, nil
}

func modelProjectCategoryGetByProjectId(projectId int, prj *ProjectCategoryItem) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	err = db.QueryRow("SELECT * FROM "+projectCategoryItemTableName+" WHERE id = ?", projectId).Scan(
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

	_, err = tx.Exec("UPDATE "+projectCategoryItemTableName+
		" SET project_name = ? , project_describe = ?, image = ?, post_priv = ? WHERE id = ?",
		prj.ProjectName,
		prj.ProjectDescribe,
		prj.Image,
		prj.PostPriv,
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
	var err error

	err = o.Begin()
	if nil != err {
		return 0, err
	}

	insertId, err := o.Insert(article)
	if nil != err {
		o.Rollback()
		return 0, err
	}

	//	inc article count
	rs := o.Raw("UPDATE "+projectCategoryItemTableName+" SET item_count = item_count + 1 WHERE id = ?", article.ProjectId)
	_, err = rs.Exec()

	if nil != err {
		o.Rollback()
		return 0, err
	}

	o.Commit()

	//	create image directory
	//	NOTE : If using cdn , remove it
	articleImagePath := "./" + kPrefixImagePath + "/article-images/" + strconv.Itoa(article.ProjectId) + "/" + strconv.FormatInt(insertId, 10)
	err = os.MkdirAll(articleImagePath, 0777)
	if nil != err {
		seelog.Error(err)
	}

	return insertId, nil
	/*db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	tx, err := db.Begin()
	if nil != err {
		return 0, err
	}

	insertRet, err := tx.Exec("INSERT INTO "+projectArticleItemTableName+` (
		project_name,
		article_title,
		article_content_html,
		article_content_markdown,
		article_author,
		like,
		post_time,
		edit_time,
		top,
		reply_author,
		reply_time,
		active_time,
		click,
		project_id
	)
	 VALUES (
		?,?,?,?,?,?,?,?,?,?,?,?,?,?
	)
	`,
		article.ProjectName,
		article.ArticleTitle,
		article.ArticleContentHtml,
		article.ArticleContentMarkdown,
		article.ArticleAuthor,
		article.Like,
		article.PostTime,
		article.EditTime,
		article.Top,
		article.ReplyAuthor,
		article.ReplyTime,
		article.ActiveTime,
		article.Click,
		article.ProjectId)

	if nil != err {
		tx.Rollback()
		return 0, err
	}

	insertId, err := insertRet.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//	inc article count
	_, err = tx.Exec("UPDATE "+projectCategoryItemTableName+" SET item_count = item_count + 1 WHERE id = ?", article.ProjectId)
	if nil != err {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return insertId, nil*/
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

func modelProjectArticleMarkPrivate(articleID int, private bool) error {
	o := orm.NewOrm()
	var article ProjectArticleItem
	article.Id = articleID
	err := o.Read(&article)
	if nil != err {
		return err
	}

	if private {
		article.ReplyTime |= ArticleFlagPrivate
	} else {
		article.ReplyTime &= (^ArticleFlagPrivate)
	}

	_, err = o.Update(&article, "reply_time")
	if nil != err {
		return err
	}
	return nil
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

func modelProjectArticleGetArticleCountByAuthor(author string) (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	counter := 0
	err = db.QueryRow("SELECT COUNT(*) FROM "+projectArticleItemTableName+" WHERE article_author = ?", author).Scan(&counter)

	return counter, err
}

func modelProjectArticleGetArticleCountAll() (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	counter := 0
	err = db.QueryRow("SELECT COUNT(*) FROM " + projectArticleItemTableName).Scan(&counter)

	return counter, err
}

func modelProjectArticleGetByAuthor(author string, page int, limit int) ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(`SELECT 
	id,
	project_name,
	article_title,
	post_time,
	reply_author,
	reply_time,
	project_id,
	active_time,
	click,
	cover_image
	 FROM `+projectArticleItemTableName+" WHERE article_author = ? ORDER BY post_time DESC LIMIT ? OFFSET ?",
		author,
		limit,
		page*limit); nil != err {
		return nil, err
	}

	//	free the conn
	defer rows.Close()

	//	get all item
	articles := make([]*ProjectArticleItem, 0, 10)
	for rows.Next() {
		item := &ProjectArticleItem{}
		if err = rows.Scan(
			&item.Id,
			&item.ProjectName,
			&item.ArticleTitle,
			&item.PostTime,
			&item.ReplyAuthor,
			&item.ReplyTime,
			&item.ProjectId,
			&item.ActiveTime,
			&item.Click,
			&item.CoverImage); nil != err {
			return nil, err
		}
		item.ArticleAuthor = author
		articles = append(articles, item)
	}

	return articles, nil
}

func modelProjectArticleDelete(articleId, projectId int) error {
	o := orm.NewOrm()
	/*c, err := o.Delete(&ProjectArticleItem{Id: articleId})
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("delete article failed")
	}
	return nil*/
	err := o.Begin()
	if nil != err {
		return err
	}

	_, err = o.Delete(&ProjectArticleItem{Id: articleId})
	if nil != err {
		o.Rollback()
		return err
	}

	//	decrease project item count
	rs := o.Raw("UPDATE "+projectCategoryItemTableName+" SET item_count = item_count - 1 WHERE id = ?", projectId)
	_, err = rs.Exec()
	if nil != err {
		o.Rollback()
		return err
	}

	o.Commit()

	//	delete project article image path
	//	NOTE : using cdn , remove it
	articleImagePath := "./" + kPrefixImagePath + "/article-images/" + strconv.Itoa(projectId) + "/" + strconv.Itoa(articleId)
	err = os.RemoveAll(articleImagePath)
	if nil != err {
		seelog.Error(err)
	}

	return nil
}

//	items, total page count
func modelProjectArticleGetTopArticlesByProjectName(project string, page int, limit int) ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(`
	SELECT id,
	article_title,
	article_author,
	post_time,
	top,
	reply_author,
	reply_time,
	active_time,
	click,
	cover_image FROM `+projectArticleItemTableName+
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
			&item.ActiveTime,
			&item.Click,
			&item.CoverImage); nil != err {
			return nil, err
		}
		item.ProjectName = project
		resultSet = append(resultSet, item)
	}

	return resultSet, nil
}

func modelProjectArticleGetTopArticleCountByProjectName(project string) (int, error) {
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
func modelProjectArticleGetArticlesByProjectName(project string, page int, limit int) ([]*ProjectArticleItem, int, error) {
	topArticles, err := modelProjectArticleGetTopArticlesByProjectName(project, page, limit)
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
			topArticleCount, err := modelProjectArticleGetTopArticleCountByProjectName(project)
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
		if rows, err = db.Query(`
		SELECT id,
		article_title,
		article_author,
		post_time,
		top,
		reply_author,
		reply_time,
		active_time,
		click,
		cover_iamge FROM `+projectArticleItemTableName+
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
				&item.ActiveTime,
				&item.Click,
				&item.CoverImage); nil != err {
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

func modelProjectArticleGetAllTopArticles(page int, limit int) ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if 0 != limit {
		if rows, err = db.Query(`
	SELECT id,
	project_id,
	project_name,
	article_title,
	article_author,
	article_content_html,
	post_time,
	top,
	reply_author,
	reply_time,
	active_time,
	click,
	cover_image FROM `+projectArticleItemTableName+
			" WHERE top = 1 ORDER BY active_time DESC LIMIT ? OFFSET ?", limit, page*limit); nil != err {
			return nil, err
		}
	} else {
		if rows, err = db.Query(`
	SELECT id,
	project_id,
	project_name,
	article_title,
	article_author,
	article_content_html,
	post_time,
	top,
	reply_author,
	reply_time,
	active_time,
	click,
	cover_image FROM ` + projectArticleItemTableName +
			" WHERE top = 1 ORDER BY active_time DESC"); nil != err {
			return nil, err
		}
	}

	//	free the conn
	defer rows.Close()

	//	get all item
	resultSet := make([]*ProjectArticleItem, 0, limit)
	for rows.Next() {
		item := &ProjectArticleItem{}
		if err = rows.Scan(
			&item.Id,
			&item.ProjectId,
			&item.ProjectName,
			&item.ArticleTitle,
			&item.ArticleAuthor,
			&item.ArticleContentHtml,
			&item.PostTime,
			&item.Top,
			&item.ReplyAuthor,
			&item.ReplyTime,
			&item.ActiveTime,
			&item.Click,
			&item.CoverImage); nil != err {
			return nil, err
		}

		resultSet = append(resultSet, item)
	}

	return resultSet, nil
}

//	using project id
//	items, total page count
func modelProjectArticleGetTopArticles(projectId int, page int, limit int) ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(`
	SELECT id,
	project_name,
	article_title,
	article_author,
	article_content_html,
	post_time,
	top,
	reply_author,
	reply_time,
	active_time,
	click,
	cover_image FROM `+projectArticleItemTableName+
		" WHERE project_id = ? AND top = 1 ORDER BY active_time DESC LIMIT ? OFFSET ?", projectId, limit, page*limit); nil != err {
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
			&item.ProjectName,
			&item.ArticleTitle,
			&item.ArticleAuthor,
			&item.ArticleContentHtml,
			&item.PostTime,
			&item.Top,
			&item.ReplyAuthor,
			&item.ReplyTime,
			&item.ActiveTime,
			&item.Click,
			&item.CoverImage); nil != err {
			return nil, err
		}
		item.ProjectId = projectId
		resultSet = append(resultSet, item)
	}

	return resultSet, nil
}

func modelProjectArticleGetTopArticleCount(projectId int) (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	articleCount := 0
	err = db.QueryRow("SELECT COUNT(*) FROM "+projectArticleItemTableName+" WHERE project_id = ? AND top = 1", projectId).Scan(&articleCount)
	if nil != err {
		return 0, err
	}

	return articleCount, nil
}

func modelProjectArticleGetTitleAndPostTime() ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}
	articles := make([]*ProjectArticleItem, 0, 32)
	rows, err := db.Query("SELECT id, article_title, post_time FROM " + projectArticleItemTableName + `
	ORDER BY post_time DESC`)
	if nil != err {
		return nil, err
	}
	for rows.Next() {
		var article ProjectArticleItem
		if err = rows.Scan(&article.Id, &article.ArticleTitle, &article.PostTime); nil != err {
			return nil, err
		}
		articles = append(articles, &article)
	}
	return articles, nil
}

//	first, get top articles
//	second, get the left articles if necessary
func modelProjectArticleGetArticles(projectId int, page int, limit int) ([]*ProjectArticleItem, int, error) {
	topArticles, err := modelProjectArticleGetTopArticles(projectId, page, limit)
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
			topArticleCount, err := modelProjectArticleGetTopArticleCount(projectId)
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
		if rows, err = db.Query(`
		SELECT id,
		project_name,
		article_title,
		article_author,
		article_content_html,
		post_time,
		top,
		reply_author,
		reply_time,
		active_time,
		click,
		cover_image FROM `+projectArticleItemTableName+
			" WHERE project_id = ? AND top = 0 ORDER BY post_time DESC LIMIT ? OFFSET ?", projectId, leftCount, offset); nil != err {
			return nil, 0, err
		}

		//	free the conn
		defer rows.Close()

		//	get all item
		for rows.Next() {
			item := &ProjectArticleItem{}
			if err = rows.Scan(
				&item.Id,
				&item.ProjectName,
				&item.ArticleTitle,
				&item.ArticleAuthor,
				&item.ArticleContentHtml,
				&item.PostTime,
				&item.Top,
				&item.ReplyAuthor,
				&item.ReplyTime,
				&item.ActiveTime,
				&item.Click,
				&item.CoverImage); nil != err {
				return nil, 0, err
			}
			item.ProjectId = projectId
			topArticles = append(topArticles, item)
		}
	}

	//	get page count
	articleCount := 0
	err = db.QueryRow("SELECT COUNT(*) FROM "+projectArticleItemTableName+" WHERE project_id = ?", projectId).Scan(&articleCount)
	if nil != err {
		return nil, 0, err
	}

	pages := (int(articleCount) + limit - 1) / limit

	return topArticles, pages, nil
}

func modelProjectArticleGetRecentNotTopArticles(page int, limit int) ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(`
		SELECT id,
		project_name,
		article_title,
		article_author,
		article_content_html,
		post_time,
		reply_author,
		reply_time,
		active_time,
		click,
		project_id,
		cover_image FROM `+projectArticleItemTableName+
		" WHERE top = 0 ORDER BY post_time DESC LIMIT ? OFFSET ?", limit, limit*page); nil != err {
		return nil, err
	}

	//	free the conn
	defer rows.Close()

	//	get all item
	articles := make([]*ProjectArticleItem, 0, 10)
	for rows.Next() {
		item := &ProjectArticleItem{}
		if err = rows.Scan(
			&item.Id,
			&item.ProjectName,
			&item.ArticleTitle,
			&item.ArticleAuthor,
			&item.ArticleContentHtml,
			&item.PostTime,
			&item.ReplyAuthor,
			&item.ReplyTime,
			&item.ActiveTime,
			&item.Click,
			&item.ProjectId,
			&item.CoverImage); nil != err {
			return nil, err
		}
		articles = append(articles, item)
	}

	return articles, nil
}

func modelProjectArticleGetRecentArticles(page int, limit int) ([]*ProjectArticleItem, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(`
		SELECT id,
		project_name,
		article_title,
		article_author,
		article_content_html,
		post_time,
		reply_author,
		reply_time,
		active_time,
		click,
		project_id,
		cover_image FROM `+projectArticleItemTableName+
		" ORDER BY post_time DESC LIMIT ? OFFSET ?", limit, limit*page); nil != err {
		return nil, err
	}

	//	free the conn
	defer rows.Close()

	//	get all item
	articles := make([]*ProjectArticleItem, 0, 10)
	for rows.Next() {
		item := &ProjectArticleItem{}
		if err = rows.Scan(
			&item.Id,
			&item.ProjectName,
			&item.ArticleTitle,
			&item.ArticleAuthor,
			&item.ArticleContentHtml,
			&item.PostTime,
			&item.ReplyAuthor,
			&item.ReplyTime,
			&item.ActiveTime,
			&item.Click,
			&item.ProjectId,
			&item.CoverImage); nil != err {
			return nil, err
		}
		articles = append(articles, item)
	}

	return articles, nil
}

func modelProjectArticleGetLastPostTime(author string) int64 {
	db, err := getRawDB()
	if nil != err {
		return 0xffffffff
	}

	var lastPostTime int64
	err = db.QueryRow("SELECT MAX(post_time) FROM "+projectArticleItemTableName+" WHERE article_author = ?", author).Scan(&lastPostTime)
	if nil != err {
		return 0
	}

	return lastPostTime
}

func modelProjectArticlesPack(dest string) (string, error) {
	u := uuid.NewV4()
	packPath := dest
	packFilename := "pack_" + u.String() + ".zip"
	packFullPath := packPath + packFilename
	err := os.MkdirAll(packPath, 0777)
	if nil != err {
		return "", err
	}

	//	create zip package file
	zipFile, err := os.Create(packFullPath)
	if nil != err {
		return "", err
	}
	defer zipFile.Close()
	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	db, err := getRawDB()
	if nil != err {
		return "", err
	}

	var rows *sql.Rows
	if rows, err = db.Query(`
		SELECT id,
		project_name,
		article_title,
		article_author,
		article_content_markdown FROM ` + projectArticleItemTableName); nil != err {
		return "", err
	}

	//	free the conn
	defer rows.Close()

	//	get all item
	for rows.Next() {
		item := &ProjectArticleItem{}
		if err = rows.Scan(
			&item.Id,
			&item.ProjectName,
			&item.ArticleTitle,
			&item.ArticleAuthor,
			&item.ArticleContentMarkdown); nil != err {
			return "", err
		}

		articlePath := item.ProjectName + "/" + item.ArticleAuthor
		markdownPath := articlePath + "/" + item.ArticleTitle + ".md"
		// Encoding to gbk
		encoder := mahonia.NewEncoder("gb18030")
		zipPath := encoder.ConvertString(markdownPath)
		writter, err := zw.Create(zipPath)
		if nil != err {
			return "", err
		}
		writter.Write([]byte(item.ArticleContentMarkdown))
	}

	return packFullPath, nil
}
