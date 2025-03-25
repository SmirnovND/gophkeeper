package pkg

// TokenExtractorFunc определяет тип функции для извлечения логина из токена
type TokenExtractorFunc func(tokenString string) (string, error)

// TokenExtractor - переменная, содержащая функцию для извлечения логина из токена
// По умолчанию использует ExtractLoginFromToken, но может быть переопределена для тестов
var TokenExtractor TokenExtractorFunc = ExtractLoginFromToken