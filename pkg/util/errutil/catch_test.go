package errutil_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
)

var _ = Describe("Catch", func() {
	Describe("CatchSimple", func() {
		Context("No error encountered", func() {
			var (
				counter int
				catcher *errutil.CatchSimple
			)
			BeforeEach(func() {
				counter = 1
				catcher = errutil.NewCatchSimple()
				for i := 0; i < 4; i++ {
					catcher.Exec(func() error {
						counter++
						return nil
					})
				}
			})
			It("Should continue to execute functions", func() {

				Expect(counter).To(Equal(5))
			})
			It("Should contain a nil error", func() {
				Expect(catcher.Error()).To(BeNil())
			})
		})
		Context("Error encountered", func() {
			var (
				counter int
				catcher *errutil.CatchSimple
			)
			BeforeEach(func() {
				counter = 1
				catcher = errutil.NewCatchSimple()
				for i := 0; i < 4; i++ {
					catcher.Exec(func() error {
						if i == 2 {
							return fmt.Errorf("encountered unknown error")
						}
						counter++
						return nil
					})
				}
			})
			It("Should stop execution", func() {
				Expect(counter).To(Equal(3))
			})
			It("Should contain a non-nil error", func() {
				Expect(catcher.Error()).ToNot(BeNil())
			})
			Describe("Reset", func() {
				It("Should reset the catcher", func() {
					catcher.Reset()
					Expect(catcher.Error()).To(BeNil())
				})
			})

		})
		Context("Aggregation", func() {
			var catcher = errutil.NewCatchSimple(errutil.WithAggregation())
			It("Should aggregate the errors", func() {
				counter := 1
				for i := 0; i < 4; i++ {
					catcher.Exec(func() error {
						counter++
						return fmt.Errorf("error encountered")
					})
				}
				Expect(counter).To(Equal(5))
				Expect(catcher.Errors()).To(HaveLen(4))
			})
		})
	})
	Describe("CatchContext", func() {
		var (
			ctx     = context.Background()
			counter int
			catcher *errutil.CatchContext
		)
		BeforeEach(func() {
			counter = 1
			catcher = errutil.NewCatchContext(ctx)
			for i := 0; i < 4; i++ {
				catcher.Exec(func(ctx context.Context) error {
					if i == 2 {
						return fmt.Errorf("encountered unknown error")
					}
					counter++
					return nil
				})
			}
		})
		It("Should stop execution", func() {
			Expect(counter).To(Equal(3))
		})
		It("Should contain a non-nil error", func() {
			Expect(catcher.Error()).ToNot(BeNil())
		})
		Describe("Reset", func() {
			It("Should reset the catcher", func() {
				catcher.Reset()
				Expect(catcher.Error()).To(BeNil())
			})
		})
	})
	Describe("Pipe Hook", func() {
		It("Should pipe errors", func() {
			pipe := make(chan error, 10)
			c := errutil.NewCatchSimple(errutil.WithHooks(errutil.NewPipeHook(pipe)))
			c.Exec(func() error {
				return errors.New("hello")
			})
			var errs []error
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				for err := range pipe {
					errs = append(errs, err)
				}
				wg.Done()
			}()
			close(pipe)
			wg.Wait()
			Expect(errs).To(HaveLen(1))
		})
	})
	Describe("With Converter", func() {
		It("Should convert the error", func() {
			cc := errutil.ConvertChain{func(err error) (error, bool) {
				if err.Error() == "not random error" {
					return errors.New("random error"), true
				}
				return nil, false
			}}
			c := errutil.NewCatchSimple(errutil.WithConvert(cc))
			c.Exec(func() error {
				return errors.New("not random error")
			})
			Expect(c.Error()).To(Equal(errors.New("random error")))
		})
	})
	Describe("AddError", func() {
		It("Should be able to take a function with arbitrary return values", func() {
			c := errutil.NewCatchSimple()
			c.AddError(func() (int, error) { return 1, nil }())
			Expect(c.Error()).To(BeNil())
		})

		It("Should do nothing when no arguments are passed", func() {
			c := errutil.NewCatchSimple()
			c.AddError()
			Expect(c.Error()).To(BeNil())
		})
		It("Should bin a non nil error", func() {
			c := errutil.NewCatchSimple()
			c.AddError(errors.New("error"))
			Expect(c.Error()).ToNot(BeNil())
		})
		Describe("Edge cases + errors", func() {
			It("Should panic when the function passed isn't called", func() {
				c := errutil.NewCatchSimple()
				Expect(func() {
					c.AddError(func() (int, error) { return 1, nil })
				}).To(PanicWith("function must be called when running AddError!"))
			})
			It("Should panic when a non-error is returned as the last value", func() {
				c := errutil.NewCatchSimple()
				Expect(func() {
					c.AddError("a", "b")
				}).To(PanicWith("catch function didn't return an error as its last value!"))
			})
		})
	})
})
