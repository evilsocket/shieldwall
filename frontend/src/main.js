import Vue from 'vue';
import App from './App.vue';
import {router} from './router';
import store from './store';
import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import VeeValidate from 'vee-validate';
import Vuex from 'vuex';
import {library} from '@fortawesome/fontawesome-svg-core';

import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';
import {faDesktop, faHome, faSignInAlt, faSignOutAlt, faTerminal, faAnchor, faArchive, faUser, faUserPlus} from '@fortawesome/free-solid-svg-icons';
import VueTimeago from 'vue-timeago'

library.add(faHome, faUser, faDesktop, faUserPlus, faTerminal, faArchive, faAnchor, faSignInAlt, faSignOutAlt);

Vue.config.productionTip = false;

Vue.use(VeeValidate);
Vue.component('font-awesome-icon', FontAwesomeIcon);

Vue.use(Vuex);

Vue.use(VueTimeago, {
    name: 'Timeago', // Component name, `Timeago` by default
    locale: 'en' // Default locale
});


new Vue({
    router,
    store,
    render: h => h(App)
}).$mount('#app');
