package paper

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"io"
	"os"
	"sync"
)

// Page allow you to read and write values to a filer.
type Page struct {
	sync.RWMutex
	file filer
	size int64
}

// NewPage creates a new page from a file.
func NewPage(path string) (*Page, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &Page{
		file: f,
		size: info.Size(),
	}, nil
}

func InitPage(f filer, size int64) *Page {
	return &Page{
		file: f,
		size: size,
	}
}

func (p *Page) Put(v interface{}) (key int64, err error) {
	var buff bytes.Buffer

	err = gob.NewEncoder(&buff).Encode(v)
	if err != nil {
		return 0, err
	}

	value := []byte(base64.StdEncoding.EncodeToString(buff.Bytes()) + "\n")

	p.Lock()
	defer p.Unlock()

	key = p.size
	p.size = p.size + int64(len(value))

	_, err = p.file.WriteAt(value, key)
	if err != nil {
		return 0, err
	}
	return key, err
}

func (p *Page) Get(key int64, into interface{}) error {
	decoder, err := p.GetDecoder(key)
	if err != nil {
		return err
	}
	return decoder(into)
}

func (p *Page) getLine(key int64) (string, int64, error) {
	p.RLock()
	defer p.RUnlock()

	var line bytes.Buffer

	b := make([]byte, 1)
	for i := key; i < p.size; i++ {
		p.file.ReadAt(b, i)
		if string(b) == "\n" {
			break
		}
		_, err := line.Write(b)
		if err != nil {
			return "", 0, err
		}
	}
	size := int64(line.Len()) + 1
	str := line.String()
	return str, size, nil
}

func (p *Page) GetDecoder(key int64) (Decoder, error) {
	line, _, err := p.getLine(key)
	if err != nil {
		return nil, err
	}
	return NewlineDecoder(line)
}

func NewlineDecoder(line string) (Decoder, error) {
	b, err := base64.StdEncoding.DecodeString(line)
	if err != nil {
		return nil, err
	}
	return gob.NewDecoder(bytes.NewReader(b)).Decode, nil
}

func (p *Page) GetSize() int64 {
	p.RLock()
	defer p.RUnlock()
	return p.size
}

func (p *Page) Walk(fn func(key int64, decoder Decoder) error) error {
	size := p.GetSize()
	for key := int64(0); key < size; {
		line, length, err := p.getLine(key)
		if err != nil {
			return err
		}
		decoder, err := NewlineDecoder(line)
		if err != nil {
			return err
		}
		err = fn(key, decoder)
		if err != nil {
			return err
		}

		key = key + length
	}
	return nil
}

type filer interface {
	io.WriterAt
	io.ReaderAt
}

type Decoder func(into interface{}) error

type Line struct {
	Key int64
	Decoder
}

func (p *Page) Iterate(ctx context.Context) <-chan Line {
	lines := make(chan Line)
	go func() {
		defer close(lines)
		p.Walk(func(key int64, decoder Decoder) error {
			select {
			case <-ctx.Done():
				return nil
			default:
				lines <- Line{
					Key:     key,
					Decoder: decoder,
				}
			}
			return nil
		})
	}()
	return lines
}
