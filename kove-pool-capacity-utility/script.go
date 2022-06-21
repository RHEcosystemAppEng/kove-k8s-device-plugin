/*
Copyright 2022 Kove IO, Inc.

This program source is the property of Kove IO, Inc. and contains
information which is confidential and proprietary to Kove IO, Inc.
No part of this source may be copied, reproduced, disclosed to third
parties, or transmitted in any form or by any means, electronic or
mechanical for any purpose without the express written consent of
Kove IO, Inc.
*/

package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	rangeMin = 512 * 1024  // 512 GiB in MiB
	rangeMax = 1024 * 1024 // 1 TiB in MiB
	rangeLen = rangeMax - rangeMin
	mib      = 1024 * 1024 // 1 MiB
)

func main() {
	// Get a random integer in the range [0, rangeLen) from the system
	// non-blocking random source
	randBig, err := rand.Int(rand.Reader, big.NewInt(rangeLen))
	if err != nil {
		// Fall back to a static value at the midpoint of the range
		randBig = big.NewInt(rangeLen >> 1)
	}
	// Shift the random integer to the defined range and convert MiB to bytes
	randPoolCapacity := (rangeMin + randBig.Int64()) * mib
	// Output the pool capacity in bytes
	fmt.Println(randPoolCapacity)
}
