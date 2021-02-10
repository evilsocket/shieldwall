import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import Docs from './views/Docs.vue';
import Login from './views/Login.vue';
import Register from './views/Register.vue';
import Verify from './views/Verify.vue';

Vue.use(Router);

export const router = new Router({
  mode: 'hash',
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/docs',
      name: 'docs',
      component: Docs
    },
    {
      path: '/login',
      component: Login
    },
    {
      path: '/register',
      component: Register
    },
    {
      path: '/verify/:code',
      component: Verify,
      params: true
    },
    {
      path: '/profile',
      name: 'profile',
      // lazy-loaded
      component: () => import('./views/Profile.vue')
    },
    {
      path: '/agents',
      name: 'agents',
      // lazy-loaded
      component: () => import('./views/Agents.vue')
    },
    {
      path: '/agents/new',
      name: 'agent-new',
      // lazy-loaded
      component: () => import('./views/Agent.vue')
    },
    {
      path: '/agent/:id',
      name: 'agent-edit',
      // lazy-loaded
      component: () => import('./views/Agent.vue')
    },
  ]
});

// router.beforeEach((to, from, next) => {
//   const publicPages = ['/login', '/register', '/home'];
//   const authRequired = !publicPages.includes(to.path);
//   const loggedIn = localStorage.getItem('user');

//   // trying to access a restricted page + not logged in
//   // redirect to login page
//   if (authRequired && !loggedIn) {
//     next('/login');
//   } else {
//     next();
//   }
// });
