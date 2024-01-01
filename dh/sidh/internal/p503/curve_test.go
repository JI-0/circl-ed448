// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package p503

import (
	"bytes"
	crand "crypto/rand"
	. "github.com/JI-0/circl-ed448/dh/sidh/internal/common"
	"io"
	"math"
	"math/rand"
	"testing"
	"time"
)

func vartimeEqProjFp2(lhs, rhs *ProjectivePoint) bool {
	var t0, t1 Fp2
	mul(&t0, &lhs.X, &rhs.Z)
	mul(&t1, &lhs.Z, &rhs.X)
	return vartimeEqFp2(&t0, &t1)
}

func toAffine(point *ProjectivePoint) *Fp2 {
	var affineX Fp2
	inv(&affineX, &point.Z)
	mul(&affineX, &affineX, &point.X)
	return &affineX
}

func Test_jInvariant(t *testing.T) {
	curve := ProjectiveCurveParameters{A: curveA, C: curveC}
	jbufRes := make([]byte, params.SharedSecretSize)
	jbufExp := make([]byte, params.SharedSecretSize)
	var jInv Fp2

	Jinvariant(&curve, &jInv)
	FromMontgomery(&jInv, &jInv)
	Fp2ToBytes(jbufRes, &jInv, params.Bytelen)

	jInv = expectedJ
	FromMontgomery(&jInv, &jInv)
	Fp2ToBytes(jbufExp, &jInv, params.Bytelen)

	if !bytes.Equal(jbufRes[:], jbufExp[:]) {
		t.Error("Computed incorrect j-invariant: found\n", jbufRes, "\nexpected\n", jbufExp)
	}
}

func TestProjectivePointVartimeEq(t *testing.T) {
	var xP ProjectivePoint

	xP = ProjectivePoint{X: affineXP, Z: params.OneFp2}
	xQ := xP

	// Scale xQ, which results in the same projective point
	mul(&xQ.X, &xQ.X, &curveA)
	mul(&xQ.Z, &xQ.Z, &curveA)
	if !vartimeEqProjFp2(&xP, &xQ) {
		t.Error("Expected the scaled point to be equal to the original")
	}
}

func TestPointMulVersusSage(t *testing.T) {
	curve := ProjectiveCurveParameters{A: curveA, C: curveC}
	cparams := CalcCurveParamsEquiv4(&curve)
	var xP ProjectivePoint

	// x 2
	xP = ProjectivePoint{X: affineXP, Z: params.OneFp2}
	Pow2k(&xP, &cparams, 1)
	afxQ := toAffine(&xP)
	if !vartimeEqFp2(afxQ, &affineXP2) {
		t.Error("\nExpected\n", affineXP2, "\nfound\n", afxQ)
	}

	// x 4
	xP = ProjectivePoint{X: affineXP, Z: params.OneFp2}
	Pow2k(&xP, &cparams, 2)
	afxQ = toAffine(&xP)
	if !vartimeEqFp2(afxQ, &affineXP4) {
		t.Error("\nExpected\n", affineXP4, "\nfound\n", afxQ)
	}
}

func TestPointMul9VersusSage(t *testing.T) {
	curve := ProjectiveCurveParameters{A: curveA, C: curveC}
	cparams := CalcCurveParamsEquiv3(&curve)
	var xP ProjectivePoint

	xP = ProjectivePoint{X: affineXP, Z: params.OneFp2}
	Pow3k(&xP, &cparams, 2)
	afxQ := toAffine(&xP)
	if !vartimeEqFp2(afxQ, &affineXP9) {
		t.Error("\nExpected\n", affineXP9, "\nfound\n", afxQ)
	}
}

func BenchmarkThreePointLadder(b *testing.B) {
	curve := ProjectiveCurveParameters{A: curveA, C: curveC}
	for n := 0; n < b.N; n++ {
		ScalarMul3Pt(&curve, &threePointLadderInputs[0], &threePointLadderInputs[1], &threePointLadderInputs[2], uint(len(scalar3Pt)*8), scalar3Pt[:])
	}
}

/* -------------------------------------------------------------------------
   Generate invalid public key points / ciphertext for test TestKEMInvalidPK
   -------------------------------------------------------------------------*/

