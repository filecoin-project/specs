package impl

import filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
import gascost "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/gascost"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import util "github.com/filecoin-project/specs/util"
import vmr "github.com/filecoin-project/specs/actors/runtime"

type Any = util.Any
type Int = util.Int
type ComputeFunctionID = vmr.ComputeFunctionID

type ComputeFunctionBody = func([]Any) Any
type ComputeFunctionGasCostFn = func([]Any) msg.GasAmount

type ComputeFunctionDef struct {
	Body      ComputeFunctionBody
	GasCostFn ComputeFunctionGasCostFn
}

var _computeFunctionDefs = map[ComputeFunctionID]ComputeFunctionDef{}

func init() {
	// VerifySignature
	_computeFunctionDefs[vmr.ComputeFunctionID_VerifySignature] = ComputeFunctionDef{
		Body: func(args []Any) Any {
			if len(args) != 3 {
				return nil
			}
			i := 0

			pk, ok := args[i].(filcrypto.PublicKey)
			i++
			if !ok {
				return nil
			}
			sig, ok := args[i].(filcrypto.Signature)
			i++
			if !ok {
				return nil
			}
			m, ok := args[i].(filcrypto.Message)
			i++
			if !ok {
				return nil
			}

			if i != len(args) {
				return nil
			}

			valid, err := filcrypto.Verify(pk, sig, m)
			if err != nil {
				return nil
			}
			return valid
		},
		GasCostFn: func(args []Any) msg.GasAmount {
			return gascost.PublicKeyCryptoOp
		},
	}
}
