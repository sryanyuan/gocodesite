package gocodecc

import (
	"math/rand"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// Get random charactors
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func articleApplyPrivate(user *WebUser, article *ProjectArticleItem) {
	if !article.IsArticlePrivate() {
		return
	}
	article.Private = true

	// Super admin can see all articles
	if user.Permission >= kPermission_SuperAdmin {
		return
	}
	// Self can see self article
	if user.Uid != 0 {
		if user.UserName == article.ArticleAuthor {
			return
		}
	}
	// Hide the article content
	article.PrivateInvisible = true
}

func articleAccessible(user *WebUser, article *ProjectArticleItem) bool {
	if user.Permission >= kPermission_SuperAdmin {
		return true
	}
	if user.UserName == article.ArticleAuthor {
		return true
	}
	return false
}
