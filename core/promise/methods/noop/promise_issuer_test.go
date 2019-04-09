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
	"errors"
	"testing"

	"github.comskytells-research/DNA/network/node/core/connection"
	"github.comskytells-research/DNA/network/node/core/promise"
	"github.comskytells-research/DNA/network/node/identity"
	"github.comskytells-research/DNA/network/node/logconfig"
	"github.comskytells-research/DNA/network/node/market"
	"github.comskytells-research/DNA/network/node/money"
	"github.com/stretchr/testify/assert"
)

var (
	providerID = identity.FromAddress("provider-id")
	proposal   = market.ServiceProposal{
		ProviderID:    providerID.Address,
		PaymentMethod: fakePaymentMethod{},
	}
)

var _ connection.PromiseIssuer = &PromiseIssuer{}

func TestPromiseIssuer_Start_SubscriptionFails(t *testing.T) {
	dialog := &fakeDialog{
		returnError: errors.New("reject subscriptions"),
	}

	logs := make([]string, 0)
	logger := logconfig.ReplaceLogger(logconfig.NewLoggerCapture(&logs))
	defer logconfig.ReplaceLogger(logger)

	issuer := &PromiseIssuer{dialog: dialog, signer: &identity.SignerFake{}}
	err := issuer.Start(proposal)
	defer issuer.Stop()

	assert.EqualError(t, err, "reject subscriptions")
	assert.Len(t, logs, 0)
}

func TestPromiseIssuer_Start_SubscriptionOfBalances(t *testing.T) {
	dialog := &fakeDialog{
		returnReceiveMessage: promise.BalanceMessage{1, true, testToken(10)},
	}

	logs := make([]string, 0)
	logger := logconfig.ReplaceLogger(logconfig.NewLoggerCapture(&logs))
	defer logconfig.ReplaceLogger(logger)

	issuer := &PromiseIssuer{dialog: dialog, signer: &identity.SignerFake{}}
	err := issuer.Start(proposal)
	assert.NoError(t, err)

	assert.Len(t, logs, 1)
	assert.Equal(t, "[promise-issuer] Promise balance notified: 1000000000TEST", logs[0])
}

func testToken(amount float64) money.Money {
	return money.NewMoney(amount, money.Currency("TEST"))
}

type fakePaymentMethod struct{}

func (fpm fakePaymentMethod) GetPrice() money.Money {
	return money.NewMoney(1111111111, money.Currency("FAKE"))
}
