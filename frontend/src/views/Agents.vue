<template>
  <div class="agents-container container-fluid table-responsive">

    <a class="btn btn-sm btn-success"
       style="margin-bottom: 10px"
       :href="'/#/agents/new'"
       v-if="agents.length"
    >
      new agent
    </a>

    <br/>

    <div class="jumbotron" v-if="!agents.length">
      No agents yet, <a href="/#/agents/new">create the first one</a>!
    </div>

    <div v-if="message" class="alert alert-success" role="alert">{{ message }}</div>
    <div v-if="error" class="alert alert-danger" role="alert">{{ error }}</div>

    <table class="table table-sm table-striped table-hover" v-if="agents.length">
      <thead class="thead-dark">
      <tr>
        <th scope="col">Created</th>
        <th scope="col" class="fit">Last Seen</th>
        <th scope="col">Name</th>
        <th scope="col">Address</th>
        <th scope="col">Version</th>
        <th scope="col">Interface</th>
        <th scope="col">Resources</th>
        <th scope="col">Rules</th>
        <th scope="col"></th>
      </tr>
      </thead>

      <tbody>
      <tr v-for="agent in agents" :key="agent.name">
        <td class="fit">
          <small>
            <timeago :datetime="agent.created_at" :auto-update="60"></timeago>
          </small>
        </td>
        <td class="fit">
          <small>
            <timeago :datetime="agent.seen_at" :auto-update="60"
                     v-if="agent.seen_at.indexOf('0001-01-01T00') == -1"></timeago>
            <span v-if="agent.seen_at.indexOf('0001-01-01T00') != -1">never</span>
          </small>
        </td>
        <td>
          <a :href="'/#/agent/' + agent.id">
            {{ agent.name }}
          </a>
        </td>
        <td class="fit">
          <small v-if="agent.address">{{ agent.address }}</small>
          <small v-if="!agent.address" class="text-muted">not seen yet</small>
        </td>
        <td class="fit">
          <small v-if="agent.user_agent">{{ agent.user_agent.replace('ShieldWall Agent ', '') }}</small>
          <small v-if="!agent.user_agent" class="text-muted">not seen yet</small>
        </td>
        <td class="fit">
          <small v-if="agent.active_interface">{{ agent.active_interface }}</small>
          <small v-if="!agent.active_interface" class="text-muted">-</small>
        </td>
        <td class="fit">
          <span v-if="agent.resources">
            <small class="resource-badge" :class="getCpuClass(agent.resources.cpu_percent)">
              CPU {{ formatPercent(agent.resources.cpu_percent) }}
            </small>
            <small class="resource-badge" :class="getMemClass(agent.resources.memory_percent)">
              MEM {{ formatPercent(agent.resources.memory_percent) }}
            </small>
          </span>
          <small v-if="!agent.resources" class="text-muted">-</small>
        </td>
        <td class="fit">
          <span v-if="agent.rules.length" class="badge badge-info">{{ agent.rules.length }}</span>
          <small v-if="!agent.rules.length" class="text-muted">none</small>
        </td>
        <td class="fit">
          <a class="btn btn-sm btn-danger" href="#" v-on:click="handleAgentDelete(agent)">
            x
          </a>
        </td>
      </tr>
      </tbody>
    </table>

  </div>
</template>

<script>
import UserService from '../services/user.service';

export default {
  name: 'Agents',

  data() {
    return {
      agents: [],
      message: '',
      error: '',
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
      return;
    }

    UserService.getUserAgents().then(
        response => {
          this.agents = response.data;
        },
        error => {
          this.error =
              (error.response && error.response.data && error.response.data.message) ||
              error.message ||
              error.toString();
          if(error.response.status === 401) {
            this.$store.dispatch('auth/logout');
            this.$router.push('/login');
          }
        }
    );
  },

  methods: {
    formatPercent(value) {
      if (value === undefined || value === null) return '-';
      return value.toFixed(1) + '%';
    },

    getCpuClass(value) {
      if (value === undefined || value === null) return '';
      if (value > 80) return 'resource-high';
      if (value > 50) return 'resource-medium';
      return 'resource-low';
    },

    getMemClass(value) {
      if (value === undefined || value === null) return '';
      if (value > 80) return 'resource-high';
      if (value > 50) return 'resource-medium';
      return 'resource-low';
    },

    handleAgentDelete(agent) {
      if (confirm("Are you sure you want to delete " + agent.name + ' ?')) {

        UserService.deleteAgent(agent.id).then(
            () => {
              for( var i in this.agents ) {
                if(this.agents[i].id === agent.id) {
                  this.agents.splice(i, 1);
                  break;
                }
              }
            },
            error => {
              this.error =
                  (error.response && error.response.data && error.response.data.error) ||
                  error.message ||
                  error.toString();
            }
        );
      }
    }
  }
};
</script>

<style scoped>
.agents-container {
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 25px;
}

.table td.fit,
.table th.fit {
  white-space: nowrap;
  width: 1%;
}

.resource-badge {
  display: inline-block;
  padding: 2px 6px;
  border-radius: 3px;
  margin-right: 4px;
  font-size: 0.75rem;
}

.resource-low {
  background-color: #28a745;
  color: white;
}

.resource-medium {
  background-color: #ffc107;
  color: black;
}

.resource-high {
  background-color: #dc3545;
  color: white;
}
</style>