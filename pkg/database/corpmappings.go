package database

import (
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type CorpMapping = struct {
	tables.CorpMapping
	CorpType tables.CorporationType
	Tag      tables.Tag
	TagType  tables.TagType
	Rules    []CorpMappingRule `json:"rules"`
}

func (db *DBAL) CorpMappingList() (corpMappingList []CorpMapping, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.CorpMappings.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.CorpMappings.Table())
	builder.Write(`WHERE TRUE`)

	query, queryArgs := builder.MustBuild()
	corpMappings, err := tables.CorpMappings.List(db, query, queryArgs...)
	if err != nil {
		return corpMappingList, err
	}
	for _, mapping := range corpMappings {
		var extendedCorpMapping = CorpMapping{
			CorpMapping: mapping,
		}
		rules, err := db.CorpMappingRuleList(CorpMappingRulesListArgs{CorpMappingID: mapping.CorpMappingID})
		if err != nil {
			return corpMappingList, err
		}
		extendedCorpMapping.Rules = rules

		corpType, err := db.CorporationTypeGet(mapping.CorpTypeID)
		if err != nil {
			return corpMappingList, err
		}
		extendedCorpMapping.CorpType = corpType

		tagType, err := db.TagTypeGet(corpType.DatasetID, mapping.TagTypeID)
		if err != nil {
			return corpMappingList, err
		}
		extendedCorpMapping.TagType = tagType

		tag, err := db.TagGet(corpType.DatasetID, mapping.TagTypeID, mapping.TagID)
		if err != nil {
			return corpMappingList, err
		}
		extendedCorpMapping.Tag = tag

		corpMappingList = append(corpMappingList, extendedCorpMapping)
	}

	return corpMappingList, nil
}

func (db *DBAL) CorpMappingUpsert(cm *tables.CorpMapping) error {
	return tables.CorpMappings.Upsert(db, cm)
}

func (db *DBAL) CorpMappingGet(corpMappingID string) (cm CorpMapping, err error) {
	mapping, err := tables.CorpMappings.Get(db, corpMappingID)
	if err != nil {
		return cm, err
	}
	cm.CorpMapping = mapping

	rules, err := db.CorpMappingRuleList(CorpMappingRulesListArgs{CorpMappingID: mapping.CorpMappingID})
	if err != nil {
		return cm, err
	}
	cm.Rules = rules

	corpType, err := db.CorporationTypeGet(mapping.CorpTypeID)
	if err != nil {
		return cm, err
	}
	cm.CorpType = corpType

	tagType, err := db.TagTypeGet(corpType.DatasetID, mapping.TagTypeID)
	if err != nil {
		return cm, err
	}
	cm.TagType = tagType

	tag, err := db.TagGet(corpType.DatasetID, mapping.TagTypeID, mapping.TagID)
	if err != nil {
		return cm, err
	}
	cm.Tag = tag

	return cm, err
}
