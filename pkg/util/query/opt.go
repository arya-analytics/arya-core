package query

type opts map[OptKey]interface{}

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

type OptConverters []OptConverter

type OptConverter func(p *Pack)

func (ocs OptConverters) Exec(p *Pack) {
	for _, oc := range ocs {
		oc(p)
	}
}
