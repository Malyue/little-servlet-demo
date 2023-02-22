import wrap from 'svelte-spa-router/wrap';

const routes = {
    '/test':wrap({
        asyncComponent: ()=> import('../views/test.svelte')
    })
}


export default routes;