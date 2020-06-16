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

package build

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running a build command", func() {
	Describe("newCmd", func() {
		It("builds and returns a cobra command", func() {
			cmd := NewCmd()
			Expect(cmd).NotTo(BeNil())

			flag := cmd.Flags().Lookup("image-build-args")
			Expect(flag).NotTo(BeNil())

			flag = cmd.Flags().Lookup("image-builder")
			Expect(flag).NotTo(BeNil())
			Expect(flag.DefValue).To(Equal("docker"))

			flag = cmd.Flags().Lookup("go-build-args")
			Expect(flag).NotTo(BeNil())
		})
	})

	Describe("validate", func() {
		It("fails if not exactly 1 arg is provided", func() {
			cmd := buildCmd{}

			err := cmd.validate([]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("an image name is required"))

			err = cmd.validate([]string{"a", "b"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("an image name is required"))
		})
		It("succeeds if exactly 1 arg is provided", func() {
			cmd := buildCmd{}

			err := cmd.validate([]string{"c"})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
