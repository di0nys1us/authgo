package model_test

import (
	"context"

	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"

	. "github.com/di0nys1us/authgo/model"
	"github.com/di0nys1us/authgo/sqlgo"
)

var _ = Describe("User", func() {

	var (
		db *sqlgo.DB
		tx *sqlgo.Tx
	)

	BeforeSuite(func() {
		var err error
		db, err = sqlgo.NewDB()

		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		var err error
		tx, err = db.Start()

		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err := tx.Rollback()

		Expect(err).To(BeNil())
	})

	AfterSuite(func() {
		err := db.Close()

		Expect(err).To(BeNil())
	})

	Context("single user", func() {
		It("should find user", func() {
			user := &User{ID: "ce274fd4-5803-11e8-8879-afa0dd22785d"}

			err := user.Find(tx)

			Expect(err).To(BeNil())
			Expect(user.FirstName).To(Equal("Erik"))

			user = &User{Email: "erik@eies.land"}

			err = user.Find(tx)

			Expect(err).To(BeNil())
			Expect(user.FirstName).To(Equal("Erik"))
		})

		It("should save user", func() {
			user := &User{
				FirstName: "Test",
				LastName:  "Test",
				Email:     "test@test",
				Password:  "secret",
				Enabled:   true,
			}

			err := user.Save(context.Background(), tx)

			Expect(err).To(BeNil())
			Expect(user.ID).To(BeUUID())
			Expect(user.Events).To(HaveLen(1))

			user = &User{Email: "test@test"}

			err = user.Find(tx)

			Expect(err).To(BeNil())
			Expect(user.FirstName).To(Equal("Test"))
			Expect(user.Events).To(HaveLen(1))

			event := user.Events[0]

			Expect(event.ID).To(BeUUID())
			Expect(event.Type).To(Equal("USER_CREATED"))
			Expect(event.Description).To(Equal("User \"test@test\" created."))
		})

		It("should delete user", func() {
			user := &User{ID: "ce274fd4-5803-11e8-8879-afa0dd22785d"}

			err := user.Delete(tx)

			Expect(err).To(BeNil())
		})
	})

	Context("multiple users", func() {
		It("should find all users", func() {
			var users Users

			err := users.FindAll(tx)

			Expect(err).To(BeNil())
			Expect(users).ToNot(BeEmpty())
		})
	})
})

func BeUUID() GomegaMatcher {
	return MatchRegexp("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
}
