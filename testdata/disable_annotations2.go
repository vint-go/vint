package fixtures

func foo1() {
	//revive:disable-next-line:lint/style/useVarNaming
	var invalid_name = 0
	var invalid_name2 = 1 //revive:disable-line:lint/style/useVarNaming
}

func foo11() {
	//revive:disable-next-line:lint/style/useVarNaming I'm an Eiffel programmer thus I like underscores
	var invalid_name = 0
	var invalid_name2 = 1 //revive:disable-line:lint/style/useVarNaming I'm an Eiffel programmer thus I like underscores
}

func foo2() {
	// 		revive:disable-next-line:lint/style/useVarNaming
	//revive:disable
	var invalid_name = 0
	var invalid_name2 = 1
}

func foo3() {
	//revive:enable
	// revive:disable-next-line:lint/style/useVarNaming
	var invalid_name = 0
	// not a valid annotation revive:disable-next-line:lint/style/useVarNaming
	var invalid_name2 = 1 // MATCH /don't use underscores in Go names; var invalid_name2 should be invalidName2/
	/* revive:disable-next-line:lint/style/useVarNaming */
	var invalid_name3 = 0 // MATCH /don't use underscores in Go names; var invalid_name3 should be invalidName3/
}

func foo2p1() {
	//revive:disable Underscores are fine
	var invalid_name = 0
	//revive:enable No! Underscores are not nice!
	var invalid_name2 = 1 // MATCH /don't use underscores in Go names; var invalid_name2 should be invalidName2/
}
