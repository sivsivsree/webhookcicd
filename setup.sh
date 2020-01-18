#!/usr/bin/env bash

function setUpAWS() {
    if which aws2 >/dev/null; then
    echo "AWS cli is installed ✅  "
else
    echo "Installing AWS CLI 🚀  "
    unamestr=`uname`
    if [[ "$unamestr" == 'Linux' ]]; then
       echo "Install CLI for Linux"
       curl "https://d1vvhvl2y92vvt.cloudfront.net/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
       unzip awscliv2.zip
       sudo ./aws/install -b
    elif [[ "$unamestr" == 'Darwin' ]]; then
       echo "Install CLI for Darwin"
       curl "https://d1vvhvl2y92vvt.cloudfront.net/awscli-exe-macos.zip" -o "awscliv2.zip"
        unzip awscliv2.zip
        sudo ./aws/install -b
    fi

    echo "AWS CLI is installed ✅"
fi


}

function showSSHLoginInfo() {
    echo "Add your id_rsa key to github before clone"
    echo "------------------------------------------"
    cat ~/.ssh/id_rsa.pub
    echo ""
    echo "Copy this key and add to github"
    echo ""
}



showSSHLoginInfo
read -p "Did you copied? [Y/N] " -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    # do dangerous stuff
    setUpAWS
    aws2 configure
    $(aws2 ecr get-login --no-include-email --region us-east-2)
fi