// Left-to-right Montgomery ladder, Algorithm 4 in Costello-Smith
// Input: ProjectivePoint P (xP, zP)
// Output: x([scalar]P), z([scalar]P)
func montgomeryLadder(cparams *ProjectiveCurveParameters, P *ProjectivePoint, scalar []uint8, random uint) ProjectivePoint {
	var R0, R2, R1 ProjectivePoint
	coefEq := CalcCurveParamsEquiv4(cparams) // for xDbl
	aPlus2Over4 := CalcAplus2Over4(cparams)  // for xDblAdd
	R0 = *P                                  // RO <- P
	R1 = *P
	Pow2k(&R1, &coefEq, 1) // R1 <- [2]P
	R2 = *P                // R2 = R1-R0 = P

	prevBit := uint8(0)
	for i := int(random); i >= 0; i-- {
		bit := (scalar[i>>3] >> (i & 7) & 1)
		swap := prevBit ^ bit
		prevBit = bit
		cswap(&R0.X, &R0.Z, &R1.X, &R1.Z, swap)
		R0, R1 = xDbladd(&R0, &R1, &R2, &aPlus2Over4)
	}
	cswap(&R0.X, &R0.Z, &R1.X, &R1.Z, prevBit)
	return R0
}

// P = P + T
// From paper https://eprint.iacr.org/2017/212.pdf
// The map tau_T: P->P+T is (X : Z) -> (Z : X) on Montgomery curves
func tauT(P *ProjectivePoint) {
	P.X, P.Z = P.Z, P.X // magic!
}

// Construct Invalid public key tuple (P,Q) such that P and Q are linearly dependent
// Simulate section 3.1.1 of paper https://eprint.iacr.org/2022/054.pdf
// We only construct point P and Q because in the attacks the third point is P-Q by construction
// and the countermeasure does not test it
// Without loss of generality, we assume the curve is the starting curve
func testInvalidPKNoneLinear(t *testing.T) {

	// Generate random scalar as secret
	secret := make([]byte, params.B.SecretByteLen)
	_, err := io.ReadFull(crand.Reader, secret)
	if err != nil {
		t.Error("Fail read random bytes")
	}

	var P, Q ProjectivePoint

	rand.Seed(time.Now().UnixNano())
	random_index := rand.Intn(int(params.B.SecretByteLen-1) * 8)

	// Set P as a point of order 3^e3
	P = ProjectivePoint{X: params.B.AffineP, Z: params.OneFp2}

	// Set Q = [k]P, where k = secret[:random_index]
	Q = montgomeryLadder(&params.InitCurve, &P, secret, uint(random_index))

	// Make sure Q is of full order 3^e_3,
	var test_Q ProjectivePoint
	test_Q = Q

	var e3 uint32
	e3_float := float64(int(params.B.SecretBitLen)+1) / math.Log2(3)
	e3 = uint32(e3_float)
	cparam_q := CalcCurveParamsEquiv3(&params.InitCurve)
	Pow3k(&test_Q, &cparam_q, e3-1)

	var test_QZ Fp2
	FromMontgomery(&test_QZ, &test_Q.Z)

	// Q are not of full order 3^e_3
	for isZero(&test_QZ) == 1 {
		rand.Seed(time.Now().UnixNano())
		random_index = rand.Intn(int(params.B.SecretByteLen-1) * 8)
		Q = montgomeryLadder(&params.InitCurve, &P, secret, uint(random_index))
		test_Q = Q
		Pow3k(&test_Q, &cparam_q, e3-1)
		FromMontgomery(&test_QZ, &test_Q.Z)
	}

	// invQz = 1/Q.Z
	var invQz Fp2
	invQz = Q.Z
	inv(&invQz, &invQz)

	mul(&P.X, &P.X, &P.Z)
	mul(&Q.X, &Q.X, &invQz)

	var xP, xQ, xQmP ProjectivePoint
	xP = ProjectivePoint{X: P.X, Z: params.OneFp2}
	xQ = ProjectivePoint{X: Q.X, Z: params.OneFp2}
	xQmP = ProjectivePoint{X: params.OneFp2, Z: params.OneFp2}

	error_verify := PublicKeyValidation(&params.InitCurve, &xP, &xQ, &xQmP, params.B.SecretBitLen)
	if error_verify == nil {
		t.Errorf("\nExpect linearly dependent ciphertext to fail, index: %v  scalar: %v ", random_index, secret)
	}
}

