import axios from 'axios';
import {API_2STEP_URL, API_LOGIN_URL, API_REGISTER_URL, API_VERIFICATION_URL} from './api';
import authHeader from './auth-header';

class AuthService {
    login(user) {
        return axios
            .post(API_LOGIN_URL, {
                email: user.email,
                password: user.password
            })
            .then(response => {
                if (response.data.token) {
                    if (!response.data.data.use_2fa) {
                        window.console.log("saving user to user", response.data);
                        localStorage.setItem('user', JSON.stringify(response.data));
                    } else {
                        window.console.log("saving user to user2fa", response.data);
                        localStorage.setItem('user2fa', JSON.stringify(response.data));
                    }
                }
                return response.data;
            });
    }

    step2(code) {
        return axios
            .post(API_2STEP_URL, {
                code: code
            }, { headers: authHeader() })
            .then(response => {
                if (response.data.token) {
                    window.console.log("finalizing user to user", response.data);
                    localStorage.setItem('user', JSON.stringify(response.data));
                    localStorage.removeItem('user2fa');
                }
                return response.data;
            });
    }

    logout() {
        localStorage.removeItem('user');
        localStorage.removeItem('user2fa');
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
