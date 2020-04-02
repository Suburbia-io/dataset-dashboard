package tables

import (
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
)

func NewCustomerForTesting(db DBi) Customer {
	return Customer{
		CustomerID: crypto.NewUUID(),
		Name:       crypto.RandAlphaNum(16),
		CreatedAt:  time.Now().Round(time.Second),
	}
}

func NewUserForTesting(db DBi) User {
	customer := NewCustomerForTesting(db)
	if err := Customers.Insert(db, &customer); err != nil {
		panic(err)
	}
	k := crypto.RandAlphaNum(32)
	email := randEmail()
	return User{
		CustomerID:   customer.CustomerID,
		UserID:       crypto.NewUUID(),
		Name:         crypto.RandAlphaNum(16),
		Email:        email,
		APIKey:       &k,
		SFTPUsername: k,
		CreatedAt:    time.Now().Round(time.Second),
		LoginToken:   &k,
	}
}

func NewDatasetForTesting(db DBi) Dataset {
	return Dataset{
		DatasetID:  crypto.NewUUID(),
		Name:       crypto.RandAlphaNum(16),
		Slug:       randSlug(),
		Manageable: true,
		CreatedAt:  time.Now().Round(time.Second),
	}
}

func NewLocationForTesting(db DBi) Location {
	ds := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &ds); err != nil {
		panic(err)
	}
	return Location{
		DatasetID:           ds.DatasetID,
		LocationHash:        randSlug() + randSlug(),
		LocationString:      "amsterdam, nl, , ,",
		GeonamesHierarchy:   []byte("{}"),
		GeonamesPostalCodes: []byte("{}"),
		CreatedAt:           time.Now().Round(time.Second),
	}
}

func NewTagTypeForTesting(db DBi) TagType {
	ds := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &ds); err != nil {
		panic(err)
	}
	return TagType{
		DatasetID:   ds.DatasetID,
		TagTypeID:   crypto.NewUUID(),
		TagType:     randSlug(),
		Description: crypto.RandAlphaNum(16),
	}
}

func NewTagForTesting(db DBi) Tag {
	tt := NewTagTypeForTesting(db)
	if err := TagTypes.Insert(db, &tt); err != nil {
		panic(err)
	}
	return Tag{
		DatasetID:   tt.DatasetID,
		TagTypeID:   tt.TagTypeID,
		TagID:       crypto.NewUUID(),
		Tag:         randSlug(),
		Description: crypto.RandAlphaNum(16),
	}
}

func NewFingerprintForTesting(db DBi) Fingerprint {
	ds := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &ds); err != nil {
		panic(err)
	}
	return Fingerprint{
		DatasetID:   ds.DatasetID,
		Fingerprint: crypto.NewUUID(),
		RawText:     crypto.RandAlphaNum(16),
		Annotations: crypto.RandAlphaNum(16),
		Count:       32,
	}
}

func NewTagAppForTesting(db DBi) TagApp {
	return TagApp{
		Name:   crypto.RandAlphaNum(16),
		Weight: 0.5,
	}
}

func NewTagAppTagForTesting(db DBi) TagAppTag {
	tag := NewTagForTesting(db)
	if err := Tags.Insert(db, &tag); err != nil {
		panic(err)
	}
	fp := Fingerprint{
		DatasetID:   tag.DatasetID,
		Fingerprint: crypto.NewUUID(),
		RawText:     crypto.RandAlphaNum(16),
		Annotations: crypto.RandAlphaNum(16),
		Count:       32,
	}
	if err := Fingerprints.Insert(db, &fp); err != nil {
		panic(err)
	}

	app := NewTagAppForTesting(db)
	if err := TagApps.Insert(db, &app); err != nil {
		panic(err)
	}

	user := NewUserForTesting(db)
	if err := Users.Insert(db, &user); err != nil {
		panic(err)
	}

	return TagAppTag{
		DatasetID:   tag.DatasetID,
		Fingerprint: fp.Fingerprint,
		TagTypeID:   tag.TagTypeID,
		TagAppID:    app.TagAppID,
		TagID:       tag.TagID,
		Confidence:  0.5,
		UpdatedAt:   time.Now().Round(time.Second),
		UserID:      user.UserID,
	}
}

