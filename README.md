The shieldwall agent embraces the zero trust principle and instruments your server firewall to block inbound
connections from every IP on any port, by default. The website and API allow you to push policies
to your agents and temporarily allow you to reach a certain port from your IP.

### Compile and Run the Agent

Requires go and make, the agent runs on the server to protect.

    cp agent.example.yaml agent.yaml

Edit the configuration then build the agent:

    make agent

Run the agent with:

    sudo ./_build/shieldwall-agent -config agent.yaml

### Compile and Run your own API

Requires go and make, the API needs to be hosted on an IP that the agents can reach.

    cp api.example.yaml api.yaml
    cp database.example.env database.env

Edit both to your needs, the API requires a database that can be started with:

    docker-compose up

One the database is running you can compile and start the API service:

    make api
    ./_build/shieldwall-api -config api.yaml

### Compile and Run the Frontend
    
Requires npm:

    cd frontend
    npm install
    npm run serve

TODO: production mode