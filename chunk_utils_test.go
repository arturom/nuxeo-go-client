package nuxeogoclient

import (
	"os"
	"testing"
)

type chunkCountTest struct {
	maxChunkSize       int64
	fileSize           int64
	expectedChunkCount int
}

func TestChunkCount(t *testing.T) {
	tests := []chunkCountTest{
		{maxChunkSize: 10, fileSize: 100, expectedChunkCount: 10},
		{maxChunkSize: 10, fileSize: 99, expectedChunkCount: 10},
		{maxChunkSize: 10, fileSize: 101, expectedChunkCount: 11},
		{maxChunkSize: 10, fileSize: 10, expectedChunkCount: 1},
		{maxChunkSize: 100, fileSize: 100, expectedChunkCount: 1},
		{maxChunkSize: 10, fileSize: 110, expectedChunkCount: 11},
	}
	for _, test := range tests {
		actual := GetChunkCount(test.maxChunkSize, test.fileSize)
		if test.expectedChunkCount != actual {
			t.Errorf("Unexpected chunk count of %d for test %d", actual, test.expectedChunkCount)
		}
	}
}

type readFileChunkTest struct {
	chunkIndex    int
	maxChunkSize  int64
	expectedValue string
}

func TestReadFileChunk(t *testing.T) {
	file, err := os.Open("TestChunk.txt")
	if err != nil {
		t.Error(err)
	}
	tests := []readFileChunkTest{
		{chunkIndex: 0, maxChunkSize: 5, expectedValue: "A1234"},
		{chunkIndex: 1, maxChunkSize: 5, expectedValue: "B1234"},
		{chunkIndex: 2, maxChunkSize: 5, expectedValue: "C1234"},
		{chunkIndex: 3, maxChunkSize: 5, expectedValue: "D12"},
		{chunkIndex: 0, maxChunkSize: 10, expectedValue: "A1234B1234"},
		{chunkIndex: 0, maxChunkSize: 50, expectedValue: "A1234B1234C1234D12"},
	}
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	fileSize := info.Size()

	for _, test := range tests {
		data, err := ReadFileChunk(file, test.maxChunkSize, fileSize, test.chunkIndex)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != test.expectedValue {
			t.Fatalf("Expected '%s' but received '%s'", test.expectedValue, data)
		}

	}
}
