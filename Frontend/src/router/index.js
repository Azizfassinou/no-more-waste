import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ClientDashboard from '../views/ClientDashboard.vue'
import MerchantDashboard from '../views/MerchantDashboard.vue'
import StaffDashboard from '../views/StaffDashboard.vue'
import AdminDashboard from '../views/AdminDashboard.vue'
import VolunteerDashboard from '../views/VolunteerDashboard.vue'


const routes = [
    {
        path: '/',
        name: 'Login',
        component: LoginView
    },
    {
        path: '/client/dashboard',
        name: 'ClientDashboard',
        component: ClientDashboard,
        meta: { requiresAuth: true, role: 'client' }
    },
    {
        path: '/merchant/dashboard',
        name: 'MerchantDashboard',
        component: MerchantDashboard,
        meta: { requiresAuth: true, role: 'merchant' }
    },
    {
        path: '/volunteer/dashboard',
        name: 'VolunteerDashboard',
        component: VolunteerDashboard,
        meta: { requiresAuth: true, role: 'volunteer' }
    },
    {
        path: '/staff/dashboard',
        name: 'StaffDashboard',
        component: StaffDashboard,
        meta: { requiresAuth: true, roles: ['staff', 'admin'] }
    },
    {
        path: '/admin/dashboard',
        name: 'AdminDashboard',
        component: AdminDashboard,
        meta: { requiresAuth: true, role: 'admin' }
    },
    {
        path: '/:catchAll(.*)',
        redirect: '/'
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    const userRole = localStorage.getItem('role')

    if (to.meta.requiresAuth) {
        if (!token) {
            next({ name: 'Login' })
        } else {
            if (to.meta.roles) {
                if (to.meta.roles.includes(userRole)) {
                    next()
                } else {
                    next({ name: 'Login' })
                }
            } else if (to.meta.role && to.meta.role !== userRole) {
                next({ name: 'Login' })
            } else {
                next()
            }
        }
    } else {
        next()
    }
})

export default router