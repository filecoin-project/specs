package abi

import (
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// PieceCID is the main reference to pieces in Filecoin. It is the CID of the piece.
type PieceCID cid.Cid

func (c PieceCID) Bytes() []byte {
	return cid.Cid(c).Bytes()
}

func PieceCIDOf(data []byte) PieceCID {
	b := cid.V1Builder{Codec: cid.DagCBOR, MhType: mh.SHA2_256}
	c, err := b.Sum(data)
	assertNoError(err)
	return PieceCID(c)
}

// PieceSize is the size of a piece, in bytes
type PieceSize struct {
	PayloadSize  int64
	OverheadSize int64
}

func (p PieceSize) Total() int64 {
	return p.PayloadSize + p.OverheadSize
}

type PieceInfo struct {
	Size     int64 // Size in nodes. For BLS12-381 (capacity 254 bits), must be >= 16. (16 * 8 = 128)
	PieceCID PieceCID
}

type PieceInfos struct {
	Items []PieceInfo
}
