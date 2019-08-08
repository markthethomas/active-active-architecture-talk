package models

import (
	"floqars/shared"

	"gopkg.in/underarmour/dynago.v2"
)

// BaseGet gets a document from a table
func BaseGet(table string, doc dynago.Document) *dynago.GetItem {
	return shared.DAL.GetItem(table, doc)
}

// BasePut gets a document from a table
func BasePut(table string, doc dynago.Document) *dynago.PutItem {
	return shared.DAL.PutItem(table, doc)
}

// BaseDelete removes a doc from a table
func BaseDelete(table string, doc dynago.Document) *dynago.DeleteItem {
	return shared.DAL.DeleteItem(table, doc)
}

// BaseBatchWrite returns the base form of a batch write
func BaseBatchWrite() *dynago.BatchWrite {
	return shared.DAL.BatchWrite()
}

// BaseQuery constructs a query for a given table
func BaseQuery(table string) *dynago.Query {
	return shared.DAL.Query(table)
}

// ToDALError casts a DAL error to its dynago type for use
// (like checking to see if something wasn't found etc.)
func ToDALError(err error) *dynago.Error {
	return err.(*dynago.Error)
}

// NotFound returns whether a result is actually not found
func NotFound(res *dynago.GetItemResult) bool {
	return len(res.Item) == 0
}
