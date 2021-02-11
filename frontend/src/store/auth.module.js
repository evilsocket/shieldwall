import AuthService from '../services/auth.service';

const user = JSON.parse(localStorage.getItem('user'));
const initialState = user
  ? { status: { loggedIn: true }, user }
  : { status: { loggedIn: false }, user: null };

export const auth = {
  namespaced: true,
  state: initialState,
  actions: {
    login({ commit }, user) {
      return AuthService.login(user).then(
        user => {
          commit('loginSuccess', user);
          return Promise.resolve(user);
        },
        error => {
          commit('loginFailure');
          return Promise.reject(error);
        }
      );
    },
    step2({ commit }, code) {
      return AuthService.step2(code).then(
          user => {
            commit('step2Success', user);
            return Promise.resolve(user);
          },
          error => {
            commit('loginFailure');
            return Promise.reject(error);
          }
      );
    },
    logout({ commit }) {
      AuthService.logout();
      commit('logout');
    },
    register({ commit }, user) {
      return AuthService.register(user).then(
        response => {
          commit('registerSuccess');
          return Promise.resolve(response.data);
        },
        error => {
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    },
    verify({ commit }, code) {
      return AuthService.verify(code).then(
          response => {
            commit('verifySuccess');
            return Promise.resolve(response.data);
          },
          error => {
            commit('verifyFailure');
            return Promise.reject(error);
          }
      );
    }
  },
  mutations: {
    loginSuccess(state, user) {
      state.status.loggedIn = !user.data.use_2fa;
      state.user = user;
    },
    step2Success(state, user) {
      state.status.loggedIn = true;
      state.user = user;
    },
    loginFailure(state) {
      state.status.loggedIn = false;
      state.user = null;
    },
    logout(state) {
      state.status.loggedIn = false;
      state.user = null;
    },
    registerSuccess(state) {
      state.status.loggedIn = false;
    },
    registerFailure(state) {
      state.status.loggedIn = false;
    },
    verifySuccess() {

    },
    verifyFailure(state) {
      state.status.loggedIn = false;
      state.user = null;
    }
  }
};
