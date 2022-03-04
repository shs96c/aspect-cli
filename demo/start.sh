#!/usr/bin/env bash

########################
# include the magic
########################
. demo-magic.sh


########################
# Configure the options
########################

#
# speed at which to simulate typing. bigger num = faster
#
# TYPE_SPEED=10

#
# custom prompt
#
# see http://www.tldp.org/HOWTO/Bash-Prompt-HOWTO/bash-prompt-escape-sequences.html for escape sequences
#
# DEMO_PROMPT="${GREEN}➜ ${CYAN}\W "

# text color
# DEMO_CMD_COLOR=$BLACK

cd ..

pei "bazel build ..."

cp -f bazel-bin/cmd/aspect/aspect_/aspect .
cp -f bazel-bin/cmd/aspect/aspect_/aspect examples/fix-visibility
cp -f bazel-bin/cmd/aspect/aspect_/aspect examples/error-augmentor
cp -f bazel-bin/cmd/aspect/aspect_/aspect examples/predefined-queries

cd demo/

pei asciinema rec -c "./demo.sh" asciinema

pei asciinema play asciinema