import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

// Router
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: () => import('./layouts/MainLayout.vue'),
      children: [
        { path: '', name: 'dashboard', component: () => import('./views/Dashboard.vue') },
        { path: 'nodes', name: 'nodes', component: () => import('./views/Nodes.vue') },
        { path: 'plugins', name: 'plugins', component: () => import('./views/Plugins.vue') },
        { path: 'playbooks', name: 'playbooks', component: () => import('./views/Playbooks.vue') },
        { path: 'users', name: 'users', component: () => import('./views/Users.vue') },
        { path: 'agents', name: 'agents', component: () => import('./views/Agents.vue') },
        { path: 'logs', name: 'logs', component: () => import('./views/Logs.vue') },
        { path: 'settings', name: 'settings', component: () => import('./views/Settings.vue') },
      ]
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('./views/Login.vue')
    }
  ]
})

// Pinia
const pinia = createPinia()

// App
const app = createApp(App)
app.use(pinia)
app.use(router)
app.mount('#app')