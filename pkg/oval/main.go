package oval

import (
	"alt-oval-scanner/pkg/repository"
	"alt-oval-scanner/pkg/utils"
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go.etcd.io/bbolt"
	"golang.org/x/xerrors"
	"io"
	"log"
	"time"
)

func NewOvalManager(branchesUrl, baseUrl, pathToDB string) (*Manager, error) {
	r := repository.NewRepository(pathToDB)
	err := r.Open()
	if err != nil {
		return nil, err
	}
	err = r.CreateBucket("meta")
	if err != nil {
		return nil, err
	}
	err = r.CreateBucket("oval")
	if err != nil {
		return nil, err
	}
	return &Manager{baseUrl: baseUrl, branchesUrl: branchesUrl, repository: r}, nil
}

func (o *Manager) Download() error {
	bytes, err := o.repository.Get("meta", []byte("fetch-time"))
	if err != nil {
		return err
	}
	if bytes != nil {
		var t time.Time
		err = json.Unmarshal(bytes, &t)
		if err != nil {
			return err
		}
		if err == nil {
			if time.Now().Sub(t) < (4 * time.Hour) {
				log.Printf("fresh data")
				log.Printf("last download: %s", t)
				return nil
			}
		}
	}
	ovalPaths, err := o.fetchBranches()
	if err != nil {
		return xerrors.Errorf("failed to get oval file paths: %w", err)
	}
	for _, ovalPath := range ovalPaths {
		log.Printf("fetch: %s....", ovalPath)
		if err = o.updateOVAL(ovalPath); err != nil {
			return xerrors.Errorf("failed to update ALT OVAL: %w", err)
		}
	}
	err = o.updateMeta()
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}

func (o *Manager) fetchBranches() ([]string, error) {
	res, err := utils.FetchURL(o.branchesUrl, "", 5)
	if err != nil {
		return nil, xerrors.Errorf("failed to fetch branches: %w", err)
	}
	var branches Branches
	err = json.Unmarshal(res, &branches)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal branches.json: %w", err)
	}
	var paths []string
	for _, b := range branches.Branches {
		paths = append(paths, fmt.Sprintf(o.baseUrl, b))
	}
	return paths, nil
}

func (o *Manager) updateOVAL(ovalPath string) error {
	res, err := utils.FetchURL(ovalPath, "", 5)
	if err != nil {
		return err
	}
	r, err := zip.NewReader(bytes.NewReader(res), int64(len(res)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		var oval OVAL
		rc, err := f.Open()
		if err != nil {
			return err
		}
		content, err := io.ReadAll(rc)
		if err != nil {
			rc.Close()
			return err
		}
		err = xml.Unmarshal(content, &oval)
		if err != nil {
			rc.Close()
			return err
		}
		err = o.saveBucket(f.Name, oval)
		if err != nil {
			return err
		}
	}
	log.Printf("save %d ovals", len(r.File))
	return nil
}

func (o *Manager) saveBucket(name string, oval OVAL) error {
	marshal, err := json.Marshal(oval)
	if err != nil {
		return err
	}
	return o.repository.Save("oval", []byte(name), marshal)
}

func (o *Manager) updateMeta() error {
	var currentTime = time.Now()
	marshal, err := json.Marshal(currentTime)
	if err != nil {
		return err
	}
	err = o.repository.Save("meta", []byte("fetch-time"), marshal)
	if err != nil {
		return err
	}
	return err
}

func (o *Manager) OVALs() ([]OVAL, error) {
	var ovals []OVAL
	err := o.repository.Pointer.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("oval"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var oval OVAL
			err := json.Unmarshal(v, &oval)
			if err != nil {
				return err
			}
			ovals = append(ovals, oval)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	o.repository.Close()
	return ovals, nil
}
