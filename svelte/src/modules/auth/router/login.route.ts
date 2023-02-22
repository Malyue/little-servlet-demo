import {wrap} from 'svelte-spa-router/wrap'
// import Login from 'src/modules/login/views/login.svelte'
// import Register from '../views/register.svelte'
import Login from '@/modules/auth/views/login.svelte'
import Register from '@/modules/auth/views/register.svelte'

const routes = {
    "/login":Login,
    '/register':Register
}


export default routes;