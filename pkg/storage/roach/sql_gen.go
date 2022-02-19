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

type sqlGen struct {
	db *bun.DB
	m  *model.Reflect
}

func (sg sqlGen) bindTableToField(tableName, fldName string) string {
	return fmt.Sprintf("%s.%s", tableName, fldName)
}

func (sg sqlGen) pkField() reflect.StructField {
	st, ok := sg.m.StructTagChain().RetrieveByFieldRole(model.PKRole)
	if !ok {
		panic("model has no pk field!")
	}
	return st.Field
}

func (sg sqlGen) pk() string {
	return fmt.Sprintf("%s = ?", sg.bindTableToField(sg.table().ModelName, sg.pkField().Name))
}

func (sg sqlGen) pks() string {
	return fmt.Sprintf("%s IN (?)", sg.bindTableToField(sg.table().ModelName, sg.pkField().Name))
}

func (sg sqlGen) fieldName(fldName string) string {
	return caseconv.PascalToSnake(fldName)
}

func (sg sqlGen) fieldNames(fldNames ...string) (sqlNames []string) {
	for _, n := range fldNames {
		sqlNames = append(sqlNames, sg.fieldName(n))
	}
	return sqlNames
}

func (sg sqlGen) table() *schema.Table {
	return sg.db.Table(sg.m.Type())
}

func (sg sqlGen) relFldEquals(fldName string) string {
	relN, baseN := model.SplitLastFieldName(fldName)
	relFldName := sg.bindTableToField(sg.relTableName(relN), sg.fieldName(baseN))
	return fmt.Sprintf("%s = ?", relFldName)

}

const nestedTableSeparator = "__"

func (sg sqlGen) relTableName(relName string) (tableName string) {
	sn := model.SplitFieldNames(relName)
	if len(sn) == 1 && sn[0] == "" {
		return sg.table().ModelName
	}
	for i := range sn {
		nRelName := strings.Join(sn[0:i+1], ".")
		table := sg.db.Table(sg.m.FieldTypeByName(nRelName))
		if i != 0 {
			tableName += nestedTableSeparator
		}
		tableName += table.ModelName
	}
	return tableName
}