func NewTagAppHistoricalTagForTesting(db DBi) TagAppHistoricalTag {
	t := NewTagAppTagForTesting(db)
	return TagAppHistoricalTag{
		DatasetID:   t.DatasetID,
		Fingerprint: t.Fingerprint,
		TagTypeID:   t.TagTypeID,
		TagAppID:    t.TagAppID,
		TagID:       t.TagID,
		Confidence:  t.Confidence,
		UpdatedAt:   time.Now().Round(time.Second),
		UserID:      t.UserID,
	}
}

func NewConsensusTagForTesting(db DBi) ConsensusTag {
	t := NewTagAppTagForTesting(db)
	return ConsensusTag{
		DatasetID:   t.DatasetID,
		Fingerprint: t.Fingerprint,
		TagTypeID:   t.TagTypeID,
		TagID:       t.TagID,
		Confidence:  t.Confidence,
		UpdatedAt:   time.Now().Round(time.Second),
	}
}

func NewCustomerDatasetForTesting(db DBi) CustomerDataset {
	customer := NewCustomerForTesting(db)
	if err := Customers.Insert(db, &customer); err != nil {
		panic(err)
	}
	dataset := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &dataset); err != nil {
		panic(err)
	}
	return CustomerDataset{
		CustomerDatasetID: crypto.NewUUID(),
		DatasetEntity:     dataset.DatasetID,
		CustomerEntity:    customer.CustomerID,
		CreatedAt:         time.Now().Round(time.Second),
	}
}

func NewAuditTrailForTesting(db DBi) AuditTrail {
	return AuditTrail{
		AuditTrailID: crypto.NewUUID(),
		Type:         "insert",
		RelatedTable: "admins",
		RelatedID:    crypto.NewUUID(),
		Payload:      "{}",
		CreatedAt:    time.Now().Round(time.Second),
	}
}

func NewSessionForTesting(db DBi) Session {
	user := NewUserForTesting(db)
	if err := Users.Insert(db, &user); err != nil {
		panic(err)
	}
	return Session{
		Token:     crypto.RandAlphaNum(32),
		UserID:    user.UserID,
		ExpiresAt: time.Now().Add(time.Hour).Round(time.Second),
	}
}

func NewCorporationTypeForTesting(db DBi) CorporationType {
	ds := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &ds); err != nil {
		panic(err)
	}

	return CorporationType{
		CorporationTypeID: crypto.NewUUID(),
		DatasetID:         ds.DatasetID,
		CorporationType:   randSlug(),
		Description:       randSlug(),
		CreatedAt:         time.Now().Round(time.Second),
	}
}

func NewCorporationForTesting(db DBi) Corporation {
	ds := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &ds); err != nil {
		panic(err)
	}
	return Corporation{
		Exchange:  "NYSE",
		Code:      crypto.NewUUID(),
		Name:      randSlug(),
		Slug:      randSlug(),
		Isin:      randSlug(),
		Cusip:     randSlug(),
		DatasetID: ds.DatasetID,
		CreatedAt: time.Now().Round(time.Second),
		UpdatedAt: time.Now().Round(time.Second),
	}
}

func NewCorpMappingForTesting(db DBi) CorpMapping {
	ct := NewCorporationTypeForTesting(db)
	if err := CorporationTypes.Insert(db, &ct); err != nil {
		panic(err)
	}
	tt := NewTagTypeForTesting(db)
	if err := TagTypes.Insert(db, &tt); err != nil {
		panic(err)
	}
	t := NewTagForTesting(db)
	if err := Tags.Insert(db, &t); err != nil {
		panic(err)
	}
	return CorpMapping{
		CorpMappingID: crypto.NewUUID(),
		CorpTypeID:    ct.CorporationTypeID,
		TagTypeID:     tt.TagTypeID,
		TagID:         t.TagID,
	}
}

func NewCorpMappingRuleForTesting(db DBi) CorpMappingRule {
	cm := NewCorpMappingForTesting(db)
	if err := CorpMappings.Insert(db, &cm); err != nil {
		panic(err)
	}
	c := NewCorporationForTesting(db)
	if err := Corporations.Insert(db, &c); err != nil {
		panic(err)
	}
	return CorpMappingRule{
		CorpMappingRuleID: crypto.NewUUID(),
		CorpMappingID:     cm.CorpMappingID,
		CorpID:            c.CorporationID,
		ExternalNotes:     "xyz",
		InternalNotes:     "abc",
		FromDate:          time.Now().Round(time.Second),
		Country:           "NL",
	}
}
