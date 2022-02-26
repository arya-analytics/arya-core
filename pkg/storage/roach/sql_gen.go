package roach

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/caseconv"
	"github.com/arya-analytics/aryacore/pkg/util/model"
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

func (sg sqlGen) relFldExp(fldName string, fldVal interface{}) (string, []interface{}) {
	relN, baseN := model.SplitLastFieldName(fldName)
	relFldName := sg.bindTableToField(sg.relTableName(relN), sg.fieldName(baseN))
	return sg.parseFldExp(relFldName, fldVal)
}

func (sg sqlGen) parseFldExp(fldName string, fldVal interface{}) (string, []interface{}) {
	exp, ok := fldVal.(model.FieldExp)
	if !ok {
		return fmt.Sprintf("%s = ?", fldName), []interface{}{fldVal}
	}
	switch exp.Op {
	case model.FieldExpOpInRange:
		return fmt.Sprintf("%s BETWEEN ? and ?", fldName), exp.Vals
	case model.FieldExpOpLessThan:
		return fmt.Sprintf("%s < ?", fldName), exp.Vals
	case model.FieldExpOpGreaterThan:
		return fmt.Sprintf("%s > ?", fldName), exp.Vals
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

func (sg sqlGen) calculate(c storage.Calculate) string {
	calcSQL := map[storage.Calculate]string{
		storage.CalculateSum:   "SUM",
		storage.CalculateAVG:   "AVG",
		storage.CalculateCount: "COUNT",
		storage.CalculateMax:   "MAX",
		storage.CalculateMin:   "MIN",
	}
	return fmt.Sprintf("%s(?)", calcSQL[c])

}
