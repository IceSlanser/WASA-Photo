import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from "../views/LoginView.vue";
import myProfileView from "../views/MyProfileView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/stream', component: HomeView},
		{path: '/users/:UID/profile', component: myProfileView},
		//{path: '/users/:UID/profile', component: GenericProfileView},
	]
})

export default router
