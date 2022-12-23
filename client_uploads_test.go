package nuxeogoclient

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestGetHandlers(t *testing.T) {
	nuxeo, err := getTestClient()
	if err != nil {
		t.Error(err)
	}

	uploads := nuxeo.Uploads()

	handlers, err := uploads.GetHandlers()
	if err != nil {
		t.Error(err)
	}

	handlerCount := len(handlers.Handlers)
	if handlerCount != 1 {
		t.Errorf("Unexpected size. Expected 1 but received %d", handlerCount)
	}

	handlerName := handlers.Handlers[0].Name
	if handlerName != "default" {
		t.Errorf("Unexpected handler name. Expected 'default' but received '%s'", handlerName)
	}

	batch, err := uploads.InitBatch(handlerName)
	if err != nil {
		t.Error(err)
	}
	if batch.Provider != handlerName {
		t.Errorf("Expected handler name to equal '%s' but received '%s'", handlerName, batch.Provider)
	}

	fmt.Println(batch.BatchId, batch.Provider)

	filePath := "TestFile.txt"
	fileUpl, err := uploads.UploadFile(batch.BatchId, 0, filePath, filePath, "text/plain")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fileUpl.BatchId)

	batchFile, err := uploads.GetUploadInfo(batch.BatchId, 0)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(batchFile.Name, batchFile.Size)

	filePath = "TestChunk.txt"
	fmt.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	fileSize := fileInfo.Size()

	chunk, err := ReadFileChunk(file, 5, fileSize, 0)
	if err != nil {
		t.Fatal(err)
	}
	chunkInfo, statusCode, err := uploads.UploadChunk(batch.BatchId, 1, 0, 4, filePath, "text/plain", 5, fileSize, bytes.NewReader(chunk))
	if err != nil {
		t.Fatal(err)
	}
	if statusCode != 202 {
		t.Fatalf("Expected status code to be 202 but it was %d", statusCode)
	}
	if chunkInfo.ChunkCount != 4 {
		t.Fatalf("Expected chunk count to be 4 but it was %d", chunkInfo.ChunkCount)
	}
	if len(chunkInfo.UploadedChunkIds) != 1 {
		t.Fatalf("Expected uploaded count to be 1 but it was %d", len(chunkInfo.UploadedChunkIds))
	}
	fmt.Println(chunkInfo)

	chunk, err = ReadFileChunk(file, 5, fileSize, 1)
	if err != nil {
		t.Fatal(err)
	}
	chunkInfo, statusCode, err = uploads.UploadChunk(batch.BatchId, 1, 1, 4, filePath, "text/plain", 5, fileSize, bytes.NewReader(chunk))
	if err != nil {
		t.Fatal(err)
	}
	if statusCode != 202 {
		t.Fatalf("Expected status code to be 202 but it was %d", statusCode)
	}
	if chunkInfo.ChunkCount != 4 {
		t.Fatalf("Expected chunk count to be 4 but it was %d", chunkInfo.ChunkCount)
	}
	if len(chunkInfo.UploadedChunkIds) != 2 {
		t.Fatalf("Expected uploaded count to be 2 but it was %d", len(chunkInfo.UploadedChunkIds))
	}
	fmt.Println(chunkInfo)

	chunk, err = ReadFileChunk(file, 5, fileSize, 2)
	if err != nil {
		t.Fatal(err)
	}
	chunkInfo, statusCode, err = uploads.UploadChunk(batch.BatchId, 1, 2, 4, filePath, "text/plain", 5, fileSize, bytes.NewReader(chunk))
	if err != nil {
		t.Fatal(err)
	}
	if statusCode != 202 {
		t.Fatalf("Expected status code to be 202 but it was %d", statusCode)
	}
	if chunkInfo.ChunkCount != 4 {
		t.Fatalf("Expected chunk count to be 4 but it was %d", chunkInfo.ChunkCount)
	}
	if len(chunkInfo.UploadedChunkIds) != 3 {
		t.Fatalf("Expected uploaded count to be 3 but it was %d", len(chunkInfo.UploadedChunkIds))
	}
	fmt.Println(chunkInfo)

	chunk, err = ReadFileChunk(file, 5, fileSize, 3)
	if err != nil {
		t.Fatal(err)
	}
	chunkInfo, statusCode, err = uploads.UploadChunk(batch.BatchId, 1, 3, 4, filePath, "text/plain", 5, fileSize, bytes.NewReader(chunk))
	if err != nil {
		t.Fatal(err)
	}
	if statusCode != 201 {
		t.Fatalf("Expected status code to be 201 but it was %d", statusCode)
	}
	if chunkInfo.ChunkCount != 4 {
		t.Fatalf("Expected chunk count to be 4 but it was %d", chunkInfo.ChunkCount)
	}
	if len(chunkInfo.UploadedChunkIds) != 4 {
		t.Fatalf("Expected uploaded count to be 4 but it was %d", len(chunkInfo.UploadedChunkIds))
	}
	fmt.Println(chunkInfo)

	batchInfo, err := uploads.GetBatchInfo(batch.BatchId)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(batchInfo)

	err = uploads.DeleteBatch(batch.BatchId)
	if err != nil {
		t.Error(err)
	}

	_, err = uploads.GetBatchInfo(batch.BatchId)
	if err == nil {
		t.Errorf("Expected an error. Batch should no longer exist exist")
	}

	err = uploads.DeleteBatch(batch.BatchId)
	if err == nil {
		t.Errorf("Expected an error. Batch should not exist")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Should had received a Request error")
	}

	if reqErr.StatusCode != 404 {
		t.Errorf("Expected error code to be '404' but it was '%d'", reqErr.StatusCode)
	}

}
