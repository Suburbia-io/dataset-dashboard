package database

import (
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type CorpMappingRule = struct {
	tables.CorpMappingRule
	Corporation Corporation
}

type CorpMappingRulesListArgs struct {
	CorpMappingID string `json:"corpMappingID"`
}

func (db *DBAL) CorpMappingRuleList(args CorpMappingRulesListArgs) (corpMappingRuleList []CorpMappingRule, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.CorpMappingRules.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.CorpMappingRules.Table())
	builder.Write(`WHERE TRUE`)

	if args.CorpMappingID != "" {
		builder.Write(`AND (corp_mapping_id = $1)`, args.CorpMappingID)
	}

	//	builder.Write(`GROUP BY country`)
	//	builder.Write(`ORDER BY from_date DESC`)

	query, queryArgs := builder.MustBuild()
	corpMappingRules, err := tables.CorpMappingRules.List(db, query, queryArgs...)

	for _, corpMappingRule := range corpMappingRules {
		corp, err := tables.Corporations.Get(db, corpMappingRule.CorpID)
		if err != nil {
			return corpMappingRuleList, err
		}

		extendedCorpMappingRule := CorpMappingRule{
			CorpMappingRule: corpMappingRule,
			Corporation:     corp,
		}

		corpMappingRuleList = append(corpMappingRuleList, extendedCorpMappingRule)
	}
	return corpMappingRuleList, nil
}

func (db *DBAL) CorpMappingRuleUpsert(cmr *tables.CorpMappingRule) error {
	return tables.CorpMappingRules.Upsert(db, cmr)
}

func (db *DBAL) CorpMappingRuleGet(corpMappingRuleID string) (rule CorpMappingRule, err error) {
	cmr, err := tables.CorpMappingRules.Get(db, corpMappingRuleID)
	if err != nil {
		return rule, err
	}
	rule.CorpMappingRule = cmr

	corp, err := tables.Corporations.Get(db, cmr.CorpID)
	if err != nil {
		return rule, err
	}

	rule.Corporation = corp
	return rule, err
}

func (db *DBAL) CorpMappingRuleDelete(corpMappingRuleID string) (err error) {
	err = tables.CorpMappingRules.Delete(db, corpMappingRuleID)
	return err
}
