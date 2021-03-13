package osext

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOSExt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YamlFile Test Suite")
}

var _ = Describe("Package osext functions", func() {
	var fileContents = `
	
	

This some text.

And More Text

`
	var fileContentsTrimmed = strings.TrimSpace(fileContents)

	var testFile *os.File

	BeforeSuite(func() {
		var err error
		testFile, err = os.CreateTemp("", "testFile*.txt")
		Expect(err).ToNot(HaveOccurred())
		Expect(testFile).ToNot(BeNil())

		err = os.WriteFile(testFile.Name(), []byte(fileContents), 0644)
		Expect(err).ToNot(HaveOccurred())
	})
	AfterSuite(func() {
		if testFile != nil {
			os.Remove(testFile.Name())
		}
	})

	Context("Checking if file(s) exist or for reading", func() {
		When("checking if a file is a directory using IsDirectory()", func() {
			It("returns true when calling IsDirectory() on a directory that does exist", func() {
				// dir := path.Dir(testFile.Name())
				dir, _ := filepath.Split(testFile.Name())
				Expect(IsDirectory(dir)).To(BeTrue())
			})
			It("returns false when calling IsDirectory() on a directory that does not exist", func() {
				// dir := path.Dir(testFile.Name())
				dir, _ := filepath.Split(testFile.Name())
				Expect(IsDirectory(dir + "-foo")).To(BeFalse())
			})
		})
		When("Checking if file exists using FileExists()", func() {
			It("returns true on a file that does exist", func() {
				Expect(FileExists(testFile.Name())).To(BeTrue())
			})
			It("returns false on a file that does not exist", func() {
				Expect(FileExists(testFile.Name() + "fpp")).To(BeFalse())
			})
			It("returns false on an empty filename", func() {
				Expect(FileExists("")).To(BeFalse())
			})
		})
		When("Checking if file or directory exists using FileOrDirectoryExists()", func() {
			It("returns true on a file that does exist", func() {
				Expect(FileOrDirectoryExists(testFile.Name())).To(BeTrue())
			})
			It("returns false on a file that does not exist", func() {
				Expect(FileOrDirectoryExists(testFile.Name() + "fpp")).To(BeFalse())
			})
			It("returns false on an empty filename", func() {
				Expect(FileOrDirectoryExists("")).To(BeFalse())
			})
		})
		When("Reading file contents as string using ReadFileAsString", func() {
			It("returns error when trying to read bytes from a file that does not exist (with and without trimming)", func() {
				_, err := ReadFile(testFile.Name()+"fpp", false)
				Expect(err).To(HaveOccurred())
				_, err = ReadFile(testFile.Name()+"fpp", true)
				Expect(err).To(HaveOccurred())
			})
			It("retrieves the exepcted contents as []byte (with and without trimming)", func() {
				var (
					bytes []byte
					err   error
				)
				bytes, err = ReadFile(testFile.Name(), false)
				Expect(err).ToNot(HaveOccurred())
				Expect(bytes).ToNot(BeEmpty())
				Expect(reflect.DeepEqual(bytes, []byte(fileContents))).To(BeTrue())

				bytes, err = ReadFile(testFile.Name(), true)
				Expect(err).ToNot(HaveOccurred())
				Expect(bytes).ToNot(BeEmpty())
				Expect(reflect.DeepEqual(bytes, []byte(fileContentsTrimmed))).To(BeTrue())
			})
			It("retrieves the expected contents as string (with and without trimming)", func() {
				var (
					readStr string
					err     error
				)
				readStr, err = ReadFileAsString(testFile.Name(), false)
				Expect(err).ToNot(HaveOccurred())
				Expect(readStr).ToNot(BeEmpty())
				Expect(readStr).To(Equal(fileContents))

				readStr, err = ReadFileAsString(testFile.Name(), true)
				Expect(err).ToNot(HaveOccurred())
				Expect(readStr).ToNot(BeEmpty())
				Expect(readStr).To(Equal(fileContentsTrimmed))
			})
			It("returns error when trying to read bytes from an empty filename (with and without trimming)", func() {
				_, err := ReadFile("", false)
				Expect(err).To(HaveOccurred())
				_, err = ReadFile("", true)
				Expect(err).To(HaveOccurred())
			})
			It("returns error when trying to read as string from an empty filename (with and without trimming)", func() {
				_, err := ReadFileAsString("", false)
				Expect(err).To(HaveOccurred())
				_, err = ReadFileAsString("", true)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	When("Creating temp file with contents", func() {
		var tmpFile *os.File
		BeforeEach(func() {

		})
		AfterEach(func() {
			if tmpFile != nil {
				os.Remove(tmpFile.Name())
			}
		})
		It("creates a file with the expected contents", func() {
			testFile, err := CreateTempWithContents("", "testFile2*.txt", []byte(fileContentsTrimmed), 0644)
			// Delete file after
			defer func() {
				if testFile != nil {
					os.Remove(testFile.Name())
				}
			}()
			Expect(err).ToNot(HaveOccurred())
			Expect(testFile).ToNot(BeNil())

			var bytes []byte
			bytes, err = os.ReadFile(testFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(bytes).ToNot(BeEmpty())
			Expect(reflect.DeepEqual(bytes, []byte(fileContentsTrimmed))).To(BeTrue())
		})
		It("returns an error when it cannot create file", func() {
			testFile, err := CreateTempWithContents("", "/foo*.txt", []byte(fileContentsTrimmed), 0644)
			defer func() {
				if testFile != nil {
					os.Remove(testFile.Name())
				}
			}()
			Expect(err).To(HaveOccurred())
		})
	})
})
