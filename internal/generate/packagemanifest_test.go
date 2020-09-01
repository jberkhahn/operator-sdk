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

package generate_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/operator-sdk/internal/generate"
)

var _ = Describe("A package manifest generator", func() {
	Describe("GeneratePackageManifest", func() {
		var (
			g                                    generate.Generator
			operatorName                         string
			outputDir                            string
			pkgManFilename                       string
			pkgManDefault                        string
			pkgManOneChannel                     string
			pkgManUpdatedOneChannel              string
			pkgManUpdatedSecondChannel           string
			pkgManUpdatedSecondChannelNewDefault string
		)
		BeforeEach(func() {
			operatorName = "memcached-operator"
			pkgManFilename = operatorName + ".package.yaml"
			outputDir = os.TempDir()
			pkgManDefault = `channels:
- currentCSV: memcached-operator.v0.0.1
  name: alpha
defaultChannel: alpha
packageName: memcached-operator
`
			pkgManOneChannel = `channels:
- currentCSV: memcached-operator.v0.0.1
  name: stable
defaultChannel: stable
packageName: memcached-operator
`
			pkgManUpdatedOneChannel = `channels:
- currentCSV: memcached-operator.v0.0.2
  name: alpha
defaultChannel: alpha
packageName: memcached-operator
`
			pkgManUpdatedSecondChannel = `channels:
- currentCSV: memcached-operator.v0.0.1
  name: alpha
- currentCSV: memcached-operator.v0.0.2
  name: stable
defaultChannel: alpha
packageName: memcached-operator
`
			pkgManUpdatedSecondChannelNewDefault = `channels:
- currentCSV: memcached-operator.v0.0.1
  name: alpha
- currentCSV: memcached-operator.v0.0.2
  name: stable
defaultChannel: stable
packageName: memcached-operator
`
		})
		Context("when writing a new package manifest", func() {
			It("writes a package manifest", func() {
				err := g.GeneratePackageManifest(operatorName, "0.0.1", outputDir)
				Expect(err).NotTo(HaveOccurred())
				file, err := ioutil.ReadFile(outputDir + string(os.PathSeparator) + pkgManFilename)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(file)).To(Equal(pkgManDefault))
			})
			It("writes a package manifest with a non-default channel", func() {
				opts := &generate.PkgOptions{
					ChannelName: "stable",
				}

				err := g.GeneratePackageManifest(operatorName, "0.0.1", outputDir, opts)
				Expect(err).NotTo(HaveOccurred())
				file, err := ioutil.ReadFile(outputDir + string(os.PathSeparator) + pkgManFilename)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(file)).To(Equal(pkgManOneChannel))
			})
		})
		Context("when updating an existing package manifest", func() {
			It("creates a new package manifest if provided an existing packagemanifest that doesn't exist", func() {
				opts := &generate.PkgOptions{
					BaseDir:     "testpotato",
					ChannelName: "stable",
				}

				err := g.GeneratePackageManifest(operatorName, "0.0.1", outputDir, opts)
				Expect(err).NotTo(HaveOccurred())
				file, err := ioutil.ReadFile(outputDir + string(os.PathSeparator) + pkgManFilename)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(file)).To(Equal(pkgManOneChannel))
			})
			It("updates an existing package manifest with a updated channel", func() {
				opts := &generate.PkgOptions{
					BaseDir:     "testdata",
					ChannelName: "alpha",
				}

				err := g.GeneratePackageManifest(operatorName, "0.0.2", outputDir, opts)
				Expect(err).NotTo(HaveOccurred())
				file, err := ioutil.ReadFile(outputDir + string(os.PathSeparator) + pkgManFilename)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(file)).To(Equal(pkgManUpdatedOneChannel))
			})
			It("updates an existing package manifest with a new channel", func() {
				opts := &generate.PkgOptions{
					BaseDir:     "testdata",
					ChannelName: "stable",
				}

				err := g.GeneratePackageManifest(operatorName, "0.0.2", outputDir, opts)
				Expect(err).NotTo(HaveOccurred())
				file, err := ioutil.ReadFile(outputDir + string(os.PathSeparator) + pkgManFilename)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(file)).To(Equal(pkgManUpdatedSecondChannel))
			})
			It("updates an existing package manifest with a new channel and an updated default channel", func() {
				opts := &generate.PkgOptions{
					BaseDir:          "testdata",
					ChannelName:      "stable",
					IsDefaultChannel: true,
				}

				err := g.GeneratePackageManifest(operatorName, "0.0.2", outputDir, opts)
				Expect(err).NotTo(HaveOccurred())
				file, err := ioutil.ReadFile(outputDir + string(os.PathSeparator) + pkgManFilename)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(file)).To(Equal(pkgManUpdatedSecondChannelNewDefault))
			})
		})
		Context("when incorrect params are provided", func() {
			It("fails if no operator name is specified", func() {
				err := g.GeneratePackageManifest("", "", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(generate.ErrNoOpName.Error()))
			})
			It("fails if no version is specified", func() {
				err := g.GeneratePackageManifest(operatorName, "", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(generate.ErrNoVersion.Error()))
			})
			It("fails if no output directory is set", func() {
				err := g.GeneratePackageManifest(operatorName, "0.0.1", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(generate.ErrNoOutputDir.Error()))
			})
		})
	})
	Describe("GetBase", func() {
		var (
			b generate.PackageManifest
		)
		BeforeEach(func() {
			b = generate.PackageManifest{}
		})
		It("returns a new blank packagemanifest", func() {
			b.PackageName = "sweetsop"

			pm, err := b.GetBase()
			Expect(err).NotTo(HaveOccurred())
			Expect(pm).NotTo(BeNil())
			Expect(pm.PackageName).To(Equal(b.PackageName))
		})
		It("reads an existing packagemanifest from disk", func() {
			b.BasePath = "testdata/memcached-operator.package.yaml"

			pm, err := b.GetBase()
			Expect(err).NotTo(HaveOccurred())
			Expect(pm).NotTo(BeNil())
			Expect(pm.PackageName).To(Equal("memcached-operator"))
			Expect(len(pm.Channels)).To(Equal(1))
			Expect(pm.Channels[0].Name).To(Equal("alpha"))
			Expect(pm.Channels[0].CurrentCSVName).To(Equal("memcached-operator.v0.0.1"))
			Expect(pm.DefaultChannelName).To(Equal("alpha"))
		})
		It("fails if provided a non-existent base path", func() {
			b.BasePath = "not-a-real-thing.yaml"

			pm, err := b.GetBase()
			Expect(pm).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error reading existing"))
		})
	})
})
