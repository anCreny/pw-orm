package pworm

import "github.com/anCreny/pw-orm/helpers"

// идентификатор экземпляра пакета в памяти работы приложения.
// Он нужен для гарантированнного разграничения условно глобальных переменных
// в powershell в рамках разных экземпляров приложений, использующих один и тот
// же powershell
var scopeID string

func Init() {
	scopeID = helpers.GenerateRandomString(10)
}

func ScopeID() string {
	return scopeID
}
