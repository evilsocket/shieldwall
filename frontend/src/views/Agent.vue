<template>
  <div class="agent-container container-fluid">

    <h2>{{ editing ? 'Edit' : 'New' }} Agent</h2>

    <br/>

    <form name="form" @submit.prevent="handleFormSubmit">
      <div class="form-group">
        <label for="name"><strong>Name</strong></label>
        <input
            v-model="agent.name"
            v-validate="'required|min:3|max:50'"
            type="text"
            class="form-control"
            placeholder="Agent name ..."
            name="name"
            id="name"
        />
        <div
            v-if="submitted && errors.has('name')"
            class="alert-danger"
        >{{ errors.first('name') }}
        </div>
      </div>

      <div class="form-group">
        <label><strong>Alert</strong></label>
        <br>

        <div class="form-group form-inline" v-if="editing">

          <label for="alert_after">
           If this agent was not active for the last
          </label>
          <select class="form-control form-control-sm" v-model.number="agent.alert_after" type="number" id="alert_after"
                  name="alert_after"
                  style="margin-left: 5px">
            <option :selected="agent.alert_after == 0" value="0">never (alert disabled)</option>
            <option :selected="agent.alert_after == 300" value="300">5 minutes</option>
            <option :selected="agent.alert_after == 600" value="600">10 minutes</option>
            <option :selected="agent.alert_after == 1800" value="1800">30 minutes</option>
            <option :selected="agent.alert_after == 3600" value="3600">1 hour</option>
            <option :selected="agent.alert_after == 10800" value="10800">3 hours</option>
            <option :selected="agent.alert_after == 21600" value="21600">6 hours</option>
          </select>
          <div
              v-if="submitted && errors.has('alert_at')"
              class="alert-danger"
          >{{ errors.first('alert_at') }}
          </div>

          <label for="alert_period" style="margin-left: 5px">
            send an email alert
          </label>
          <select class="form-control form-control-sm" v-model.number="agent.alert_period" type="number" id="alert_period"
                  name="alert_period"
                  style="margin-left: 5px">
            <option :selected="agent.alert_period == 0" value="0">once</option>
            <option :selected="agent.alert_period == 1800" value="1800">every 30 minutes</option>
            <option :selected="agent.alert_period == 3600" value="1800">every hour</option>
            <option :selected="agent.alert_period == 10800" value="10800">every 3 hours</option>
            <option :selected="agent.alert_period == 21600" value="21600">every 6 hours</option>
          </select>
          <div
              v-if="submitted && errors.has('alert_period')"
              class="alert-danger"
          >{{ errors.first('alert_period') }}
          </div>

        </div>
      </div>

      <span v-if="editing">
        <div class="form-group">
          <label for="token"><strong>Token</strong></label>
          <input
              v-model="agent.token"
              type="text"
              class="form-control"
              name="token"
              id="token"
              readonly="true"
          />
          <div
              v-if="submitted && errors.has('token')"
              class="alert-danger"
          >{{ errors.first('token') }}
          </div>
        </div>

        <div class="form-group">
          <label for="token"><strong>agent.yaml</strong></label>
          <small style="margin-left: 10px">
            Use this as the agent configuration file.
          </small>
          <pre class="config-file form-control">
# where to store the lists
data: '/var/lib/shieldwall/'
# path to iptables
iptables: '/sbin/iptables'
# check for newer versions and self update the agent
update: true

# api configuration
api:
  # api server to use
  server: 'https://shieldwall.me'
  # authentication token
  token: '{{ agent.token }}'
  # api polling period in seconds
  period: 10
  # api timeout in seconds or 0 for no timeout
  timeout: 0

# list of ip addresses to always allow just in case
allow:
  - '127.0.0.1'

