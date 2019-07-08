set h=%time:~0,2%
set h=%h: =0%
set ver=%date:~2,2%.%date:~5,2%.%date:~8,2%.%h%%time:~3,2%

go build -ldflags "-X main._version_=%ver%" 
copy grest.exe %GOROOT%/bin
pause