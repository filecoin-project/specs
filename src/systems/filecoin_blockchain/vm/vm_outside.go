package vm

type VMInterpreter interface {
  ApplyMessage(inTree InputTree, msg Message) (outTree StateTree, ret MessageReceipt)
}

type InvocationInput struct {
  InTree    StateTree
  VMContext VMContext
  FromActor Actor
  ToActor   Actor
  Method    ActorMethod
  Params    Params
  Value     TokenAmount
  GasLimit  GasAmount
  // GasPrice  GasPrice
}

type InvocationOutput struct {
  OutTree     StateTree
  ExitCode    UVarint
  ReturnValue Bytes
  GasUsed     GasAmount
}

func (vmi *VMInterpreter) ApplyMessage(inTree InputTree, msg Message, minerAddr Address) (outTree StateTree, ret MessageReceipt) {

  compTree := inTree
  fromActor, found := compTree.GetActor(msg.From)
  if !found {
    Fatal("no such from actor")
  }

  // make sure fromActor has enough money to run the max invocation
  maxGasCost := (msg.GasLimit * msg.GasPrice)
  totalCost := msg.Value + maxGasCost
  if fromActor.Balance < totalCost {
    Fatal("not enough funds")
  }

  // make sure this is the right message order for fromActor
  // (this is protection against replay attacks, and useful sequencing)
  if msg.Nonce() != fromActor.SeqNo+1 {
    Fatal("invalid nonce")
  }

  // may return a different tree on succeess.
  // this MUST get rolled back if the invocation fails.
  var toActor Actor
  compTree, toActor = vmi.getOrCreateAccountActor(compTree, msgTo)

  // deduct maximum expenditure gas funds first
  compTree = vmi.deductFunds(compTree, fromActor, maxGasCost)

  // transfer funds fromActor -> toActor
  // (yes deductions can be combined, spelled out here for clarity)
  compTree = vmi.deductFunds(compTree, fromActor, msg.Value)
  compTree = vmi.depositFunds(compTree, toActor, msg.Value)

  // perform the method call to the actor
  // TODO: eval if we should lift gas tracking and calc to the beginning of invocation
  // (ie, include account creation, gas accounting itself)
  out := vmi.MethodDispatch(InvocationInput{
    InTree:    compTree,
    VMContext: makeVMContext(compTree, msg),
    FromActor: fromActor,
    ToActor:   toActor,
    Method:    msg.Method,
    Params:    msg.Params,
    Value:     msg.Value,
    GasLimit:  msg.GasLimit,
  })

  var outTree StateTree
  if out.ExitCode != 0 {
    // error -- revert all state changes -- ie drop updates. burn used gas.
    outTree = inTree // wipe!
    outTree = vmi.deductFunds(outTree, fromActor, vmi.calcGas(out.GasUsed, msg.GasPrice))

  } else {
    // success -- refund unused gas
    outTree = out.OutTree // take the state from the invocation output
    refundGas := msg.GasLimit - out.GasUsed
    outTree = vmi.depositFunds(outTree, msg.From, vmi.calcGas(refundGas, msg.GasPrice))
    outTree = vmi.incrementActorSeqNo(fromActor)
  }

  // reward miner gas fees
  outTree = vmi.depositFunds(outTree, minerAddr, vmi.calcGas(out.GasUsed, msg.GasPrice))

  return outTree, MessageReceipt{
    ExitCode: out.ExitCode,
    Return:   out.ReturnValue,
    GasUsed:  out.GasUsed,
  }
}

func (_ *VMInterpreter) dispatchMethod(input InvocationInput) InvocationOutput {
  if input.Method == 0 {
    // just sending money. move along.
    return InvocationOutput{
      OutTree:     input.InTree,
      GasUsed:     vmi.gasForSendingMoney(),
      ExitCode:    0,
      ReturnValue: nil,
    }
  }

  //TODO: actually invoke the funtion here.
  // put any vtable lookups in this function.
  output := input.ToActor.Call(input.StateTree, input.Method, input.Params)
}

func (_ *VMInterpreter) invokeMethod(vmctx VMContext, inTree StateTree, toActor Actor, method ActorMethod, params Params) {

}

func (_ *VMInterpreter) deductFunds(inTree StateTree, a Actor, amt TokenAmount) (outTree StateTree) {
  panic("todo")
}

func (_ *VMInterpreter) depositFunds(inTree StateTree, a Actor, amt TokenAmount) (outTree StateTree) {
  panic("todo")
}

func (_ *VMInterpreter) calcGas(gasUsed GasAmount, gasPrice GasPrice) {
  return gasUsed * gasPrice
}

func (vmi *VMInterpreter) getOrCreateAccountActor(inTree StateTree, addr Address) (outTree StateTree, _ Actor) {

  toActor, found := inTree.GetActor(msg.To)
  if found {
    return inTree, toActor
  }

  switch addr.Type() {
  case BLS:
    return vmi.newBLSAccountActor(inTree, addr)
  case Secp256k1:
    return vmi.newSecp256k1AccountActor(inTree, addr)
  case ID:
    Fatal("no actor with given ID")
  case Actor:
    Fatal("no such actor")
  }
}

func (_ *VMInterpreter) newBLSAccountActor(inTree StateTree, addr Address) (outTree StateTree, _ Actor) {
  panic("todo")
}

func (_ *VMInterpreter) newSecp256k1AccountActor(inTree StateTree, addr Address) (outTree StateTree, _ Actor) {
  panic("todo")
}

//------------
