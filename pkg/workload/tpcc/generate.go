// Copyright 2017 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.

package tpcc

import (
	"math/rand"
	"strconv"
)

// These constants are all set by the spec - they're not knobs. Don't change
// them.
const (
	numItems                 = 100000
	numDistrictsPerWarehouse = 10
	numStockPerWarehouse     = 100000
	numCustomersPerDistrict  = 3000
	numCustomersPerWarehouse = numCustomersPerDistrict * numDistrictsPerWarehouse
	numOrdersPerDistrict     = numCustomersPerDistrict
	numOrdersPerWarehouse    = numOrdersPerDistrict * numDistrictsPerWarehouse
	numNewOrdersPerDistrict  = 900
	numNewOrdersPerWarehouse = numNewOrdersPerDistrict * numDistrictsPerWarehouse
	minOrderLinesPerOrder    = 5
	maxOrderLinesPerOrder    = 15

	// TODO(dan): This doesn't match the spec, it should be rand 5-15 per order.
	hackOrderLinesPerOrder     = (minOrderLinesPerOrder + maxOrderLinesPerOrder) / 2
	hackOrderLinesPerDistrict  = hackOrderLinesPerOrder * numOrdersPerDistrict
	hackOrderLinesPerWarehouse = hackOrderLinesPerDistrict * numDistrictsPerWarehouse

	originalString = "ORIGINAL"
	wYtd           = 300000.00
	ytd            = 30000.00
	nextOrderID    = 3001
	creditLimit    = 50000.00
	balance        = -10.00
	ytdPayment     = 10.00
	middleName     = "OE"
	paymentCount   = 1
	deliveryCount  = 0
	goodCredit     = "GC"
	badCredit      = "BC"
)

func (w *tpcc) tpccItemInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	iID := rowIdx + 1

	return []interface{}{
		iID,
		randInt(rng, 1, 10000),                           // im_id: "Image ID associated to Item"
		randAString(rng, 14, 24),                         // name
		float64(randInt(rng, 100, 10000)) / float64(100), // price
		randOriginalString(rng),
	}
}

func (w *tpcc) tpccWarehouseInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	wID := rowIdx // warehouse ids are 0-indexed. every other table is 1-indexed

	return []interface{}{
		wID,
		strconv.Itoa(randInt(rng, 6, 10)),  // name
		strconv.Itoa(randInt(rng, 10, 20)), // street_1
		strconv.Itoa(randInt(rng, 10, 20)), // street_2
		strconv.Itoa(randInt(rng, 10, 20)), // city
		randState(rng),
		randZip(rng),
		randTax(rng),
		wYtd,
	}
}

func (w *tpcc) tpccStockInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	sID := (rowIdx % numStockPerWarehouse) + 1
	wID := (rowIdx / numStockPerWarehouse)

	return []interface{}{
		sID, wID,
		randInt(rng, 10, 100),    // quantity
		randAString(rng, 24, 24), // dist_01
		randAString(rng, 24, 24), // dist_02
		randAString(rng, 24, 24), // dist_03
		randAString(rng, 24, 24), // dist_04
		randAString(rng, 24, 24), // dist_05
		randAString(rng, 24, 24), // dist_06
		randAString(rng, 24, 24), // dist_07
		randAString(rng, 24, 24), // dist_08
		randAString(rng, 24, 24), // dist_09
		randAString(rng, 24, 24), // dist_10
		0, // ytd
		0, // order_cnt
		0, // remote_cnt
		randOriginalString(rng), // data
	}
}

func (w *tpcc) tpccDistrictInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	dID := (rowIdx % numDistrictsPerWarehouse) + 1
	wID := (rowIdx / numDistrictsPerWarehouse)

	return []interface{}{
		dID,
		wID,
		randAString(rng, 6, 10),  // name
		randAString(rng, 10, 20), // street 1
		randAString(rng, 10, 20), // street 2
		randAString(rng, 10, 20), // city
		randState(rng),
		randZip(rng),
		randTax(rng),
		ytd,
		nextOrderID,
	}
}

