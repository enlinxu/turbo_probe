#!/bin/bash

options=" --v=4 --logtostderr "
options="$options --protocolVersion=6.1.0-SNAPSHOT"
options="$options --probeType=my-probe"
options="$options --serverHost=https://localhost:9400"
options="$options --turboUser=administrator"
options="$options --turboPasswd=a"

set -x
_output/turbo $options 
