package growibackup

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

// Revision is growy entry struct
type Revision struct {
	ID            string    `json:"_id"`
	Format        string    `json:"format"`
	CreatedAt     time.Time `json:"createdAt"`
	Path          string    `json:"path"`
	Body          string    `json:"body"`
	Author        string    `json:"author"`
	HasDiffToPrev bool      `json:"hasDiffToPrev"`
	V             int       `json:"__v"`
}

// Revisions is as equal as revision.json
type Revisions []Revision

func convertStructFromJSON(jsonPATH string) (*Revisions, error) {
	file, err := ioutil.ReadFile(jsonPATH)
	if err != nil {
		return nil, xerrors.Errorf("[error] failed open %s, %+w", jsonPATH, err)
	}

	revs := Revisions{}
	err = json.Unmarshal([]byte(file), &revs)

	if err != nil {
		return nil, xerrors.Errorf("[error] failed unmarshal json %s, %+w", jsonPATH, err)
	}

	return &revs, nil
}

func (revs *Revisions) createUniqRevisions() *Revisions {
	path2Revision := make(map[string]Revision)
	for _, rev := range *revs {
		if prevRev, ok := path2Revision[rev.Path]; ok {
			if prevRev.CreatedAt.After(rev.CreatedAt) {
				path2Revision[rev.Path] = rev
			}
			continue
		}
		path2Revision[rev.Path] = rev
	}

	newRevs := make(Revisions, 0, len(path2Revision))
	for _, rev := range path2Revision {
		newRevs = append(newRevs, rev)
	}
	return &newRevs
}

func (revs *Revisions) backup(outputPATH string) error {
	var sbuilder strings.Builder
	path2Exists := make(map[string]bool)

	oPATH, err := getAbsRootPATH(outputPATH)
	if err != nil {
		return xerrors.Errorf("[error] failed get abs at %s %+w", outputPATH, err)
	}

	for _, rev := range *revs {
		sbuilder.WriteString(rev.Path)
		sbuilder.WriteString(".md")
		path := sbuilder.String()
		sbuilder.Reset()

		mdPATH := filepath.Join(oPATH, path)
		dirPATH := filepath.Dir(mdPATH)
		if _, ok := path2Exists[dirPATH]; !ok {
			err := checkAfterMkdir(dirPATH)
			if err != nil {
				return err
			}
		}

		f, err := os.Create(mdPATH)
		if err != nil {
			return xerrors.Errorf("[error] failed create %s, %+w", mdPATH, err)
		}
		w := bufio.NewWriter(f)
		w.WriteString(rev.Body)
		w.Flush()
		f.Close()
	}
	return nil
}

func checkAfterMkdir(dirpath string) error {
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		err := os.MkdirAll(dirpath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAbsRootPATH(outputPath string) (string, error) {
	if !filepath.IsAbs(outputPath) {
		oPATH, err := filepath.Abs(outputPath)
		if err != nil {
			return "", err
		}
		return oPATH, nil
	}
	return outputPath, nil
}
