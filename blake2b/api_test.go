package blake2b

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"
)

// Tests Blocksize() declared in hash.Hash
func TestBlockSize(t *testing.T) {
	h, err := New(&Params{HashSize: Size})
	if err != nil {
		t.Fatalf("Could not create blake2b instance: %s", err)
	}
	if bs := h.BlockSize(); bs != BlockSize || bs != 128 {
		t.Fatalf("BlockSize() returned: %d - but expected: %d", bs, 128)
	}
}

// Tests Size() declared in hash.Hash
func TestSize(t *testing.T) {
	h, err := New(&Params{HashSize: Size})
	if err != nil {
		t.Fatalf("Could not create blake2b instance: %s", err)
	}
	if s := h.Size(); s != Size || s != 64 {
		t.Fatalf("Size() returned: %d - but expected: %d", s, 64)
	}

	h, err = New(&Params{HashSize: 32})
	if err != nil {
		t.Fatalf("Could not create blake2b instance: %s", err)
	}
	if s := h.Size(); s != Size/2 || s != 32 {
		t.Fatalf("Size() returned: %d - but expected: %d", s, 32)
	}
}

// Tests Reset() declared in hash.Hash
func TestReset(t *testing.T) {
	h, err := New(&Params{HashSize: Size})
	if err != nil {
		t.Fatalf("Could not create blake2b instance: %s", err)
	}
	b, ok := h.(*hashFunc)
	if !ok {
		t.Fatal("Impossible situation: New returns no blake2b struct")
	}
	orig := *b // copy

	var randData [BlockSize]byte
	if _, err := rand.Read(randData[:]); err != nil {
		t.Fatalf("Failed to read random bytes form crypto/rand: %s", err)
	}

	b.Write(randData[:])
	b.Reset()

	if b.hsize != orig.hsize {
		t.Fatalf("Reseted hsize field: %d - but expected: %d", b.hsize, orig.hsize)
	}
	if b.keyed != orig.keyed {
		t.Fatalf("Reseted keyed field: %v - but expected: %v", b.keyed, orig.keyed)
	}
	if b.ctr != orig.ctr {
		t.Fatalf("Reseted ctr field: %v - but expected: %v", b.ctr, orig.ctr)
	}
	if b.off != orig.off {
		t.Fatalf("Reseted off field: %d - but expected: %d", b.off, orig.off)
	}
	if b.buf != orig.buf {
		t.Fatalf("Reseted buf field %v - but expected %v", b.buf, orig.buf)
	}
	if b.key != orig.key {
		t.Fatalf("Reseted key field: %v - but expected: %v", b.key, orig.key)
	}
	if b.hVal != orig.hVal {
		t.Fatalf("Reseted hVal field: %v - but expected: %v", b.hVal, orig.hVal)
	}
	if b.initVal != orig.initVal {
		t.Fatalf("Reseted initVal field: %v - but expected: %v", b.initVal, orig.initVal)
	}
}

// Tests Write(p []byte) declared in hash.Hash
func TestWrite(t *testing.T) {
	h, err := New(&Params{})
	if err != nil {
		t.Fatalf("Failed to create instance of blake2b - Cause: %s", err)
	}
	n, err := h.Write(nil)
	if n != 0 || err != nil {
		t.Fatalf("Failed to process nil slice: Processed bytes: %d - Returned error: %s", n, err)
	}
	n, err = h.Write(make([]byte, h.Size()))
	if n != h.Size() || err != nil {
		t.Fatalf("Failed to process 0-slice with length %d: Processed bytes: %d - Returned error: %s", h.Size(), n, err)
	}
	n, err = h.Write(make([]byte, h.BlockSize()))
	if n != h.BlockSize() || err != nil {
		t.Fatalf("Failed to process 0-slice with length %d: Processed bytes: %d - Returned error: %s", h.BlockSize(), n, err)
	}
	n, err = h.Write(make([]byte, 211)) // 211 = (2*3*5*7)+1 is prime
	if n != 211 || err != nil {
		t.Fatalf("Failed to process 0-slice with length %d: Processed bytes: %d - Returned error: %s", 211, n, err)
	}
}

