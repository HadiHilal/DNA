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

package connection

import (
	"github.comskytells-research/DNA/network/node/core/connection"
	wg "github.comskytells-research/DNA/network/node/services/wireguard"
	"github.comskytells-research/DNA/network/node/services/wireguard/key"
)

// Factory is the wireguard connection factory
type Factory struct{}

// Create creates a new wireguard connection
func (f *Factory) Create(stateChannel connection.StateChannel, statisticsChannel connection.StatisticsChannel) (connection.Connection, error) {
	privateKey, err := key.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	config := wg.ServiceConfig{}
	config.Consumer.PrivateKey = privateKey

	return &Connection{
		stopChannel:       make(chan struct{}),
		stateChannel:      stateChannel,
		statisticsChannel: statisticsChannel,
		config:            config,
	}, nil
}

// NewConnectionCreator creates wireguard connections
func NewConnectionCreator() connection.Factory {
	return &Factory{}
}
