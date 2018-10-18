#!/bin/bash
#RUN_MODE=develop
#APP_NAME=grestcmd
#SVC_NAME=$APP_NAME-$RUN_MODE
#PRJ_PATH=/data/projects/$RUN_MODE/$APP_NAME
#BIN_NAME=$APP_NAME-$RUN_MODE
#RUN_PATH=/data/GORUN/$RUN_MODE/$APP_NAME

ver=`date '+%Y.%m.%d.%H%M'`

#mkdir -p $RUN_PATH

#export GOROOT=/data/GO1.4.1
#export GOPATH=/data/projects/develop/gopath-develop
#export GOPATH=$GOPATH:$PRJ_PATH

#cd $PRJ_PATH
$GOROOT/bin/go build -ldflags "-X main._version_=$ver" 
#-o $PRJ_PATH/bin/$BIN_NAME "src/$APP_NAME.go"

#cd $PRJ_PATH/bin
#/usr/bin/supervisorctl stop $SVC_NAME
#cp * $RUN_PATH -Rf
#cd $RUN_PATH
#/usr/bin/supervisorctl start $SVC_NAME