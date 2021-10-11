import { config } from "../config/config.js";

export function getEndpointPrefix() {
  let prefix = `http://${config.server.ip}:${config.server.port}`;
  console.log("[getEndpointPrefix::prefix]: ", prefix);
  return prefix;
}
