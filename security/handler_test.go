package security_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/di0nys1us/httpgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/di0nys1us/authgo/security"
)

type findSubjectByEmailFunc func(email string) (Subject, error)

func (fn findSubjectByEmailFunc) FindSubjectByEmail(email string) (Subject, error) {
	return fn(email)
}

var _ = Describe("Handler", func() {

	var (
		security = New(findSubjectByEmailFunc(func(email string) (subj Subject, err error) { return }))
	)

	Describe("Authenticate", func() {
		It("should return 401", func() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			httpgo.ErrorHandlerFunc(security.Authenticate).ServeHTTP(w, r)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
