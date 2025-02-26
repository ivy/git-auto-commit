package template_test

import (
	"io"
	"strings"
	"testing"

	"github.com/ivy/git-auto-commit/template"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestTemplate is the entry point for the Ginkgo test suite.
func TestTemplate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Template Suite")
}

var _ = Describe("Engine", func() {
	var (
		engine *template.Engine
	)

	BeforeEach(func() {
		// Create a new template engine before each test
		engine = template.New()
	})

	Context("New()", func() {
		It("should create a new template engine without error", func() {
			Expect(engine).NotTo(BeNil())
		})
	})

	Context("Lookup(name)", func() {
		When("the template exists", func() {
			It("should return the template without error", func() {
				tmpl, lookupErr := engine.Lookup("prompt/commit.tmpl") // Replace with your actual embedded template name
				Expect(lookupErr).NotTo(HaveOccurred())
				Expect(tmpl).NotTo(BeNil())
			})
		})

		When("the template does not exist", func() {
			It("should return an error indicating the template was not found", func() {
				tmpl, lookupErr := engine.Lookup("nonexistent.tmpl")
				Expect(lookupErr).To(HaveOccurred())
				Expect(tmpl).To(BeNil())
			})
		})
	})

	Context("RenderBytes(name, data)", func() {
		When("the template exists", func() {
			It("should render the template to bytes without error", func() {
				bytesOutput, renderErr := engine.RenderBytes(
					"prompt/commit.tmpl",
					map[string]string{
						"Message": "Ginkgo",
					},
				)
				Expect(renderErr).NotTo(HaveOccurred())
				Expect(bytesOutput).NotTo(BeEmpty())
				// You can also check if the string contains expected output
				Expect(string(bytesOutput)).To(ContainSubstring("Ginkgo"))
			})
		})

		When("the template does not exist", func() {
			It("should return an error", func() {
				bytesOutput, renderErr := engine.RenderBytes("nonexistent.tmpl", nil)
				Expect(renderErr).To(HaveOccurred())
				Expect(bytesOutput).To(BeEmpty())
			})
		})
	})

	Context("RenderString(name, data)", func() {
		When("the template exists", func() {
			It("should render the template to string without error", func() {
				strOutput, renderErr := engine.RenderString(
					"prompt/commit.tmpl",
					map[string]string{
						"Message": "Gomega",
					},
				)
				Expect(renderErr).NotTo(HaveOccurred())
				Expect(strOutput).NotTo(BeEmpty())
				Expect(strOutput).To(ContainSubstring("Gomega"))
			})
		})

		When("the template does not exist", func() {
			It("should return an error", func() {
				strOutput, renderErr := engine.RenderString("nonexistent.tmpl", nil)
				Expect(renderErr).To(HaveOccurred())
				Expect(strOutput).To(BeEmpty())
			})
		})
	})

	Context("Render(name, data)", func() {
		When("the template exists", func() {
			It("should render the template and return an io.Reader", func() {
				reader, renderErr := engine.Render("prompt/commit.tmpl", map[string]string{
					"Message": "ReaderTest",
				})
				Expect(renderErr).NotTo(HaveOccurred())
				Expect(reader).NotTo(BeNil())

				// Verify the contents of the io.Reader
				buf := new(strings.Builder)
				_, copyErr := io.Copy(buf, reader)
				Expect(copyErr).NotTo(HaveOccurred())
				Expect(buf.String()).To(ContainSubstring("ReaderTest"))
			})
		})

		When("the template does not exist", func() {
			It("should return an error and a nil reader", func() {
				reader, renderErr := engine.Render("nonexistent.tmpl", nil)
				Expect(renderErr).To(HaveOccurred())
				Expect(reader).To(BeNil())
			})
		})
	})
})
