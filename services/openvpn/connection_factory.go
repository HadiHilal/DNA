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

package openvpn

import (
	"time"

	"github.com/skytells-research/DNA/network/go-openvpn/openvpn"
	"github.com/skytells-research/DNA/network/go-openvpn/openvpn/management"
	"github.com/skytells-research/DNA/network/go-openvpn/openvpn/middlewares/client/auth"
	openvpn_bytescount "github.com/skytells-research/DNA/network/go-openvpn/openvpn/middlewares/client/bytescount"
	"github.com/skytells-research/DNA/network/go-openvpn/openvpn/middlewares/state"
	"github.com/skytells-research/DNA/network/node/core/connection"
	"github.com/skytells-research/DNA/network/node/core/ip"
	"github.com/skytells-research/DNA/network/node/core/location"
	"github.com/skytells-research/DNA/network/node/identity"
	"github.com/skytells-research/DNA/network/node/services/openvpn/middlewares/client/bytescount"
	openvpn_session "github.com/skytells-research/DNA/network/node/services/openvpn/session"
	"github.com/skytells-research/DNA/network/node/session"
)

// ProcessBasedConnectionFactory represents a factory for creating process-based openvpn connections
type ProcessBasedConnectionFactory struct {
	openvpnBinary         string
	configDirectory       string
	runtimeDirectory      string
	originalLocationCache location.Cache
	signerFactory         identity.SignerFactory
	ipResolver            ip.Resolver
	natPinger             NATPinger
}

// NewProcessBasedConnectionFactory creates a new ProcessBasedConnectionFactory
func NewProcessBasedConnectionFactory(
	openvpnBinary, configDirectory, runtimeDirectory string,
	originalLocationCache location.Cache,
	signerFactory identity.SignerFactory,
	resolver ip.Resolver,
	natPinger NATPinger,
) *ProcessBasedConnectionFactory {
	return &ProcessBasedConnectionFactory{
		openvpnBinary:         openvpnBinary,
		configDirectory:       configDirectory,
		runtimeDirectory:      runtimeDirectory,
		originalLocationCache: originalLocationCache,
		signerFactory:         signerFactory,
		ipResolver:            resolver,
		natPinger:             natPinger,
	}
}

func (op *ProcessBasedConnectionFactory) newAuthMiddleware(sessionID session.ID, signer identity.Signer) management.Middleware {
	credentialsProvider := openvpn_session.SignatureCredentialsProvider(sessionID, signer)
	return auth.NewMiddleware(credentialsProvider)
}

func (op *ProcessBasedConnectionFactory) newBytecountMiddleware(statisticsChannel connection.StatisticsChannel) management.Middleware {
	statsSaver := bytescount.NewSessionStatsSaver(statisticsChannel)
	return openvpn_bytescount.NewMiddleware(statsSaver, 1*time.Second)
}

func (op *ProcessBasedConnectionFactory) newStateMiddleware(session session.ID, signer identity.Signer, connectionOptions connection.ConnectOptions, stateChannel connection.StateChannel) management.Middleware {
	stateCallback := GetStateCallback(stateChannel)
	return state.NewMiddleware(stateCallback)
}

// Create creates a new openvpn connection
func (op *ProcessBasedConnectionFactory) Create(stateChannel connection.StateChannel, statisticsChannel connection.StatisticsChannel) (connection.Connection, error) {
	procFactory := func(options connection.ConnectOptions) (openvpn.Process, *ClientConfig, error) {
		vpnClientConfig, err := NewClientConfigFromSession(options.SessionConfig, op.configDirectory, op.runtimeDirectory)
		if err != nil {
			return nil, nil, err
		}

		signer := op.signerFactory(options.ConsumerID)

		stateMiddleware := op.newStateMiddleware(options.SessionID, signer, options, stateChannel)
		authMiddleware := op.newAuthMiddleware(options.SessionID, signer)
		byteCountMiddleware := op.newBytecountMiddleware(statisticsChannel)
		proc := openvpn.CreateNewProcess(op.openvpnBinary, vpnClientConfig.GenericConfig, stateMiddleware, byteCountMiddleware, authMiddleware)
		return proc, vpnClientConfig, nil
	}

	return &Client{
		processFactory: procFactory,
		ipResolver:     op.ipResolver,
		natPinger:      op.natPinger,
	}, nil
}
