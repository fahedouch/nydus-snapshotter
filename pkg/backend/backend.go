/*
 * Copyright (c) 2022. Nydus Developers. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package backend

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/v2/core/content"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	BackendTypeOSS     = "oss"
	BackendTypeS3      = "s3"
	BackendTypeLocalFS = "localfs"
)

var (
	// We always use multipart upload for backend, and limit the
	// multipart chunk size to 500MB by default.
	MultipartChunkSize int64 = 500 * 1024 * 1024
)

// Backend uploads blobs generated by nydus-image builder to a backend storage.
type Backend interface {
	// Push pushes specified blob file to remote storage backend.
	Push(ctx context.Context, cs content.Store, desc ocispec.Descriptor) error
	// Check checks whether a blob exists in remote storage backend,
	// blob exists -> return (blobPath, nil)
	// blob not exists -> return ("", err)
	Check(blobDigest digest.Digest) (string, error)
	// Type returns backend type name.
	Type() string
}

// Nydus driver majorly works for registry backend, which means blob is stored in
// registry as per OCI distribution specification. But nydus can also make OSS or
// other storage services as backend storage. Pass config as byte slice here because
// we haven't find a way to represent all backend config at the same time.
func NewBackend(_type string, config []byte, forcePush bool) (Backend, error) {
	switch _type {
	case BackendTypeOSS:
		return newOSSBackend(config, forcePush)
	case BackendTypeS3:
		return newS3Backend(config, forcePush)
	case BackendTypeLocalFS:
		return newLocalFSBackend(config, forcePush)
	default:
		return nil, fmt.Errorf("unsupported backend type %s", _type)
	}
}
