#!/bin/bash

options=" --v=3 --logtostderr "
options="$options --serverHost=https://localhost:9400"

_output/turbo $options 
