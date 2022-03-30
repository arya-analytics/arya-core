package route_test

import (
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

type DeltaInletTest struct {
	data chan int
	err  chan error
}

func (d *DeltaInletTest) Stream() chan int {
	return d.data
}

func (d *DeltaInletTest) Errors() chan error {
	return d.err
}

func (d *DeltaInletTest) Update(ctx route.DeltaContext[int, int]) {}

type DeltaOutletTest struct {
	data []int
	err  []error
}

func (d *DeltaOutletTest) Context() int {
	return 0
}

func (d *DeltaOutletTest) Send(v int) {
	d.data = append(d.data, v)
}

func (d *DeltaOutletTest) SendError(err error) {
	d.err = append(d.err, err)
}

var _ = Describe("Delta", func() {
	Describe("NewDelta", func() {
		It("should create a new delta", func() {
			delta := route.NewDelta[int, int](
				&DeltaInletTest{
					data: make(chan int),
					err:  make(chan error),
				},
			)
			go delta.Start()
			delta.AddOutlet(&DeltaOutletTest{})
			Expect(delta).ToNot(BeNil())
		})
	})
	Describe("Streaming Data", func() {
		It("Should streamq data correctly", func() {
			inlet := &DeltaInletTest{
				data: make(chan int),
				err:  make(chan error),
			}
			delta := route.NewDelta[int, int](inlet)
			go delta.Start()
			oOne := &DeltaOutletTest{}
			oTwo := &DeltaOutletTest{}
			delta.AddOutlet(oOne)
			delta.AddOutlet(oTwo)
			for i := 0; i < 20; i++ {
				inlet.Stream() <- i
			}
			time.Sleep(1 * time.Millisecond)
			Expect(oOne.data).To(HaveLen(20))
			Expect(oTwo.data).To(HaveLen(20))
		})
	})
	Describe("Streaming errors", func() {
		It("Should streamq errors correctly", func() {
			inlet := &DeltaInletTest{
				data: make(chan int),
				err:  make(chan error),
			}
			delta := route.NewDelta[int, int](inlet)
			go delta.Start()
			oOne := &DeltaOutletTest{}
			oTwo := &DeltaOutletTest{}
			delta.AddOutlet(oOne)
			delta.AddOutlet(oTwo)
			inlet.Errors() <- errors.New("Hello")
			Expect(oOne.err).To(HaveLen(1))
			Expect(oTwo.err).To(HaveLen(1))
		})
	})
	Describe("Removing an Outlet", func() {
		It("Should remove the outlet correctly", func() {
			inlet := &DeltaInletTest{
				data: make(chan int),
				err:  make(chan error),
			}
			delta := route.NewDelta[int, int](inlet)
			go delta.Start()
			oOne := &DeltaOutletTest{}
			oTwo := &DeltaOutletTest{}
			delta.AddOutlet(oOne)
			delta.AddOutlet(oTwo)
			delta.RemoveOutlet(oOne)
			for i := 0; i < 20; i++ {
				inlet.Stream() <- i
			}
			time.Sleep(1 * time.Millisecond)
			Expect(oOne.data).To(HaveLen(0))
			Expect(oTwo.data).To(HaveLen(20))
		})
	})
})
