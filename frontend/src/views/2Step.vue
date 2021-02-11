<template>
  <div class="col-md-12">
    <div class="card card-container">
      <img
          id="profile-img"
          src="/logo.png"
          class="profile-img-card"
      />

      <form name="form" @submit.prevent="handleFormSubmit">
        <div class="form-group">
          <label for="code" style="text-align: center">Verification Code</label>
          <input
              v-validate="'required'"
              type="text"
              class="form-control form-control-lg"
              name="code"
              id="code"
              v-model="code"
          />
          <div
              v-if="errors.has('code')"
              class="alert alert-danger"
          >{{ errors.first('code') }}
          </div>
        </div>

        <div class="form-group">
          <button class="btn btn-primary btn-block" :disabled="loading" v-if="dev">
            <span v-show="loading" class="spinner-border spinner-border-sm"></span>
            <span>Login</span>
          </button>
        </div>

        <div class="form-group">
          <div v-if="message" class="alert alert-danger" role="alert">{{ message }}</div>
        </div>

      </form>
    </div>

  </div>
</template>

<script>
import {API_DEV} from '../services/api';

export default {
  name: 'LoginVerification',

  data() {
    return {
      loading: false,
      message: '',
      dev: API_DEV,
      code: '',
    };
  },
  computed: {
    loggedIn() {
      return this.$store.state.auth.status.loggedIn;
    },
    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  mounted() {
    if (this.loggedIn || !this.currentUser) {
      this.$router.push('/login');
    }

    window.console.log("2step.loggedIn = ",this.loggedIn);
    window.console.log("2step.currentUser = ",this.currentUser);
  },
  methods: {
    handleFormSubmit() {
      this.loading = true;
      this.$validator.validateAll().then(isValid => {
        if (!isValid) {
          this.loading = false;
          return;
        }

        if (this.code) {
          this.$store.dispatch('auth/step2', this.code).then(
              () => {
                this.$router.push('/agents');
              },
              error => {
                this.loading = false;
                this.message =
                    (error.response && error.response.data && error.response.data.error) ||
                    error.message ||
                    error.toString();
              }
          );
        }
      });
    }
  }
};
</script>

<style scoped>
label {
  display: block;
  margin-top: 10px;
}

.card-container.card {
  max-width: 350px !important;
  padding: 40px 40px;
}

.card {
  background-color: #f7f7f7;
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 50px;
  -moz-border-radius: 2px;
  -webkit-border-radius: 2px;
  border-radius: 2px;
  -moz-box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
  -webkit-box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
  box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
}

.profile-img-card {
  width: 96px;
  height: 96px;
  margin: 0 auto 10px;
  display: block;
}
</style>