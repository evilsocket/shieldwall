<template>
  <div class="agent-container container-fluid">

    <h2>{{ editing ? 'Edit' : 'Create' }} Agent</h2>

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
      </span>

      <div class="form-group">
        <h3>Rules
          <a class="btn btn-sm btn-success"
             href="#"
             v-on:click="handleRuleAdd()">
            +
          </a>
        </h3>

        <table class="table table-sm table-striped table-hover" id="rules" v-if="agent.rules.length">
          <thead class="thead-dark">
          <tr>
            <th scope="col">Type</th>
            <th scope="col">Address</th>
            <th scope="col">Proto</th>
            <th scope="col">Ports</th>
            <th scope="col"></th>
          </tr>
          </thead>

          <tbody>
          <tr v-for="(rule, index) in agent.rules" :key="`rule-${index}`">
            <td class="fit">
              <select class="form-control">
                <option :selected="rule.type == 'allow'" value="allow">Allow</option>
                <option :selected="rule.type == 'block'" value="block">Block</option>
              </select>
            </td>
            <td class="fit">
              <input
                  v-model="rule.address"
                  v-validate="'required'"
                  type="text"
                  class="form-control"
                  name="address"
              />
            </td>
            <td class="fit">
              <select class="form-control">
                <option :selected="rule.protocol == 'tcp'" value="tcp">TCP</option>
                <option :selected="rule.protocol == 'udp'" value="udp">UDP</option>
                <option :selected="rule.protocol == 'all'" value="all">All</option>
              </select>
            </td>
            <td class="fit">
              <input
                  v-model="rule.ports"
                  v-validate="'required'"
                  type="text"
                  class="form-control"
                  name="ports"
              />
            </td>

            <td class="fit">
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
          ["1:65535"]
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
    if (!this.currentUser) {
      this.$router.push('/login');
      return;
    }

    if(this.$route.params.id) {
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
            this.$router.push('/agents');
            return;
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
          () => {
            this.$router.push('/agents');
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
          ["443", "80"]
      ));
    }
  }
};
</script>

<style scoped>
.agent-container {
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 25px;
}

</style>