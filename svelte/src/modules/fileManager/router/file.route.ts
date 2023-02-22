import File from '@/modules/fileManager/views/fileManage.svelte'
import wrap from 'svelte-spa-router/wrap';

const routes = {
    // '/fileManage':wrap({
    //     asyncComponent:()=>
    // }),
    '/fileManager':File
}


export default routes;