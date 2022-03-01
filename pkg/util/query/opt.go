package query

type opts map[optKey]interface{}

type optKey string

const (
	pkOptKey          optKey = "pk"
	whereFieldsOptKey optKey = "wFld"
	relationOptKey    optKey = "rel"
	fieldsOptKey      optKey = "fld"
	calculateOptKey   optKey = "calc"
	bulkUpdateOptKey  optKey = "bulkU"
)
