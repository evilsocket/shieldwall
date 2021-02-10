<template>
  <div class="docs-container container-fluid">
    <h3>Agent Installation</h3>

    <p>For the time being compiling the agent from sources is the only option, soon statically precompiled binaries will
    be available, hang tight.</p>

    <p>The following steps require the go compiler and make:</p>

    <pre class="command">go get github.com/evilsocket/shieldwall
cd $GOPATH/src/github.com/evilsocket/shieldwall
make agent
sudo make install_agent</pre>

    <p>Now edit the file <strong>/etc/shieldwall/config.yaml</strong> and set the agent token. You can find the token
      by clicking on the agent name in the <a href="/#/agents">agents page</a> once you registered and logged in.
      <br/>
      <br/>
      The final step is restarting the service:</p>

    <pre class="command">sudo service shieldwall-agent restart</pre>

    <p>The agent will start polling the API and enabling the addresses and ports that you'll configure on this website. You can check the logs in /var/log//var/log/shieldwall.log</p>

    <br/>

    <h3>Hosting your own API</h3>

    <p>Please refer to the <a href="https://github.com/evilsocket/shieldwall/blob/main/README.md" target="_blank">
      README</a>.</p>

    <br/>

    <h3>Help</h3>

    <p>Found a bug? The documentation could be better? You have an idea for a new feature? Good news, this project is
      open source! You are more than welcome to <a href="https://github.com/evilsocket/shieldwall" target="_blank">fork the
        repository</a> and
    contribute with <a target="_blank" href="https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request">a pull request</a>.</p>
  </div>
</template>

<script>

export default {
  name: 'Docs',

  mounted() {
  }
};
</script>

<style scoped>
.docs-container {
  padding: 20px 25px 30px;
  margin: 25px auto;
}

.command {
  padding: 5px;
  height: fit-content;
  background-color: #212529;
  color: white;
}
</style>