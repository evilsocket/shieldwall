The shieldwall agent embraces the zero trust principle and instruments your server firewall to block inbound
connections from every IP on any port, by default. The website and API allow you to push policies
to your agents and temporarily allow you to reach a certain port from your IP.

**This project and its documentation are work in progress**

https://shieldwall.me/

### Compile and Run the Agent

Requires go and make, the agent runs on the server to protect.

    cp agent.example.yaml agent.yaml

Edit the configuration then build the agent:

    make agent

To install it as a systemd service:

    sudo make install_agent

To run it manually:

    sudo ./_build/shieldwall-agent -config agent.yaml

### Compile and Run your own API + Frontend

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

### TODO


* Precompile agent and make installation process easier.
* deb, rpm, etc for the agent?
* Agent self update?
* Upload rules json to bucket?
