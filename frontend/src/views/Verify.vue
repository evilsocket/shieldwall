<template>
    <div class="container">
      <div v-if="status" class="alert" role="alert">{{status}}</div>
      <div v-if="message" class="alert alert-success" role="alert">{{message}}</div>
      <div v-if="error" class="alert alert-danger" role="alert">{{error}}</div>
    </div>
</template>

<script>

export default {
  name: 'Agents',
  data() {
    return {
      status: '',
      message: '',
      error: ''
    };
  },
  mounted() {
    // this.message = "mounted";
    // alert("code = " + this.$route.params.code);
    this.status = 'checking code ...';

    this.$store.dispatch('auth/verify', this.$route.params.code).then(
        () => {
          this.status = '';
          this.message = 'account verified';
          this.$router.push('/login');
        },
        error => {
          console.log(error.response);
          this.status = '';
          this.error =
              (error.response && error.response.data && error.response.data.error) ||
              error.message ||
              error.toString();
        }
    );
  }
};
</script>

<style scoped>
.container {
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 50px;
}
</style>