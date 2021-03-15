<template>
  <div id="app">
    <nav class="navbar navbar-expand navbar-dark bg-dark">
      <a class="navbar-brand" href="/">shieldwall</a>

      <div class="navbar-nav mr-auto">
        <li class="nav-item">
          <a href="https://www.evilsocket.net/2021/02/13/Hide-your-servers-in-plain-sight-presenting-ShieldWall/" target="_blank" class="nav-link">
            <font-awesome-icon icon="anchor" /> <span class="d-none d-sm-inline-block">Rationale</span>
          </a>
        </li>

        <li class="nav-item">
          <a href="https://github.com/evilsocket/shieldwall/wiki" target="_blank" class="nav-link">
            <font-awesome-icon icon="archive" /> <span class="d-none d-sm-inline-block">Docs</span>
          </a>
        </li>
      </div>

      <div v-if="!currentUser && !loggedIn" class="navbar-nav ml-auto">
        <li class="nav-item">
          <router-link to="/register" class="nav-link">
            <font-awesome-icon icon="user" /> <span class="d-none d-sm-inline-block">Sign Up</span>
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/login" class="nav-link">
            <font-awesome-icon icon="sign-in-alt" /> <span class="d-none d-sm-inline-block">Login</span>
          </router-link>
        </li>
      </div>

      <div v-if="loggedIn && currentUser" class="navbar-nav ml-auto">
        <li class="nav-item">
          <router-link to="/agents" class="nav-link">
            <font-awesome-icon icon="terminal" /> <span class="d-none d-sm-inline-block">Agents</span>
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/profile" class="nav-link">
            <font-awesome-icon icon="user" /> <span class="d-none d-sm-inline-block">Profile</span>
          </router-link>
        </li>
        <li class="nav-item">
          <a class="nav-link" href @click.prevent="logOut">
            <font-awesome-icon icon="sign-out-alt" /> <span class="d-none d-sm-inline-block">Logout</span>
          </a>
        </li>
      </div>

      <div v-if="currentUser && !loggedIn" class="navbar-nav ml-auto">
        <li class="nav-item">
          <a class="nav-link" href @click.prevent="logOut">
            <font-awesome-icon icon="sign-out-alt" /> <span class="d-none d-sm-inline-block">Logout</span>
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
