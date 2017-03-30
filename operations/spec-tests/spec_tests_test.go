package spec_tests_test

import (
	"os/exec"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("cf-d opsfiles", func() {
	var boshLiteManifest []byte
	var singleAZManifest []byte
	// var singleAZStruct interface

	Describe("Applying single-az ops file to bosh-lite manifest", func() {
		BeforeEach(func() {
			// Setup a bosh-lite manifest
			command := exec.Command("bosh2",
				"interpolate",
				"cf-deployment.yml",
				"-o", "operations/bosh-lite.yml")
			command.Dir = "../../"

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Eventually(session.Out, 30).Should(Say("name"))
			boshLiteManifest = session.Out.Contents()

			// Setup a bosh-lite-and-single-az manifest
			applySingleAZCommand := exec.Command("bosh2",
				"interpolate",
				"cf-deployment.yml",
				"-o", "operations/bosh-lite.yml",
				"-o", "operations/scale-to-one-az.yml")
			applySingleAZCommand.Dir = "../.."

			session, err = gexec.Start(applySingleAZCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 30).Should(gexec.Exit(0))
			Eventually(session.Out, 30).Should(Say("name"))
			singleAZManifest = session.Out.Contents()

			err = yaml.Unmarshal(singleAZManifest, singleAZStruct)
		})
		It("Doesn't change it", func() {
			Expect(len(boshLiteManifest)).NotTo(Equal(0))
			Expect(len(singleAZManifest)).NotTo(Equal(0))

			Expect(boshLiteManifest).To(Equal(singleAZManifest))

		})
	})
})
