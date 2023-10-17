staticcheck ./cmd/yt-dlp-web
@if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
mkdir build
set "GOOS=linux"   & go build -o ./build/ ./cmd/yt-dlp-web
set "GOOS=windows" & go build -o ./build/ ./cmd/yt-dlp-web
@if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
