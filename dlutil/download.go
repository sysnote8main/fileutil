package dlutil

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)

func DownloadParallel(urls []string, concurrentSize int64, folderPath string) {
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(concurrentSize)
	for _, u := range urls {
		wg.Add(1)
		go downloadFromURL(u, &wg, s, folderPath)
	}
	// TODO returns wg and wait outside
	wg.Wait()
}

func downloadFromURL(_url string, wg *sync.WaitGroup, s *semaphore.Weighted, folderPath string) (*int64, error) {
	defer wg.Done()
	if err := s.Acquire(context.Background(), 1); err != nil {
		return nil, err
	}
	defer s.Release(1)

	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	path := u.Path
	fileName := extractFileNameFromURL(path)
	file, err := os.Create(folderPath + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resp, err := http.Get(_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}
	return &size, nil
}

func extractFileNameFromURL(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
