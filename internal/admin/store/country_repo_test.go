package adminstore_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/google/go-cmp/cmp"
)

func TestCountryRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewCountryRepo(tx)

	countries := []admin.CountryAttrs{
		{
			HumanName:  "Japan",
			Alpha2Code: "JP",
			Alpha3Code: "JPN",
		},
		{
			HumanName:  "China",
			Alpha2Code: "CN",
			Alpha3Code: "CHN",
		},
	}

	want := make([]admin.Country, len(countries))
	for i, attrs := range countries {
		country, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, country, nil)
		}

		want[i] = *country
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestCountryRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewCountryRepo(tx)

	attrs := &admin.CountryAttrs{
		HumanName:  "Japan",
		Alpha2Code: "JP",
		Alpha3Code: "JPN",
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestCountryRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewCountryRepo(tx)

	attrs := admin.CountryAttrs{
		HumanName:  "Japan",
		Alpha2Code: "JP",
		Alpha3Code: "JPX",
	}

	country, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, country, nil)
	}

	want := country
	want.Alpha3Code = "JPN"

	updateParams := &admin.CountryAttrs{
		Alpha3Code: want.Alpha3Code,
	}
	got, err := repo.Update(context.Background(), country.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestCountryRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewCountryRepo(tx)

	attrs := &admin.CountryAttrs{
		HumanName:  "Japan",
		Alpha2Code: "JP",
		Alpha3Code: "JPX",
	}
	country, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, country, nil)
	}

	err = repo.Delete(context.Background(), country.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", country.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), country.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", country.ID, got, err, nil, "record not found")
	}
}
