import {createWebHistory, createRouter} from 'vue-router'

const home = () => import('./components/home-app.vue')
const user = () => import('./components/user-app.vue')
const module = () => import('./components/module-app.vue')
const task = () => import('./components/task-app.vue')
const space = () => import('./components/space-app.vue')
const rawtext = () => import('./components/rawtext-app.vue')
const dashboard = () => import('./components/dashboard-app.vue')
const database = () => import('./components/database-app.vue')
const rawdata = () => import('./components/rawdata-app.vue')
const problem = () => import('./components/problem-app.vue')
const notification = () => import('./components/notification-app.vue')

const routes = [
    {
        path: "/",
        name: "Home",
        component: home,
    },
    {
        path: "/user",
        name: "Пользователи",
        component: user,
    },
    {
        path: "/rawdata",
        name: "Данные по метрикам",
        component: rawdata,
    },
    {
        path: "/module",
        name: "Модули",
        component: module,
    },
    {
        path: "/task",
        name: "Задачи",
        component: task,
    },
    {
        path: "/space",
        name: "Пространства (клиенты)",
        component: space,
    },
    {
        path: "/rawtext",
        name: "Текстовые данные",
        component: rawtext,
    },
    {
        path: "/dashboard",
        name: "Дашборды",
        component: dashboard,
    },
    {
        path: "/problem",
        name: "Проблемы",
        component: problem,
    },
    {
        path: "/database",
        name: "Состояние БД",
        component: database,
    },
    {
        path: "/notification",
        name: "Уведомления",
        component: notification,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;