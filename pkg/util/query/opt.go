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
	optRequired bool
}

func newOptRetrieveOpts(opts ...OptRetrieveOpt) *optRetrieveOpts {
	ret := &optRetrieveOpts{}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

type OptRetrieveOpt func(o *optRetrieveOpts)

// RequireOpt is passed as an option to a RetrieveOpt function that requires the option to be present.
// Panics if the option is not present.
//
// Example:
// 		pkc, _ := query.RetrievePKOpt(p, query.RequireOpt())
//
//  The function will panic if the pk option was not set on the query.
func RequireOpt() OptRetrieveOpt {
	return func(o *optRetrieveOpts) {
		o.optRequired = true
	}
}
