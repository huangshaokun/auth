for /f "delims=" %%t in ('go version') do set goversion=%%t

set hour=%time:~,2%
if "%time:~,1%"==" " set hour=0%time:~1,1%

for /F %%i in ('git rev-parse --short HEAD') do ( set gitversion=%%i)

set GOOS=linux
set GOARCH=amd64

go build -i -v -a -ldflags "-s -w -X 'git.algor.tech/algor/algorlib/version.VERSION=1.1.1' -X 'git.algor.tech/algor/algorlib/version.BUILDTIME=%date:~0,4%-%date:~5,2%-%date:~8,2% %hour%:%time:~3,2%:%time:~6,2%' -X 'git.algor.tech/algor/algorlib/version.GOVERSION=%goversion%' -X 'git.algor.tech/algor/algorlib/version.GITHASH=%gitversion%'"