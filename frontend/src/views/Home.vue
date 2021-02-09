<template>
  <div class="home-container container-fluid">
    <div class="jumbotron">
      <h1 class="display-4">
        <img src="/logo.png" width="50px"/>
        shieldwall!
      </h1>
      <br/>
      <p class="lead">
        Shieldwall is a free service that helps you secure your most private servers.
      </p>
      <hr class="my-4">
      <p>Every agent connects to our API
        and fetch a list of the IP addresses you allowed to navigate the host and block the inbound traffic from
        every other source.</p>
      <p class="lead" v-if="!currentUser">
        <br/>
        <a class="btn btn-primary btn" href="/register" role="button">
          Create a Free Account
        </a>
      </p>
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
</style>