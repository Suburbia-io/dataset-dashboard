package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validators"
)

type BrandCategory struct {
	CategoryID string     `json:"categoryID" db:"category_id"`
	ParentID   *string    `json:"parentID" db:"parent_id"`
	DatasetID  string     `json:"datasetID" db:"dataset_entity"`
	Dataset    Dataset    `json:"dataset"`
	Name       string     `json:"name" db:"name"`
	ArchivedAt *time.Time `json:"archivedAt" db:"archived_at"`
}

func (db *DBAL) BrandCategoryCreate(name string, parentID *string, dataset Dataset) (cat BrandCategory, err error) {
	name, err = validators.Slug(name)
	if err != nil {
		return cat, err
	}

	cat.CategoryID = crypto.NewUUID()
	cat.Name = name
	cat.ParentID = parentID
	cat.DatasetID = dataset.DatasetID
	return cat, db.insert("brand_categories", cat)
}

func (db *DBAL) BrandCategoryGet(categoryID string) (cat BrandCategory, err error) {
	if err := validators.UUIDDeleteMe(categoryID); err != nil {
		return cat, errors.DBNotFound
	}
	return cat, db.getBy(&cat, "brand_categories", map[string]interface{}{
		"category_id": categoryID,
	})
}

func (db *DBAL) BrandCategorySetName(categoryID, name string) (err error) {
	if err := validators.UUIDDeleteMe(categoryID); err != nil {
		return errors.DBNotFound
	}
	if name, err = validators.Slug(name); err != nil {
		return err
	}
	return db.update(
		"brand_categories",
		map[string]interface{}{"name": name},
		map[string]interface{}{"category_id": categoryID})
}

func (db *DBAL) BrandCategoryArchive(categoryID string) (err error) {
	if err := validators.UUIDDeleteMe(categoryID); err != nil {
		return errors.DBNotFound
	}
	return db.update(
		"brand_categories",
		map[string]interface{}{"archived_at": time.Now()},
		map[string]interface{}{"category_id": categoryID})
}

func (db *DBAL) BrandCategoryUnArchive(categoryID string) (err error) {
	if err := validators.UUIDDeleteMe(categoryID); err != nil {
		return errors.DBNotFound
	}
	var t *time.Time
	return db.update(
		"brand_categories",
		map[string]interface{}{"archived_at": t},
		map[string]interface{}{"category_id": categoryID})
}

type BrandCategoryListArgs struct {
	Archived    bool
	NotArchived bool
	DatasetID   string
}

func (db *DBAL) BrandCategoryList(args BrandCategoryListArgs) (
	l []BrandCategory, err error,
) {
	return db.brandCategoryList(args, nil)
}

func (db *DBAL) BrandCategoryListChildren(parentID string, args BrandCategoryListArgs) (
	l []BrandCategory, err error,
) {
	return db.brandCategoryList(args, &parentID)
}

func (db *DBAL) brandCategoryList(args BrandCategoryListArgs, parentID *string) (
	l []BrandCategory, err error,
) {
	ands := []qb.Builder{}
	queryArgs := []interface{}{}

	if parentID != nil {
		if *parentID == "" {
			ands = append(ands, qb.New("(parent_id IS NULL)"))
		} else {
			ands = append(ands, qb.New("(parent_id=$1)", *parentID))
		}
	}

	if args.Archived && !args.NotArchived {
		// archived only
		ands = append(ands, qb.New("(archived_at IS NOT NULL)"))
	} else if !args.Archived && args.NotArchived {
		// not archived only
		ands = append(ands, qb.New("(archived_at IS NULL)"))
	} else if !args.Archived && args.NotArchived {
		return l, nil
	}

	ands = append(ands, qb.New("(dataset_entity=$1)", args.DatasetID))

	sel := qb.New(`SELECT ` +
		` category_id,` +
		` parent_id,` +
		` name,` +
		` archived_at, ` +
		` dataset_entity ` +
		`FROM brand_categories`)
	where := qb.WhereAnd(ands...)
	orderBy := qb.New("ORDER BY name ASC")
	query, queryArgs := qb.Join(" ", sel, where, orderBy).MustBuild()

	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		// TODO: Alert
		fmt.Println(err)
		return nil, errors.Unexpected.WithMsg("Failed to list brand categories: %v", err)
	}

	for rows.Next() {
		row := BrandCategory{}
		if err := rows.Scan(
			&row.CategoryID,
			&row.ParentID,
			&row.Name,
			&row.ArchivedAt,
			&row.DatasetID,
		); err != nil {
			// TODO: Alert
			fmt.Println(err)
			return nil, errors.Unexpected.WithMsg("Failed to list brand categories: %v", err)
		}

		row.Dataset, err = db.DatasetGet(row.DatasetID)
		if err != nil {
			return l, err
		}

		l = append(l, row)
	}

	return l, nil
}

func (db *DBAL) BrandCategoryExport(w io.Writer, datasetID string) (err error) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	andWheres := []qb.Builder{}
	andWheres = append(andWheres, qb.New("(archived_at IS NULL)"))

	if datasetID != "" {
		andWheres = append(andWheres, qb.New("(dataset_entity = $1)", datasetID))
	}

	sel := qb.New(`SELECT ` +
		`category_id, ` +
		` name, ` +
		` parent_id ` +
		`FROM brand_categories `)
	where := qb.WhereAnd(andWheres...)
	orderBy := qb.New("ORDER BY name ASC")
	query, queryArgs := qb.Join(" ", sel, where, orderBy).MustBuild()

	rows, err := db.Query(query, queryArgs...)

	if err != nil {
		return errors.UnexpectedError(err, "Failed listing brand categories")
	}
	defer rows.Close()

	err = writer.Write([]string{
		"category_id",
		"name",
		"parent_id",
	})
	if err != nil {
		return errors.UnexpectedError(err, "Failed writing brandcategory header")
	}

	categoryID := ""
	name := ""
	var parentID *string
	for rows.Next() {
		if err := rows.Scan(
			&categoryID,
			&name,
			&parentID,
		); err != nil {
			return errors.UnexpectedError(err, "Failed scanning brandcategory")
		}

		parentIDStr := ""
		if parentID != nil {
			parentIDStr = *parentID
		}
		err := writer.Write([]string{
			categoryID,
			name,
			parentIDStr,
		})
		if err != nil {
			return errors.UnexpectedError(err, "Failed writing brandcategory row")
		}
	}

	if err := rows.Err(); err != nil {
		return errors.UnexpectedError(err, "Failed iterating brandcategory rows")
	}

	return nil
}
