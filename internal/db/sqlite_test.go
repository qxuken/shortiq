package db_test

import (
	"slices"
	"testing"

	"github.com/qxuken/short/internal/config"
	mdb "github.com/qxuken/short/internal/db"
)

func createmainDb() *mdb.UrlStore {
	return mdb.ConnectUrlStore(&config.Config{}, ":memory:")
}

func createauxDb() *mdb.TrackingStore {
	return mdb.ConnectTrackingStore(&config.Config{}, ":memory:")
}

func TestSqlite3KVs(t *testing.T) {
	db := createmainDb()
	kvs := []mdb.LinkItem{{"testurl", "testshort"}, {"UpercaseUrl", "UpercaseSHORT"}}

	for _, kv := range kvs {
		err := db.SetLink(kv.RedirectUrl, kv.ShortUrl)
		if err != nil {
			t.Fatal(err)
		}
		v, err := db.GetLink(kv.ShortUrl)
		if err != nil {
			t.Fatal(err)
		}
		if v != kv.RedirectUrl {
			t.Fatalf("Value missmatch for %s, expected %s but found %s", kv.ShortUrl, kv.RedirectUrl, v)
		}
	}

	r, err := db.GetLinks()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(kvs, r) {
		t.Fatalf("Value missmatch, got response (len = %v, val = %v)", len(r), r)
	}
}

func TestSqlite3EmptyKey(t *testing.T) {
	db := createmainDb()
	v, err := db.GetLink("empty")
	if err == nil || v != "" {
		t.Fatal("Found value where it shouldnt be")
	}
}

func TestSqlite3LogVisit(t *testing.T) {
	tv := []mdb.AnalyticsItem{{"s", "c1", "o1", "192.168.0.1", 1}, {"s", "c2", "o2", "192.168.0.2", 2}, {"s2", "c3", "o3", "192.168.0.2", 3}}
	db := createauxDb()
	for _, v := range tv {
		err := db.LogVisit(v)
		if err != nil {
			t.Fatal(err)
		}
	}
	r, err := db.GetLinkAnalytics()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(tv, r) {
		t.Fatalf("Value missmatch, got response (len = %v, val = %v)", len(r), r)
	}
}

func TestTrackingStoreGetTrackingTotalsEmpty(t *testing.T) {
	db := createauxDb()
	totals, err := db.GetTrackingTotals()
	if err != nil {
		t.Fatal(err)
	}
	if totals.TotalClicks != 0 || totals.UniqueVisitors != 0 {
		t.Fatalf("expected zero totals, got %+v", totals)
	}
}

func TestTrackingStoreGetTrackingTotals(t *testing.T) {
	db := createauxDb()
	visits := []mdb.AnalyticsItem{
		{"link1", "c1", "o1", "192.168.0.1", 1},
		{"link1", "c2", "o1", "192.168.0.1", 2},
		{"link1", "c3", "o2", "192.168.0.2", 3},
		{"link2", "c4", "o3", "192.168.0.3", 4},
	}
	for _, v := range visits {
		if err := db.LogVisit(v); err != nil {
			t.Fatal(err)
		}
	}
	totals, err := db.GetTrackingTotals()
	if err != nil {
		t.Fatal(err)
	}
	if totals.TotalClicks != 4 {
		t.Fatalf("expected TotalClicks=4, got %d", totals.TotalClicks)
	}
	if totals.UniqueVisitors != 3 {
		t.Fatalf("expected UniqueVisitors=3, got %d", totals.UniqueVisitors)
	}
}

func TestTrackingStoreGetAllLinksTrafficStatsEmpty(t *testing.T) {
	db := createauxDb()
	stats, err := db.GetAllLinksTrafficStats()
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 0 {
		t.Fatalf("expected empty stats, got %d items", len(stats))
	}
}

func TestTrackingStoreGetAllLinksTrafficStats(t *testing.T) {
	db := createauxDb()
	visits := []mdb.AnalyticsItem{
		{"link1", "c1", "o1", "192.168.0.1", 1},
		{"link1", "c2", "o1", "192.168.0.1", 2},
		{"link1", "c3", "o2", "192.168.0.2", 3},
		{"link2", "c4", "o3", "192.168.0.3", 4},
	}
	for _, v := range visits {
		if err := db.LogVisit(v); err != nil {
			t.Fatal(err)
		}
	}
	stats, err := db.GetAllLinksTrafficStats()
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 links, got %d", len(stats))
	}
	byUrl := make(map[string]mdb.LinkTrafficStats, 2)
	for _, s := range stats {
		byUrl[s.ShortUrl] = s
	}
	link1 := byUrl["link1"]
	if link1.TotalClicks != 3 {
		t.Fatalf("expected link1 TotalClicks=3, got %d", link1.TotalClicks)
	}
	if link1.UniqueVisitors != 2 {
		t.Fatalf("expected link1 UniqueVisitors=2, got %d", link1.UniqueVisitors)
	}
	link2 := byUrl["link2"]
	if link2.TotalClicks != 1 {
		t.Fatalf("expected link2 TotalClicks=1, got %d", link2.TotalClicks)
	}
	if link2.UniqueVisitors != 1 {
		t.Fatalf("expected link2 UniqueVisitors=1, got %d", link2.UniqueVisitors)
	}
}

func TestComposedStatsSeparateStores(t *testing.T) {
	mainDb := createmainDb()
	auxDb := createauxDb()

	if err := mainDb.SetLink("https://example.com", "abc"); err != nil {
		t.Fatal(err)
	}
	if err := mainDb.SetLink("https://test.com", "xyz"); err != nil {
		t.Fatal(err)
	}

	if err := auxDb.LogVisit(mdb.AnalyticsItem{"abc", "c1", "o1", "10.0.0.1", 1}); err != nil {
		t.Fatal(err)
	}
	if err := auxDb.LogVisit(mdb.AnalyticsItem{"abc", "c2", "o2", "10.0.0.2", 2}); err != nil {
		t.Fatal(err)
	}

	links, err := mainDb.GetLinks()
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}

	trafficStats, err := auxDb.GetAllLinksTrafficStats()
	if err != nil {
		t.Fatal(err)
	}
	trafficMap := make(map[string]mdb.LinkTrafficStats, len(trafficStats))
	for _, ts := range trafficStats {
		trafficMap[ts.ShortUrl] = ts
	}

	for _, link := range links {
		ts := trafficMap[link.ShortUrl]
		switch link.ShortUrl {
		case "abc":
			if ts.TotalClicks != 2 {
				t.Fatalf("expected abc TotalClicks=2, got %d", ts.TotalClicks)
			}
		case "xyz":
			if ts.TotalClicks != 0 {
				t.Fatalf("expected xyz TotalClicks=0, got %d", ts.TotalClicks)
			}
		}
	}

	totals, err := auxDb.GetTrackingTotals()
	if err != nil {
		t.Fatal(err)
	}
	if totals.TotalClicks != 2 {
		t.Fatalf("expected TotalClicks=2, got %d", totals.TotalClicks)
	}
	if totals.UniqueVisitors != 2 {
		t.Fatalf("expected UniqueVisitors=2, got %d", totals.UniqueVisitors)
	}
}
