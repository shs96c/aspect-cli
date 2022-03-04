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

PROMPT_TIMEOUT=1

cd ..

# hide the evidence
clear

########################
# aspect demo
########################

p "aspect"
./aspect

sleep 5
clear

p "aspect info"
./aspect info

sleep 5
clear

p "aspect help"
./aspect help

sleep 5
clear

p "aspect version"
./aspect version

sleep 5
clear

pe "cd examples/fix-visibility"

p "aspect build ..."
(sleep 7; echo y) | ./aspect build ...

sleep 2

pe "git diff bar"

p "aspect build ..."
./aspect build ...

sleep 5
clear

pe "cd ../error-augmentor"

p "aspect build ..."
./aspect build ...

sleep 5
clear

# pe "cd ../predefined-queries"

# pe "(sleep 5; echo y) | ./aspect query"

# sleep 5
# clear
