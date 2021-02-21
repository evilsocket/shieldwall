<template>
  <div class="profile-container container-fluid">

    <form name="form" @submit.prevent="handlePasswordChange">
      <div class="form-group">
        <label for="created_at"><strong>Account Creation</strong></label>
        <span
            class="form-control"
            name="created_at"
            id="created_at"
            readonly="true"
        >
          <timeago :datetime="currentUser.data.created_at " :auto-update="60"></timeago>
        </span>
      </div>
      <div class="form-group">
        <label for="email"><strong>Email</strong></label>
        <span
            class="form-control"
            name="email"
            id="email"
            readonly="true"
        >{{ currentUser.data.email }}</span>
      </div>
      <div class="form-group">
        <label for="address"><strong>Address</strong></label>
        <span
            class="form-control"
            name="address"
            id="address"
            readonly="true"
        >{{ currentUser.data.address }}</span>
      </div>
      <div class="form-group">
        <label for="password"><strong>New Password</strong></label>
        <input
            v-model="new_password"
            v-validate="'min:8|max:64'"
            type="password"
            class="form-control"
            name="password"
            id="password"
        />
        <div
            v-if="errors.has('password')"
            class="alert alert-danger"
        >{{ errors.first('password') }}
        </div>
      </div>

      <div class="form-group form-check">
        <input class="form-check-input" type="checkbox" value="1" id="use_2fa" v-model="use_2fa">
        <label class="form-check-label" for="use_2fa">
          Two-step verification via email
        </label>
      </div>

      <div class="form-group">
        <button class="btn btn-primary btn-block" :disabled="loading">
          <span v-show="loading" class="spinner-border spinner-border-sm"></span>
          <span>Save</span>
        </button>
      </div>
      <div class="form-group">
        <div v-if="message" class="alert alert-success" role="alert">{{message}}</div>
        <div v-if="error" class="alert alert-danger" role="alert">{{error}}</div>
      </div>

    </form>

  </div>
</template>

<script>
import UserService from '../services/user.service';

export default {
  name: 'Profile',

  data() {
    return {
      loading: false,
      message: '',
      error: '',
      use_2fa: this.$store.state.auth.user.data.use_2fa,
      new_password: ''
    };
  },

  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    }
  },

  mounted() {
    if (!this.$store.state.auth.status.loggedIn) {
      this.$router.push('/login');
    }
  },

  methods: {
    handlePasswordChange() {
      this.loading = true;
      this.$validator.validateAll().then(isValid => {
        if (!isValid) {
          this.loading = false;
          return;
        }

        if (this.new_password.length >= 8 || this.use_2fa !== this.currentUser.data.use_2fa) {
          UserService.update(this.new_password, this.use_2fa).then(
              response => {
                this.loading = false;
                this.message = 'Profile updated.';
                localStorage.setItem('user', JSON.stringify(response.data));
              },
              error => {
                this.loading = false;
                this.message =
                    (error.response && error.response.data && error.response.data.error) ||
                    error.message ||
                    error.toString();

                if(error.response.status === 401) {
                  this.$store.dispatch('auth/logout');
                  this.$router.push('/login');
                }
              }
          );
        } else {
          this.loading = false;
        }
      });
    }
  }
};
</script>

<style scoped>
.profile-container {
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 25px;
}
</style>