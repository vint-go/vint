package fixtures

import "os"

func mkdirOverlyPermissive() error {
	return os.Mkdir("/tmp/myapp", 0777) // MATCH /directory permission 0777 is more permissive than 0750/
}

func mkdirAllWorldReadable() error {
	return os.MkdirAll("/var/myapp/data", 0755) // MATCH /directory permission 0755 is more permissive than 0750/
}

func mkdirGroupWritable() error {
	return os.Mkdir("/tmp/shared", 0770) // MATCH /directory permission 0770 is more permissive than 0750/
}

func mkdirAllWorldExecutable() error {
	return os.MkdirAll("/opt/myapp", 0751) // MATCH /directory permission 0751 is more permissive than 0750/
}

func mkdirRestrictive() error {
	return os.Mkdir("/tmp/myapp", 0750)
}

func mkdirAllOwnerOnly() error {
	return os.MkdirAll("/var/myapp/data", 0700)
}

func mkdirMinimal() error {
	return os.Mkdir("/tmp/private", 0500)
}

func mkdirAllGroupRead() error {
	return os.MkdirAll("/opt/app", 0740)
}
