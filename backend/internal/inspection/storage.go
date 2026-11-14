package inspection

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MaxImageSize = 5 << 20

type Storage struct {
	dir       string
	container *url.URL
	client    *http.Client
}

func NewStorage(dir, containerURL string) (*Storage, error) {
	s := &Storage{dir: dir, client: &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	if containerURL == "" {
		if dir == "" {
			return nil, errors.New("inspection storage directory is required")
		}
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, err
		}
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() || info.Mode().Perm()&0077 != 0 {
			return nil, errors.New("inspection storage directory must be private")
		}
		return s, nil
	}
	u, err := url.Parse(containerURL)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.Fragment != "" || u.Query().Get("sig") == "" {
		return nil, errors.New("invalid Azure container SAS URL")
	}
	s.container = u
	return s, nil
}

func (s *Storage) blobURL(key string) string {
	u := *s.container
	u.Path = strings.TrimRight(u.Path, "/") + "/" + key
	return u.String()
}

func (s *Storage) Put(ctx context.Context, key, contentType string, data []byte) error {
	if s.container == nil {
		f, err := os.OpenFile(filepath.Join(s.dir, key), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		_, writeErr := f.Write(data)
		closeErr := f.Close()
		if writeErr != nil || closeErr != nil {
			_ = os.Remove(filepath.Join(s.dir, key))
			return errors.Join(writeErr, closeErr)
		}
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.blobURL(key), bytes.NewReader(data))
	if err != nil {
		return errors.New("cannot create storage request")
	}
	req.Header.Set("x-ms-blob-type", "BlockBlob")
	req.Header.Set("x-ms-version", "2023-11-03")
	req.Header.Set("If-None-Match", "*")
	req.Header.Set("Content-Type", contentType)
	return s.do(req)
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	if s.container == nil {
		return os.Remove(filepath.Join(s.dir, key))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.blobURL(key), nil)
	if err != nil {
		return errors.New("cannot create storage request")
	}
	req.Header.Set("x-ms-version", "2023-11-03")
	return s.do(req)
}

func (s *Storage) Get(ctx context.Context, key string) ([]byte, error) {
	if s.container == nil {
		file, err := os.Open(filepath.Join(s.dir, key))
		if err != nil {
			return nil, err
		}
		defer file.Close()
		data, err := io.ReadAll(io.LimitReader(file, MaxImageSize+1))
		if err != nil || len(data) > MaxImageSize {
			return nil, errors.New("storage response invalid")
		}
		return data, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.blobURL(key), nil)
	if err != nil {
		return nil, errors.New("cannot create storage request")
	}
	req.Header.Set("x-ms-version", "2023-11-03")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.New("storage request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("storage returned status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxImageSize+1))
	if err != nil || len(data) > MaxImageSize {
		return nil, errors.New("storage response invalid")
	}
	return data, nil
}

func (s *Storage) do(req *http.Request) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return errors.New("storage request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("storage returned status %d", resp.StatusCode)
	}
	return nil
}
