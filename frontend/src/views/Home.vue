<template>
  <div class="home-container container-fluid">
      <h3>Home</h3>
      <p>{{content}}</p>
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