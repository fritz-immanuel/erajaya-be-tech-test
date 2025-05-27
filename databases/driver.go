package databases

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/source"
)

// RiceBoxSource represents the golang-migrate data source from go.rice Box
type RiceBoxSource struct {
	lock       sync.Mutex
	box        *rice.Box
	migrations *source.Migrations
}

func (s *RiceBoxSource) loadMigrations() (*source.Migrations, error) {
	migrations := source.NewMigrations()
	err := s.box.Walk("", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("failed to access a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		migration, err := source.Parse(path)
		if err != nil {
			return err
		}
		migrations.Append(migration)
		return nil
	})
	return migrations, err
}

// PopulateMigrations populates all migration files from the go.rice box
func (s *RiceBoxSource) PopulateMigrations(box *rice.Box) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.box = box
	migrations, err := s.loadMigrations()
	if err != nil {
		return err
	}
	s.migrations = migrations
	return nil
}

// Open implements the golang-migrate source driver Open interface
func (s *RiceBoxSource) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// Close implements the golang-migrate source driver Close interface
func (s *RiceBoxSource) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.migrations = nil
	return nil
}

// First implements the golang-migrate source driver First interface
func (s *RiceBoxSource) First() (version uint, err error) {
	v, ok := s.migrations.First()
	if !ok {
		return 0, os.ErrNotExist
	}
	return v, nil
}

// Prev implements the golang-migrate source driver Prev interface
func (s *RiceBoxSource) Prev(version uint) (prevVersion uint, err error) {
	v, ok := s.migrations.Prev(version)
	if !ok {
		return 0, os.ErrNotExist
	}
	return v, nil
}

// Next implements the golang-migrate source driver Next interface
func (s *RiceBoxSource) Next(version uint) (nextVersion uint, err error) {
	v, ok := s.migrations.Next(version)
	if !ok {
		return 0, os.ErrNotExist
	}
	return v, nil
}

// ReadUp implements the golang-migrate source driver ReadUp interface
func (s *RiceBoxSource) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	migration, ok := s.migrations.Up(version)
	if !ok {
		return nil, "", os.ErrNotExist
	}
	b, err := s.box.Bytes(migration.Raw)
	if err != nil {
		return nil, "", err
	}
	return ioutil.NopCloser(bytes.NewBuffer(b)),
		migration.Identifier,
		nil
}

// ReadDown implements the golang-migrate source driver ReadDown interface
func (s *RiceBoxSource) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	migration, ok := s.migrations.Down(version)
	if !ok {
		return nil, "", migrate.ErrNilVersion
	}
	b, err := s.box.Bytes(migration.Raw)
	if err != nil {
		return nil, "", err
	}
	return ioutil.NopCloser(bytes.NewBuffer(b)),
		migration.Identifier,
		nil
}
