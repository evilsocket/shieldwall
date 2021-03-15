import axios from 'axios';
import authHeader from './auth-header';
import {API_AGENTS_URL, API_USER_UPDATE_URL, API_NEW_AGENT_URL, API_SUBNETS_URL} from "@/services/api";

class UserService {
  getUserAgents() {
    return axios.get(API_AGENTS_URL, { headers: authHeader() });
  }

  update(new_password, use_2fa) {
    return axios.post(API_USER_UPDATE_URL, {
      password: new_password,
      use_2fa: use_2fa
    }, { headers: authHeader() });
  }

  getAgent(id) {
    return axios.get(API_AGENTS_URL + '/' + id, { headers: authHeader() });
  }

  deleteAgent(id) {
    return axios.delete(API_AGENTS_URL + '/' + id, { headers: authHeader() });
  }

  createAgent(agent) {
    return axios.put(API_NEW_AGENT_URL, agent, { headers: authHeader() });
  }

  updateAgent(agent) {
    return axios.put(API_AGENTS_URL + '/' + agent.id, agent, { headers: authHeader() });
  }

  getCloudFlareSubnets() {
    return axios.get(API_SUBNETS_URL + '/cloudflare', { headers: authHeader() });
  }
}

export default new UserService();
