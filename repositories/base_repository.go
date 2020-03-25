package repositories

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
	"reflect"
)

type IRepository interface {
	Find(outEntities interface{}, query interface{}, args ...interface{}) error
	FindAll(outEntities interface{}) error
	FindPage(outEntities interface{}, pageIndex, pageSize uint64, orderBy []string, query interface{}, args ...interface{}) (uint64, error)
	First(outEntity interface{}) error
	FindOne(outEntity, query interface{}, args ...interface{}) error
	Count(entity interface{}) (uint64, error)
	Delete(entity interface{}) error
	Create(entity interface{}) error
	Update(entity interface{}) error
	Clear(entity interface{}) error
	Exec(sql string, args ...interface{})
	ClearAssociation(entity interface{}, related string)
	ReplaceAssociation(entity interface{}, related string, list interface{})
}

type TableNameAble interface {
	TableName() string
}

type baseRepository struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
}

func (b baseRepository) ClearAssociation(entity interface{}, related string) {
	b.db.Model(entity).Association(related).Clear()
}

func (b baseRepository) ReplaceAssociation(entity interface{}, related string, list interface{}) {
	b.db.Model(entity).Association(related).Replace(list)
}

func (b baseRepository) getTx(entity interface{}) (*gorm.DB, error) {
	var (
		tableNameAble TableNameAble
		ok            bool
	)

	// 通过接口断言判断是否实现了TableName接口
	if tableNameAble, ok = entity.(TableNameAble); !ok {
		// 没有实现TableName判断是否Slice，再通过反射获取TableName
		// type Result []*Item{}
		// result := &Result{}
		resultType := reflect.TypeOf(entity)
		if resultType.Kind() != reflect.Ptr {
			//fmt.Print("result is not a pointer")
			return nil, errors.New("result is not a pointer")
		}

		sliceType := resultType.Elem()
		if sliceType.Kind() != reflect.Slice {
			return nil, errors.New("result doesn't point to a slice")
		}
		// *Item
		itemPtrType := sliceType.Elem()
		// Item
		itemType := itemPtrType.Elem()

		elemValue := reflect.New(itemType)
		elemValueType := reflect.TypeOf(elemValue)
		tableNameAbleType := reflect.TypeOf((*TableNameAble)(nil)).Elem()

		if elemValueType.Implements(tableNameAbleType) {
			fmt.Print("neither the query nor result implement TableNameAble")
			return nil, errors.New("neither the query nor result implement TableNameAble")
		}

		tableNameAble = elemValue.Interface().(TableNameAble)
	}
	tableName := tableNameAble.TableName()
	tx := b.db.Table(tableName)
	tx = tx.Set("gorm:auto_preload", true) // 自动加载关联表
	tx = tx.Debug()
	err := tx.Error
	return tx, err
}

func (b baseRepository) Exec(sql string, args ...interface{}) {
	b.db.Exec(sql, args...)
}

func (b baseRepository) FindOne(outEntity, query interface{}, args ...interface{}) error {
	tx, err := b.getTx(outEntity)
	if err != nil {
		return err
	}
	err = tx.Where(query, args...).Take(outEntity).Error
	return err
}

func (b baseRepository) Find(outEntities interface{}, query interface{}, args ...interface{}) error {

	tx, err := b.getTx(outEntities)
	if err != nil {
		return err
	}
	return tx.Where(query, args...).Find(outEntities).Error

}

func (b baseRepository) FindAll(outEntities interface{}) error {

	tx, err := b.getTx(outEntities)
	if err != nil {
		return err
	}
	return tx.Find(outEntities).Error

}

func (b baseRepository) FindPage(outEntities interface{}, pageIndex, pageSize uint64, orderBy []string, query interface{}, args ...interface{}) (uint64, error) {
	var totalCount uint64
	tx, err := b.getTx(outEntities)
	if err != nil {
		return totalCount, err
	}
	if orderBy != nil {
		for _, orderCon := range orderBy {
			tx = tx.Order(orderCon)
		}
	}
	if query != nil {
		tx = tx.Where(query, args...)
	}
	if pageIndex > 0 {
		tx = tx.Offset((pageIndex - 1) * pageSize)
	}
	if pageSize > 0 {
		tx = tx.Limit(pageSize)
	}

	err = tx.Find(outEntities).Count(&totalCount).Error
	return totalCount, err
}

func (b baseRepository) First(outEntity interface{}) error {
	tx, err := b.getTx(outEntity)
	if err != nil {
		return err
	}
	return tx.First(outEntity).Error
}

func (b baseRepository) Count(entity interface{}) (uint64, error) {
	var totalCount uint64
	tx, err := b.getTx(entity)
	if err != nil {
		return totalCount, err
	}
	err = tx.Count(&totalCount).Error
	return totalCount, err
}

func (b baseRepository) Delete(entity interface{}) error {
	tx, err := b.getTx(entity)
	if err != nil {
		return err
	}
	return tx.Delete(entity).Error
}

func (b baseRepository) Create(entity interface{}) error {
	tx, err := b.getTx(entity)
	if err != nil {
		return err
	}
	if tx.NewRecord(entity) {
		return tx.Create(entity).Error
	}
	return errors.New("对象已存在")
}

func (b baseRepository) Update(entity interface{}) error {
	if !b.db.NewRecord(entity) {
		return b.db.Model(entity).Update(entity).Error
	}
	return errors.New("对象不存在")
}

func (b baseRepository) Clear(entity interface{}) error {
	b.db.Unscoped().Delete(entity)
	return nil
}
