/*
 * Copyright (C) 2019 Skytells, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package location

import (
	"github.com/skytells-research/DNA/network/node/core/location/gendb"
	"github.com/oschwald/geoip2-golang"
)

//go:generate go run generator/generator.go --dbname db/GeoLite2-Country.mmdb --output gendb --compress

// NewBuiltInResolver returns new db resolver initialized from built in data
func NewBuiltInResolver() Resolver {
	dbBytes, err := gendb.LoadData()
	if err != nil {
		return NewFailingResolver(err)
	}

	dbReader, err := geoip2.FromBytes(dbBytes)
	if err != nil {
		return NewFailingResolver(err)
	}
	return &DbResolver{dbReader: dbReader}
}
