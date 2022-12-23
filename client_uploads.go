package nuxeogoclient

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type Handler struct {
	Name string `json:"name"`
}

type GetHandlersResponse struct {
	Handlers []Handler `json:"handlers"`
}

type Batch struct {
	Provider string `json:"provider"`
	BatchId  string `json:"batchId"`
}

type UploadedFile struct {
	// {"uploaded":"true","fileIdx":"0","uploadType":"normal","uploadedSize":"-1","batchId":"batchId-abf927ff-74dc-49cd-a935-2e6369554c66"}
	Uploaded     bool   `json:"uploaded"`
	BatchId      string `json:"batchId"`
	FileIdx      int    `json:"fileIdx"`
	UploadType   string `json:"uploadType"`
	UploadedSize int64  `json:"uploadedSize"`
}

type BatchInfo struct {
}

type BatchFileInfo struct {
	// {"size":11,"name":"TestFile.txt","uploadType":"normal"}
	Size       int64  `json:"size"`
	Name       string `json:"name"`
	UploadType string `json:"uploadType"`
}

type UploadedChunk struct {
	/*
		{
			"chunkCount":2,
			"uploaded":"true",
			"fileIdx":"1",
			"uploadType":"chunked",
			"uploadedSize":"5",
			"uploadedChunkIds":[0],
			"batchId":"batchId-cc4365be-531f-4fc9-b148-3fe4d9eb4bea"
		}
	*/
	ChunkCount       int    `json:"chunkCount"`
	Uploaded         bool   `json:"uploaded"`
	FileIdx          int    `json:"fileIdx"`
	UploadType       string `json:"uploadType"`
	UploadedSize     int64  `json:"uploadedSize"`
	UploadedChunkIds []int  `json:"uploadedChunkIds"`
	BatchId          string `json:"batchId"`
}

type UploadsClient interface {
	GetHandlers() (*GetHandlersResponse, error)
	InitBatch(handlerName string) (*Batch, error)
	UploadChunk(batchId string, fileIdx, chunkIdx, chunkCount int, filename, fileType string, chunkSize, fileSize int64, body io.Reader) (*UploadedChunk, int, error)
	UploadFile(batchId string, fileIdx int, filePath, fileName, fileType string) (*UploadedFile, error)
	GetUploadInfo(batchId string, fileIdx int) (*BatchFileInfo, error)
	GetBatchInfo(batchId string) ([]BatchFileInfo, error)
	DeleteBatch(batchId string) error
}

type uploadsClient struct {
	nuxeo client
}

func newUploadsClient(nuxeo client) UploadsClient {
	c := new(uploadsClient)
	c.nuxeo = nuxeo
	return c
}

func (c uploadsClient) GetHandlers() (*GetHandlersResponse, error) {
	handlers := new(GetHandlersResponse)
	_, err := c.nuxeo.GetJson("/api/v1/upload/handlers", handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func (c uploadsClient) InitBatch(handlerName string) (*Batch, error) {
	batch := new(Batch)
	path := fmt.Sprintf("/api/v1/upload/new/%s", handlerName)
	_, err := c.nuxeo.PostJson(path, nil, batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (c uploadsClient) UploadReader(batchId string, fileIdx int, filename, fileType string, contentLength int64, body io.Reader) (*UploadedFile, error) {
	path := fmt.Sprintf("/api/v1/upload/%s/%d", batchId, fileIdx)
	req, err := c.nuxeo.initRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-File-Name", filename)
	req.Header.Set("X-File-Type", fileType)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", contentLength))

	resp, err := c.nuxeo.sendRequest(req)
	if err != nil {
		return nil, err
	}

	f := new(UploadedFile)
	unmarshallJSONResponse(resp, f)
	return f, nil
}

func (c uploadsClient) UploadFile(batchId string, fileIdx int, filePath, fileName, fileType string) (*UploadedFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileName == "" {
		fileName = path.Base(filePath)
	}
	return c.UploadReader(batchId, fileIdx, fileName, fileType, info.Size(), bufio.NewReader(file))
}

func (c uploadsClient) UploadChunk(batchId string, fileIdx, chunkIdx, chunkCount int, filename, fileType string, chunkSize, fileSize int64, body io.Reader) (*UploadedChunk, int, error) {
	path := fmt.Sprintf("/api/v1/upload/%s/%d", batchId, fileIdx)
	req, err := c.nuxeo.initRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("X-Upload-Type", "chunked")
	req.Header.Set("X-Upload-Chunk-Index", fmt.Sprintf("%d", chunkIdx))
	req.Header.Set("X-Upload-Chunk-Count", fmt.Sprintf("%d", chunkCount))

	req.Header.Set("X-File-Name", filename)
	req.Header.Set("X-File-Type", fileType)
	req.Header.Set("X-File-Size", fmt.Sprintf("%d", fileSize))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", chunkSize))

	resp, err := c.nuxeo.sendRequest(req)
	if err != nil {
		if resp != nil {
			fmt.Println(readStringResponse(resp))
		}
		return nil, 0, err
	}

	f := new(UploadedChunk)
	unmarshallJSONResponse(resp, f)
	return f, resp.StatusCode, nil
}

func (c uploadsClient) GetUploadInfo(batchId string, fileIdx int) (*BatchFileInfo, error) {
	file := new(BatchFileInfo)
	path := fmt.Sprintf("/api/v1/upload/%s/%d", batchId, fileIdx)
	_, err := c.nuxeo.GetJson(path, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c uploadsClient) GetBatchInfo(batchId string) ([]BatchFileInfo, error) {
	batch := make([]BatchFileInfo, 0)
	path := fmt.Sprintf("/api/v1/upload/%s", batchId)
	_, err := c.nuxeo.GetJson(path, &batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (c uploadsClient) DeleteBatch(batchId string) error {
	path := fmt.Sprintf("/api/v1/upload/%s", batchId)
	resp, err := c.nuxeo.Delete(path)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status, resp.StatusCode)
	return nil
}
