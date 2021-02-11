export default class User {
  constructor(email, password) {
    this.email = email;
    this.password = password;
    this.use_2fa = false;
    this.user_agent = '';
  }
}