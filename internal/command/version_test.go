package command

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestCommand_VersionCmd(t *testing.T) {
	// Сохраняем оригинальные значения
	origVersion := version
	origBuildDate := buildDate
	
	// Устанавливаем тестовые значения
	SetVersionInfo("1.0.0", "2023-01-01")
	
	// Создаем команду
	cmd := &Command{}
	versionCmd := cmd.VersionCmd()
	
	// Перенаправляем stdout для проверки вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Выполняем команду
	versionCmd.Run(versionCmd, []string{})
	
	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout
	
	// Читаем вывод
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	
	// Проверяем вывод
	assert.Contains(t, output, "Версия: 1.0.0")
	assert.Contains(t, output, "Дата сборки: 2023-01-01")
	
	// Восстанавливаем оригинальные значения
	version = origVersion
	buildDate = origBuildDate
}