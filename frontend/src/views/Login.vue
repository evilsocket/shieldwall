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
          <label for="email">Email</label>
          <input
              v-model="user.email"
              v-validate="'required|email|max:50'"
              type="text"
              class="form-control"
              name="email"
          />
          <div
              v-if="errors.has('email')"
              class="alert alert-danger"
          >{{ errors.first('email') }}
          </div>
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input
              v-model="user.password"
              v-validate="'required'"
              type="password"
              class="form-control"
              name="password"
          />
          <div
              v-if="errors.has('password')"
              class="alert alert-danger"
          >{{ errors.first('email') }}
          </div>
        </div>
        <div class="form-group">
          <vue-recaptcha
              ref="recaptcha"
              v-if="!dev"
              @verify="onCaptchaVerified"
              @expired="onCaptchaExpired"
              sitekey="6LewaVIaAAAAAFn37I4KpU4OOKcjJBh_D0GXB8gC">
            <button class="btn btn-primary btn-block" :disabled="loading">
              <span v-show="loading" class="spinner-border spinner-border-sm"></span>
              <span>Login</span>
            </button>
          </vue-recaptcha>
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

    <center>
      Don't have an account? <br/>
      <a href="/#/register">Sign up for free!</a>
    </center>
  </div>
</template>

<script>
import User from '../models/user';
import VueRecaptcha from 'vue-recaptcha';
import {API_DEV} from '../services/api';

export default {
  name: 'Login',
  components: { VueRecaptcha },

  data() {
    return {
      user: new User('', ''),
      loading: false,
      message: '',
      dev: API_DEV
    };
  },
  computed: {
    loggedIn() {
      return this.$store.state.auth.status.loggedIn;
    }
  },
  created() {
    if (this.loggedIn) {
      this.$router.push('/');
    }
  },
  methods: {
    handleFormSubmit() {
      if(this.dev === true) {
        this.handleLogin();
      }
    },

    handleLogin() {
      this.loading = true;
      this.$validator.validateAll().then(isValid => {
        if (!isValid) {
          this.loading = false;
          return;
        }

        if (this.user.email && this.user.password) {
          this.$store.dispatch('auth/login', this.user).then(
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
    },

    onCaptchaExpired: function () {
      this.$refs.recaptcha.reset();
    },

    onCaptchaVerified: function () {
      this.handleLogin();
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