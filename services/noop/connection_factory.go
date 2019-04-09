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

package noop

import (
	"github.comskytells-research/DNA/network/node/core/connection"
)

// ConnectionFactory is the noop connection factory
type ConnectionFactory struct{}

// Create creates a new noop connnection
func (cf *ConnectionFactory) Create(stateChannel connection.StateChannel, statisticsChannel connection.StatisticsChannel) (connection.Connection, error) {
	return &Connection{
		stateChannel: stateChannel,
	}, nil
}

// NewConnectionCreator creates noop connections
func NewConnectionCreator() connection.Factory {
	return &ConnectionFactory{}
}
