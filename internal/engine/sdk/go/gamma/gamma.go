/**
 * Copyright 2019 The Vearch Authors.
 *
 * This source code is licensed under the Apache License, Version 2.0 license
 * found in the LICENSE file in the root directory of this source tree.
 */

package gamma

/*
#cgo CFLAGS : -I../../../c_api
#cgo LDFLAGS: -L../../../engine/build -lgamma

#include "gamma_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

func Init(config *Config) unsafe.Pointer {
	var buffer []byte
	config.Serialize(&buffer)
	return C.Init((*C.char)(unsafe.Pointer(&buffer[0])), C.int(len(buffer)))
}

func Close(engine unsafe.Pointer) int {
	return int(C.Close(engine))
}

func CreateTable(engine unsafe.Pointer, table *Table) int {
	var buffer []byte
	table.Serialize(&buffer)
	return int(C.CreateTable(engine, (*C.char)(unsafe.Pointer(&buffer[0])), C.int(len(buffer))))
}

func AddOrUpdateDoc(engine unsafe.Pointer, buffer []byte) int {
	return int(C.AddOrUpdateDoc(engine, (*C.char)(unsafe.Pointer(&buffer[0])), C.int(len(buffer))))
}

func AddOrUpdateDocs(engine unsafe.Pointer, buffer [][]byte) BatchResult {
	num := len(buffer)
	C.AddOrUpdateDocsNum(engine, C.int(num))
	for i, b := range buffer {
		C.PrepareDocs(engine, (*C.char)(unsafe.Pointer(&(b[0]))), C.int(i))
	}
	var CBuffer *C.char
	zero := 0
	length := &zero
	C.AddOrUpdateDocsFinish(engine, C.int(num), (**C.char)(unsafe.Pointer(&CBuffer)), (*C.int)(unsafe.Pointer(length)))
	defer C.free(unsafe.Pointer(CBuffer))

	var result BatchResult
	buffer2 := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	result.DeSerialize(buffer2)
	return result
}

func DeleteDoc(engine unsafe.Pointer, docID []byte) int {
	return int(C.DeleteDoc(engine, (*C.char)(unsafe.Pointer(&docID[0])), C.int(len(docID))))
}

func GetEngineStatus(engine unsafe.Pointer, status *EngineStatus) {
	if engine == nil || status == nil {
		return
	}
	var CBuffer *C.char
	zero := 0
	length := &zero
	C.GetEngineStatus(engine, (**C.char)(unsafe.Pointer(&CBuffer)), (*C.int)(unsafe.Pointer(length)))
	defer C.free(unsafe.Pointer(CBuffer))
	buffer := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	status.DeSerialize(buffer)
}

func GetEngineMemoryInfo(engine unsafe.Pointer, status *MemoryInfo) {
	if engine == nil || status == nil {
		return
	}
	var CBuffer *C.char
	zero := 0
	length := &zero
	C.GetMemoryInfo(engine, (**C.char)(unsafe.Pointer(&CBuffer)), (*C.int)(unsafe.Pointer(length)))
	defer C.free(unsafe.Pointer(CBuffer))
	buffer := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	status.DeSerialize(buffer)
}

func GetDocByID(engine unsafe.Pointer, docID []byte, doc *Doc) int {
	var CBuffer *C.char
	zero := 0
	length := &zero
	ret := int(C.GetDocByID(engine,
		(*C.char)(unsafe.Pointer(&docID[0])),
		C.int(len(docID)),
		(**C.char)(unsafe.Pointer(&CBuffer)),
		(*C.int)(unsafe.Pointer(length))))
	defer C.free(unsafe.Pointer(CBuffer))
	buffer := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	doc.DeSerialize(buffer)
	return ret
}

func GetDocByDocID(engine unsafe.Pointer, docID int, next bool, doc *Doc) int {
	var CBuffer *C.char
	zero := 0
	length := &zero

	cNext := 0
	if next {
		cNext = 1
	}
	ret := int(C.GetDocByDocID(engine,
		C.int(docID),
		C.char(cNext),
		(**C.char)(unsafe.Pointer(&CBuffer)),
		(*C.int)(unsafe.Pointer(length))))
	defer C.free(unsafe.Pointer(CBuffer))
	buffer := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	doc.DeSerialize(buffer)
	return ret
}

func BuildIndex(engine unsafe.Pointer) int {
	return int(C.BuildIndex(engine))
}

func RebuildIndex(engine unsafe.Pointer, drop_before_rebuild int, limit_cpu int, describe int) int {
	return int(C.RebuildIndex(engine, C.int(drop_before_rebuild), C.int(limit_cpu), C.int(describe)))
}

func Dump(engine unsafe.Pointer) int {
	return int(C.Dump(engine))
}

func Load(engine unsafe.Pointer) int {
	return int(C.Load(engine))
}

/*func Search(engine unsafe.Pointer, request *Request, response *Response) int {
	var buffer []byte
	request.Serialize(&buffer)

	var CBuffer *C.char
	zero := 0
	length := &zero

	ret := int(C.Search(engine,
		(*C.char)(unsafe.Pointer(&buffer[0])), C.int(len(buffer)),
		(**C.char)(unsafe.Pointer(&CBuffer)),
		(*C.int)(unsafe.Pointer(length))))
	defer C.free(unsafe.Pointer(CBuffer))
	res := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	response.DeSerialize(res)
	return ret
}*/

func Search(engine unsafe.Pointer, reqByte []byte) (int, []byte) {
	var CBuffer *C.char
	zero := 0
	length := &zero

	ret := int(C.Search(engine,
		(*C.char)(unsafe.Pointer(&reqByte[0])), C.int(len(reqByte)),
		(**C.char)(unsafe.Pointer(&CBuffer)),
		(*C.int)(unsafe.Pointer(length))))
	defer C.free(unsafe.Pointer(CBuffer))
	respByte := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	return ret, respByte
}

func SetEngineCfg(engine unsafe.Pointer, config *Config) int {
	var buffer []byte
	config.Serialize(&buffer)
	ret := int(C.SetConfig(engine, (*C.char)(unsafe.Pointer(&buffer[0])), C.int(len(buffer))))
	return ret
}

func GetEngineCfg(engine unsafe.Pointer, config *Config) {
	var CBuffer *C.char
	zero := 0
	length := &zero
	C.GetConfig(engine, (**C.char)(unsafe.Pointer(&CBuffer)), (*C.int)(unsafe.Pointer(length)))
	defer C.free(unsafe.Pointer(CBuffer))
	buffer := C.GoBytes(unsafe.Pointer(CBuffer), C.int(*length))
	config.DeSerialize(buffer)
}