// Construct Invalid public key tuple (P,Q) such that Q = [k]P + T, where k is random and T is the point of order 2.
// Simulate HB and section 3.1.2 of paper https://eprint.iacr.org/2022/054.pdf
// We only construct point P and Q because in the attacks the third point is P-Q by construction
// and the countermeasure does not test it
// Without loss of generality, we assume the curve is the starting curve
func testInvalidPKT(t *testing.T) {

	// Generate random scalar as secret
	secret := make([]byte, params.B.SecretByteLen)
	_, err := io.ReadFull(crand.Reader, secret)
	if err != nil {
		t.Error("Fail read random bytes")
	}

	var P, Q ProjectivePoint

	rand.Seed(time.Now().UnixNano())
	random_index := rand.Intn(int(params.B.SecretByteLen-1) * 8)

	// Set P as a point of order 3^e3
	P = ProjectivePoint{X: params.B.AffineP, Z: params.OneFp2}

	// Set Q = [k]P, where k = secret[:random_index]
	Q = montgomeryLadder(&params.InitCurve, &P, secret, uint(random_index))
	// Q = [k]P + T
	tauT(&Q)

	var invQz Fp2
	invQz = Q.Z
	inv(&invQz, &invQz)

	mul(&P.X, &P.X, &P.Z)
	mul(&Q.X, &Q.X, &invQz)

	var xP, xQ, xQmP ProjectivePoint
	xP = ProjectivePoint{X: P.X, Z: params.OneFp2}
	xQ = ProjectivePoint{X: Q.X, Z: params.OneFp2}
	xQmP = ProjectivePoint{X: params.OneFp2, Z: params.OneFp2}

	error_verify := PublicKeyValidation(&params.InitCurve, &xP, &xQ, &xQmP, params.B.SecretBitLen)
	if error_verify == nil {
		t.Errorf("\nExpect ciphertext involve point T to fail, index: %v  scalar: %v ", random_index, secret)
	}
}

// Construct Invalid public key tuple (P,Q) such that P and Q are in E[2^e2]
// Simulate section 3.2 of paper https://eprint.iacr.org/2022/054.pdf
// We only construct point P and Q because in the attacks the third point is P-Q by construction
// and the countermeasure does not test it
// Without loss of generality, we assume the curve is the starting curve
func testInvalidPKOrder2(t *testing.T) {

	// Generate random scalar as secret
	secret := make([]byte, params.B.SecretByteLen)
	_, err := io.ReadFull(crand.Reader, secret)
	if err != nil {
		t.Error("Fail read random bytes")
	}

	var P, Q ProjectivePoint

	P = ProjectivePoint{X: params.A.AffineP, Z: params.OneFp2}
	Q = ProjectivePoint{X: params.A.AffineQ, Z: params.OneFp2}

	rand.Seed(time.Now().UnixNano())
	random_index_p := rand.Intn(int(params.A.SecretByteLen-1) * 8)
	random_index_q := rand.Intn(int(params.A.SecretByteLen-1) * 8)

	P = montgomeryLadder(&params.InitCurve, &P, secret, uint(random_index_p))
	Q = montgomeryLadder(&params.InitCurve, &Q, secret, uint(random_index_q))

	var invQz, invPz Fp2
	invQz = Q.Z
	invPz = P.Z
	inv(&invQz, &invQz)
	inv(&invPz, &invPz)

	mul(&P.X, &P.X, &invPz)
	mul(&Q.X, &Q.X, &invQz)

	var xP, xQ, xQmP ProjectivePoint
	xP = ProjectivePoint{X: P.X, Z: params.OneFp2}
	xQ = ProjectivePoint{X: Q.X, Z: params.OneFp2}
	xQmP = ProjectivePoint{X: params.OneFp2, Z: params.OneFp2}

	error_verify := PublicKeyValidation(&params.InitCurve, &xP, &xQ, &xQmP, params.B.SecretBitLen)
	if error_verify == nil {
		t.Errorf("\nExpect ciphertext in torsion E[2^e2] to fail, index_p: %v  index_q: %v  scalar: %v ", random_index_p, random_index_q, secret)
	}

}

