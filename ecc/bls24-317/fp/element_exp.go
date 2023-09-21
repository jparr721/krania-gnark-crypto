// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by jparr721/krania-gnark-crypto DO NOT EDIT

package fp

// expBySqrtExp is equivalent to z.Exp(x, 41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b95976acaab)
//
// uses github.com/mmcloughlin/addchain v0.4.0 to generate a shorter addition chain
func (z *Element) expBySqrtExp(x Element) *Element {
	// addition chain:
	//
	//	_10      = 2*1
	//	_11      = 1 + _10
	//	_101     = _10 + _11
	//	_111     = _10 + _101
	//	_1001    = _10 + _111
	//	_1011    = _10 + _1001
	//	_1101    = _10 + _1011
	//	_1111    = _10 + _1101
	//	_10000   = 1 + _1111
	//	_10001   = 1 + _10000
	//	_10011   = _10 + _10001
	//	_10101   = _10 + _10011
	//	_10111   = _10 + _10101
	//	_11001   = _10 + _10111
	//	_11101   = _1101 + _10000
	//	_111010  = 2*_11101
	//	_111111  = _101 + _111010
	//	_1000000 = 1 + _111111
	//	i37      = ((_1000000 << 3 + _1011) << 8 + _11001) << 6
	//	i53      = ((_10001 + i37) << 8 + _10011) << 5 + _1111
	//	i74      = ((i53 << 3 + _11) << 10 + _10001) << 6
	//	i86      = ((_1001 + i74) << 6 + _11001) << 3 + _111
	//	i109     = ((i86 << 5 + _101) << 9 + _111111) << 7
	//	i123     = ((_1011 + i109) << 2 + 1) << 9 + _1011
	//	i141     = ((i123 << 7 + _111111) << 2 + _11) << 7
	//	i160     = ((_11101 + i141) << 9 + _111) << 7 + _10001
	//	i177     = ((i160 << 5 + _1101) << 6 + _1101) << 4
	//	i194     = ((_11 + i177) << 8 + _1111) << 6 + _1101
	//	i217     = ((i194 << 9 + _10011) << 6 + _111) << 6
	//	i232     = ((_1101 + i217) << 8 + _10001) << 4 + _1011
	//	i251     = ((i232 << 5 + _1011) << 5 + _1111) << 7
	//	i262     = ((_11001 + i251) << 5 + _11001) << 3 + _111
	//	i283     = ((i262 << 3 + 1) << 8 + _1101) << 8
	//	i296     = ((_1001 + i283) << 3 + 1) << 7 + _111111
	//	i316     = ((i296 << 7 + _111111) << 6 + _10001) << 5
	//	i332     = ((_10101 + i316) << 6 + _10001) << 7 + _10111
	//	i351     = ((i332 << 7 + _10101) << 4 + _1001) << 6
	//	i365     = ((_11101 + i351) << 5 + _10101) << 6 + _11001
	//	return     (i365 << 6 + _10101) << 5 + _1011
	//
	// Operations: 310 squares 68 multiplies

	// Allocate Temporaries.
	var (
		t0  = new(Element)
		t1  = new(Element)
		t2  = new(Element)
		t3  = new(Element)
		t4  = new(Element)
		t5  = new(Element)
		t6  = new(Element)
		t7  = new(Element)
		t8  = new(Element)
		t9  = new(Element)
		t10 = new(Element)
		t11 = new(Element)
		t12 = new(Element)
		t13 = new(Element)
	)

	// var t0,t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11,t12,t13 Element
	// Step 1: t1 = x^0x2
	t1.Square(&x)

	// Step 2: t11 = x^0x3
	t11.Mul(&x, t1)

	// Step 3: t12 = x^0x5
	t12.Mul(t1, t11)

	// Step 4: t8 = x^0x7
	t8.Mul(t1, t12)

	// Step 5: t3 = x^0x9
	t3.Mul(t1, t8)

	// Step 6: z = x^0xb
	z.Mul(t1, t3)

	// Step 7: t7 = x^0xd
	t7.Mul(t1, z)

	// Step 8: t9 = x^0xf
	t9.Mul(t1, t7)

	// Step 9: t2 = x^0x10
	t2.Mul(&x, t9)

	// Step 10: t5 = x^0x11
	t5.Mul(&x, t2)

	// Step 11: t10 = x^0x13
	t10.Mul(t1, t5)

	// Step 12: t0 = x^0x15
	t0.Mul(t1, t10)

	// Step 13: t4 = x^0x17
	t4.Mul(t1, t0)

	// Step 14: t1 = x^0x19
	t1.Mul(t1, t4)

	// Step 15: t2 = x^0x1d
	t2.Mul(t7, t2)

	// Step 16: t6 = x^0x3a
	t6.Square(t2)

	// Step 17: t6 = x^0x3f
	t6.Mul(t12, t6)

	// Step 18: t13 = x^0x40
	t13.Mul(&x, t6)

	// Step 21: t13 = x^0x200
	for s := 0; s < 3; s++ {
		t13.Square(t13)
	}

	// Step 22: t13 = x^0x20b
	t13.Mul(z, t13)

	// Step 30: t13 = x^0x20b00
	for s := 0; s < 8; s++ {
		t13.Square(t13)
	}

	// Step 31: t13 = x^0x20b19
	t13.Mul(t1, t13)

	// Step 37: t13 = x^0x82c640
	for s := 0; s < 6; s++ {
		t13.Square(t13)
	}

	// Step 38: t13 = x^0x82c651
	t13.Mul(t5, t13)

	// Step 46: t13 = x^0x82c65100
	for s := 0; s < 8; s++ {
		t13.Square(t13)
	}

	// Step 47: t13 = x^0x82c65113
	t13.Mul(t10, t13)

	// Step 52: t13 = x^0x1058ca2260
	for s := 0; s < 5; s++ {
		t13.Square(t13)
	}

	// Step 53: t13 = x^0x1058ca226f
	t13.Mul(t9, t13)

	// Step 56: t13 = x^0x82c6511378
	for s := 0; s < 3; s++ {
		t13.Square(t13)
	}

	// Step 57: t13 = x^0x82c651137b
	t13.Mul(t11, t13)

	// Step 67: t13 = x^0x20b19444dec00
	for s := 0; s < 10; s++ {
		t13.Square(t13)
	}

	// Step 68: t13 = x^0x20b19444dec11
	t13.Mul(t5, t13)

	// Step 74: t13 = x^0x82c651137b0440
	for s := 0; s < 6; s++ {
		t13.Square(t13)
	}

	// Step 75: t13 = x^0x82c651137b0449
	t13.Mul(t3, t13)

	// Step 81: t13 = x^0x20b19444dec11240
	for s := 0; s < 6; s++ {
		t13.Square(t13)
	}

	// Step 82: t13 = x^0x20b19444dec11259
	t13.Mul(t1, t13)

	// Step 85: t13 = x^0x1058ca226f60892c8
	for s := 0; s < 3; s++ {
		t13.Square(t13)
	}

	// Step 86: t13 = x^0x1058ca226f60892cf
	t13.Mul(t8, t13)

	// Step 91: t13 = x^0x20b19444dec11259e0
	for s := 0; s < 5; s++ {
		t13.Square(t13)
	}

	// Step 92: t12 = x^0x20b19444dec11259e5
	t12.Mul(t12, t13)

	// Step 101: t12 = x^0x41632889bd8224b3ca00
	for s := 0; s < 9; s++ {
		t12.Square(t12)
	}

	// Step 102: t12 = x^0x41632889bd8224b3ca3f
	t12.Mul(t6, t12)

	// Step 109: t12 = x^0x20b19444dec11259e51f80
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 110: t12 = x^0x20b19444dec11259e51f8b
	t12.Mul(z, t12)

	// Step 112: t12 = x^0x82c651137b044967947e2c
	for s := 0; s < 2; s++ {
		t12.Square(t12)
	}

	// Step 113: t12 = x^0x82c651137b044967947e2d
	t12.Mul(&x, t12)

	// Step 122: t12 = x^0x1058ca226f60892cf28fc5a00
	for s := 0; s < 9; s++ {
		t12.Square(t12)
	}

	// Step 123: t12 = x^0x1058ca226f60892cf28fc5a0b
	t12.Mul(z, t12)

	// Step 130: t12 = x^0x82c651137b044967947e2d0580
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 131: t12 = x^0x82c651137b044967947e2d05bf
	t12.Mul(t6, t12)

	// Step 133: t12 = x^0x20b19444dec11259e51f8b416fc
	for s := 0; s < 2; s++ {
		t12.Square(t12)
	}

	// Step 134: t12 = x^0x20b19444dec11259e51f8b416ff
	t12.Mul(t11, t12)

	// Step 141: t12 = x^0x1058ca226f60892cf28fc5a0b7f80
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 142: t12 = x^0x1058ca226f60892cf28fc5a0b7f9d
	t12.Mul(t2, t12)

	// Step 151: t12 = x^0x20b19444dec11259e51f8b416ff3a00
	for s := 0; s < 9; s++ {
		t12.Square(t12)
	}

	// Step 152: t12 = x^0x20b19444dec11259e51f8b416ff3a07
	t12.Mul(t8, t12)

	// Step 159: t12 = x^0x1058ca226f60892cf28fc5a0b7f9d0380
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 160: t12 = x^0x1058ca226f60892cf28fc5a0b7f9d0391
	t12.Mul(t5, t12)

	// Step 165: t12 = x^0x20b19444dec11259e51f8b416ff3a07220
	for s := 0; s < 5; s++ {
		t12.Square(t12)
	}

	// Step 166: t12 = x^0x20b19444dec11259e51f8b416ff3a0722d
	t12.Mul(t7, t12)

	// Step 172: t12 = x^0x82c651137b044967947e2d05bfce81c8b40
	for s := 0; s < 6; s++ {
		t12.Square(t12)
	}

	// Step 173: t12 = x^0x82c651137b044967947e2d05bfce81c8b4d
	t12.Mul(t7, t12)

	// Step 177: t12 = x^0x82c651137b044967947e2d05bfce81c8b4d0
	for s := 0; s < 4; s++ {
		t12.Square(t12)
	}

	// Step 178: t11 = x^0x82c651137b044967947e2d05bfce81c8b4d3
	t11.Mul(t11, t12)

	// Step 186: t11 = x^0x82c651137b044967947e2d05bfce81c8b4d300
	for s := 0; s < 8; s++ {
		t11.Square(t11)
	}

	// Step 187: t11 = x^0x82c651137b044967947e2d05bfce81c8b4d30f
	t11.Mul(t9, t11)

	// Step 193: t11 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3c0
	for s := 0; s < 6; s++ {
		t11.Square(t11)
	}

	// Step 194: t11 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd
	t11.Mul(t7, t11)

	// Step 203: t11 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a00
	for s := 0; s < 9; s++ {
		t11.Square(t11)
	}

	// Step 204: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a13
	t10.Mul(t10, t11)

	// Step 210: t10 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c0
	for s := 0; s < 6; s++ {
		t10.Square(t10)
	}

	// Step 211: t10 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c7
	t10.Mul(t8, t10)

	// Step 217: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131c0
	for s := 0; s < 6; s++ {
		t10.Square(t10)
	}

	// Step 218: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd
	t10.Mul(t7, t10)

	// Step 226: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd00
	for s := 0; s < 8; s++ {
		t10.Square(t10)
	}

	// Step 227: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11
	t10.Mul(t5, t10)

	// Step 231: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd110
	for s := 0; s < 4; s++ {
		t10.Square(t10)
	}

	// Step 232: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b
	t10.Mul(z, t10)

	// Step 237: t10 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a2360
	for s := 0; s < 5; s++ {
		t10.Square(t10)
	}

	// Step 238: t10 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b
	t10.Mul(z, t10)

	// Step 243: t10 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d60
	for s := 0; s < 5; s++ {
		t10.Square(t10)
	}

	// Step 244: t9 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f
	t9.Mul(t9, t10)

	// Step 251: t9 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b780
	for s := 0; s < 7; s++ {
		t9.Square(t9)
	}

	// Step 252: t9 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799
	t9.Mul(t1, t9)

	// Step 257: t9 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f320
	for s := 0; s < 5; s++ {
		t9.Square(t9)
	}

	// Step 258: t9 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339
	t9.Mul(t1, t9)

	// Step 261: t9 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799c8
	for s := 0; s < 3; s++ {
		t9.Square(t9)
	}

	// Step 262: t8 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf
	t8.Mul(t8, t9)

	// Step 265: t8 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce78
	for s := 0; s < 3; s++ {
		t8.Square(t8)
	}

	// Step 266: t8 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce79
	t8.Mul(&x, t8)

	// Step 274: t8 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce7900
	for s := 0; s < 8; s++ {
		t8.Square(t8)
	}

	// Step 275: t7 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d
	t7.Mul(t7, t8)

	// Step 283: t7 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d00
	for s := 0; s < 8; s++ {
		t7.Square(t7)
	}

	// Step 284: t7 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d09
	t7.Mul(t3, t7)

	// Step 287: t7 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c86848
	for s := 0; s < 3; s++ {
		t7.Square(t7)
	}

	// Step 288: t7 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c86849
	t7.Mul(&x, t7)

	// Step 295: t7 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e4342480
	for s := 0; s < 7; s++ {
		t7.Square(t7)
	}

	// Step 296: t7 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf
	t7.Mul(t6, t7)

	// Step 303: t7 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125f80
	for s := 0; s < 7; s++ {
		t7.Square(t7)
	}

	// Step 304: t6 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf
	t6.Mul(t6, t7)

	// Step 310: t6 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efc0
	for s := 0; s < 6; s++ {
		t6.Square(t6)
	}

	// Step 311: t6 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1
	t6.Mul(t5, t6)

	// Step 316: t6 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa20
	for s := 0; s < 5; s++ {
		t6.Square(t6)
	}

	// Step 317: t6 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa35
	t6.Mul(t0, t6)

	// Step 323: t6 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d40
	for s := 0; s < 6; s++ {
		t6.Square(t6)
	}

	// Step 324: t5 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d51
	t5.Mul(t5, t6)

	// Step 331: t5 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a880
	for s := 0; s < 7; s++ {
		t5.Square(t5)
	}

	// Step 332: t4 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a897
	t4.Mul(t4, t5)

	// Step 339: t4 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b80
	for s := 0; s < 7; s++ {
		t4.Square(t4)
	}

	// Step 340: t4 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b95
	t4.Mul(t0, t4)

	// Step 344: t4 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b950
	for s := 0; s < 4; s++ {
		t4.Square(t4)
	}

	// Step 345: t3 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b959
	t3.Mul(t3, t4)

	// Step 351: t3 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d512e5640
	for s := 0; s < 6; s++ {
		t3.Square(t3)
	}

	// Step 352: t2 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d512e565d
	t2.Mul(t2, t3)

	// Step 357: t2 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacba0
	for s := 0; s < 5; s++ {
		t2.Square(t2)
	}

	// Step 358: t2 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacbb5
	t2.Mul(t0, t2)

	// Step 364: t2 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a8972b2ed40
	for s := 0; s < 6; s++ {
		t2.Square(t2)
	}

	// Step 365: t1 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a8972b2ed59
	t1.Mul(t1, t2)

	// Step 371: t1 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacbb5640
	for s := 0; s < 6; s++ {
		t1.Square(t1)
	}

	// Step 372: t0 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacbb5655
	t0.Mul(t0, t1)

	// Step 377: t0 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b95976acaa0
	for s := 0; s < 5; s++ {
		t0.Square(t0)
	}

	// Step 378: z = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b95976acaab
	z.Mul(z, t0)

	return z
}

