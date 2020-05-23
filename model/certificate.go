package model

import (
	"github.com/jinzhu/gorm"
)

type Certificate struct {
	GModel   gorm.Model `gorm:"embedded"`
	UID      string     `gorm:"column:uid;type:char(32);unique;not null"`
	Space    string     `gorm:"column:space;type:varchar(32);not null"`
	Key      string     `gorm:"column:key;type:char(32);not null"`
	Consumer string     `gorm:"column:consumer;type:varchar(128);not null"`
	Content  string     `gorm:"column:content;type:TEXT;not null"`
}

func (Certificate) TableName() string {
	return "msa_license_certificate"
}

type CertificateQuery struct {
	Space    string
	Consumer string
	Key      string
}

type CertificateDAO struct {
}

func NewCertificateDAO() *CertificateDAO {
	return &CertificateDAO{}
}

func (CertificateDAO) Insert(_cer Certificate) error {
	db, err := openSqlDB()
	if nil != err {
		return err
	}
	defer closeSqlDB(db)

	return db.Create(&_cer).Error
}

func (CertificateDAO) Find(_uid string) (Certificate, error) {
	var cer Certificate
	db, err := openSqlDB()
	if nil != err {
		return cer, err
	}
	defer closeSqlDB(db)

	res := db.Where("uid = ?", _uid).First(&cer)
	if res.RecordNotFound() {
		return Certificate{}, nil
	}
	return cer, err
}

func (CertificateDAO) Query(_query CertificateQuery) ([]*Certificate, error) {
	var cers []*Certificate
	db, err := openSqlDB()
	if nil != err {
		return nil, err
	}
	defer closeSqlDB(db)

	db = db.Model(&Certificate{})
	blankQuery := true
	if "" != _query.Space {
		db = db.Where("space = ?", _query.Space)
		blankQuery = false
	}
	if "" != _query.Consumer {
		db = db.Where("consumer = ?", _query.Consumer)
		blankQuery = false
	}
	if "" != _query.Key {
		db = db.Where("key = ?", _query.Key)
		blankQuery = false
	}

	if blankQuery {
		return make([]*Certificate, 0), nil
	}

	res := db.Find(&cers)
	return cers, res.Error
}

func (CertificateDAO) Count(_query CertificateQuery) (int, error) {
	count := 0
	db, err := openSqlDB()
	if nil != err {
		return count, err
	}
	defer closeSqlDB(db)

	db = db.Model(&Certificate{})

	if "" != _query.Space {
		db = db.Where("space = ?", _query.Space)
	}
	if "" != _query.Consumer {
		db = db.Where("consumer = ?", _query.Consumer)
	}
	if "" != _query.Key {
		db = db.Where("key = ?", _query.Key)
	}

	res := db.Count(&count)
	return count, res.Error
}
