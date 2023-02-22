
const modules:any = import.meta.glob('../modules/**/*.route.ts',{
    eager:true,
})

let routes = {};

Object.keys(modules).forEach((key)=>{
    const route = modules[key].default
    let r = Object.assign(route,routes)
    routes = r;
})

export default routes;