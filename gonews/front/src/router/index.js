import Vue from 'vue'
import VueRouter from 'vue-router'
import ItemListView from '../views/ItemList.vue'

Vue.use(VueRouter)

export default new VueRouter({
    mode: 'history',
    fallback: false,
    scrollBehavior: () => ({ y: 0 }),
    routes: [
        { path: '/news/:page(\\d+)?', component: ItemListView },
        { path: '/', component: ItemListView }
    ]
})