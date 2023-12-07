package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type part struct {
	index      int
	start      int64
	end        int64
	data       *bytes.Buffer
	lock       *sync.RWMutex
	current    int64
	total      int64
	retryCount int
	done       bool
}

func newPart(index int, start, end int64, retry int) *part {
	return &part{
		index:      index,
		start:      start,
		end:        end,
		data:       bytes.NewBuffer(nil),
		total:      end - start,
		retryCount: retry,
		lock:       &sync.RWMutex{},
	}
}

func (this *part) Write(body []byte) (int, error) {
	n := len(body)
	this.lock.Lock()
	defer this.lock.Unlock()
	this.current += int64(n)
	return n, nil
}

func (this *part) Download(url string) error {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	r.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", this.start, this.end))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	_, err = io.Copy(this.data, io.TeeReader(resp.Body, this))
	defer resp.Body.Close()
	return err
}

type Downloader struct {
	url       string
	directory string
	filename  string
	routine   int

	size     int64
	parts    []*part
	wg       *sync.WaitGroup
	retry    int
	progress bool
}

type downloaderOption func(*Downloader)

func DownloaderWithFilename(name string) downloaderOption {
	return func(d *Downloader) {
		d.filename = name
	}
}

func DownloaderWithDirectory(directory string) downloaderOption {
	return func(d *Downloader) {
		d.directory = directory
	}
}

func DownloaderWithRoutine(routine int) downloaderOption {
	return func(d *Downloader) {
		d.routine = routine
	}
}

func DownloaderWithProgress(p bool) downloaderOption {
	return func(d *Downloader) {
		d.progress = p
	}
}

func DownloaderWithRetry(n int) downloaderOption {
	return func(d *Downloader) {
		d.retry = n
	}
}

func NewDownloader(url string, options ...downloaderOption) (*Downloader, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	d := &Downloader{
		url:       url,
		directory: pwd,
		filename:  filepath.Base(url),
		routine:   runtime.NumCPU(),
		retry:     3,
		wg:        new(sync.WaitGroup),
	}

	for _, option := range options {
		option(d)
	}
	return d, nil
}

func (this *Downloader) Head() error {
	resp, err := http.Head(this.url)
	if err != nil {
		return err
	}

	if resp.Header.Get("Accept-Ranges") != "bytes" {
		this.routine = 1
		return errors.New("Unsupported")
	}

	this.size = resp.ContentLength

	return nil
}

func (this *Downloader) Download() *Downloader {
	chunkSize := int64(math.Ceil(float64(this.size) / float64(this.routine)))
	this.parts = make([]*part, this.routine)

	for i := 0; i < len(this.parts); i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if i == len(this.parts)-1 {
			end = this.size
		}

		part := newPart(i, start, end, this.retry)
		this.wg.Add(1)
		go func() {
			defer this.wg.Done()
			for {
				err := part.Download(this.url)
				if err == nil {
					part.done = true
					break
				}
				if this.retry < part.retryCount {
					part.retryCount += 1
					log.Printf("err: %s, retry: %d", err.Error(), part.retryCount+1)
					break
				}
			}
		}()

		this.parts[i] = part
	}

	if this.progress {
		this.wg.Add(1)
		go this.printProgress()
	}

	this.wg.Wait()
	return this
}

func (this *Downloader) Merge() error {
	file, err := os.OpenFile(filepath.Join(this.directory, this.filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	for i := range this.parts {
		if !this.parts[i].done {
			return fmt.Errorf("Partition %d download failed", i+1)
		}
		_, err := file.Write(this.parts[i].data.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Downloader) printProgress() {
	defer this.wg.Done()
	for {
		var current int64
		for i := range this.parts {
			this.parts[i].lock.RLock()
			current += this.parts[i].current
			this.parts[i].lock.RUnlock()
		}

		fmt.Printf("\r%s", strings.Repeat(" ", 34))
		fmt.Printf("\rcurrent: %d total: %d downloading... %.2f%%", current, this.size, float64(current*10000/int64(this.size))/100)
		if current >= this.size {
			fmt.Println()
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
