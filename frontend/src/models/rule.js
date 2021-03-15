export default class Rule {
  constructor(type, address, protocol, ports, ttl, comment) {
    this.type = type;
    this.address = address;
    this.protocol = protocol;
    this.ports = ports;
    this.ttl = ttl;
    this.comment = comment;
  }
}