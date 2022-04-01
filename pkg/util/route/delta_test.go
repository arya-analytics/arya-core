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

func (d *DeltaInletTest) Next() <-chan int {
	return d.data
}

func (d *DeltaInletTest) Errors() <-chan error {
	return d.err
}

func (d *DeltaInletTest) Update(ctx route.DeltaContext[int, int]) {}

type DeltaOutletTest struct {
	errC chan error
	valC chan int
	data []int
	err  []error
}

func (d *DeltaOutletTest) Start() {
	go func() {
		for err := range d.errC {
			d.err = append(d.err, err)
		}
	}()
	go func() {
		for v := range d.valC {
			d.data = append(d.data, v)
		}
	}()
}

func (d *DeltaOutletTest) Context() int {
	return 0
}

func (d *DeltaOutletTest) Send() chan<- int {
	return d.valC
}

func (d *DeltaOutletTest) SendError() chan<- error {
	return d.errC
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
			oOne := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oOne.Start()
			oTwo := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oTwo.Start()
			delta.AddOutlet(oOne)
			delta.AddOutlet(oTwo)
			for i := 0; i < 20; i++ {
				inlet.data <- i
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
			oOne := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oOne.Start()
			oTwo := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oTwo.Start()
			delta.AddOutlet(oOne)
			delta.AddOutlet(oTwo)
			inlet.err <- errors.New("Hello")
			time.Sleep(1 * time.Millisecond)
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
			oOne := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oOne.Start()
			oTwo := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oTwo.Start()
			delta.AddOutlet(oOne)
			delta.AddOutlet(oTwo)
			delta.RemoveOutlet(oOne)
			for i := 0; i < 20; i++ {
				inlet.data <- i
			}
			time.Sleep(1 * time.Millisecond)
			Expect(oOne.data).To(HaveLen(0))
			Expect(oTwo.data).To(HaveLen(20))
		})
	})
	Describe("Delta with data rate", func() {
		It("Should run the delta at the correct rate", func() {
			inlet := &DeltaInletTest{
				data: make(chan int),
				err:  make(chan error),
			}
			delta := route.NewDelta[int, int](inlet, route.WithDataRate(1000))
			go delta.Start()
			oOne := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oOne.Start()
			oTwo := &DeltaOutletTest{errC: make(chan error), valC: make(chan int)}
			oTwo.Start()
			delta.AddOutlet(oOne)
			delta.AddOutlet(oTwo)
			t := time.NewTimer(10 * time.Millisecond)
		o:
			for {
				select {
				case inlet.data <- 1:
				case <-t.C:
					break o
				}
			}
			Expect(oOne.data).To(HaveLen(10))
			Expect(oTwo.data).To(HaveLen(10))
		})
	})
})
