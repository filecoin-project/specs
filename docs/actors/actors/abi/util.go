package abi

func assertNoError(e error) {
	if e != nil {
		panic(e.Error())
	}
}
