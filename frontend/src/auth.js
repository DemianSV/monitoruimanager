import axios from 'axios';

let accessToken = null;
let refreshPromise = null;
let onAuthenticationExpired = null;
let authenticationExpired = false;

function isApiRequest(url = '') {
    try {
        const requestUrl = new URL(url, window.location.origin);
        return requestUrl.origin === window.location.origin && requestUrl.pathname.startsWith('/api/');
    } catch {
        return false;
    }
}

function setAccessToken(token) {
    if (typeof token !== 'string' || token.length === 0) {
        throw new Error('The authentication response did not contain an access token.');
    }
    accessToken = token;
    authenticationExpired = false;
}

export function clearAccessToken() {
    accessToken = null;
}

export async function refreshAccessToken() {
    if (!refreshPromise) {
        refreshPromise = axios.post('/auth/refresh', null, { withCredentials: true })
            .then((response) => {
                setAccessToken(response.data?.access_token);
                return accessToken;
            })
            .finally(() => {
                refreshPromise = null;
            });
    }
    return refreshPromise;
}

export async function login(username, password) {
    const form = new URLSearchParams();
    form.set('username', username);
    form.set('password', password);

    const response = await axios.post('/auth', form, { withCredentials: true });
    setAccessToken(response.data?.access_token);
}

export async function logout() {
    await axios.post('/auth/logout', null, { withCredentials: true });
    clearAccessToken();
}

export function installAuthenticationInterceptors(onExpired) {
    onAuthenticationExpired = onExpired;

    axios.interceptors.request.use((config) => {
        if (isApiRequest(config.url) && accessToken) {
            config.headers.Authorization = `Bearer ${accessToken}`;
        }
        return config;
    });

    axios.interceptors.response.use(
        (response) => response,
        async (error) => {
            const request = error.config;
            if (error.response?.status !== 401 || !request || !isApiRequest(request.url)) {
                return Promise.reject(error);
            }
            if (request._authRetry) {
                clearAccessToken();
                if (!authenticationExpired) {
                    authenticationExpired = true;
                    onAuthenticationExpired?.();
                }
                return Promise.reject(error);
            }

            request._authRetry = true;
            try {
                await refreshAccessToken();
                return axios.request(request);
            } catch (refreshError) {
                clearAccessToken();
                if (!authenticationExpired) {
                    authenticationExpired = true;
                    onAuthenticationExpired?.();
                }
                return Promise.reject(refreshError);
            }
        },
    );
}
