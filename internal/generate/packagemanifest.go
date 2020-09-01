// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	apimanifests "github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	genutil "github.com/operator-framework/operator-sdk/internal/generate/internal"
)

const (
	// File extension for all PackageManifests written by Generator.
	packageManifestFileExt = ".package.yaml"
)

// PackageManifest configures the PackageManifest that GetBase() returns.
type PackageManifest struct {
	PackageName string
	BasePath    string
}

var (
	// User-facing errors.

	// ErrNoVersion if no version # has been provided
	ErrNoVersion = errors.New("version must be set")

	// Internal errors.

	// ErrNoOpName if the operator name has not been set
	ErrNoOpName = genutil.InternalError("operator name must be set")
	// ErrNoOutputDir if the directory to write the package manifest to has not been set
	ErrNoOutputDir = genutil.InternalError("output directory must be set")
)

type PkgOptions struct {
	// BaseDir is a directory to look for an existing base package manifest
	// to update.
	BaseDir string
	// ChannelName is operator's PackageManifest channel. If a new PackageManifest is generated
	// or ChannelName is the only channel in the generated PackageManifest,
	// this channel will be set to the PackageManifest's default.
	ChannelName string
	// IsDefaultChannel determines whether ChannelName should be the default channel in the
	// generated PackageManifest. If true, ChannelName will be the PackageManifest's default channel.
	// Setting this field is only necessary when more than one channel exists.
	IsDefaultChannel bool
}

// Generate configures the Generator with opts then runs it.
func (g Generator) GeneratePackageManifest(operatorName, version, outputDir string, opts ...*PkgOptions) error {
	if operatorName == "" {
		return ErrNoOpName
	} else if version == "" {
		return ErrNoVersion
	} else if outputDir == "" {
		return ErrNoOutputDir
	}
	myOpts := &PkgOptions{}
	if len(opts) > 0 {
		myOpts = opts[0]
	}

	pkg, err := g.generatePackageManifest(operatorName, version, myOpts)
	if err != nil {
		return err
	}

	outputWriter, err := genutil.Open(outputDir, makePkgManFileName(operatorName))
	if err != nil {
		return err
	}

	return genutil.WriteYAML(outputWriter, pkg)
}

// generatePackageManifest takes the input and generates the populated package manifest object
func (g *Generator) generatePackageManifest(operatorName, version string, opts *PkgOptions) (*apimanifests.PackageManifest, error) {
	b := PackageManifest{
		PackageName: operatorName,
	}
	if opts.BaseDir != "" {
		basePath := filepath.Join(opts.BaseDir, makePkgManFileName(operatorName))
		if genutil.IsNotExist(basePath) {
			basePath = ""
		}
		b.BasePath = basePath
	}
	base, err := b.GetBase()
	if err != nil {
		return nil, fmt.Errorf("error getting PackageManifest base: %v", err)
	}

	csvName := genutil.MakeCSVName(operatorName, version)
	if opts.ChannelName != "" {
		setChannels(base, opts.ChannelName, csvName)
		sortChannelsByName(base)
		if opts.IsDefaultChannel || len(base.Channels) == 1 {
			base.DefaultChannelName = opts.ChannelName
		}
	} else if len(base.Channels) == 0 {
		setChannels(base, "alpha", csvName)
		base.DefaultChannelName = "alpha"
	}

	if err = validatePackageManifest(base); err != nil {
		return nil, err
	}

	return base, nil
}

// makePkgManFileName will return the file name of a PackageManifest.
func makePkgManFileName(operatorName string) string {
	return strings.ToLower(operatorName) + packageManifestFileExt
}

// sortChannelsByName sorts pkg.Channels by each element's name.
func sortChannelsByName(pkg *apimanifests.PackageManifest) {
	sort.Slice(pkg.Channels, func(i int, j int) bool {
		return pkg.Channels[i].Name < pkg.Channels[j].Name
	})
}

// validatePackageManifest will validate pkg and log warnings and errors.
// If a validation error is encountered, an error is returned.
func validatePackageManifest(pkg *apimanifests.PackageManifest) error {
	if pkg == nil {
		return errors.New("empty PackageManifest")
	}

	hasErrors := false
	results := validation.PackageManifestValidator.Validate(pkg)
	for _, r := range results {
		for _, e := range r.Errors {
			log.Errorf("PackageManifest validation: [%s] %s", e.Type, e.Detail)
		}
		for _, w := range r.Warnings {
			log.Warnf("PackageManifest validation: [%s] %s", w.Type, w.Detail)
		}
		if r.HasError() {
			hasErrors = true
		}
	}

	if hasErrors {
		return errors.New("invalid generated PackageManifest")
	}

	return nil
}

// setChannels checks for duplicate channels in pkg and sets the default channel if possible.
func setChannels(pkg *apimanifests.PackageManifest, channelName, csvName string) {
	channelIdx := -1
	for i, channel := range pkg.Channels {
		if channel.Name == channelName {
			pkg.Channels[i].CurrentCSVName = csvName
			channelIdx = i
			break
		}
	}
	if channelIdx == -1 {
		pkg.Channels = append(pkg.Channels, apimanifests.PackageChannel{
			Name:           channelName,
			CurrentCSVName: csvName,
		})
	}
}

// GetBase returns a base PackageManifest, populated either with default
// values or, if b.BasePath is set, bytes from disk.
func (b PackageManifest) GetBase() (base *apimanifests.PackageManifest, err error) {
	if b.BasePath != "" {
		if base, err = readPackageManifestBase(b.BasePath); err != nil {
			return nil, fmt.Errorf("error reading existing PackageManifest base %s: %v", b.BasePath, err)
		}
	} else {
		base = b.makeNewBase()
	}

	return base, nil
}

// makeNewBase returns a base makeNewBase to modify.
func (b PackageManifest) makeNewBase() *apimanifests.PackageManifest {
	return &apimanifests.PackageManifest{
		PackageName: b.PackageName,
	}
}

// readPackageManifestBase returns the PackageManifest base at path.
// If no base is found, readPackageManifestBase returns an error.
func readPackageManifestBase(path string) (*apimanifests.PackageManifest, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pkg := &apimanifests.PackageManifest{}
	if err := yaml.Unmarshal(b, pkg); err != nil {
		return nil, fmt.Errorf("error unmarshalling PackageManifest from %s: %w", path, err)
	}
	if pkg.PackageName == "" {
		return nil, fmt.Errorf("no PackageManifest in %s", path)
	}
	return pkg, nil
}
