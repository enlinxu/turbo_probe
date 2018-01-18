#!/bin/bash

options=" --v=3 --logtostderr "
options="$options --protocolVersion=6.1.0-SNAPSHOT"
options="$options --probeType=my-probe"
options="$options --serverHost=https://localhost:9400"

set -x
_output/turbo $options 
