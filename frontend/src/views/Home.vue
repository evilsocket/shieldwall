<template>
  <div class="home-container container-fluid">
    <div class="jumbotron">
      <h1 class="display-4">
        <img src="/logo.png" width="50px"/>
        shieldwall!
      </h1>
      <br/>
      <p class="lead">
        Shieldwall is a free service you can use in order to secure your most private servers. Every agent connects to our API and fetch a list of the IP addresses you allowed to navigate the host and block the inbound traffic from every other source.
      </p>
      <hr class="my-4">
      <p>Start using it for free!</p>
      <p class="lead">
        <a class="btn btn-primary btn-lg" href="/register" role="button">
          Create an Account
        </a>
      </p>
    </div>
  </div>
</template>

<script>
// import UserService from '../services/user.service';

export default {
  name: 'Home',
  data() {
    return {
      content:
          "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
    };
  },
  mounted() {
    /*
    UserService.getPublicContent().then(
      response => {
        this.content = response.data;
      },
      error => {
        this.content =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString();
      }
    );
     */
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