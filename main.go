package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/deislabs/ratify/pkg/common"
	"github.com/deislabs/ratify/pkg/ocispecs"
	"github.com/deislabs/ratify/pkg/referrerstore"
	"github.com/deislabs/ratify/pkg/referrerstore/plugin/skel"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
)

func main() {
	skel.PluginMain("filesystem", "1.0.0", ListReferrers, GetBlobContent, GetReferenceManifest, GetSubjectDescriptor, []string{"1.0.0"})
}

// FileSystem store local subjectReference.
var localsubjectReference = common.Reference{
	Path:     "/home/vscode/.ratify/plugins/data/sha256:8fef74aa37a48e8c9e911f5ca1d10809d01660ec5a42bb3514bb75e508d0276d.json",
	Digest:   "sha256:8fef74aa37a48e8c9e911f5ca1d10809d01660ec5a42bb3514bb75e508d0276d",
	Tag:      "",
	Original: "/home/vscode/.ratify/plugins/data/sha256:8fef74aa37a48e8c9e911f5ca1d10809d01660ec5a42bb3514bb75e508d0276d.json",
}

var localArtifactReference = common.Reference{
	Path:     "/home/vscode/.ratify/plugins/data/references/referenceManifest1.json",
	Digest:   "/home/vscode/.ratify/plugins/data/references/referenceManifest1.json",
	Tag:      "",
	Original: "/home/vscode/.ratify/plugins/data/references/referenceManifest1.json",
}

func ListReferrers(args *skel.CmdArgs, subjectReference common.Reference, artifactTypes []string, nextToken string, subjectDesc *ocispecs.SubjectDescriptor) (*referrerstore.ListReferrersResult, error) {

	subjectReference = localsubjectReference

	var referrers []ocispecs.ReferenceDescriptor

	directory := filepath.Dir("/home/vscode/.ratify/plugins/data/references/")

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	// Go through the directory to retrieve all reference json files

	// ! this is the implementation for GetReferenceManifest + appending them to the referrers slice. (should refactor and call GetReferenceManifest)
	for _, referenceFile := range files {

		//! Just assume we have the right file types.
		reference := handleJSONFiles(fmt.Sprintf("%s/%s", directory, referenceFile.Name()))

		for _, referenceSubject := range reference.Subjects {

			if subjectReference.Digest == referenceSubject.Digest {
				logrus.Info("Found reference: ", referenceFile.Name())
				referrers = append(referrers, ocispecs.ReferenceDescriptor{
					ArtifactType: reference.ArtifactType,
					Descriptor: v1.Descriptor{
						MediaType: referenceSubject.MediaType,
						Digest:    referenceSubject.Digest,
						Size:      referenceSubject.Size,
					},
				})
				break
			}
		}
	}

	return &referrerstore.ListReferrersResult{
		Referrers: referrers,
		NextToken: "",
	}, nil
}

func handleJSONFiles(filePath string) ocispecs.ReferenceManifest {

	readFile, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	var fileContent ocispecs.ReferenceManifest
	err = json.Unmarshal(readFile, &fileContent)

	if err != nil {
		log.Fatal(err)
	}

	return fileContent
}

func GetBlobContent(args *skel.CmdArgs, subjectReference common.Reference, digest digest.Digest) ([]byte, error) {

	artifactPath := subjectReference.Digest.String()

	if digest != "" {
		artifactPath = digest.String()
	}

	// get content of artifact file
	content, err := os.ReadFile(artifactPath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func GetSubjectDescriptor(args *skel.CmdArgs, subjectReference common.Reference) (*ocispecs.SubjectDescriptor, error) {

	subjectReference = localsubjectReference

	dig := subjectReference.Digest

	if dig == "" {
		dig = digest.FromString(subjectReference.Tag)
	}

	return &ocispecs.SubjectDescriptor{Descriptor: v1.Descriptor{Digest: dig}}, nil
}

func GetReferenceManifest(args *skel.CmdArgs, subjectReference common.Reference, digest digest.Digest) (ocispecs.ReferenceManifest, error) {

	subjectReference = localArtifactReference

	directory := filepath.Dir("/home/vscode/.ratify/plugins/filesystem_data/references/")

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	var referenceManifest ocispecs.ReferenceManifest

	for _, referenceFile := range files {
		reference := handleJSONFiles(fmt.Sprintf("%s/%s", directory, referenceFile.Name()))
		for _, artifactSubject := range reference.Subjects {
			if artifactSubject.Digest == subjectReference.Digest {
				fmt.Printf("Found reference: %s \n", referenceFile.Name())
				referenceManifest = reference
				break
			}
		}
	}
	return referenceManifest, nil
}
