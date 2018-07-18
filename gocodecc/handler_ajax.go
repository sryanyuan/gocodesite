package gocodecc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"regexp"

	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

type AjaxResult struct {
	Result    int    `json:"Result"`
	Msg       string `json:"Msg"`
	CaptchaId string `json:"CaptchaId"`
}

type ArticleImageUploadResult struct {
	Success int    `json:"success"`
	Url     string `json:"url"`
	Message string `json:"message"`
}

var successData = []byte("success")
var projectArticleReg = regexp.MustCompile("^/project/\\d+/article/(\\d+)$")

func ajaxHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	action := vars["action"]
	var result AjaxResult
	result.Result = -1

	//	for article upload image
	var uploadResult ArticleImageUploadResult
	autoRender := true

	defer func() {
		if action == "upload" {
			//	need present result
			redirectUrl := ""
			if 0 == result.Result {
				redirectUrl = fmt.Sprintf("/common/message?text=&result=&title=上传成功")
			} else {
				redirectUrl = fmt.Sprintf("/common/message?text=%s&result=1&title=上传失败", result.Msg)
			}
			ctx.Redirect(redirectUrl, http.StatusFound)
		} else if action == "article_submit" ||
			action == "article_edit" {
			if 0 != result.Result {
				//	new captcha
				result.CaptchaId = captcha.NewLen(4)
			}
			renderJson(ctx, &result)
		} else if action == "article_image_upload" {
			ctx.RenderJson(&uploadResult)
		} else {
			if autoRender {
				renderJson(ctx, &result)
			}
		}
	}()

	switch action {
	case "project_create":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			//	check project name and project describe
			ctx.r.ParseForm()
			defer ctx.r.Body.Close()
			projectName := ctx.r.Form.Get("project[name]")
			projectDescribe := ctx.r.Form.Get("project[describe]")
			projectImage := ctx.r.Form.Get("project[image]")
			//	check with auth
			auth, err := strconv.Atoi(ctx.r.Form.Get("dst"))
			if nil != err {
				result.Msg = "Invalid auth select"
				return
			}
			if auth != kPermission_User &&
				auth != kPermission_SuperAdmin {
				result.Msg = "Invalid auth select"
				return
			}

			if len(projectName) == 0 ||
				len(projectDescribe) == 0 ||
				len(projectName) >= kCategoryNameLimit ||
				len(projectDescribe) >= kCategoryDescribeLimit {
				result.Msg = "invalid project name or project describe"
				return
			}

			var project ProjectCategoryItem
			project.Author = ctx.user.NickName
			project.Image = projectImage
			project.ProjectName = projectName
			project.ProjectDescribe = projectDescribe
			project.PostPriv = uint32(auth)
			err = modelProjectCategoryAdd(&project)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "project_edit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			//	check project name and project describe
			ctx.r.ParseForm()
			defer ctx.r.Body.Close()

			var err error
			projectName := ctx.r.Form.Get("project[name]")
			projectDescribe := ctx.r.Form.Get("project[describe]")
			projectImage := ctx.r.Form.Get("project[image]")
			projectId, err := strconv.Atoi(ctx.r.Form.Get("project[id]"))
			//	check with auth
			auth, err := strconv.Atoi(ctx.r.Form.Get("dst"))
			if nil != err {
				result.Msg = "Invalid auth select"
				return
			}
			if auth != kPermission_User &&
				auth != kPermission_SuperAdmin {
				result.Msg = "Invalid auth select"
				return
			}

			if len(projectName) == 0 ||
				len(projectDescribe) == 0 ||
				len(projectName) >= kCategoryNameLimit ||
				len(projectDescribe) >= kCategoryDescribeLimit ||
				nil != err ||
				0 == projectId {
				result.Msg = "invalid project name or project describe"
				return
			}

			//	get the original item
			var originPrj ProjectCategoryItem
			if err := modelProjectCategoryGetByProjectId(projectId, &originPrj); nil != err {
				result.Msg = err.Error()
				return
			}

			if originPrj.ProjectName == projectName &&
				originPrj.ProjectDescribe == projectDescribe &&
				originPrj.Image == projectImage &&
				originPrj.PostPriv == uint32(auth) {
				return
			}

			var newPrj ProjectCategoryItem
			newPrj = originPrj
			newPrj.ProjectName = projectName
			newPrj.ProjectDescribe = projectDescribe
			newPrj.Image = projectImage
			newPrj.PostPriv = uint32(auth)
			err = modelProjectCategoryUpdateProject(&originPrj, &newPrj)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "project_delete":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			ctx.r.ParseForm()
			projectId, err := strconv.Atoi(ctx.r.Form.Get("project[id]"))
			ctx.r.Body.Close()

			if projectId == 0 ||
				nil != err {
				result.Msg = "invalid project name"
				return
			}

			err = modelProjectCategoryRemove(projectId)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "article_submit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			defer ctx.r.Body.Close()
			projectId, err := strconv.Atoi(ctx.r.Form.Get("projectid"))
			if projectId == 0 ||
				nil != err {
				result.Msg = "invalid project"
				return
			}
			//	check captcha
			if !captcha.VerifyString(ctx.r.Form.Get("captchaid"), ctx.r.Form.Get("captchaSolution")) {
				result.Msg = "验证码错误"
				return
			}
			//	check auth
			var prj ProjectCategoryItem
			if err := modelProjectCategoryGetByProjectId(projectId, &prj); nil != err {
				result.Msg = err.Error()
				return
			}
			//	check auth
			if ctx.user.Permission < prj.PostPriv &&
				ctx.user.NickName != prj.Author {
				result.Msg = "permission denied"
				return
			}
			//	check post time
			if ctx.user.Permission < kPermission_Admin {
				lastPostTime := modelProjectArticleGetLastPostTime(ctx.user.UserName)
				tmNow := time.Now().Unix()
				if tmNow-lastPostTime < kMemberPostInterval {
					nextPostTime := lastPostTime + kMemberPostInterval - tmNow
					result.Msg = "离下一次发帖时间还有" + strconv.FormatInt(nextPostTime, 10) + "秒"
					return
				}
			}
			//	check valid
			title := ctx.r.Form.Get("title")
			if len(title) >= kArticleTitleLimit {
				result.Msg = "标题长度太长了"
				return
			}
			if len(title) == 0 {
				result.Msg = "请输入标题"
				return
			}
			contentHtml := ctx.r.Form.Get("editormd-html-code")
			//contentHtml = strings.Replace(contentHtml, "<pre>", `<pre class="prettyprint linenums">`, -1)
			if len(contentHtml) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentHtml) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
			if len(contentMarkdown) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentMarkdown) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			coverImage := ctx.r.Form.Get("coverImage")

			//	do post
			var postArticle ProjectArticleItem
			postArticle.ActiveTime = time.Now().Unix()
			postArticle.PostTime = time.Now().Unix()
			postArticle.ArticleTitle = title
			postArticle.ArticleAuthor = ctx.user.NickName
			postArticle.ArticleContentHtml = contentHtml
			postArticle.ArticleContentMarkdown = contentMarkdown
			postArticle.ProjectName = prj.ProjectName
			postArticle.ProjectId = prj.Id
			postArticle.CoverImage = coverImage
			articleId, err := modelProjectArticleNewArticle(&postArticle)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%d/article/%d", projectId, articleId)
		}
	case "article_edit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			defer ctx.r.Body.Close()

			//	check captcha
			if !captcha.VerifyString(ctx.r.Form.Get("captchaid"), ctx.r.Form.Get("captchaSolution")) {
				result.Msg = "验证码错误"
				return
			}
			projectId, _ := strconv.Atoi(ctx.r.Form.Get("projectId"))
			if projectId == 0 {
				result.Msg = "invalid project"
				return
			}
			articleId, _ := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if 0 == articleId {
				result.Msg = "invalid articleId"
				return
			}
			coverImage := ctx.r.Form.Get("coverImage")

			//	check auth
			article, err := modelProjectArticleGet(articleId)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			if article.ArticleAuthor != ctx.user.NickName {
				if ctx.user.Permission < kPermission_SuperAdmin {
					result.Msg = "access denied"
					return
				}
			}
			//	check valid
			title := ctx.r.Form.Get("title")
			if len(title) >= kArticleTitleLimit {
				result.Msg = "标题长度太长了"
				return
			}
			if len(title) == 0 {
				result.Msg = "请输入标题"
				return
			}
			contentHtml := ctx.r.Form.Get("editormd-html-code")
			//contentHtml = strings.Replace(contentHtml, "<pre>", `<pre class="prettyprint linenums">`, -1)
			if len(contentHtml) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentHtml) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
			if len(contentMarkdown) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentMarkdown) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}

			if len(coverImage) == 0 &&
				len(article.CoverImage) == 0 {
				//	find the first image label and use it
				coverImage = getOneImageFromHtml(contentHtml)
				if len(coverImage) != 0 {
					coverImage = filepath.Base(coverImage)
					extType := strings.ToLower(filepath.Ext(coverImage))
					switch extType {
					case ".jpg":
						fallthrough
					case ".jpeg":
						fallthrough
					case ".png":
						fallthrough
					case ".gif":
						fallthrough
					case ".webp":
						{
							//	nothing
						}
					default:
						{
							extType = ""
						}
					}
					if len(extType) == 0 {
						//	invalid image extension
						coverImage = ""
					}
				}
			}

			//	do post
			colsEdit := []string{"active_time", "edit_time"}
			article.ActiveTime = time.Now().Unix()
			article.EditTime = time.Now().Unix()
			if article.ArticleTitle != title {
				article.ArticleTitle = title
				colsEdit = append(colsEdit, "article_title")
			}
			if article.ArticleContentHtml != contentHtml {
				article.ArticleContentHtml = contentHtml
				article.ArticleContentMarkdown = contentMarkdown
				colsEdit = append(colsEdit, "article_content_html")
				colsEdit = append(colsEdit, "article_content_markdown")
			}
			if article.CoverImage != coverImage {
				article.CoverImage = coverImage
				colsEdit = append(colsEdit, "cover_image")
			}
			_, err = modelProjectArticleEditArticle(article, colsEdit)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%d/article/%d", projectId, articleId)
		}
	case "article_delete":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			ctx.r.Body.Close()
			if err != nil ||
				0 == articleId {
				result.Msg = "invalid articleId"
				return
			}

			//	get article
			article, err := modelProjectArticleGet(articleId)
			if nil != err {
				result.Msg = "invalid article"
				return
			}

			//	must be superadmin
			if ctx.user.Permission <= kPermission_Admin {
				result.Msg = "access denied"
				return
			}

			err = modelProjectArticleDelete(articleId, article.ProjectId)
			if nil != err {
				result.Msg = "delete article failed"
				return
			}

			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%d/page/1", article.ProjectId)
		}
	case "article_top":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "access denied"
				return
			}

			ctx.r.ParseForm()
			defer ctx.r.Body.Close()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if err != nil ||
				0 == articleId {
				result.Msg = "invalid articleId"
				return
			}
			top, err := strconv.Atoi(ctx.r.Form.Get("top"))
			if err != nil {
				result.Msg = "invalid top"
				return
			}

			doTop := true
			if 0 == top {
				doTop = false
			}

			err = modelProjectArticleSetTop(articleId, doTop)
			if nil != err {
				result.Msg = "set top failed"
				return
			}

			//	done
			result.Result = 0
		}
	case "article_mark_private":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if err != nil ||
				0 == articleId {
				result.Msg = "invalid articleId"
				return
			}
			privateStr := ctx.r.Form.Get("private")
			private := true
			if "" == privateStr {
				// Clear private flag
				private = false
			}
			// Check can modify private flag
			if ctx.user.Uid == 0 {
				result.Msg = "permission denied"
				return
			}
			article, err := modelProjectArticleGet(articleId)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			if ctx.user.Permission < kPermission_SuperAdmin &&
				ctx.user.UserName != article.ArticleAuthor {
				result.Msg = "permission denied"
				return
			}
			if err = modelProjectArticleMarkPrivate(articleId, private); nil != err {
				result.Msg = err.Error()
				return
			}

			result.Result = 0
		}
	case "article_image_upload":
		{
			ctx.r.ParseForm()
			projectId, err := strconv.Atoi(ctx.r.Form.Get("projectId"))
			if err != nil ||
				0 == projectId {
				uploadResult.Message = "非法的参数"
				return
			}
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if nil != err ||
				0 == articleId {
				uploadResult.Message = "非法的参数"
				return
			}

			//	create directory
			articleImagePath := "." + kPrefixImagePath + "/article-images/" + strconv.Itoa(projectId) + "/" + strconv.Itoa(articleId)
			err = os.MkdirAll(articleImagePath, 0777)
			if nil != err {
				uploadResult.Message = err.Error()
				return
			}

			file, header, err := ctx.r.FormFile("editormd-image-file")
			if nil != err {
				panic(err)
				return
			}
			defer file.Close()

			// 检查是否是jpg或png文件
			uploadFileType := header.Header["Content-Type"][0]

			filenameExtension := ""
			if uploadFileType == "image/jpeg" {
				filenameExtension = ".jpg"
			} else if uploadFileType == "image/png" {
				filenameExtension = ".png"
			} else if uploadFileType == "image/gif" {
				filenameExtension = ".gif"
			}

			if filenameExtension == "" {
				uploadResult.Message = "不支持的文件格式，请上传 jpg/png/gif 图片"
				return
			}

			//	copy to dest directory
			uploadImagePath := articleImagePath + "/" + header.Filename
			f, err := os.OpenFile(uploadImagePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				uploadResult.Message = err.Error()
				return
			}
			defer f.Close()
			io.Copy(f, file)
			uploadResult.Success = 1
			uploadResult.Url = strings.Trim(uploadImagePath, ".")
		}
	case "account_verify":
		{
			if ctx.r.Method != "GET" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			account := ctx.r.Form.Get("account")
			password := ctx.r.Form.Get("password")

			if len(account) == 0 ||
				len(password) == 0 ||
				len(account) > 20 ||
				len(password) > 100 {
				result.Msg = "invalid input"
				return
			}

			user := modelWebUserGetUserByUserName(account)
			if nil == user {
				result.Msg = "user not exists"
				result.Result = -2
				return
			}

			if password != user.PassToken {
				result.Msg = "invalid password"
				result.Result = -3
				return
			}

			//	done
			result.Result = 0
		}
	case "upload":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "Invalid method"
				return
			}
			if ctx.user.Permission < kPermission_SuperAdmin {
				//result.Msg = "access denied"
				result.Msg = kErrMsg_AccessDenied
				return
			}

			//	1 MB
			var fileSizeLimit int64 = 1 * 1024 * 1024
			ctx.r.ParseMultipartForm(fileSizeLimit)
			file, handler, err := ctx.r.FormFile("uploadfile")
			path := strings.Trim(ctx.r.Form.Get("path"), "/")
			path = strings.Trim(path, "\\")
			if len(path) != 0 {
				path += "/"
			}
			if nil != err {
				result.Msg = err.Error()
				return
			}
			defer file.Close()

			fileSize := int64(0)
			if statInterface, ok := file.(FileStat); ok {
				fileInfo, _ := statInterface.Stat()
				fileSize = fileInfo.Size()
			}
			if 0 == fileSize {
				if sizeInterface, ok := file.(FileSize); ok {
					fileSize = sizeInterface.Size()
				}
			}

			if fileSize > fileSizeLimit {
				result.Msg = "文件大小超过限制"
				return
			}

			//	check with path
			pathSel := ctx.r.Form.Get("dst")
			pathBase := ""
			if pathSel == "static" {
				pathBase = "./static/"
			} else if pathSel == "tpl" {
				pathBase = "./template/"
			} else if pathSel == "resume" {
				pathBase = ctx.config.ResumeFile
			} else {
				result.Msg = "Invalid file type"
				return
			}

			var f *os.File
			if pathSel == "resume" {
				f, err = os.OpenFile(pathBase, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			} else {
				f, err = os.OpenFile(pathBase+path+handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			}
			if err != nil {
				result.Msg = err.Error()
				return
			}
			defer f.Close()
			io.Copy(f, file)
			result.Result = 0

			// Clear file cache
			if pathSel == "resume" {
				readFileLock.Lock()
				delete(readFileCacheMap, pathSel)
				readFileLock.Unlock()
			}
		}
	case "reply_add":
		{
			var err error
			if ctx.r.Method != "POST" {
				result.Msg = "Invalid method"
				return
			}

			ctx.r.ParseForm()

			if !captcha.VerifyString(ctx.r.Form.Get("captchaid"), ctx.r.Form.Get("captchaSolution")) {
				result.Msg = "验证码错误"
				result.CaptchaId = captcha.NewLen(4)
				return
			}

			url := ctx.r.Form.Get("uri")
			seelog.Debugf("reply_add url %s", url)
			var article *ProjectArticleItem
			// Check if url is valid, maybe project article or guestbook
			selfComment := false
			if url == "/guestbook" {
				// ok
				if ctx.user.Uid == 1 {
					selfComment = true
				}
			} else {
				substrs := projectArticleReg.FindStringSubmatch(url)
				if nil == substrs ||
					len(substrs) != 2 {
					result.Msg = "非法的URL"
					result.CaptchaId = captcha.NewLen(4)
					return
				}
				articleID, err := strconv.Atoi(substrs[1])
				if nil != err {
					result.Msg = err.Error()
					return
				}
				// Just check if exists
				article, err = modelProjectArticleGet(articleID)
				if nil != err {
					result.Msg = err.Error()
					result.CaptchaId = captcha.NewLen(4)
					return
				}
				if nil == article {
					result.Msg = "非法的URL"
					result.CaptchaId = captcha.NewLen(4)
					return
				}
				// Get article author
				author := modelWebUserGetUserByUserName(article.ArticleAuthor)
				if nil == author {
					ctx.RenderMessagePage("错误", "无作者的文章", false)
					return
				}
				if author.Uid == ctx.user.Uid {
					selfComment = true
				}
			}

			comment := ctx.r.PostForm.Get("content")
			if len(comment) == 0 {
				result.Msg = "请输入留言内容"
				return
			}

			user := ctx.user
			if ctx.user.Uid == 0 {
				// If role is guest, need mail info
				mail := ctx.r.Form.Get("mail")
				if len(mail) == 0 {
					result.Msg = "请留下邮箱信息"
					return
				}
				if len(mail) > 21 {
					result.Msg = "游客信息过长，请限制在21个字符以内"
					return
				}
				user = &WebUser{}
				user.UserName = mail
			}

			replyID, err := modelReplyNew(url, user, comment)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			// Self to self comment will not send message tip
			if !selfComment {
				// Add message tip, if message is reply for guestbook, tip will send to admin
				receiverUid := uint32(1)
				if url != "/guestbook" {
					articleAuthor := modelWebUserGetUserByUserName(article.ArticleAuthor)
					if nil != articleAuthor {
						receiverUid = articleAuthor.Uid
					} else {
						receiverUid = 0
					}
				}
				if 0 != receiverUid {
					err = modelMessageNew(receiverUid, MessageTypeComment, comment, url, user, int(replyID))
					if nil != err {
						seelog.Errorf("Add comment message error:%s", err.Error())
					}
				} else {
					seelog.Errorf("Get receiver uid failed, url %s", url)
				}
			}
			// Find all people mentioned
			if user.Uid != 0 {
				substrs := mentionPeopleReg.FindAllString(comment, -1)
				if nil != substrs {
					for _, v := range substrs {
						username := strings.TrimSpace(v[1:])
						atUser := modelWebUserGetUserByUserName(username)
						if nil == atUser {
							continue
						}
						if atUser.Uid == user.Uid {
							continue
						}
						if err = modelMessageNew(atUser.Uid, MessageTypeReply, "", url, user, int(replyID)); nil != err {
							seelog.Errorf("Add reply message error:%s", err.Error())
						}
						seelog.Debugf("Add reply message for user %s success", username)
					}
				}
			}

			result.Result = 0
		}
	case "reply_del":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "Invalid method"
				return
			}
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = kErrMsg_AccessDenied
				return
			}

			ctx.r.ParseForm()

			replyIDStr := ctx.r.Form.Get("replyId")
			replyID, err := strconv.Atoi(replyIDStr)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			seelog.Debug(replyID)
			if err = modelReplyMarkDelete(replyID); nil != err {
				result.Msg = err.Error()
				return
			}
			result.Result = 0
		}
	case "message_get_count":
		{
			if ctx.r.Method != "GET" {
				result.Msg = "Invalid method"
				return
			}
			if ctx.user.Uid == 0 {
				result.Result = 0
				return
			}

			cnt, err := modelMessageGetCountByReceiver(ctx.user.Uid)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			result.Result = 0
			result.Msg = strconv.Itoa(cnt)
		}
	case "message_get":
		{
			if ctx.r.Method != "GET" {
				result.Msg = "Invalid method"
				return
			}
			if ctx.user.Uid == 0 {
				result.Result = 0
				return
			}

			messages, err := modelMessageGetByReceiver(ctx.user.Uid, 0, 8)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			for _, m := range messages {
				if m.Url == "/guestbook" {
					m.Title = "留言板"
				} else {
					substrs := projectArticleReg.FindStringSubmatch(m.Url)
					if nil == substrs ||
						2 != len(substrs) {
						m.Title = "N/A"
						continue
					}
					articleId, err := strconv.Atoi(substrs[1])
					if nil != err {
						seelog.Error("Fetch url articleid error:", err)
						m.Title = "N/A"
						continue
					}
					article, err := modelProjectArticleGet(articleId)
					if nil != err {
						seelog.Errorf("Get article %d error:%s", articleId, err.Error())
						m.Title = "N/A"
						continue
					}
					m.Title = article.ArticleTitle
				}
			}

			jsonBytes, err := json.Marshal(messages)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			result.Result = 0
			result.Msg = string(jsonBytes)
		}
	case "message_read":
		{
			if ctx.r.Method != "GET" {
				result.Msg = "Invalid method"
				return
			}
			if ctx.user.Uid == 0 {
				result.Msg = "Access denied"
				return
			}

			ctx.r.ParseForm()
			messageIDStr := ctx.r.Form.Get("message")
			messageID, err := strconv.Atoi(messageIDStr)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			if err = modelMessageMarkRead(ctx.user.Uid, messageID); nil != err {
				result.Msg = err.Error()
				return
			}

			result.Result = 0
		}
	case "zfbqr_pay":
		{
			ctx.r.ParseForm()
			donateAccount := ctx.r.Form.Get("user[account]")
			if donateAccount == "" {
				result.Msg = "请输入账户"
				return
			}
			donateValueStr := ctx.r.Form.Get("user[value]")
			if "" == donateValueStr {
				result.Msg = "请输入点数"
				return
			}
			donateValue, err := strconv.Atoi(donateValueStr)
			if nil != err {
				result.Msg = "错误的金额格式"
				return
			}
			if donateValue != 10 {
				//result.Msg = "金额必须为10"
				//return
			}
			if donateValue < 10 ||
				donateValue > 500 {
				result.Msg = "点数范围(10-500)"
				return
			}
			payMethodStr := ctx.r.Form.Get("paymethod")
			payMethod := payMethodAlipayQR
			if "" != payMethodStr {
				payMethod, err = strconv.Atoi(payMethodStr)
				if nil != err {
					result.Msg = "无效的支付方式"
					return
				}
			}
			if payMethod != payMethodWxQR &&
				payMethod != payMethodAlipayQR &&
				payMethod != payMethodUnion {
				result.Msg = "无效的支付方式"
				return
			}
			if payMethod == payMethodWxQR {
				if donateValue < 100 {
					result.Msg = "微信支付仅支持100元以上金额，小金额请使用支付宝"
					return
				}
			}

			seelog.Infof("Request to create order, account=%v, value=%v, pm=%v, debug=%v",
				donateAccount, donateValue, payMethod, ctx.config.Debug)
			orderInfo, err := createDonateOrder(donateAccount, donateValue, payMethod, ctx.config.Debug)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			// Request for payment url when paymethod is union pay
			if payMethodUnion == payMethod {
				qrURL, err := requestForPaymentURL(orderInfo, ctx.config)
				if nil != err {
					result.Msg = err.Error()
					return
				}
				orderInfo.QRUrl = qrURL
			}

			// Append a iframe into front
			seelog.Info("Order info ", orderInfo)
			jsonBytes, _ := json.Marshal(orderInfo)
			result.Result = 0
			result.Msg = string(jsonBytes)
		}
	case "zfbqr_pay_confirm":
		{
			ctx.r.ParseForm()
			ctx.r.ParseMultipartForm(10 * 1024)

			orderID := getFormValueAllMethod(ctx.r, "addnum")
			apikey := getFormValueAllMethod(ctx.r, "apikey")
			totalStr := getFormValueAllMethod(ctx.r, "total")
			uid := getFormValueAllMethod(ctx.r, "uid")
			seelog.Infof("Confirm order with orderID %s, apikey %s, total %s, uid %s", orderID, apikey, totalStr, uid)

			if "" == orderID {
				result.Msg = "Invalid order id"
				return
			}

			if "" == apikey {
				result.Msg = "Invalid apikey"
				return
			}

			if "" == totalStr {
				result.Msg = "Invalid total"
				return
			}

			if "" == uid {
				result.Msg = "Invalid uid"
				return
			}

			totalF, err := strconv.ParseFloat(totalStr, 32)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			err = confirmDonateOrder(uid, orderID, apikey, totalF)
			if nil != err {
				result.Msg = err.Error()
				seelog.Errorf("Confirm failed by remote server, error = %v", err)
				return
			}
			seelog.Infof("Confirm done by remote server, order id = %v", orderID)

			autoRender = false
			ctx.w.Write(successData)
		}
	case "pushmessage":
		{
			ctx.r.ParseForm()
			title := ctx.r.FormValue("title")
			body := ctx.r.FormValue("body")
			err := PushMessage(title, body)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			result.Result = 0
		}
	default:
		{
			result.Msg = "invalid ajax request"
		}
	}
}

func getFormValueAllMethod(r *http.Request, key string) string {
	value := r.FormValue(key)
	if "" != value {
		return value
	}

	value = r.PostFormValue(key)
	if "" != value {
		return value
	}

	return ""
}