// Tests Sum(b []byte) declared in hash.Hash
func TestSum(t *testing.T) {
	h, err := New(&Params{HashSize: 32})
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}
	var one = [1]byte{1}

	h.Sum(nil)
	h.Write(make([]byte, BlockSize))
	h.Write(one[:])

	sum1 := h.Sum(nil)
	sum2 := Sum256(append(make([]byte, BlockSize), one[:]...))

	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
}

// Tests New(p *Params) declared here (blake2b)
func TestNew(t *testing.T) {
	p := &Params{}
	_, err := New(p)
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	p.HashSize = 80 // invalid but verify should adjust this
	_, err = New(p)
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	p.Key = make([]byte, Size)
	_, err = New(p)
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	p.Key = make([]byte, Size+1)
	_, err = New(p)
	if err == nil {
		t.Fatalf("Verification of key parameter failed: Accepted illegal keysize: %d", Size+1)
	}
	p.Key = nil

	p.Salt = make([]byte, saltSize)
	_, err = New(p)
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	p.Salt = make([]byte, saltSize+1)
	_, err = New(p)
	if err == nil {
		t.Fatalf("Verification of salt parameter failed: Accepted illegal saltsize: %d", saltSize+1)
	}
	p.Salt = nil
}

// Tests Sum256(msg []byte) declared here (blake2b)
func TestSum256(t *testing.T) {
	h, err := New(&Params{HashSize: 32})
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	h.Write(nil)
	sum1 := h.Sum(nil)
	sum2 := Sum256(nil)
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
	h.Reset()

	h.Write(make([]byte, 1))
	sum1 = h.Sum(nil)
	sum2 = Sum256(make([]byte, 1))
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
	h.Reset()

	h.Write(make([]byte, BlockSize+1))
	sum1 = h.Sum(nil)
	sum2 = Sum256(make([]byte, BlockSize+1))
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
}

// Tests Sum512(msg []byte) declared here (blake2b)
func TestSum512(t *testing.T) {
	h, err := New(&Params{HashSize: Size})
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	h.Write(nil)
	sum1 := h.Sum(nil)
	sum2 := Sum512(nil)
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
	h.Reset()

	h.Write(make([]byte, 1))
	sum1 = h.Sum(nil)
	sum2 = Sum512(make([]byte, 1))
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
	h.Reset()

	h.Write(make([]byte, BlockSize+1))
	sum1 = h.Sum(nil)
	sum2 = Sum512(make([]byte, BlockSize+1))
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
}

// Tests Sum(msg []byte. p *Params) declared here (blake2b)
func TestSumFunc(t *testing.T) {
	p := &Params{}
	h, err := New(p)
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	h.Write(nil)
	sum1 := h.Sum(nil)
	sum2, err := Sum(nil, p)
	if err != nil {
		t.Fatalf("Failed to calculate the sum: %s", err)
	}
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}

	p.HashSize = 48
	h, err = New(p)
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	h.Write(make([]byte, 1))
	sum1 = h.Sum(nil)
	sum2, err = Sum(make([]byte, 1), p)
	if err != nil {
		t.Fatalf("Failed to calculate the sum: %s", err)
	}
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}

	p.Salt = make([]byte, 14)
	h, err = New(p)
	if err != nil {
		t.Fatalf("Failed to create blake2b instance: %s", err)
	}

	h.Write(make([]byte, BlockSize+1))
	sum1 = h.Sum(nil)
	sum2, err = Sum(make([]byte, BlockSize+1), p)
	if err != nil {
		t.Fatalf("Failed to calculate the sum: %s", err)
	}
	if !bytes.Equal(sum1, sum2[:]) {
		t.Fatalf("Hash does not match:\nFound:    %s\nExpected: %s", hex.EncodeToString(sum1), hex.EncodeToString(sum2[:]))
	}
}
