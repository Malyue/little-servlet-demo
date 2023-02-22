<script lang="ts">
  import { onMount } from 'svelte';
  import Router, { push } from 'svelte-spa-router';
  import routes from './router/index';
  import Header from '../src/lib/Header.svelte'


  // export let flag:boolean = false
  onMount(()=>{
    console.log(localStorage.getItem("userToken"))
    if(localStorage.getItem("userToken")===''){
      push('/login')
    }else{
      if (location.hash === ''){
        push('/userManager')
      }
    }
  })
</script>


<svelte:head>
    <!-- Fonts -->
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/icon?family=Material+Icons"
  />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,600,700"
  />

  <!-- Material Typography -->
  <link
    rel="stylesheet"
    href="https://unpkg.com/@material/typography@14.0.0/dist/mdc.typography.css"
  />

  <!-- SMUI -->
  <link rel="stylesheet" href="https://unpkg.com/svelte-material-ui/bare.css" />
</svelte:head>


<div>
  <main>
    <Router {routes} />
  </main>
</div>



<style>
</style>
