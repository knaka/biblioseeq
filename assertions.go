// todo: remove me

package common

// ensure checks if the value is available. If err is not nil, it panics.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func ensure[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// assert checks that err is nil. If err is not nil, it panics.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func assert[T any](first T, rest ...any) {
	if len(rest) > 0 {
		if err, ok := (rest[len(rest)-1]).(error); ok && err != nil {
			panic(err)
		}
	}
	if err, ok := any(first).(error); ok && err != nil {
		panic(err)
	}
}

// ignore ignores errors explicitly.
func ignore[T any](T, ...any) {}
