#!/bin/bash
#
# Agent I/O preparation "Step 1".
# Creates a raw instance that can be saved as an image.
# Users and tools can repeatedly deploy this image to create agents.
#
# MUST be run as root.
#

#
# Part 1: Runtime environment (Go)
#
# Tested with ubuntu-12.04.2-server-amd64.iso. 
# Other Ubuntu and Debian installations may also work well.
#

apt-get update
apt-get install golang -y
apt-get install git -y
apt-get install bzr -y
apt-get install mercurial -y

#
# Part 2: Third-party Agent I/O components
#
# This installs nginx, mongodb, postfix, dovecot
# and any other third-party necessities not previously
# installed.
#

apt-get install nginx -y
apt-get install unzip -y

# mail setup
# aptitude remove exim4 && aptitude install postfix && postfix stop
# aptitude install dovecot-core dovecot-imapd
# aptitude install dovecot-common # might not be necessary

# this is installed earlier, but we don't need to keep it
# apt-get uninstall whoopsie

# get mongodb from the official mongodb repository, we want 2.6 or later
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | 
tee /etc/apt/sources.list.d/mongodb.list
apt-get update
apt-get install mongodb-org

#
# Part 3: Install CONTROL, the Agent I/O monitor
# 
# This web service provides an API for remotely administering an agent.
# It also includes tools that configure nginx as a container for agent HTTP services.
#

#adduser --system -disabled-login control 
adduser --system control 
addgroup control

# replace /home/control with the control repository
rm -rf /home/control
cp -r control /home/control

cd /home/control
mkdir -p nginx/logs
mkdir -p var
mkdir -p workers
mkdir -p go/bin go/pkg go/src
chown -R control /home/control
chgrp -R control /home/control

cp upstart/agentio-control.conf /etc/init
initctl start agentio-control

#
# That's it! Now save your image or remotely configure it (Step 2).
#
echo "Agent preparation is complete."

