package filproofs

import (
	file "github.com/filecoin-project/specs/systems/filecoin_files/file"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
/// Binary Merkle-tree generation

// BuildTree<T>
func BuildTree_T(data []byte) (util.T, file.Path) {
	// Plan: define this in terms of BinaryHash_T, then copy-paste changes into T-specific specializations, for now.

	// Nodes are always the digest size so data cannot be compressed to digest for storage.
	nodeSize := DigestSize_T()

	// TODO: Fail if len(dat) is not a power of 2 and a multiple of the node size.

	rows := [][]byte{data}

	for row := []byte{}; len(row) > nodeSize; {
		for i := 0; i < len(data); i += 2 * nodeSize {
			left := data[i : i+nodeSize]
			right := data[i+nodeSize : i+2*nodeSize]

			hashed := BinaryHash_T(left, right)

			row = append(row, AsBytes_T(hashed)...)
		}
		rows = append(rows, row)
	}

	// Last row is the root
	root := rows[len(rows)-1]

	if len(root) != nodeSize {
		panic("math failed us")
	}

	var filePath file.Path // TODO: dump tree to file.
	// NOTE: merkle tree file layout is illustrative, not prescriptive.

	// TODO: Check above more carefully. It's just an untested sketch for the moment.
	return fromBytes_T(root), filePath
}

// BuildTree<PedersenHash>
func BuildTree_PedersenHash(data []byte) (PedersenHash, file.Path) {
	return PedersenHash{}, file.Path("") // FIXME
}

//  BuildTree<SHA256Hash>
func BuildTree_SHA256Hash(data []byte) (SHA256Hash, file.Path) {
	return []byte{}, file.Path("") // FIXME
}
