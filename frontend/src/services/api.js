var apiDev, apiProto, apiHost, apiPort;

if (process.env.NODE_ENV === 'development') {
    apiDev = true;
    apiProto = 'http';
    apiHost = 'localhost';
    apiPort = 8080;
} else {
    apiDev = false;
    apiProto = 'https';
    apiHost = 'shieldwall.me';
    apiPort = 443;
}

export const API_DEV = apiDev;

export const API_PROTO = apiProto;
export const API_HOST = apiHost;
export const API_PORT = apiPort;

export const API_BASE_URL = API_PROTO + '://' + API_HOST + ':' + API_PORT + '/api/v1';

export const API_LOGIN_URL = API_BASE_URL + '/user/login';
export const API_2STEP_URL = API_BASE_URL + '/user/2step';
export const API_REGISTER_URL = API_BASE_URL + '/user/register';
export const API_VERIFICATION_URL = API_BASE_URL + '/user/verify';

export const API_USER_UPDATE_URL = API_BASE_URL + '/user/';

export const API_AGENTS_URL = API_BASE_URL + '/user/agents';
export const API_NEW_AGENT_URL = API_AGENTS_URL + '/new';

export const API_SUBNETS_URL = API_BASE_URL + '/subnets';
