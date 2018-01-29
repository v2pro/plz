package counselor

import (
	"bytes"
	"github.com/v2pro/plz/countlog"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type objectId struct {
	namespace  string
	objectName string
}

func (oid *objectId) String() string {
	return oid.namespace + "/" + oid.objectName
}

type itemName string
type itemDataFormat string
type items map[itemName]interface{}

var dataFormatPrefix = []byte("data_format:")
var store = map[objectId]items{}
var parsers = map[itemDataFormat]Parser{}
var parsersMutex = &sync.Mutex{}

type Parser interface {
	Parse(data []byte) (interface{}, error)
}

func RegisterParser(dataFormat string, parser Parser) {
	parsersMutex.Lock()
	defer parsersMutex.Unlock()
	parsers[itemDataFormat(dataFormat)] = parser
}

func RegisterParserByFunc(dataFormat string, f func(data []byte) (interface{}, error)) {
	parsersMutex.Lock()
	defer parsersMutex.Unlock()
	parsers[itemDataFormat(dataFormat)] = &funcParser{f}
}

type funcParser struct {
	f func(data []byte) (interface{}, error)
}

func (parser *funcParser) Parse(data []byte) (interface{}, error) {
	return parser.f(data)
}

func getParser(dataFormat string) Parser {
	parsersMutex.Lock()
	defer parsersMutex.Unlock()
	return parsers[itemDataFormat(dataFormat)]
}

var SourceDir = os.Getenv("HOME") + "/conf"

func getItem(objectId objectId, _itemName itemName) interface{} {
	items := store[objectId]
	if items == nil {
		items = map[itemName]interface{}{}
		store[objectId] = items
	}
	item := items[_itemName]
	if item != nil {
		return item
	}
	data := loadItem(objectId, _itemName)
	if data == nil {
		return nil
	}
	return parseItem(objectId, _itemName, data)
}

func loadItem(objectId objectId, itemName itemName) []byte {
	filename := path.Join(SourceDir, objectId.namespace, objectId.objectName, string(itemName))
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		countlog.Error("event!counselor.failed to load item",
			"err", err, "objectId", objectId, "itemName", itemName)
		return nil
	}
	return data
}

func parseItem(objectId objectId, itemName itemName, data []byte) interface{} {
	if bytes.Index(data, dataFormatPrefix) != 0 {
		countlog.Error("event!counselor.item data missing data format",
			"data", data, "objectId", objectId, "itemName", itemName)
		return nil
	}
	newlinePos := bytes.IndexByte(data, '\n')
	if newlinePos == -1 {
		countlog.Error("event!counselor.item data missing newline",
			"data", data, "objectId", objectId, "itemName", itemName)
		return nil
	}
	dataFormat := string(data[len(dataFormatPrefix):newlinePos])
	parser := getParser(dataFormat)
	if parser == nil {
		countlog.Error("event!counselor.no parser for the data format",
			"data", data, "objectId", objectId, "itemName", itemName, "dataFormat", dataFormat)
		return nil
	}
	parsed, err := parser.Parse(data[newlinePos+1:])
	if err != nil {
		countlog.Error("event!counselor.parse failed",
			"err", err, "data", data, "objectId", objectId, "itemName", itemName, "dataFormat", dataFormat)
		return nil
	}
	return parsed
}
