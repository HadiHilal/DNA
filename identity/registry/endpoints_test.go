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

package registry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/julienschmidt/httprouter"
	"github.com/skytells-research/DNA/network/node/identity"
	payments_identity "github.com/skytells-research/DNA/network/payments/identity"
	"github.com/skytells-research/DNA/network/payments/registry"
	"github.com/stretchr/testify/assert"
)

const (
	testPublicKeyPart1 = "0xFA001122334455667788990011223344556677889900112233445566778899AF"
	testPublicKeyPart2 = "0xDE001122334455667788990011223344556677889900112233445566778899AD"
)

func TestIdentityRegistrationEndpointReturnsRegistrationData(t *testing.T) {

	mockedDataProvider := &mockRegistrationDataProvider{}
	mockedDataProvider.RegistrationData = &registry.RegistrationData{
		PublicKey: registry.PublicKeyParts{
			Part1: common.FromHex(testPublicKeyPart1),
			Part2: common.FromHex(testPublicKeyPart2),
		},
		Signature: &payments_identity.DecomposedSignature{
			R: [32]byte{1},
			S: [32]byte{2},
			V: 27,
		},
	}

	mockedStatusProvider := &mockRegistrationStatus{
		Registered: false,
	}

	endpoint := newRegistrationEndpoint(mockedDataProvider, mockedStatusProvider)

	req, err := http.NewRequest(
		http.MethodGet,
		"/notimportant",
		nil,
	)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()

	endpoint.IdentityRegistrationData(
		resp,
		req,
		httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: "0x1231323131",
			},
		},
	)

	assert.Equal(t, identity.FromAddress("0x1231323131"), mockedDataProvider.RecordedIdentity)

	assert.JSONEq(
		t,
		`{
			"registered" : false,
            "publicKey": {
				"part1" : "0xfa001122334455667788990011223344556677889900112233445566778899af",
				"part2" : "0xde001122334455667788990011223344556677889900112233445566778899ad"
			},
			"signature": {
				"r": "0x0100000000000000000000000000000000000000000000000000000000000000",
				"s": "0x0200000000000000000000000000000000000000000000000000000000000000",
				"v": 27
			}
        }`,
		resp.Body.String(),
	)

}

type mockRegistrationStatus struct {
	Registered bool
}

func (m *mockRegistrationStatus) IsRegistered(id identity.Identity) (bool, error) {
	return m.Registered, nil
}

func (m *mockRegistrationStatus) SubscribeToRegistrationEvent(id identity.Identity) (
	registeredEvent chan RegistrationEvent,
	unsubscribe func(),
) {
	return nil, nil
}

type mockRegistrationDataProvider struct {
	RegistrationData *registry.RegistrationData
	RecordedIdentity identity.Identity
}

func (m *mockRegistrationDataProvider) ProvideRegistrationData(id identity.Identity) (*registry.RegistrationData, error) {
	m.RecordedIdentity = id
	return m.RegistrationData, nil
}
