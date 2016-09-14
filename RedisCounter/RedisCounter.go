package cacheCounter

import (
	"fmt"
	"log"

	"github.com/dereking/grest/cache"
)

func Get(dbTableName string, dbField string, dbRecordId int64) (uint64, error) {
	var val uint64
	err := cache.GetCacheInstance().Get(fmt.Sprintf("%s_%s_%d", dbTableName, dbField, dbRecordId), &val)
	if err != nil {
		//log.Println("cacheCounter.get",dbTableName,dbField,dbRecordId ,"err:",err)
		return 0, err
	}
	return val, nil
}

func Increment(dbTableName string, dbField string, dbRecordId int64) (newVal uint64, err error) {
	//cache.Increment
	newVal, err = cache.GetCacheInstance().IncrementAuto(fmt.Sprintf("%s_%s_%d", dbTableName, dbField, dbRecordId), 1)
	if err != nil {
		log.Println("cacheCounter.Increment err:", err)
		return
	}
	return
}

func Decrement(dbTableName string, dbField string, dbRecordId int64) (newVal uint64, err error) {
	newVal, err = cache.GetCacheInstance().DecrementAuto(fmt.Sprintf("%s_%s_%d", dbTableName, dbField, dbRecordId), 1)
	if err != nil {
		log.Println("cacheCounter.Decrement err:", err)
		return
	}
	return
}
