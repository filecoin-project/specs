package codeGen

type Message interface {
	To()        Address
	From()      Address
	Nonce()     Word
	Value()     BigInt
	GasPrice()  BigInt
	GasLimit()  BigInt
	Method()    Word
	Params()    []byte
}

type MessageReceipt interface {
	ExitCode()  byte
	Return()    []byte
	GasUsed()   BigInt
}

type Block interface {
	MinerAddress()               Address
	Tickets()                    []Ticket
	ParentTipset()               Tipset
	Weight()                     BlockWeight
	Height()                     Word
	StateTree()                  StateTree
	Messages()                   []Message
	BLSAggregate()               Signature
	MessageReceipts()            []MessageReceipt
	Timestamp()                  Timestamp
	BlockSig()                   Signature

	SerializeSigned()            []byte
	ComputeUnsignedFingerprint() []byte
}

type BlockWeight BigInt

func (block *BlockI) ComputeUnsignedFingerprint() []byte {
	return Hash(HashRole_BlockSig, block.SerializeUnsigned());	
}

func (block *BlockI) SerializeUnsigned() []byte {
	panic("TODO")
}

func (block *BlockI) SerializeSigned() []byte {
	panic("TODO")
}

////////////////////
// Implementation //
////////////////////

type BlockI struct {
	minerAddress    Address
	tickets         []Ticket
	parentTipset    Tipset
	weight          BigInt
	height          Word
	stateTree       StateTree
	messages        []Message
	blsAggregate    Signature
	messageReceipts []MessageReceipt
	timestamp       Timestamp
	blockSig        Signature
}

func (block *BlockI) MinerAddress() Address {
	return block.minerAddress
}

func (block *BlockI) Tickets() []Ticket {
	return block.tickets
}

func (block *BlockI) ParentTipset() Tipset {
	panic("TODO")
}

func (block *BlockI) Weight() BlockWeight {
	panic("TODO")
}

func (block *BlockI) Height() Word {
	panic("TODO")
}

func (block *BlockI) StateTree() StateTree {
	panic("TODO")
}

func (block *BlockI) Messages() []Message {
	panic("TODO")
}

func (block *BlockI) BLSAggregate() Signature {
	return block.blsAggregate
}

func (block *BlockI) MessageReceipts() []MessageReceipt {
	return block.messageReceipts
}

func (block *BlockI) Timestamp() Timestamp {
	panic("TODO")
}

func (block *BlockI) BlockSig() Signature {
	panic("TODO")
}
