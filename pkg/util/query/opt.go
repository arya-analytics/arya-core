package query

type opts map[OptKey]interface{}

// OptKey is a unique key for a specified option in a query.
// If you're creating a new option, please be careful not to duplicate any of the OptKey already set.
type OptKey string

const (
	pkOptKey          OptKey = "pk"
	whereFieldsOptKey OptKey = "wFld"
	relationOptKey    OptKey = "rel"
	fieldsOptKey      OptKey = "fld"
	calculateOptKey   OptKey = "calc"
	bulkUpdateOptKey  OptKey = "bulkU"
	orderOptKey       OptKey = "order"
	limitOptKey       OptKey = "limit"
	verifyOptKey      OptKey = "verify"
	memoOptKey        OptKey = "memo"
)

// OptConvertChain wraps a slice of OptConvert and provides an Exec function to run them in sequence.
type OptConvertChain []OptConvert

// OptConvert is a simple utility function that allows a package implementing a query runner (Execute) to convert
// the options in a provided query.
type OptConvert func(p *Pack)

// Exec executes all OptConvert in the chain.
func (ocs OptConvertChain) Exec(p *Pack) {
	for _, oc := range ocs {
		oc(p)
	}
}

type optRetrieveOpts struct {
	panicIfNotPresent bool
}

func newOptRetrieveOpts(opts ...OptRetrieveOpt) *optRetrieveOpts {
	ret := &optRetrieveOpts{}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

type OptRetrieveOpt func(o *optRetrieveOpts)

func PanicIfOptNotPresent() OptRetrieveOpt {
	return func(o *optRetrieveOpts) {
		o.panicIfNotPresent = true
	}
}
