package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
	"testing"

	"io/ioutil"
	"os"
	"runtime"

	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	commandPath string
	homeDir     string
	server      *Server
)

var _ = BeforeEach(func() {
	var err error
	homeDir, err = ioutil.TempDir("", "cm-test")
	Expect(err).NotTo(HaveOccurred())

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", homeDir)
	} else {
		os.Setenv("HOME", homeDir)
	}

	server = NewServer()

	runCommand("api", server.URL())
})

var _ = AfterEach(func() {
	server.Close()
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/pivotal-cf/cm-cli")
	Expect(err).NotTo(HaveOccurred())
	return []byte(path)
}, func(data []byte) {
	commandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})

func runCommand(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}
