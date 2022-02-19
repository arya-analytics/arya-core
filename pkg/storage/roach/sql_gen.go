package roach

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
	"reflect"
)

type SQLGen struct {
	db *bun.DB
	m  *model.Reflect
}

func (sg SQLGen) bindModelNameToCol(column string) string {
	return fmt.Sprintf("%s.%s", sg.table().ModelName, column)
}

func (sg SQLGen) pkField() reflect.StructField {
	st, ok := sg.m.StructTagChain().RetrieveByFieldRole(model.PKRole)
	if !ok {
		panic("model has no pk field!")
	}
	return st.Field
}

func (sg SQLGen) pk() string {
	return fmt.Sprintf("%s = ?", sg.bindModelNameToCol(sg.pkField().Name))
}

func (sg SQLGen) pks() string {
	return fmt.Sprintf("%s IN (?)", sg.bindModelNameToCol(sg.pkField().Name))
}

func (sg SQLGen) fieldNameToSQL(fldName string) string {
	return caseconv.PascalToSnake(fldName)
}

func (sg SQLGen) fieldNamesToSQL(fldNames ...string) (sqlNames []string) {
	for _, n := range fldNames {
		sqlNames = append(sqlNames, sg.fieldNameToSQL(n))
	}
	return sqlNames
}

func (sg SQLGen) fieldEquals(fldName string) string {
	return fmt.Sprintf("%s = ?", caseconv.PascalToSnake(fldName))
}

func (sg SQLGen) table() *schema.Table {
	return sg.db.Table(sg.m.Type())
}
