import { createApp, version } from 'vue'
import { createStore } from 'vuex'
import { createI18n } from 'vue-i18n'
import '@/style.css';

import axios from 'axios';

import login from '@/login-app.vue';
import router from '@/router.js';
import { messages } from '@/locales.js';
import { refreshAccessToken, installAuthenticationInterceptors } from '@/auth.js';

const versionApp = "2.0.0";

console.log("FrontEnd version: " + versionApp);
console.log("Vue version: " + version);

installAuthenticationInterceptors(() => window.location.reload());

const i18n = createI18n({
    legacy: false,
    locale: window.localStorage.getItem('locale') || 'ru',
    fallbackLocale: 'en',
    messages,
});

const store = createStore({
    state () {
        return {
            mode: window.localStorage.getItem('theme') === 'dark',
            locale: window.localStorage.getItem('locale') || 'ru'
        }
    },
    mutations: {
        changemode (state, value) {
            state.mode = value;
        },
        changelocale (state, value) {
            state.locale = value;
        },
    }
})

async function startApp() {
    const isAuthenticated = await restoreSession();
    const appComponent = isAuthenticated
        ? (await import('@/main-app.vue')).default
        : login;
    const vueApp = createApp(appComponent);
    vueApp.config.globalProperties.axios = axios;
    vueApp.config.globalProperties.$i18n = {
        get locale() {
            return i18n.global.locale.value;
        },
        set locale(value) {
            i18n.global.locale.value = value;
        },
        get availableLocales() {
            return i18n.global.availableLocales;
        },
    };
    const app = vueApp.use(store)
        .use(router)
        .use(i18n)
        .mount('#main');

    if (isAuthenticated) {
        app.versionInfo();
        app.userInfo();
    }

    setTimeout(function() {
        app.showPreloader = false;
    }, 600);
}

async function restoreSession() {
    try {
        await refreshAccessToken();
        return true;
    } catch {
        return false;
    }
}

startApp();
