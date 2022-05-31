//import {createRouter, createWebHistory} from 'vue-router';
import {createRouter, createWebHashHistory} from 'vue-router';
import LoginPage from "@/components/login/LoginPage";
import DashboardPage from "@/components/Dashboard";
import HomePage from "@/components/home/HomePage";
import ControlPanel from "@/components/control/ControlPanel";
//storage
import publicStorage from "@/public.storage";
import utils from "@/utils";
import UserPanel from "@/components/users/UserPanel";
import OverviewPage from "@/components/overview/OverviewPage";
import CertPanel from "@/components/certification/CertPanel";

const routers = [
    {
        path: '/dashboard',
        name: 'dashboard',
        component: DashboardPage,
        meta: {
            requireLogin: true
        },
        children: [
            {
                path: '/dashboard/overview',
                name: 'dashboard_overview',
                component: OverviewPage,
                meta: {
                    requireLogin: true
                },
            },
            {
                path: '/dashboard/home',
                name: 'dashboard_home',
                component: HomePage,
                meta: {
                    requireLogin: true
                },
            },
            {
                path: '/dashboard/control',
                name: 'dashboard_control',
                component: ControlPanel,
                meta: {
                    requireLogin: true
                },
            },
            {
                path: '/dashboard/users',
                name: 'dashboard_users',
                component: UserPanel,
                meta: {
                    requireLogin: true
                },
            },
            {
                path: '/dashboard/cert',
                name: 'dashboard_cert',
                component: CertPanel,
                meta: {
                    requireLogin: true
                }
            }
        ]
    },
    {
        path: '/login',
        name: 'login',
        component: LoginPage,
        meta: {
            requireLogin: false
        },
    }
]

const router = createRouter({
    //history: createWebHistory(),
    history: createWebHashHistory(),
    routes: routers,
})

router.beforeEach((to, from, next) => {
    if (to.meta.requireLogin) {
        publicStorage.Load()
        if (publicStorage.User === undefined || publicStorage.User.token === "") {
            utils.Warning("禁止访问", "请登录后重试")
        } else {
            if (to.name !== "login") {
                next()
            }
        }
    } else {
        next()
    }
})

export default router