package voicebank

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// Source is the folder or the zip file a voicebank is read from. The files
// are never changed.
type Source struct {
	Name  string   // the folder's or zip's name, without .zip
	Files []string // slash paths of the files, sorted
	open  func(name string) (io.ReadCloser, error)
	close func() error
	// keys finds files whose names differ in case or Unicode
	// normalization, as macOS writes them, from the oto.ini's.
	keys map[string]string
}

// ErrNotVoicebank is returned for paths that are neither a folder nor a
// zip file.
var ErrNotVoicebank = errors.New("not a folder or zip file")

// OpenSource opens a voicebank's folder, read recursively, or zip file,
// whose names not marked UTF-8 are decoded as Shift-JIS.
func OpenSource(p string) (*Source, error) {
	info, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return openDir(p)
	}
	if strings.EqualFold(filepath.Ext(p), ".zip") {
		return openZip(p)
	}
	return nil, ErrNotVoicebank
}

func openDir(dir string) (*Source, error) {
	s := &Source{Name: filepath.Base(dir), keys: map[string]string{}, close: func() error { return nil }}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		s.add(filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.open = func(name string) (io.ReadCloser, error) { return os.Open(filepath.Join(dir, filepath.FromSlash(name))) }
	slices.Sort(s.Files)
	return s, nil
}

func openZip(p string) (*Source, error) {
	r, err := zip.OpenReader(p)
	if err != nil {
		return nil, fmt.Errorf("cannot open the zip: %w", err)
	}
	s := &Source{Name: strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)), keys: map[string]string{}, close: r.Close}
	files := map[string]*zip.File{}
	for _, f := range r.File {
		name := f.Name
		if f.NonUTF8 {
			name = Decode([]byte(name))
		}
		name = path.Clean(strings.TrimLeft(strings.ReplaceAll(name, `\`, "/"), "/"))
		if f.FileInfo().IsDir() || name == "." || strings.HasPrefix(name, "__MACOSX/") {
			continue
		}
		if _, dup := files[name]; !dup {
			files[name] = f
			s.add(name)
		}
	}
	s.open = func(name string) (io.ReadCloser, error) {
		f := files[name]
		if f == nil {
			return nil, fs.ErrNotExist
		}
		return f.Open()
	}
	slices.Sort(s.Files)
	return s, nil
}

func (s *Source) add(name string) {
	s.Files = append(s.Files, name)
	if k := key(name); s.keys[k] == "" {
		s.keys[k] = name
	}
}

func key(name string) string { return norm.NFC.String(strings.ToLower(name)) }

// Find returns the name of the file name refers to, in its case and
// normalization or not, and whether there is one.
func (s *Source) Find(name string) (string, bool) {
	f, ok := s.keys[key(name)]
	return f, ok
}

// ReadFile reads the file name refers to (see Find).
func (s *Source) ReadFile(name string) ([]byte, error) {
	f, ok := s.Find(name)
	if !ok {
		return nil, fmt.Errorf("%s: %w", name, fs.ErrNotExist)
	}
	r, err := s.open(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// Close closes a zip file.
func (s *Source) Close() error { return s.close() }
