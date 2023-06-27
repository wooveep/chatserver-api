/*
 * @Author: cloudyi.li
 * @Date: 2023-06-15 10:49:51
 * @LastEditTime: 2023-06-15 13:36:08
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/geo.go
 */
package search

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"context"

	"googlemaps.github.io/maps"
)

func GeoSearch(address string) (float64, float64, string) {
	googlecfg := config.AppConfig.GoogelConfig
	c, err := maps.NewClient(maps.WithAPIKey(googlecfg.ApiKey))
	if err != nil {
		logger.Errorf("fatal error: %s", err)
	}
	r := &maps.GeocodingRequest{
		Address: address,
	}
	result, err := c.Geocode(context.Background(), r)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
	}
	if len(result) == 0 {
		return 0, 0, ""
	}
	location := result[0].Geometry.Location
	fmtaddress := result[0].FormattedAddress
	return location.Lat, location.Lng, fmtaddress
}
