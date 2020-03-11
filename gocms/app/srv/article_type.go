package srv

import (
	"github.com/davidddw/gopj/gocms/app/dao"
	"github.com/davidddw/gopj/gocms/app/domain"
)

type ArticleTypeSrv struct {
}

func NewArticleTypeSrv() *ArticleTypeSrv {
	return &ArticleTypeSrv{}
}

func (this *ArticleTypeSrv) Edit(typeName string, articleType *domain.ArticleType) (err error) {
	err = dao.NewArticleTypeDao().Update(typeName, articleType)
	if err != nil {
		return
	}
	if typeName != articleType.Name {
		// update article
		err = dao.NewArticleDao().UpdateTypeName(typeName, articleType.Name)
	}
	return
}

func (this *ArticleTypeSrv) Delete(typeName string) (err error) {
	err = dao.NewArticleDao().DeleteByTypeName(typeName)
	if err != nil {
		return
	}
	err = dao.NewArticleTypeDao().Delete(typeName)
	return
}
