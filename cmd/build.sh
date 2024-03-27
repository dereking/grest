#!/bin/bash 

ver=`date '+%Y.%m.%d.%H%M'`
out="grest"
 


if [ "$(uname)" == "Darwin" ]; then
  # Mac OS X  
  out="grest" 
elif [ "$(uname -s |grep 'Linux')" != "" ]; then
  # GNU/Linux  
  out="grest" 
elif [ "$(uname -s |grep 'MINGW')" != "" ]; then
  # Windows NT 
  out="$out.exe"
fi  
 

#cd $PRJ_PATH
go build -ldflags "-X main._version_='$ver'" -o $out
if [ -f "$GOROOT/bin/$out" ]; then 
  mv "$GOROOT/bin/$out" "$GOROOT/bin/$out.$ver.pre"
fi

 