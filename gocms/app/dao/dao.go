package dao

import (
	"log"

	"github.com/davidddw/gopj/gocms/app/conf"
	"gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session

func init() {
	var err error
	mgoSession, err = mgo.Dial(conf.Conf.Db.Host) //connect database
	if err != nil {
		log.Fatal(err)
	}
}

func GetSession() *mgo.Session {
	return mgoSession.Copy()
}

func getCollect(session *mgo.Session, collectionName string) *mgo.Collection {
	return session.DB(conf.Conf.Db.DbName).C(collectionName)
}
