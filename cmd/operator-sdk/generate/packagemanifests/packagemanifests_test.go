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

package packagemanifests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPackagemanifests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Packagemanifests Suite")
}

var _ = Describe("Running a generate packagemanifest command", func() {
	Describe("NewCmd", func() {
		It("Builds and returns a cobra command", func() {
			cmd := NewCmd()
			Expect(*cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("packagemanifests"))
			Expect(cmd.Short).To(Equal("Generates a package manifests format"))

			flag := cmd.Flags().Lookup("kustomize")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Generate kustomize bases"))

			flag = cmd.Flags().Lookup("manifests")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Generate package manifests"))

			flag = cmd.Flags().Lookup("stdout")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Write package to stdout"))

			flag = cmd.Flags().Lookup("operator-name")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Name of the packaged operator"))

			flag = cmd.Flags().Lookup("version")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Shorthand).To(Equal("v"))
			Expect(flag.Usage).To(ContainSubstring("Semantic version of the packaged operator"))

			flag = cmd.Flags().Lookup("input-dir")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Directory to read existing package manifests from. " +
				"This directory is the parent of individual versioned package directories, and different from --manifest-root"))

			flag = cmd.Flags().Lookup("output-dir")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Directory in which to write package manifests"))

			flag = cmd.Flags().Lookup("manifest-root")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Hidden).To(BeTrue())
			Expect(flag.Usage).To(ContainSubstring("Root directory for operator manifests such as " +
				"Deployments and RBAC, ex. 'deploy' or 'config'. This directory is different from that passed to --input-dir"))

			flag = cmd.Flags().Lookup("apis-dir")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Root directory for API type defintions"))

			flag = cmd.Flags().Lookup("crds-dir")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Root directory for CustomResoureDefinition manifests"))

			flag = cmd.Flags().Lookup("channel")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("Channel name for the generated package"))

			flag = cmd.Flags().Lookup("default-channel")
			Expect(flag).NotTo(BeNil())
			Expect(flag.DefValue).To(Equal("false"))
			Expect(flag.Usage).To(ContainSubstring("Use the channel passed to --channel " +
				"as the package manifest file's default channel"))

			flag = cmd.Flags().Lookup("update-crds")
			Expect(flag).NotTo(BeNil())
			Expect(flag.DefValue).To(Equal("false"))
			Expect(flag.Usage).To(ContainSubstring("Update CustomResoureDefinition manifests " +
				"in this package"))

			flag = cmd.Flags().Lookup("quiet")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Shorthand).To(Equal("q"))
			Expect(flag.DefValue).To(Equal("false"))
			Expect(flag.Usage).To(ContainSubstring("Run in quiet mode"))

			flag = cmd.Flags().Lookup("interactive")
			Expect(flag).NotTo(BeNil())
			Expect(flag.DefValue).To(Equal("false"))
			Expect(flag.Usage).To(ContainSubstring("When set or no package base exists, an interactive " +
				"command prompt will be presented to accept package ClusterServiceVersion metadata"))
		})
	})
	Describe("validateManifests", func() {
		PIt("validates the fields to generate a manifest", func() {
			cmd := packagemanifestsCmd{}

			err := cmd.validateManifests()
			Expect(err).NotTo(HaveOccurred())
		})
		It("fails without a valid version", func() {
			cmd := packagemanifestsCmd{}

			err := cmd.validateManifests()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("is not a valid semantic version: Version string empty"))
		})
		It("fails if manifest-root is not set", func() {
			cmd := packagemanifestsCmd{
				version: "1.0.0",
			}

			err := cmd.validateManifests()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("--manifest-root must be set if not reading from stdin"))
		})
	})
})
