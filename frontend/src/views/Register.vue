<template>
  <div class="col-md-12">
    <div class="card card-container">
      <img
          id="profile-img"
          src="/logo.png"
          class="profile-img-card"
      />
      <form name="form">
        <div v-if="!successful">
          <div class="form-group">
            <label for="email">Email</label>
            <input
                v-model="user.email"
                v-validate="'required|email|max:50'"
                type="email"
                class="form-control"
                name="email"
            />
            <div
                v-if="submitted && errors.has('email')"
                class="alert-danger"
            >{{ errors.first('email') }}
            </div>
          </div>
          <div class="form-group">
            <label for="password">Password</label>
            <input
                v-model="user.password"
                v-validate="'required|min:8|max:64'"
                type="password"
                class="form-control"
                name="password"
            />
            <div
                v-if="submitted && errors.has('password')"
                class="alert-danger"
            >{{ errors.first('password') }}
            </div>
          </div>

          <div class="form-group">
            <vue-recaptcha
                ref="recaptcha"
                @verify="onCaptchaVerified"
                @expired="onCaptchaExpired"
                sitekey="6LewaVIaAAAAAFn37I4KpU4OOKcjJBh_D0GXB8gC">
              <button class="btn btn-primary btn-block">Sign Up</button>
            </vue-recaptcha>
          </div>
        </div>
      </form>

      <div
          v-if="message"
          class="alert"
          :class="successful ? 'alert-success' : 'alert-danger'"
      >{{ message }}
      </div>
    </div>

    <center>
      Already have an account? <br/>
      <a href="/login">Login in here!</a>
    </center>
  </div>
</template>

<script>
import User from '../models/user';
import VueRecaptcha from 'vue-recaptcha';

export default {
  name: 'Register',
  components: { VueRecaptcha },

  data() {
    return {
      user: new User('', ''),
      submitted: false,
      successful: false,
      message: ''
    };
  },
  computed: {
    loggedIn() {
      return this.$store.state.auth.status.loggedIn;
    }
  },
  mounted() {
    if (this.loggedIn) {
      this.$router.push('/agents');
    }
  },
  methods: {
    handleRegister() {
      this.message = '';
      this.submitted = true;
      this.$validator.validate().then(isValid => {
        if (isValid) {
          this.$store.dispatch('auth/register', this.user).then(
              data => {
                this.message = data;
                this.successful = true;
              },
              error => {
                this.message =
                    (error.response && error.response.data && error.response.data.error) ||
                    error.message ||
                    error.toString();
                this.successful = false;
              }
          );
        }
      });
    },

    onCaptchaExpired: function () {
      this.$refs.recaptcha.reset();
    },

    onCaptchaVerified: function () {
      this.handleRegister();
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