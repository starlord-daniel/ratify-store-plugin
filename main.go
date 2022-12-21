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

type conf struct {
	Name       string `json:"name"`
	FolderPath string `json:"folderPath"`
}

func ListReferrers(args *skel.CmdArgs, subjectReference common.Reference, artifactTypes []string, nextToken string, subjectDesc *ocispecs.SubjectDescriptor) (*referrerstore.ListReferrersResult, error) {

	localFileReference, err := GetSubjectDescriptor(args, subjectReference)
	if err != nil {
		log.Fatal(err)
	}
	var referrers []ocispecs.ReferenceDescriptor

	// ! Should be passed in as store args in ratify config.json
	localconfig := localFileConfig(args)

	directory := filepath.Join(localconfig.FolderPath, "/references/")

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

			if localFileReference.Digest == referenceSubject.Digest {
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

func localFileConfig(args *skel.CmdArgs) *conf {
	// already validated before this point, thank you skel.
	var config conf
	json.Unmarshal(args.StdinData, &config)
	// verify the filepath was provided.
	if config.FolderPath == "" {
		log.Fatalf("missing folderPath value in config.json for the filesystem store")
	}
	return &config
}

// called from verifier to retrieve referrer content
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

	localConfig := localFileConfig(args)

	path := filepath.Join(localConfig.FolderPath, subjectReference.Original)

	file, err := os.ReadFile(path)

	if err != nil {
		logrus.Fatal(err, "Manifest File Read Error")
	}

	type subjectConfig struct {
		MediaType string
		Size      int
		Digest    digest.Digest // this really should be a string, I think.
	}

	type SubjectManifest struct {
		MediaType string          `json:"mediaType"`
		Config    subjectConfig   `json:"config"`
		Layers    []subjectConfig `json:"layers"`
	}

	var fileContent SubjectManifest
	err = json.Unmarshal(file, &fileContent)
	if err != nil {
		logrus.Fatal(err, "Manifest File Content Error")
	}

	return &ocispecs.SubjectDescriptor{Descriptor: v1.Descriptor{Digest: fileContent.Config.Digest}}, nil
}

func GetReferenceManifest(args *skel.CmdArgs, subjectReference common.Reference, digest digest.Digest) (ocispecs.ReferenceManifest, error) {

	localConfig := localFileConfig(args)
	directory := filepath.Join(localConfig.FolderPath, "/references")

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
