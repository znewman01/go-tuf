package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type HTTPRemoteOptions struct {
	MetadataPath string
	TargetsPath  string
	UserAgent    string
	Retries      *HTTPRemoteRetries
}

type HTTPRemoteRetries struct {
	Delay time.Duration
	Total time.Duration
}

var DefaultHTTPRetries = &HTTPRemoteRetries{
	Delay: time.Second,
	Total: 10 * time.Second,
}

func HTTPRemoteStore(baseURL string, opts *HTTPRemoteOptions, client *http.Client) (RemoteStore, error) {
	if !strings.HasPrefix(baseURL, "http") {
		return nil, ErrInvalidURL{baseURL}
	}
	if opts == nil {
		opts = &HTTPRemoteOptions{}
	}
	if opts.TargetsPath == "" {
		opts.TargetsPath = "targets"
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &httpRemoteStore{baseURL, opts, client}, nil
}

type httpRemoteStore struct {
	baseURL string
	opts    *HTTPRemoteOptions
	cli     *http.Client
}

func (h *httpRemoteStore) GetMeta(name string) (io.ReadCloser, int64, error) {
	return h.get(path.Join(h.opts.MetadataPath, name))
}

func (h *httpRemoteStore) GetTarget(name string) (io.ReadCloser, int64, error) {
	return h.get(path.Join(h.opts.TargetsPath, name))
}

func (h *httpRemoteStore) get(s string) (io.ReadCloser, int64, error) {
	u := h.url(s)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, 0, err
	}
	if h.opts.UserAgent != "" {
		req.Header.Set("User-Agent", h.opts.UserAgent)
	}
	var res *http.Response
	if r := h.opts.Retries; r != nil {
		for start := time.Now(); time.Since(start) < r.Total; time.Sleep(r.Delay) {
			res, err = h.cli.Do(req)
			if err == nil && (res.StatusCode < 500 || res.StatusCode > 599) {
				break
			}
		}
	} else {
		res, err = h.cli.Do(req)
	}
	if err != nil {
		return nil, 0, err
	}

	if res.StatusCode == http.StatusNotFound {
		res.Body.Close()
		return nil, 0, ErrNotFound{s}
	} else if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, 0, &url.Error{
			Op:  "GET",
			URL: u,
			Err: fmt.Errorf("unexpected HTTP status %d", res.StatusCode),
		}
	}

	size, err := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 0)
	if err != nil {
		return res.Body, -1, nil
	}
	return res.Body, size, nil
}

func (h *httpRemoteStore) url(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return h.baseURL + path
}

type fileRemoteStore struct {
	path string // TODO: use DirFS
}

func FileRemoteStore(url string) (RemoteStore, error) {
	parts := strings.SplitN(url, "://", 2)
	if parts[0] != "file" {
		return nil, ErrInvalidURL{url}
	}
	path := parts[1]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrNotFound{path}
	}
	return &fileRemoteStore{
		path: path,
	}, nil
}

func (f *fileRemoteStore) GetMeta(name string) (stream io.ReadCloser, size int64, err error) {
	file, err := os.Open(filepath.Join(f.path, name))
	if os.IsNotExist(err) {
		return nil, 0, ErrNotFound{name}
	} else if err != nil {
		return nil, 0, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}
	return file, info.Size(), nil
}

func (f *fileRemoteStore) GetTarget(path string) (stream io.ReadCloser, size int64, err error) {
	file, err := os.Open(filepath.Join(f.path, "targets", path))
	if os.IsNotExist(err) {
		return nil, 0, ErrNotFound{path}
	} else if err != nil {
		return nil, 0, err
	}
	if err != nil {
		return nil, 0, ErrNotFound{path}
	}
	info, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}
	return file, info.Size(), nil
}

func RemoteStoreFromURL(url string) (RemoteStore, error) {
	parts := strings.SplitN(url, "://", 2)
	switch parts[0] {
	case "http":
		return HTTPRemoteStore(url, nil, nil)
	case "file":
		return FileRemoteStore(url)
	default:
		return nil, ErrInvalidURL{url}
	}
}
