export default class Agent {
  constructor(name, rules) {
    this.name = name;
    this.rules = rules;
    this.alert_after = 0;
    this.alert_period = 0;
  }
}