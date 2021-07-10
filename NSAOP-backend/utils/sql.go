package utils

import (
	"fmt"
	"nsaop/model"
	"strings"

	"gorm.io/gorm"
)

func GenerateStatusClause(tx *gorm.DB, status []string) (*gorm.DB, error) {
	clause := model.DB
	for _, stat := range status {
		if model.ValidStatus(stat) {
			clause = clause.Or("status = ?", stat)
		} else {
			return nil, fmt.Errorf("invalid status")
		}
	}
	return tx.Where(clause), nil
}

func GenerateSearchCLause(tx *gorm.DB, search string, inwhich []string) (*gorm.DB, error) {
	kwargs := strings.Fields(search)
	clause := model.DB
	for _, kwarg := range kwargs {
		for _, in := range inwhich {
			kw := kwarg
			kw = strings.Replace(kw, `\`, `\\`, -1)
			kw = strings.Replace(kw, `%`, `\%`, -1)
			kw = strings.Replace(kw, `_`, `\_`, -1)
			clause = clause.Or(in+" LIKE ?", "%"+kw+"%")
		}
	}
	return tx.Where(clause), nil
}