# log dropped packets to syslog
drops:
  log: true
  limit: '10/min'
  prefix: 'shieldwall-dropped'
  level: 4</pre>
          <div
              v-if="submitted && errors.has('token')"
              class="alert-danger"
          >{{ errors.first('token') }}
          </div>
        </div>
      </span>

      <div class="form-group table-responsive">
        <h3>Rules
          <a class="btn btn-sm btn-success"
             href="#"
             v-on:click.prevent="handleRuleAdd()">
            +
          </a>

          <a class="btn btn-sm"
             href="#"
             title="Create rules for CloudFlare subnets."
             v-on:click.prevent="handleCf()">
            <img style="height: 60px; margin-top: 4px" src="/cf.png"/>
          </a>

          <span class="float-right d-inline">
            <small style="font-size: 0.9rem">
              <a href="https://github.com/evilsocket/shieldwall/wiki/Rules" target="_blank">Help</a>
            </small>
          </span>
        </h3>

        <table class="table table-sm table-striped table-hover" id="rules" v-if="agent.rules.length">
          <thead class="thead-dark">
          <tr>
            <th scope="col">Type</th>
            <th scope="col">Address</th>
            <th scope="col">Proto</th>
            <th scope="col">Ports</th>
            <th scope="col">Expires</th>
            <th scope="col">Comment</th>
            <th scope="col"></th>
          </tr>
          </thead>

          <tbody>
          <tr v-for="(rule, index) in agent.rules" :key="`rule-${index}`">

            <td class="fit input-group-sm">
              <select class="form-control" v-model="rule.type">
                <option :selected="rule.type == 'allow'" value="allow">Allow</option>
                <option :selected="rule.type == 'block'" value="block">Block</option>
              </select>
            </td>

            <td class="fit input-group-sm">
              <input
                  v-model="rule.address"
                  v-validate="'required'"
                  type="text"
                  class="form-control"
                  name="address"
              />
            </td>

            <td class="fit input-group-sm">
              <select class="form-control" v-model="rule.protocol">
                <option :selected="rule.protocol == 'tcp'" value="tcp">TCP</option>
                <option :selected="rule.protocol == 'udp'" value="udp">UDP</option>
                <option :selected="rule.protocol == 'all'" value="all">All</option>
              </select>
            </td>

            <td class="fit input-group-sm">
              <input
                  v-model="rule.ports"
                  v-validate="'required'"
                  type="text"
                  class="form-control"
                  name="ports"
              />
            </td>

            <td class="fit input-group-sm">
              <select class="form-control" v-model.number="rule.ttl" type="number">
                <option :selected="rule.ttl == 0" value=0>Never</option>
                <option :selected="rule.ttl == 3" value=3>3 Seconds</option>
                <option :selected="rule.ttl == 300" value=300>5 Minutes</option>
                <option :selected="rule.ttl == 600" value=600>10 Minutes</option>
                <option :selected="rule.ttl == 900" value=900>15 Minutes</option>
                <option :selected="rule.ttl == 1800" value=1800>30 Minutes</option>
                <option :selected="rule.ttl == 3600" value=3600>1 Hour</option>
                <option :selected="rule.ttl == 43200" value=43200>12 Hours</option>
                <option :selected="rule.ttl == 86400" value=86400>24 Hours</option>
              </select>
            </td>
            <td class="fit input-group-sm">
              <input
                  v-model="rule.comment"
                  type="text"
                  class="form-control input-small"
                  name="comment"
              />
            </td>

            <td class="fit input-group-sm">
              <a class="btn btn-sm btn-danger" href="#" v-on:click="handleRuleDelete(agent.rules.indexOf(rule))">
                x
              </a>
            </td>
          </tr>
          </tbody>
        </table>

      </div>

      <div class="form-group">
        <button class="btn btn-primary btn-block" :disabled="loading">
          <span v-show="loading" class="spinner-border spinner-border-sm"></span>
          <span>{{ editing ? 'Save' : 'Create' }}</span>
        </button>
      </div>
      <div class="form-group">
        <div v-if="message" class="alert alert-success" role="alert">{{ message }}</div>
        <div v-if="error" class="alert alert-danger" role="alert">{{ error }}</div>
      </div>

    </form>

  </div>
</template>


<script>
import UserService from '../services/user.service';
import Agent from "@/models/agent";
import Rule from "@/models/rule";

export default {
  name: 'Agent',

  data() {
    return {
      agent: new Agent('', [new Rule(
          "allow",
          this.$store.state.auth.user.address,
          "all",
          ["1:65535"],
          0
      )]),
      loading: false,
      submitted: false,
      editing: false,
      error: '',
      message: ''
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

    if (this.$route.params.id) {
      this.editing = true;
      UserService.getAgent(this.$route.params.id).then(
          response => {
            this.agent = response.data;
          },
          error => {
            this.error =
                (error.response && error.response.data && error.response.data.error) ||
                error.error ||
                error.toString();

            if (error.response.status === 401) {
              this.$store.dispatch('auth/logout');
              this.$router.push('/login');
            } else {
              this.$router.push('/agents');
            }
          }
      );
    }
  },

  methods: {
    handleFormSubmit() {
      // fix ports
      for (let i in this.agent.rules) {
        if (typeof (this.agent.rules[i].ports) == 'string') {
          this.agent.rules[i].ports = this.agent.rules[i].ports.split(',').map(e => e.trim());
        }
      }

      let method = this.editing ? UserService.updateAgent : UserService.createAgent;

      method(this.agent).then(
          response => {
            if (this.editing) {
              this.$router.push('/agents');
            } else {
              // force reload
              window.location.href = '/#/agent/' + response.data.id;
              window.location.reload();
            }
          },
          error => {
            this.error =
                (error.response && error.response.data && error.response.data.error) ||
                error.error ||
                error.toString();
          }
      );
    },

    handleRuleDelete(idx) {
      this.agent.rules.splice(idx, 1);
    },

    handleRuleAdd() {
      this.agent.rules.push(new Rule(
          "allow",
          this.$store.state.auth.user.address,
          "all",
          ["443", "80", "22"],
          43200,
          'Client IP on ' + (new Date())
      ));
    },

    handleCf() {
    UserService.getCloudFlareSubnets().then(
          response => {
            for(let i in response.data) {
              let subnet = response.data[i];
              this.agent.rules.push(new Rule(
                  "allow",
                  subnet,
                  "tcp",
                  ["443", "80"],
                  0,
                  "CloudFlare"
              ));
            }
          },
          error => {
            this.error =
                (error.response && error.response.data && error.response.data.error) ||
                error.error ||
                error.toString();
          }
      );
    }
  }
};
</script>

<style scoped>
.agent-container {
  padding: 20px 25px 30px;
  margin: 25px auto;
}

.config-file {
  font-family: monospace;
  height: fit-content;
  background-color: #212529;
  color: white;
}
</style>