// expByLegendreExp is equivalent to z.Exp(x, 82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a8972b2ed59555)
//
// uses github.com/mmcloughlin/addchain v0.4.0 to generate a shorter addition chain
func (z *Element) expByLegendreExp(x Element) *Element {
	// addition chain:
	//
	//	_10      = 2*1
	//	_11      = 1 + _10
	//	_101     = _10 + _11
	//	_111     = _10 + _101
	//	_1001    = _10 + _111
	//	_1011    = _10 + _1001
	//	_1101    = _10 + _1011
	//	_1111    = _10 + _1101
	//	_10000   = 1 + _1111
	//	_10001   = 1 + _10000
	//	_10011   = _10 + _10001
	//	_10101   = _10 + _10011
	//	_10111   = _10 + _10101
	//	_11001   = _10 + _10111
	//	_11101   = _1101 + _10000
	//	_111010  = 2*_11101
	//	_111111  = _101 + _111010
	//	_1000000 = 1 + _111111
	//	i37      = ((_1000000 << 3 + _1011) << 8 + _11001) << 6
	//	i53      = ((_10001 + i37) << 8 + _10011) << 5 + _1111
	//	i74      = ((i53 << 3 + _11) << 10 + _10001) << 6
	//	i86      = ((_1001 + i74) << 6 + _11001) << 3 + _111
	//	i109     = ((i86 << 5 + _101) << 9 + _111111) << 7
	//	i123     = ((_1011 + i109) << 2 + 1) << 9 + _1011
	//	i141     = ((i123 << 7 + _111111) << 2 + _11) << 7
	//	i160     = ((_11101 + i141) << 9 + _111) << 7 + _10001
	//	i177     = ((i160 << 5 + _1101) << 6 + _1101) << 4
	//	i194     = ((_11 + i177) << 8 + _1111) << 6 + _1101
	//	i217     = ((i194 << 9 + _10011) << 6 + _111) << 6
	//	i232     = ((_1101 + i217) << 8 + _10001) << 4 + _1011
	//	i251     = ((i232 << 5 + _1011) << 5 + _1111) << 7
	//	i262     = ((_11001 + i251) << 5 + _11001) << 3 + _111
	//	i283     = ((i262 << 3 + 1) << 8 + _1101) << 8
	//	i296     = ((_1001 + i283) << 3 + 1) << 7 + _111111
	//	i316     = ((i296 << 7 + _111111) << 6 + _10001) << 5
	//	i332     = ((_10101 + i316) << 6 + _10001) << 7 + _10111
	//	i351     = ((i332 << 7 + _10101) << 4 + _1001) << 6
	//	i365     = ((_11101 + i351) << 5 + _10101) << 6 + _11001
	//	return     (i365 << 6 + _10101) << 6 + _10101
	//
	// Operations: 311 squares 68 multiplies

	// Allocate Temporaries.
	var (
		t0  = new(Element)
		t1  = new(Element)
		t2  = new(Element)
		t3  = new(Element)
		t4  = new(Element)
		t5  = new(Element)
		t6  = new(Element)
		t7  = new(Element)
		t8  = new(Element)
		t9  = new(Element)
		t10 = new(Element)
		t11 = new(Element)
		t12 = new(Element)
		t13 = new(Element)
	)

	// var t0,t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11,t12,t13 Element
	// Step 1: t0 = x^0x2
	t0.Square(&x)

	// Step 2: t11 = x^0x3
	t11.Mul(&x, t0)

	// Step 3: t12 = x^0x5
	t12.Mul(t0, t11)

	// Step 4: t7 = x^0x7
	t7.Mul(t0, t12)

	// Step 5: t2 = x^0x9
	t2.Mul(t0, t7)

	// Step 6: t9 = x^0xb
	t9.Mul(t0, t2)

	// Step 7: t6 = x^0xd
	t6.Mul(t0, t9)

	// Step 8: t8 = x^0xf
	t8.Mul(t0, t6)

	// Step 9: t1 = x^0x10
	t1.Mul(&x, t8)

	// Step 10: t4 = x^0x11
	t4.Mul(&x, t1)

	// Step 11: t10 = x^0x13
	t10.Mul(t0, t4)

	// Step 12: z = x^0x15
	z.Mul(t0, t10)

	// Step 13: t3 = x^0x17
	t3.Mul(t0, z)

	// Step 14: t0 = x^0x19
	t0.Mul(t0, t3)

	// Step 15: t1 = x^0x1d
	t1.Mul(t6, t1)

	// Step 16: t5 = x^0x3a
	t5.Square(t1)

	// Step 17: t5 = x^0x3f
	t5.Mul(t12, t5)

	// Step 18: t13 = x^0x40
	t13.Mul(&x, t5)

	// Step 21: t13 = x^0x200
	for s := 0; s < 3; s++ {
		t13.Square(t13)
	}

	// Step 22: t13 = x^0x20b
	t13.Mul(t9, t13)

	// Step 30: t13 = x^0x20b00
	for s := 0; s < 8; s++ {
		t13.Square(t13)
	}

	// Step 31: t13 = x^0x20b19
	t13.Mul(t0, t13)

	// Step 37: t13 = x^0x82c640
	for s := 0; s < 6; s++ {
		t13.Square(t13)
	}

	// Step 38: t13 = x^0x82c651
	t13.Mul(t4, t13)

	// Step 46: t13 = x^0x82c65100
	for s := 0; s < 8; s++ {
		t13.Square(t13)
	}

	// Step 47: t13 = x^0x82c65113
	t13.Mul(t10, t13)

	// Step 52: t13 = x^0x1058ca2260
	for s := 0; s < 5; s++ {
		t13.Square(t13)
	}

	// Step 53: t13 = x^0x1058ca226f
	t13.Mul(t8, t13)

	// Step 56: t13 = x^0x82c6511378
	for s := 0; s < 3; s++ {
		t13.Square(t13)
	}

	// Step 57: t13 = x^0x82c651137b
	t13.Mul(t11, t13)

	// Step 67: t13 = x^0x20b19444dec00
	for s := 0; s < 10; s++ {
		t13.Square(t13)
	}

	// Step 68: t13 = x^0x20b19444dec11
	t13.Mul(t4, t13)

	// Step 74: t13 = x^0x82c651137b0440
	for s := 0; s < 6; s++ {
		t13.Square(t13)
	}

	// Step 75: t13 = x^0x82c651137b0449
	t13.Mul(t2, t13)

	// Step 81: t13 = x^0x20b19444dec11240
	for s := 0; s < 6; s++ {
		t13.Square(t13)
	}

	// Step 82: t13 = x^0x20b19444dec11259
	t13.Mul(t0, t13)

	// Step 85: t13 = x^0x1058ca226f60892c8
	for s := 0; s < 3; s++ {
		t13.Square(t13)
	}

	// Step 86: t13 = x^0x1058ca226f60892cf
	t13.Mul(t7, t13)

	// Step 91: t13 = x^0x20b19444dec11259e0
	for s := 0; s < 5; s++ {
		t13.Square(t13)
	}

	// Step 92: t12 = x^0x20b19444dec11259e5
	t12.Mul(t12, t13)

	// Step 101: t12 = x^0x41632889bd8224b3ca00
	for s := 0; s < 9; s++ {
		t12.Square(t12)
	}

	// Step 102: t12 = x^0x41632889bd8224b3ca3f
	t12.Mul(t5, t12)

	// Step 109: t12 = x^0x20b19444dec11259e51f80
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 110: t12 = x^0x20b19444dec11259e51f8b
	t12.Mul(t9, t12)

	// Step 112: t12 = x^0x82c651137b044967947e2c
	for s := 0; s < 2; s++ {
		t12.Square(t12)
	}

	// Step 113: t12 = x^0x82c651137b044967947e2d
	t12.Mul(&x, t12)

	// Step 122: t12 = x^0x1058ca226f60892cf28fc5a00
	for s := 0; s < 9; s++ {
		t12.Square(t12)
	}

	// Step 123: t12 = x^0x1058ca226f60892cf28fc5a0b
	t12.Mul(t9, t12)

	// Step 130: t12 = x^0x82c651137b044967947e2d0580
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 131: t12 = x^0x82c651137b044967947e2d05bf
	t12.Mul(t5, t12)

	// Step 133: t12 = x^0x20b19444dec11259e51f8b416fc
	for s := 0; s < 2; s++ {
		t12.Square(t12)
	}

	// Step 134: t12 = x^0x20b19444dec11259e51f8b416ff
	t12.Mul(t11, t12)

	// Step 141: t12 = x^0x1058ca226f60892cf28fc5a0b7f80
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 142: t12 = x^0x1058ca226f60892cf28fc5a0b7f9d
	t12.Mul(t1, t12)

	// Step 151: t12 = x^0x20b19444dec11259e51f8b416ff3a00
	for s := 0; s < 9; s++ {
		t12.Square(t12)
	}

	// Step 152: t12 = x^0x20b19444dec11259e51f8b416ff3a07
	t12.Mul(t7, t12)

	// Step 159: t12 = x^0x1058ca226f60892cf28fc5a0b7f9d0380
	for s := 0; s < 7; s++ {
		t12.Square(t12)
	}

	// Step 160: t12 = x^0x1058ca226f60892cf28fc5a0b7f9d0391
	t12.Mul(t4, t12)

	// Step 165: t12 = x^0x20b19444dec11259e51f8b416ff3a07220
	for s := 0; s < 5; s++ {
		t12.Square(t12)
	}

	// Step 166: t12 = x^0x20b19444dec11259e51f8b416ff3a0722d
	t12.Mul(t6, t12)

	// Step 172: t12 = x^0x82c651137b044967947e2d05bfce81c8b40
	for s := 0; s < 6; s++ {
		t12.Square(t12)
	}

	// Step 173: t12 = x^0x82c651137b044967947e2d05bfce81c8b4d
	t12.Mul(t6, t12)

	// Step 177: t12 = x^0x82c651137b044967947e2d05bfce81c8b4d0
	for s := 0; s < 4; s++ {
		t12.Square(t12)
	}

	// Step 178: t11 = x^0x82c651137b044967947e2d05bfce81c8b4d3
	t11.Mul(t11, t12)

	// Step 186: t11 = x^0x82c651137b044967947e2d05bfce81c8b4d300
	for s := 0; s < 8; s++ {
		t11.Square(t11)
	}

	// Step 187: t11 = x^0x82c651137b044967947e2d05bfce81c8b4d30f
	t11.Mul(t8, t11)

	// Step 193: t11 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3c0
	for s := 0; s < 6; s++ {
		t11.Square(t11)
	}

	// Step 194: t11 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd
	t11.Mul(t6, t11)

	// Step 203: t11 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a00
	for s := 0; s < 9; s++ {
		t11.Square(t11)
	}

	// Step 204: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a13
	t10.Mul(t10, t11)

	// Step 210: t10 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c0
	for s := 0; s < 6; s++ {
		t10.Square(t10)
	}

	// Step 211: t10 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c7
	t10.Mul(t7, t10)

	// Step 217: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131c0
	for s := 0; s < 6; s++ {
		t10.Square(t10)
	}

	// Step 218: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd
	t10.Mul(t6, t10)

	// Step 226: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd00
	for s := 0; s < 8; s++ {
		t10.Square(t10)
	}

	// Step 227: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11
	t10.Mul(t4, t10)

	// Step 231: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd110
	for s := 0; s < 4; s++ {
		t10.Square(t10)
	}

	// Step 232: t10 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b
	t10.Mul(t9, t10)

	// Step 237: t10 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a2360
	for s := 0; s < 5; s++ {
		t10.Square(t10)
	}

	// Step 238: t9 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b
	t9.Mul(t9, t10)

	// Step 243: t9 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d60
	for s := 0; s < 5; s++ {
		t9.Square(t9)
	}

	// Step 244: t8 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f
	t8.Mul(t8, t9)

	// Step 251: t8 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b780
	for s := 0; s < 7; s++ {
		t8.Square(t8)
	}

	// Step 252: t8 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799
	t8.Mul(t0, t8)

	// Step 257: t8 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f320
	for s := 0; s < 5; s++ {
		t8.Square(t8)
	}

	// Step 258: t8 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339
	t8.Mul(t0, t8)

	// Step 261: t8 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799c8
	for s := 0; s < 3; s++ {
		t8.Square(t8)
	}

	// Step 262: t7 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf
	t7.Mul(t7, t8)

	// Step 265: t7 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce78
	for s := 0; s < 3; s++ {
		t7.Square(t7)
	}

	// Step 266: t7 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce79
	t7.Mul(&x, t7)

	// Step 274: t7 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce7900
	for s := 0; s < 8; s++ {
		t7.Square(t7)
	}

	// Step 275: t6 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d
	t6.Mul(t6, t7)

	// Step 283: t6 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d00
	for s := 0; s < 8; s++ {
		t6.Square(t6)
	}

	// Step 284: t6 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d09
	t6.Mul(t2, t6)

	// Step 287: t6 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c86848
	for s := 0; s < 3; s++ {
		t6.Square(t6)
	}

	// Step 288: t6 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c86849
	t6.Mul(&x, t6)

	// Step 295: t6 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e4342480
	for s := 0; s < 7; s++ {
		t6.Square(t6)
	}

	// Step 296: t6 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf
	t6.Mul(t5, t6)

	// Step 303: t6 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125f80
	for s := 0; s < 7; s++ {
		t6.Square(t6)
	}

	// Step 304: t5 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf
	t5.Mul(t5, t6)

	// Step 310: t5 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efc0
	for s := 0; s < 6; s++ {
		t5.Square(t5)
	}

	// Step 311: t5 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1
	t5.Mul(t4, t5)

	// Step 316: t5 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa20
	for s := 0; s < 5; s++ {
		t5.Square(t5)
	}

	// Step 317: t5 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa35
	t5.Mul(z, t5)

	// Step 323: t5 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d40
	for s := 0; s < 6; s++ {
		t5.Square(t5)
	}

	// Step 324: t4 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d51
	t4.Mul(t4, t5)

	// Step 331: t4 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a880
	for s := 0; s < 7; s++ {
		t4.Square(t4)
	}

	// Step 332: t3 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a897
	t3.Mul(t3, t4)

	// Step 339: t3 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b80
	for s := 0; s < 7; s++ {
		t3.Square(t3)
	}

	// Step 340: t3 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b95
	t3.Mul(z, t3)

	// Step 344: t3 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b950
	for s := 0; s < 4; s++ {
		t3.Square(t3)
	}

	// Step 345: t2 = x^0x41632889bd8224b3ca3f1682dfe740e45a69879a131cd11b5bcce790d092fdfa3544b959
	t2.Mul(t2, t3)

	// Step 351: t2 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d512e5640
	for s := 0; s < 6; s++ {
		t2.Square(t2)
	}

	// Step 352: t1 = x^0x1058ca226f60892cf28fc5a0b7f9d039169a61e684c73446d6f339e43424bf7e8d512e565d
	t1.Mul(t1, t2)

	// Step 357: t1 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacba0
	for s := 0; s < 5; s++ {
		t1.Square(t1)
	}

	// Step 358: t1 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacbb5
	t1.Mul(z, t1)

	// Step 364: t1 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a8972b2ed40
	for s := 0; s < 6; s++ {
		t1.Square(t1)
	}

	// Step 365: t0 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a8972b2ed59
	t0.Mul(t0, t1)

	// Step 371: t0 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacbb5640
	for s := 0; s < 6; s++ {
		t0.Square(t0)
	}

	// Step 372: t0 = x^0x20b19444dec11259e51f8b416ff3a0722d34c3cd098e688dade673c868497efd1aa25cacbb5655
	t0.Mul(z, t0)

	// Step 378: t0 = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a8972b2ed59540
	for s := 0; s < 6; s++ {
		t0.Square(t0)
	}

	// Step 379: z = x^0x82c651137b044967947e2d05bfce81c8b4d30f342639a236b799cf21a125fbf46a8972b2ed59555
	z.Mul(z, t0)

	return z
}
