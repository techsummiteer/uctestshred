package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

// Task represents a shred task
type Task struct {
	Offset    int64
	ChunkSize int64
}

const (
	CHUNK_SIZE = 128 * 1024
	ITERATIONS = 3
)

func Shred(path string) error {
	//Open the file for write-only
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	//Get the size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	//We will use K as a chunk
	chunkSize := int64(CHUNK_SIZE)
	numChunks := (fileSize + chunkSize - 1) / chunkSize

	//Wait group for sync of goroutines
	var wg sync.WaitGroup
	var shredErr error

	//We get the cancellation context which we monitor with Done
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := make(chan error, 1)
	go func() {
		if err := <-errChan; err != nil {
			cancel()
		}
	}()
	//We get the worker threads to parallelize
	var concurrencyLimit int
	numCPU := runtime.NumCPU()
	if numChunks > int64(numCPU) {
		concurrencyLimit = numCPU
	} else {
		concurrencyLimit = int(numChunks)
	}

	tasks := make(chan Task, numChunks)
	workSize := fileSize
	for chunk := int64(0); chunk < numChunks; chunk++ {
		if workSize < chunkSize {
			tasks <- Task{Offset: chunk * chunkSize, ChunkSize: workSize}
		} else {
			workSize -= chunkSize
			tasks <- Task{Offset: chunk * chunkSize, ChunkSize: chunkSize}
		}
	}
	close(tasks)
	//Create the worker threads if chunks are greater than 1
	if numChunks > 1 {
		for i := 0; i < concurrencyLimit; i++ {
			wg.Add(1)
			go worker(i, ctx, file, tasks, &wg, &shredErr, errChan)
		}
	} else {
		wg.Add(1)
		worker(0, ctx, file, tasks, &wg, &shredErr, errChan)
	}
	wg.Wait()
	if shredErr != nil {
		return shredErr
	}
	file.Sync()
	file.Close()
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

func worker(id int, ctx context.Context, file *os.File, tasks <-chan Task, wg *sync.WaitGroup,
	shredErr *error, errChan chan error) {
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			break
		}
		for i := 0; i < ITERATIONS; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				if err := shredChunkBase(file, task.Offset, task.ChunkSize); err != nil && *shredErr == nil {
					*shredErr = err
					errChan <- fmt.Errorf(err.Error())
					return
				}
			}
		}
	}
	return
}

func shredChunkBase(file *os.File, offset, chunkSize int64) error {
	// Seek to the correct offset in the file
	_, err := file.Seek(offset, 0)
	if err != nil {
		return err
	}
	// Create a buffered writer
	writer := bufio.NewWriterSize(file, int(chunkSize))
	// Copy random data to the buffer
	if _, err := io.CopyN(writer, rand.Reader, chunkSize); err != nil {
		return err
	}
	// Flush the buffer to ensure all data is written
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}
