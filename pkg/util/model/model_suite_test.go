package model_test

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model Suite")
}

type ExampleModel struct {
	ID   uuid.UUID
	Name string
}

// Getting the type of single model.
// Both the chain and struct return the same underlying type.
func ExampleReflect_Type() {
	r := model.NewReflect(&ExampleModel{})
	fmt.Println(r.Type().Name())
	rChain := model.NewReflect(&[]*ExampleModel{})
	fmt.Println(rChain.Type().Name())
	// Output:
	// ExampleModel
	// ExampleModel
}

func ExampleReflect_StructValue() {
	r := model.NewReflect(&ExampleModel{Name: "Hello"})
	fmt.Println(r.StructValue())
	// Output:
	// {00000000-0000-0000-0000-000000000000 Hello}
}
