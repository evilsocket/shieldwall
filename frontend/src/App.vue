<template>
  <div id="app">
    <nav class="navbar navbar-expand navbar-dark bg-dark">
      <a class="navbar-brand" href="/">shieldwall</a>

      <div class="navbar-nav mr-auto">
        <li class="nav-item">
          <router-link to="/docs" class="nav-link">Docs</router-link>
        </li>
      </div>

      <div v-if="!currentUser && !loggedIn" class="navbar-nav ml-auto">
        <li class="nav-item">
          <router-link to="/register" class="nav-link">
            <font-awesome-icon icon="user" /> Sign Up
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/login" class="nav-link">
            <font-awesome-icon icon="sign-in-alt" /> Login
          </router-link>
        </li>
      </div>

      <div v-if="loggedIn && currentUser" class="navbar-nav ml-auto">
        <li class="nav-item">
          <router-link to="/agents" class="nav-link">
            <font-awesome-icon icon="desktop" />
            Agents
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/profile" class="nav-link">
            <font-awesome-icon icon="user" />
            Profile
          </router-link>
        </li>
        <li class="nav-item">
          <a class="nav-link" href @click.prevent="logOut">
            <font-awesome-icon icon="sign-out-alt" />Logout
          </a>
        </li>
      </div>

      <div v-if="currentUser && !loggedIn" class="navbar-nav ml-auto">
        <li class="nav-item">
          <a class="nav-link" href @click.prevent="logOut">
            <font-awesome-icon icon="sign-out-alt" />Logout
          </a>
        </li>
      </div>
    </nav>

    <div class="container-fluid">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  mounted() {
    window.console.log("app.loggedIn = ",this.loggedIn);
    window.console.log("app.currentUser = ",this.currentUser);
  },
  computed: {
    loggedIn() {
      return this.$store.state.auth.status.loggedIn;
    },

    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  methods: {
    logOut() {
      this.$store.dispatch('auth/logout');
      this.$router.push('/login');
    }
  }
};
</script>
