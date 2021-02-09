export default class Rule {
  constructor(type, address, protocol, ports) {
    this.type = type;
    this.address = address;
    this.protocol = protocol;
    this.ports = ports;
  }
}