func (w *tpcc) tpccCustomerInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	cID := (rowIdx % numCustomersPerDistrict) + 1
	dID := ((rowIdx / numCustomersPerDistrict) % numDistrictsPerWarehouse) + 1
	wID := (rowIdx / numCustomersPerWarehouse)

	// 10% of the customer rows have bad credit.
	// See section 4.3, under the CUSTOMER table population section.
	credit := goodCredit
	if rng.Intn(9) == 0 {
		// Poor 10% :(
		credit = badCredit
	}
	var lastName string
	// The first 1000 customers get a last name generated according to their id;
	// the rest get an NURand generated last name.
	if cID <= 1000 {
		lastName = randCLastSyllables(cID - 1)
	} else {
		lastName = randCLast(rng)
	}

	return []interface{}{
		cID, dID, wID,
		randAString(rng, 8, 16), // first name
		middleName,
		lastName,
		randAString(rng, 10, 20), // street 1
		randAString(rng, 10, 20), // street 2
		randAString(rng, 10, 20), // city name
		randState(rng),
		randZip(rng),
		randNString(rng, 16, 16), // phone number
		w.nowString,
		credit,
		creditLimit,
		float64(randInt(rng, 0, 5000)) / float64(10000.0), // discount
		balance,
		ytdPayment,
		paymentCount,
		deliveryCount,
		randAString(rng, 300, 500), // data
	}
}

func (w *tpcc) tpccHistoryInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	cID := (rowIdx % numCustomersPerDistrict) + 1
	dID := ((rowIdx / numCustomersPerDistrict) % numDistrictsPerWarehouse) + 1
	wID := (rowIdx / numCustomersPerWarehouse)

	return []interface{}{
		cID, dID, wID, dID, wID, w.nowString, 10.00, randAString(rng, 12, 24),
	}
}

func (w *tpcc) tpccOrderInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	oID := (rowIdx % numOrdersPerDistrict) + 1
	dID := ((rowIdx / numOrdersPerDistrict) % numDistrictsPerWarehouse) + 1
	wID := (rowIdx / numOrdersPerWarehouse)

	// We need a random permutation of customers that stable for all orders in a
	// district, so use the district ID to seed the random permuation.
	// TODO(dan): Cache this in w.
	var randomCIDs [numCustomersPerDistrict]int
	for i, cID := range rand.New(rand.NewSource(int64(dID))).Perm(numCustomersPerDistrict) {
		randomCIDs[i] = cID + 1
	}
	cID := randomCIDs[oID-1]

	// TODO(dan): This doesn't match the spec, it should be rand 5-15 per order.
	olCnt := hackOrderLinesPerOrder

	var carrierID interface{}
	if oID < 2101 {
		carrierID = strconv.Itoa(randInt(rng, 1, 10))
	}

	return []interface{}{
		oID, dID, wID, cID, w.nowString, carrierID, olCnt, 1,
	}
}

func (w *tpcc) tpccNewOrderInitialRow(rowIdx int) []interface{} {
	// The last numNewOrdersPerDistrict orders have entries in new orders.
	const firstNewOrderOffset = numOrdersPerDistrict - numNewOrdersPerDistrict
	oID := (rowIdx % numNewOrdersPerDistrict) + firstNewOrderOffset + 1
	dID := ((rowIdx / numNewOrdersPerDistrict) % numDistrictsPerWarehouse) + 1
	wID := (rowIdx / numNewOrdersPerWarehouse)

	return []interface{}{
		oID, dID, wID,
	}
}

func (w *tpcc) tpccOrderLineInitialRow(rowIdx int) []interface{} {
	rng := rand.New(rand.NewSource(w.seed + int64(rowIdx)))

	olNumber := (rowIdx % hackOrderLinesPerOrder) + 1
	oID := ((rowIdx / hackOrderLinesPerOrder) % numOrdersPerDistrict) + 1
	dID := ((rowIdx / hackOrderLinesPerDistrict) % numDistrictsPerWarehouse) + 1
	wID := (rowIdx / hackOrderLinesPerWarehouse)

	var amount float64
	var deliveryD interface{}
	if oID < 2101 {
		amount = 0
		deliveryD = w.nowString
	} else {
		amount = float64(randInt(rng, 1, 999999)) / 100.0
	}

	return []interface{}{
		oID, dID, wID, olNumber,
		randInt(rng, 1, 100000), // ol_i_id
		wID, // supply_w_id
		deliveryD,
		5, // quantity
		amount,
		randAString(rng, 24, 24),
	}
}
