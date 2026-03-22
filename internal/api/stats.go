package api

import (
	"net/http"

	"github.com/go-chi/render"
	mdb "github.com/qxuken/short/internal/db"
)

type LinkStatsResponse struct {
	ShortUrl       string             `json:"short_url"`
	RedirectUrl    string             `json:"redirect_url"`
	TotalClicks    int                `json:"total_clicks"`
	UniqueVisitors int                `json:"unique_visitors"`
	TopCountries   []mdb.CountryStats `json:"top_countries"`
	TopReferers    []mdb.RefererStats `json:"top_referers"`
	DailyClicks    []mdb.DailyStats   `json:"daily_clicks"`
}

type AllLinksStatsResponse struct {
	TotalLinks     int                 `json:"total_links"`
	TotalClicks    int                 `json:"total_clicks"`
	UniqueVisitors int                 `json:"unique_visitors"`
	TopCountries   []mdb.CountryStats  `json:"top_countries"`
	TopReferers    []mdb.RefererStats  `json:"top_referers"`
	DailyClicks    []mdb.DailyStats    `json:"daily_clicks"`
	LinksStats     []LinkStatsResponse `json:"links_stats"`
}

func GetAllStats(mainDb mdb.MainDb, auxDb mdb.AuxiliaryDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingTotals, err := auxDb.GetTrackingTotals()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		links, err := mainDb.GetLinks()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		trafficStats, err := auxDb.GetAllLinksTrafficStats()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		topCountries, err := auxDb.GetAllCountryStats()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		topReferers, err := auxDb.GetAllRefererStats()
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		dailyClicks, err := auxDb.GetAllDailyClicks(30)
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		trafficMap := make(map[string]mdb.LinkTrafficStats, len(trafficStats))
		for _, ts := range trafficStats {
			trafficMap[ts.ShortUrl] = ts
		}

		linksStatsResp := make([]LinkStatsResponse, 0, len(links))
		for _, link := range links {
			ts := trafficMap[link.ShortUrl]
			countryStats, _ := auxDb.GetCountryStats(link.ShortUrl)
			refererStats, _ := auxDb.GetRefererStats(link.ShortUrl)
			daily, _ := auxDb.GetDailyClicks(link.ShortUrl, 30)

			linksStatsResp = append(linksStatsResp, LinkStatsResponse{
				ShortUrl:       link.ShortUrl,
				RedirectUrl:    link.RedirectUrl,
				TotalClicks:    ts.TotalClicks,
				UniqueVisitors: ts.UniqueVisitors,
				TopCountries:   countryStats,
				TopReferers:    refererStats,
				DailyClicks:    daily,
			})
		}

		resp := AllLinksStatsResponse{
			TotalLinks:     len(links),
			TotalClicks:    trackingTotals.TotalClicks,
			UniqueVisitors: trackingTotals.UniqueVisitors,
			TopCountries:   topCountries,
			TopReferers:    topReferers,
			DailyClicks:    dailyClicks,
			LinksStats:     linksStatsResp,
		}

		render.JSON(w, r, resp)
	}
}

func GetLinkStats(shortUrl string, auxDb mdb.AuxiliaryDB, mainDb mdb.MainDb) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingStats, err := auxDb.GetLinkStats(shortUrl)
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		redirectUrl, err := mainDb.GetLink(shortUrl)
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		topCountries, err := auxDb.GetCountryStats(shortUrl)
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		topReferers, err := auxDb.GetRefererStats(shortUrl)
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		dailyClicks, err := auxDb.GetDailyClicks(shortUrl, 30)
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}

		resp := LinkStatsResponse{
			ShortUrl:       shortUrl,
			RedirectUrl:    redirectUrl,
			TotalClicks:    trackingStats.TotalClicks,
			UniqueVisitors: trackingStats.UniqueVisitors,
			TopCountries:   topCountries,
			TopReferers:    topReferers,
			DailyClicks:    dailyClicks,
		}

		render.JSON(w, r, resp)
	}
}
