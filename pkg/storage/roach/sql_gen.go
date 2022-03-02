package roach

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
	"reflect"
	"strings"
)

type sqlGen struct {
	db *bun.DB
	m  *model.Reflect
}

// |||| PRIMARY KEY ||||

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

// |||| FIELDS ||||

func (sg sqlGen) fieldName(fldName string) string {
	return caseconv.PascalToSnake(fldName)
}

func (sg sqlGen) fieldNames(fldNames ...string) (sqlNames []string) {
	for _, n := range fldNames {
		sqlNames = append(sqlNames, sg.fieldName(n))
	}
	return sqlNames
}

func (sg sqlGen) relFldName(fldName string) string {
	relN, baseN := model.SplitLastFieldName(fldName)
	return sg.bindTableToField(sg.relTableName(relN), sg.fieldName(baseN))
}

func (sg sqlGen) relFldExp(fldName string, fldVal interface{}) (string, []interface{}) {
	return sg.parseFldExp(sg.relFldName(fldName), fldVal)
}

func (sg sqlGen) parseFldExp(fldName string, fldVal interface{}) (string, []interface{}) {
	exp, ok := fldVal.(query.FieldExp)
	if !ok {
		return fmt.Sprintf("%s = ?", fldName), []interface{}{fldVal}
	}
	switch exp.Op {
	case query.FieldOpInRange:
		return fmt.Sprintf("%s BETWEEN ? and ?", fldName), exp.Vals
	case query.FieldOpLessThan:
		return fmt.Sprintf("%s < ?", fldName), exp.Vals
	case query.FieldOpGreaterThan:
		return fmt.Sprintf("%s > (?)", fldName), exp.Vals
	case query.FieldOpIn:
		return fmt.Sprintf("%s IN (?)", fldName), []interface{}{bun.In(exp.Vals)}
	default:
		log.Warnf("roach sql gen could not parse expression opt %s. attempting equality", exp.Op)
		return fmt.Sprintf("%s = ?", fldName), []interface{}{exp}
	}

}

// |||| TABLE ||||

func (sg sqlGen) bindTableToField(tableName, fldName string) string {
	return fmt.Sprintf("%s.%s", tableName, fldName)
}

func (sg sqlGen) table() *schema.Table {
	return sg.db.Table(sg.m.Type())
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

// |||| CALCULATIONS ||||

func (sg sqlGen) calc(c query.Calc) string {
	calcSQL := map[query.Calc]string{
		query.CalcSum:   "SUM",
		query.CalcAVG:   "AVG",
		query.CalcCount: "COUNT",
		query.CalcMax:   "MAX",
		query.CalcMin:   "MIN",
	}
	return fmt.Sprintf("%s(?)", calcSQL[c])

}

// |||| ORDER ||||

const (
	orderSQLASC = "ASC"
	orderSQLDSC = "DESC"
)

func (sg sqlGen) order(o query.Order, fld string) string {
	var orderSQL string
	if o == query.OrderASC {
		orderSQL = orderSQLASC
	} else {
		orderSQL = orderSQLDSC
	}
	return fmt.Sprintf("%s %s", sg.relFldName(fld), orderSQL)
}
