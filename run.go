package growibackup

import "golang.org/x/xerrors"

func Run(jsonPATH, outputRootPATH string) error {
	revisions, err := convertStructFromJSON(jsonPATH)
	if err != nil {
		return xerrors.Errorf("[error] fialed convert struct from json %+w", err)
	}
	revisions = revisions.createUniqRevisions()
	return revisions.backup(outputRootPATH)
}
