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

package promise

import "github.com/skytells-research/DNA/network/node/money"

// SignedPromise represents payment promise signed by issuer
type SignedPromise struct {
	Promise         Promise
	IssuerSignature Signature
}

// Promise represents payment promise between two parties
type Promise struct {
	SerialNumber int    `storm:"id"`
	IssuerID     string `storm:"index"`
	BenefiterID  string `storm:"index"`
	Amount       money.Money
}

// Signature represents some data signed with a key
type Signature string
