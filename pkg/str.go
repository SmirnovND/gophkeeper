package pkg

func GetExtensionByPath(filePath string) string {
	// Извлекаем расширение файла из пути
	extension := ""
	lastDotIndex := -1
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '.' {
			lastDotIndex = i
			break
		}
		if filePath[i] == '/' || filePath[i] == '\\' {
			break
		}
	}
	if lastDotIndex != -1 {
		extension = filePath[lastDotIndex+1:]
	}

	return extension
}
