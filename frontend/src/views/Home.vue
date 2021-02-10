<template>
  <div class="home-container container-fluid">
    <div class="jumbotron">
      <div style="text-align: center;">
        <h1 class="display-4">
          <img src="/logo.png" width="50px"/>
          shieldwall
          <small class="beta">beta</small>
        </h1>
        <br/>

        <p class="lead">
          Helps you secure your most private servers!
        </p>
      </div>

      <hr class="my-4">

      <p style="text-align: center;">
        The shieldwall agent embraces the zero trust principle and instruments your server firewall to block inbound
        connections from every IP on any port, by default. The website and API allow you to push policies
        to your agents and temporarily allow you to reach a certain port from your IP.
      </p>

      <div class="container">
        <div class="row">
          <div class="col">
            <div class="card text-center">
              <div class="card-header">
                Free API
              </div>
              <div class="card-body">
                <p class="card-text">Create a free account and use this server.</p>
                <a href="/#/register" class="btn btn-primary" v-if="!currentUser">Sign Up</a>
                &nbsp;
                <a href="/#/docs" class="btn btn-info" >Quick Start</a>
              </div>

            </div>
          </div>
          <div class="col">
            <div class="card text-center">
              <div class="card-header">
                Host it Yourself
              </div>
              <div class="card-body">
                <p class="card-text">Host the backend and API service yourself for a more discreet setup.</p>
                <a href="https://github.com/evilsocket/shieldwall" target="_blank" class="btn btn-success">
                  Code Repository
                </a>
              </div>

            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script>

export default {
  name: 'Home',
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  mounted() {

  }
};
</script>

<style scoped>
.home-container {
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 25px;
}

.beta {
  font-size: 15px;
  margin-left: -10px;
}

.btn {
  margin-top: 5px;
}

.card {
  margin-top: 10px;
}
</style>