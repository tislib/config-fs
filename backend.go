package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Backend struct {
}

func (b *Backend) init() {

}

func (b *Backend) runReadOperation(database string, collection string, connection string, location string) {
	location, err := filepath.Abs(location)

	if err != nil {
		log.Print(err)
		os.Exit(10)
		return
	}

	mongoInstance := new(MongoConn)
	mongoInstance.Init(connection)

	db := mongoInstance.client.Database(database)

	list, err := db.ListCollectionNames(context.TODO(), bson.M{})

	dieIfError(err)

	connectionPattern, err := regexp.Compile(collection)

	var list2 []string

	for _, item := range list {
		if connectionPattern.Match([]byte(item)) {
			list2 = append(list2, item)
		}
	}

	for _, item := range list2 {
		b.runReadOperationItem(db, item, location+"/"+item)
	}
}

func (b *Backend) runReadOperationItem(db *mongo.Database, collection string, location string) {
	err := os.Mkdir(location, 0777)

	if err != nil && !strings.HasSuffix(err.Error(), "file exists") {
		dieIfError(err)
	}

	col := db.Collection(collection)

	cursor, err := col.Find(context.TODO(), bson.M{})

	dieIfError(err)

	var recordIds []string

	for cursor.Next(context.TODO()) {
		record := bson.M{}

		err = cursor.Decode(record)

		dieIfError(err)

		recordId := record["_id"].(string)

		delete(record, "_id")
		delete(record, "_class")

		recordIds = append(recordIds, recordId)

		data, err := yaml.Marshal(record)

		dieIfError(err)

		file, err := os.OpenFile(location+"/"+recordId+".yml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)

		dieIfError(err)

		_, err = file.Write(data)

		dieIfError(err)

		err = file.Close()

		dieIfError(err)

		log.Printf("Updating file: %s", recordId+".yml")
	}

	dieIfError(cursor.Close(context.TODO()))

	dieIfError(filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		item := path[len(location):]
		if item == "" || !strings.HasSuffix(item, ".yml") {
			return nil
		}

		item = item[1 : len(item)-4]

		if !contains(recordIds, item) {
			err = os.RemoveAll(location + "/" + item)
			dieIfError(err)
			log.Printf("Removing file: %s", item)
		}

		return nil
	}))
}

func (b *Backend) runWriteOperation(database string, collection string, connection string, location string) {
	location, err := filepath.Abs(location)

	if err != nil {
		log.Print(err)
		os.Exit(10)
		return
	}

	mongoInstance := new(MongoConn)
	mongoInstance.Init(connection)

	db := mongoInstance.client.Database(database)

	list, err := db.ListCollectionNames(context.TODO(), bson.M{})

	dieIfError(err)

	connectionPattern, err := regexp.Compile(collection)

	var list2 []string

	for _, item := range list {
		if connectionPattern.Match([]byte(item)) {
			list2 = append(list2, item)
		}
	}

	for _, item := range list2 {
		b.runWriteOperationItem(db, item, location+"/"+item)
	}
}

func (b *Backend) runWriteOperationItem(db *mongo.Database, collection string, location string) {
	col := db.Collection(collection)

	var recordIds []string

	dieIfError(filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		item := path[len(location):]
		if item == "" || !strings.HasSuffix(item, ".yml") {
			return nil
		}

		item = item[1 : len(item)-4]

		recordIds = append(recordIds, item)

		content, err := ioutil.ReadFile(path)

		var record = new(interface{})
		dieIfError(err)

		err = yaml.Unmarshal(content, record)

		dieIfError(err)

		filter := bson.M{"_id": item}
		update := bson.M{"$set": convert(record)}
		opts := new(options.UpdateOptions)
		opts.Upsert = new(bool)
		*opts.Upsert = true

		result, err := col.UpdateOne(context.TODO(), filter, update, opts)

		dieIfError(err)

		if result.UpsertedCount > 0 {
			log.Printf("Record updated: %s", item)
		}

		return nil
	}))

	cursor, err := col.Find(context.TODO(), bson.M{})

	dieIfError(err)

	for cursor.Next(context.TODO()) {
		record := bson.M{}

		err = cursor.Decode(record)

		dieIfError(err)

		item := record["_id"].(string)

		if !contains(recordIds, item) {
			log.Printf("Deleting record : %s", item)
			_, err = col.DeleteOne(context.TODO(), bson.M{"_id": item})
			dieIfError(err)
		}
	}

	cursor.Close(context.TODO())
}

func dieIfError(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(10)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case *interface{}:
		return convert(*i.(*interface{}))
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
