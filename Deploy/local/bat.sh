#!/bin/bash

export BasePath=/works/goProject/src/github.com/vmware/harbor
${BasePath}/Deploy/local/jsminify.sh ${BasePath}/views/sections/script-include.htm ${BasePath}/static/resources/js/harbor.app.min.js
