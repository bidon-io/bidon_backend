package adminstore_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/resource"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
)

func TestDemandSourceRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceRepo(tx)

	sources := []admin.DemandSourceAttrs{
		{
			HumanName: "Applovin",
			ApiKey:    "applovin",
		},
		{
			HumanName: "Admob",
			ApiKey:    "admob",
		},
	}

	items := make([]admin.DemandSource, len(sources))
	for i, attrs := range sources {
		source, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; items %T, %v", &attrs, nil, err, source, nil)
		}

		items[i] = *source
	}

	want := &resource.Collection[admin.DemandSource]{
		Items: items,
		Meta:  resource.CollectionMeta{TotalCount: int64(len(items))},
	}

	got, err := repo.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestDemandSourceRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceRepo(tx)

	attrs := &admin.DemandSourceAttrs{
		HumanName: "Applovin",
		ApiKey:    "asdf",
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

func TestDemandSourceRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceRepo(tx)

	attrs := admin.DemandSourceAttrs{
		HumanName: "Applovin",
		ApiKey:    "asdf",
	}

	source, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, source, nil)
	}

	want := source
	want.ApiKey = "fdsa"

	updateParams := &admin.DemandSourceAttrs{
		ApiKey: want.ApiKey,
	}
	got, err := repo.Update(context.Background(), source.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestDemandSourceRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceRepo(tx)

	attrs := &admin.DemandSourceAttrs{
		HumanName: "Applovin",
		ApiKey:    "asdf",
	}
	source, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, source, nil)
	}

	err = repo.Delete(context.Background(), source.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", source.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), source.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", source.ID, got, err, nil, "record not found")
	}
}
