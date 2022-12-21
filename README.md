# Filesystem referrer store

## Terminology

| Term | Description |
| ---- | ----- |
| [Referrer Store](https://pkg.go.dev/github.com/deislabs/ratify/pkg/referrerstore#ReferrerStore) | ReferrerStore is an interface that defines methods to query the graph of supply chain content including its related content. Take a look at [the interface definition in api.go](https://github.com/deislabs/ratify/tree/main/pkg/referrerstore/api.go) |
| [Reference Manifest](https://pkg.go.dev/github.com/deislabs/ratify/pkg/ocispecs#ReferenceManifest) | The manifest of a referrer. The referrer manifest defines the connection to a subject manifest via its digest. The digest acts as a unique id of the subject, that the referrer is storing. |
| [Subject Reference](https://pkg.go.dev/github.com/deislabs/ratify/pkg/common#Reference) | A reference to the subject. This contains the unique digest to look for, the path for validation, a tag and a original. |

## Description

This is a sample implementation of a referrer store that stores the referrer in files on the filesystem.

The job of this and any other referrer store is to provide methods defined in the [ReferrerStore interface](https://github.com/deislabs/ratify/pkg/referrerstore/api.go) that are used by referrers like [NotaryV2](https://github.com/deislabs/ratify/tree/main/pkg/verifier/notaryv2/notaryv2.go) to get data to verify.

The important files are:

- The referrer manifests: `data/references/referenceManifest1.json` and `data/references/referenceManifest1.json`
- The subject manifest: `data/subjectManifest.json`
- The content of the artifact: `data/artifactContent.json`
- The code that defines the filesystem storage: `main.go`

## Concept

The `main.go` file contains all methods, that this implementation of the referrerstore needs to provide to the calling application (a verifier).

These methods are defined in the [ReferrerStore interface](https://github.com/deislabs/ratify/tree/main/pkg/referrerstore/api.go) and are:

### ListReferrers

The `ListReferrers` method gets the subject reference as an input and takes a look at its digest.

The method looks for all references it can find in the filesystem and returns all of them, that reference the subject references digest.

The output is a [ListReferrerResult](https://pkg.go.dev/github.com/deislabs/ratify/pkg/referrerstore#ListReferrersResult)

### GetBlobContent

The `GetBlobContent` method returns the blob with the given digest.

It takes the `subjectReference`  and `digest` as an input. The `digest` here is the name of the artifact that is stored in the `ReferenceManifest`.

In this filesystem implementation, ... (Logic)

This method returns the contents of the referenced blob (artifact) as a byte array.

### GetReferenceManifest

The `GetReferenceManifest` method returns the reference artifact manifest as given by the descriptor.

It takes in a subject reference a returns a single reference with the same digest as itself.

### GetSubjectDescriptor

The `GetSubjectDescriptor` method returns the descriptor for the given subject. This descriptor is an [oci.Descriptor object](https://pkg.go.dev/github.com/opencontainers/image-spec/specs-go/v1#Descriptor).

As an input we get the `subjectReference` of type [common Reference](https://pkg.go.dev/github.com/deislabs/ratify/pkg/common#Reference).

From this `subjectReference` the digest is extracted and used to build the Descriptor object, which is then returned.

## Running/Testing

To run/test the behaviour of your storage, you need to build it first.
The plugin itself is not part of ratify (called out-of-tree) and therefor need to be built separately from it.

To do this, run this go command from the root of the repo:

```bash
go build
```

The ratify project can be found and [cloned from GitHub](https://github.com/deislabs/ratify).
To use this plugin with ratify, you'd need to move the executable to your ratify directory, into the `ratify/.ratify/plugins` folder.

>info: Currently, the data for the filesystem storage has to be copied to the users `.ratify` folder.

A valid structure (in the [main ratify repo](https://github.com/deislabs/ratify)) can look as follows:

```bash
.ratify/plugins/
├── filesystem
└── data
    ├── references
    │   ├── referenceManifest1.json
    │   └── referenceManifest2.json
    └── manifest.json
    └── artifactContent.json
```

Now you need to update the `config.json` in the `.ratify` folder to include your storage. An updated value can look like this:

```json
"store": {
    "version": "1.0.0",
    "plugins":
    [
        {
            "name": "filesystem",
            "folderPath": "<folder path /data/ folder >"
        }
    ]
    },
```

After this is done, you can run ratify by using the debug configuration `Verify` as defined in the `.vscode/launch.json` by pressing F5 in VS Code (in the ratify repo). When prompted provide the subject manifest file name e.g. `manifest.json` that resides in the folder path provided to the store plugin.

## Sample

To test, there is an included `artifactContent.json` which the `licensechecker` should pick up if included in the config.

```json
{
    "verifier": {
    "version": "1.0.0",
    "plugins": [
    {
        "name": "licensechecker",
        "artifactTypes": "application/vnd.ratify.spdx.v0",
        "allowedLicenses": [
          "MIT"
        ]
    }]
    }
}
```

The resulting output should demonstrate the feasibility of local filesystem store.
