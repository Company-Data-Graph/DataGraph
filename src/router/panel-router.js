
import SomethingWentWrong from "@/pages/SomethingWentWrong";
import {createRouter, createWebHistory} from "vue-router";
import HomePage from "@/pages/HomePage";
import CompanyPage from "@/pages/CompanyPage";
import ProductPage from "@/pages/ProductPage.vue";
import LinePage from "@/pages/LinePage.vue";

const routes = [
    {
        path: '/',
        component: HomePage
    },
    {
        path: '/company/:id',
        component: CompanyPage
    },
    {
        path: '/product/:id',
        component: ProductPage
    },
    {
        path: '/line',
        component: LinePage
    },
    {
        path: '/:pathMatch(.*)*',
        component: SomethingWentWrong
    }
]

const router = createRouter({
    routes,
    history: createWebHistory()
})

export default router