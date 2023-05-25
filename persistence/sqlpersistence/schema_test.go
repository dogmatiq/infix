package sqlpersistence_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dogmatiq/sqltest"
	. "github.com/dogmatiq/verity/persistence/sqlpersistence"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("creating and dropping schema", func() {
	for _, pair := range sqltest.CompatiblePairs() {
		pair := pair // capture loop variable

		When(
			fmt.Sprintf(
				"using %s with the '%s' driver",
				pair.Product.Name(),
				pair.Driver.Name(),
			),
			func() {
				var (
					ctx      context.Context
					database *sqltest.Database
					db       *sql.DB
				)

				BeforeEach(func() {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
					DeferCleanup(cancel)

					var err error
					database, err = sqltest.NewDatabase(ctx, pair.Driver, pair.Product)
					Expect(err).ShouldNot(HaveOccurred())
					DeferCleanup(database.Close)

					db, err = database.Open()
					Expect(err).ShouldNot(HaveOccurred())
				})

				Describe("func CreateSchema()", func() {
					It("can be called when the schema already exists", func() {
						err := CreateSchema(ctx, db)
						Expect(err).ShouldNot(HaveOccurred())

						err = CreateSchema(ctx, db)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				Describe("func DropSchema()", func() {
					It("can be called when the schema does not exist", func() {
						err := DropSchema(ctx, db)
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("can be called when the schema has already been dropped", func() {
						err := CreateSchema(ctx, db)
						Expect(err).ShouldNot(HaveOccurred())

						err = DropSchema(ctx, db)
						Expect(err).ShouldNot(HaveOccurred())

						err = DropSchema(ctx, db)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			},
		)
	}
})
