package util

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type section map[string]string

type File struct {
	sections map[string]section
}

func NewFileReader(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewIniFileReader(f), nil
}

func NewIniFileReader(f io.Reader) *File {
	m := make(map[string]section)
	r := bufio.NewReader(f)
	sec := ""
	var line string
	var err error
	for err == nil {
		line, err = r.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" || line[0] == ';' {
			continue
		}
		if line[0] == '[' && line[len(line)-1] == ']' {
			sec = line[1 : len(line)-1]
			_, ok := m[sec]
			if !ok {
				m[sec] = make(section)
			}
			continue
		}
		if sec == "" {
			continue
		}
		pair := strings.SplitN(line, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key := strings.TrimSpace(pair[0])
		val := strings.TrimSpace(pair[1])
		if key == "" || val == "" {
			continue
		}
		m[sec][key] = val
	}
	return &File{m}
}

func (p *File) Sections() []string {
	s := make([]string, len(p.sections))
	i := 0
	for k, _ := range p.sections {
		s[i] = k
		i++
	}
	return s
}

func (p *File) HasSection(section string) bool {
	_, ok := p.sections[section]
	return ok
}

func (p *File) Keys(sec string) []string {
	m, ok := p.sections[sec]
	if !ok {
		return nil
	}
	keys := make([]string, len(m))
	i := 0
	for key, _ := range m {
		keys[i] = key
		i++
	}
	return keys
}

func (p *File) GetString(sec, key, def string) string {
	m, ok := p.sections[sec]
	if !ok {
		return def
	}
	v, ok := m[key]
	if !ok {
		return def
	}
	return v
}

func (p *File) GetInt(sec, key string, def int) int {
	m, ok := p.sections[sec]
	if !ok {
		return def
	}
	v, ok := m[key]
	if !ok {
		return def
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return def
	}
	return int(i)
}

func (p *File) GetUint(sec, key string, def uint) uint {
	m, ok := p.sections[sec]
	if !ok {
		return def
	}
	v, ok := m[key]
	if !ok {
		return def
	}
	i, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return def
	}
	return uint(i)
}

func (p *File) GetBool(sec, key string, def bool) bool {
	m, ok := p.sections[sec]
	if !ok {
		return def
	}
	v, ok := m[key]
	if !ok {
		return def
	}
	return v != "0"
}