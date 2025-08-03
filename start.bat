@echo off
echo 🚀 Iniciando Tivix Performance Tracker Backend...
echo.

:: Verificar se o Go está instalado
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go não está instalado ou não está no PATH
    echo Por favor, instale o Go em https://golang.org/dl/
    pause
    exit /b 1
)

:: Verificar se as dependências estão instaladas
if not exist "go.sum" (
    echo 📦 Instalando dependências...
    go mod download
    if %errorlevel% neq 0 (
        echo ❌ Erro ao baixar dependências
        pause
        exit /b 1
    )
)

:: Verificar se o arquivo .env existe
if not exist "../.env" (
    echo ⚠️  Arquivo .env não encontrado na raiz do projeto
    echo Por favor, configure as variáveis de ambiente no arquivo .env
    pause
    exit /b 1
)

echo ✅ Tudo configurado!
echo 🏃‍♂️ Executando servidor...
echo.
echo 🌐 Servidor será iniciado em: http://localhost:8080
echo 📖 Health check: http://localhost:8080/health
echo 📚 API Base URL: http://localhost:8080/api/v1
echo.
echo Para parar o servidor, pressione Ctrl+C
echo.

:: Executar o servidor
go run main.go
