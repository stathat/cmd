package db

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/stathat/cmd/stathat/config"
	"github.com/stathat/cmd/stathat/intr"

	homedir "github.com/mitchellh/go-homedir"
)

type Store struct {
	IndexID   map[string]intr.Stat
	IndexName map[string]intr.Stat
	UpdatedAt time.Time
	filename  string
}

func New() (*Store, error) {
	hash := sha256.Sum256([]byte(config.AccessKey()))
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(home, ".stathat", fmt.Sprintf("%x.db", hash[0:4]))
	sf, err := load(filename)
	if err != nil {
		if config.Debug("db") {
			log.Printf("ignoring load %s error: %s", filename, err)
		}
	} else {
		return sf, nil
	}
	s := &Store{filename: filename}
	return s, nil
}

func (s *Store) update() error {
	if time.Since(s.UpdatedAt) < 5*time.Minute {
		return nil
	}
	stats, err := intr.StatList()
	if err != nil {
		return err
	}
	s.IndexID = make(map[string]intr.Stat)
	s.IndexName = make(map[string]intr.Stat)
	for _, stat := range stats {
		s.IndexID[stat.ID] = stat
		s.IndexName[stat.Name] = stat
	}
	s.UpdatedAt = time.Now()
	return s.save()
}

func (s *Store) Count() int {
	return len(s.IndexID)
}

func (s *Store) LookupID(id string) (intr.Stat, bool) {
	s.update()
	stat, ok := s.IndexID[id]
	return stat, ok
}

func (s *Store) LookupName(name string) (intr.Stat, bool) {
	s.update()
	stat, ok := s.IndexName[name]
	return stat, ok
}

func (s *Store) Lookup(query string) (intr.Stat, bool) {
	s.update()
	stat, ok := s.LookupID(query)
	if ok {
		if config.Debug("db") {
			log.Printf("id %s => %s/%s", query, stat.ID, stat.Name)
		}
		return stat, ok
	}
	stat, ok = s.LookupName(query)
	if ok && config.Debug("db") {
		log.Printf("name %s => %s/%s", query, stat.ID, stat.Name)
	}
	return stat, ok
}

func load(filename string) (*Store, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)

	var ts Store
	if err := dec.Decode(&ts); err != nil {
		return nil, err
	}
	ts.filename = filename
	if config.Debug("db") {
		log.Printf("Loaded data from %s\n", ts.filename)
	}

	return &ts, nil
}

func (s *Store) save() error {
	h, err := homedir.Dir()
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile(filepath.Join(h, ".stathat"), "swap")
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(f)
	if err := enc.Encode(s); err != nil {
		return err
	}
	f.Close()
	if config.Debug("db") {
		log.Printf("Saved data to %s", s.filename)
	}
	return os.Rename(f.Name(), s.filename)
}
