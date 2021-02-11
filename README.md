<p align="center">
  <img alt="ShieldWall" src="https://shieldwall.me/logo.png" height="140" />
  <p align="center">
    <a href="https://github.com/evilsocket/shieldwall/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/evilsocket/shieldwall.svg?style=flat-square"></a>
    <a href="https://github.com/evilsocket/shieldwall/blob/master/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
  </p>
</p>

ShieldWall embraces the zero-trust principle and instruments your server firewall to block inbound connections from every IP on any port, by default. The website allows you to push policies to your agents and temporarily unlock certain ports from your IP.

### Installing the agent

Download the [latest precompiled release](https://github.com/evilsocket/shieldwall/releases/latest) and install it 
with (adjust url to the latest version and your server architecture).

    mkdir /tmp/sw
    cd /tmp/sw
    wget https://github.com/evilsocket/shieldwall/releases/download/v1.0.0/shieldwall-agent_1.0.0_linux_arm64.tar.gz
    tar xvf shieldwall-agent_1.0.0_linux_arm64.tar.gz
    sudo ./install.sh

The agent is now installed as a systemd service, but it is not yet started nor enabled for autostart. You will first 
need to register an account on https://shieldwall.me/ and then edit the `/etc/shieldwall/config.yaml` configuration 
file, making sure it matches what you see on the agent page.

**It is very important that you double check the configuration before the next step, if the agent can't authenticate 
because of a wrong token, you will be locked out by the firewall and unable to log back.** 

You can now enable the service and start it. If configured so, it will automatically download and install its updates 
from github:

    sudo systemctl enable shieldwall-agent
    sudo service shieldwall-agent start    

Log into your https://shieldwall.me/ account to control the agent.

## Compile the agent from sources

Requires go and make, clone the repo and then:
    
    cd /path/to/repo
    make agent

To install it as a systemd service:

    sudo make install_agent

To run it manually:

    sudo ./_build/shieldwall-agent -config agent.yaml

## Compile and run your own API + frontend

Requires go and make. The API needs to be hosted on an IP that the agents can reach.

    cp api.example.yaml api.yaml

Edit both to your needs, the API requires a database that can be started with:

    cp database.example.env database.env
    docker-compose up

One the database is running you can compile and start the API service:

    make api

To install it as a systemd service:

    sudo make install_api

To run it manually:

    ./_build/shieldwall-api -config api.yaml

## Notes on future ideas

* Upload rules json to bucket?

## License

Released under the GPL3 license.