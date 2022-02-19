package roach

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
	"reflect"
	"strings"
)

type SQLGen struct {
	db *bun.DB
	m  *model.Reflect
}

func (sg SQLGen) bindTableToField(tableName, fldName string) string {
	return fmt.Sprintf("%s.%s", tableName, fldName)
}

func (sg SQLGen) pkField() reflect.StructField {
	st, ok := sg.m.StructTagChain().RetrieveByFieldRole(model.PKRole)
	if !ok {
		panic("model has no pk field!")
	}
	return st.Field
}

func (sg SQLGen) pk() string {
	return fmt.Sprintf("%s = ?", sg.bindTableToField(sg.table().ModelName, sg.pkField().Name))
}

func (sg SQLGen) pks() string {
	return fmt.Sprintf("%s IN (?)", sg.bindTableToField(sg.table().ModelName, sg.pkField().Name))
}

func (sg SQLGen) fieldName(fldName string) string {
	return caseconv.PascalToSnake(fldName)
}

func (sg SQLGen) fieldNames(fldNames ...string) (sqlNames []string) {
	for _, n := range fldNames {
		sqlNames = append(sqlNames, sg.fieldName(n))
	}
	return sqlNames
}

func (sg SQLGen) fieldEquals(fldName string) string {
	return fmt.Sprintf("%s = ?", fldName)
}

func (sg SQLGen) table() *schema.Table {
	return sg.db.Table(sg.m.Type())
}

func (sg SQLGen) relFldEquals(relName, fldName string) string {
	return sg.fieldEquals(sg.bindTableToField(sg.relTableName(relName), sg.fieldName(fldName)))

}

func (sg SQLGen) relTableName(relName string) (tableName string) {
	sn := model.SplitFieldNames(relName)
	for i := range sn {
		nRelName := strings.Join(sn[0:i+1], ".")
		table := sg.db.Table(sg.m.FieldTypeByName(nRelName))
		if i != 0 {
			tableName += "__"
		}
		tableName += table.ModelName
	}
	return tableName
}
