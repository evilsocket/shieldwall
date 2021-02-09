import axios from 'axios';
import {API_LOGIN_URL, API_REGISTER_URL, API_VERIFICATION_URL} from './api';

class AuthService {
  login(user) {
    return axios
      .post(API_LOGIN_URL, {
        email: user.email,
        password: user.password
      })
      .then(response => {
        if (response.data.token) {
          localStorage.setItem('user', JSON.stringify(response.data));
        }
        return response.data;
      });
  }

  logout() {
    localStorage.removeItem('user');
  }

  register(user) {
    return axios.post(API_REGISTER_URL, {
      email: user.email,
      password: user.password
    });
  }

  verify(code) {
    return axios.get(API_VERIFICATION_URL + '/' + code);
  }
}

export default new AuthService();
