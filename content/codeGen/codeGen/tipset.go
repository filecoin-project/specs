package codeGen

type Tipset interface {
	Blocks() []Block
	ElectionSeed() ElectionSeed

	StateTree() StateTree
}

type ElectionSeed interface {
	ToBytesInternal() []byte
	FromBytesInternal([]byte) ElectionSeed
}

func (tipset *TipsetI) ElectionSeed() ElectionSeed {
	var ret []byte

	for _, currBlock := range tipset.Blocks() {
		for _, currTicket := range currBlock.Tickets() {
			currSeed := Hash(
				HashRole_ElectionSeedFromVRFOutput,
				currTicket.VRFResult().ToBytesInternal(),
			)
			if ret == nil || CompareBytesStrict(currSeed, ret) < 0 {
				ret = currSeed
			}
		}
	}

	Assert(ret != nil)
	return ElectionSeed.FromBytesInternal(nil, ret)
}

func (tipset *TipsetI) StateTree() StateTree {
	var currTree StateTree = nil
	for _, block := range tipset.Blocks() {
		if currTree == nil {
			currTree = block.StateTree()
		} else {
			Assert(block.StateTree().CID().Equals(currTree.CID()))
		}
	}
	Assert(currTree != nil)
	for _, block := range tipset.Blocks() {
		currTree = UpdateStateTree(currTree, block)
	}
	return currTree
}

func UpdateStateTree(tree StateTree, block Block) StateTree {
	panic("TODO")
}

type TipsetI struct {
	blocks []Block
}

func (tipset *TipsetI) Blocks() []Block {
	return tipset.blocks
}

func (__ *ElectionSeedI) FromBytesInternal(data []byte) ElectionSeed {
	Assert(__ == nil)

	return &ElectionSeedI{
		data,
	}
}

func (seed *ElectionSeedI) ToBytesInternal() []byte {
	return seed.data
}

type ElectionSeedI struct {
	data []byte
}