// Construct Invalid public key tuple (P,Q) such that P and Q are in E[3^e3] but not of full order 3^e3
// Simulate section 3.1.1 of paper https://eprint.iacr.org/2022/054.pdf
// We only construct point P and Q because in the attacks the third point is P-Q by construction
// and the countermeasure does not test it
// Without loss of generality, we assume the curve is the starting curve
func testInvalidPKFullOrder(t *testing.T) {

	var P, Q ProjectivePoint

	P = ProjectivePoint{X: params.B.AffineP, Z: params.OneFp2}
	Q = ProjectivePoint{X: params.B.AffineQ, Z: params.OneFp2}

	var e3 uint32
	e3_float := float64(int(params.B.SecretBitLen)+1) / math.Log2(3)
	e3 = uint32(e3_float)

	rand.Seed(time.Now().UnixNano())
	random_index_p := rand.Intn(int(e3))
	random_index_q := rand.Intn(int(e3))

	cparam_q := CalcCurveParamsEquiv3(&params.InitCurve)
	Pow3k(&P, &cparam_q, uint32(random_index_p))
	Pow3k(&Q, &cparam_q, uint32(random_index_q))

	var invQz, invPz Fp2
	invQz = Q.Z
	invPz = P.Z
	inv(&invQz, &invQz)
	inv(&invPz, &invPz)

	mul(&P.X, &P.X, &invPz)
	mul(&Q.X, &Q.X, &invQz)

	var xP, xQ, xQmP ProjectivePoint
	xP = ProjectivePoint{X: P.X, Z: params.OneFp2}
	xQ = ProjectivePoint{X: Q.X, Z: params.OneFp2}
	xQmP = ProjectivePoint{X: params.OneFp2, Z: params.OneFp2}

	error_verify := PublicKeyValidation(&params.InitCurve, &xP, &xQ, &xQmP, params.B.SecretBitLen)
	if error_verify == nil {
		t.Errorf("\nExpect ciphertext not of full order to fail, index_p: %v  index_q: %v  ", random_index_p, random_index_q)
	}

}

// A trivial test case not covered by paper https://eprint.iacr.org/2022/054.pdf and HB
// Countermeasure in https://eprint.iacr.org/2022/054.pdf only cares about P and Q
// But if PmQ is point T or O, that can also lead to recovery of the first bit
func testInvalidPmQ(t *testing.T) {

	var zero Fp2
	var xP, xQ, xQmP ProjectivePoint
	xP = ProjectivePoint{X: params.A.AffineP, Z: params.OneFp2}
	xQ = ProjectivePoint{X: params.A.AffineQ, Z: params.OneFp2}
	xQmP = ProjectivePoint{X: zero, Z: params.OneFp2}

	error_verify := PublicKeyValidation(&params.InitCurve, &xP, &xQ, &xQmP, params.B.SecretBitLen)
	if error_verify == nil {
		t.Errorf("\nExpect PmQ as T to fail\n")
	}

}

// Test valid ciphertext
// Where P, Q are linearly independent points of correct order 3^e3 in E[3^e3]
func testValidPQ(t *testing.T) {

	var xP, xQ, xQmP ProjectivePoint
	xP = ProjectivePoint{X: params.B.AffineP, Z: params.OneFp2}
	xQ = ProjectivePoint{X: params.B.AffineQ, Z: params.OneFp2}
	xQmP = ProjectivePoint{X: params.OneFp2, Z: params.OneFp2}

	error_verify := PublicKeyValidation(&params.InitCurve, &xP, &xQ, &xQmP, params.B.SecretBitLen)
	if error_verify != nil {
		t.Errorf("\nExpect correct ciphertext to not fail\n")
	}

}

/* -------------------------------------------------------------------------
   Public key / Ciphertext validation against attacks proposed in paper https://eprint.iacr.org/2022/054.pdf and HB
   -------------------------------------------------------------------------*/

func TestInvalidPK(t *testing.T) {

	t.Run("InvalidPmQ", testInvalidPmQ)
	t.Run("InvalidPKNoneLinear", testInvalidPKNoneLinear)
	t.Run("InvalidPKT", testInvalidPKT)
	t.Run("InvalidPKOrder2", testInvalidPKOrder2)
	t.Run("InvalidPKFullOrder", testInvalidPKFullOrder)
	t.Run("ValidPQ", testValidPQ)

}
