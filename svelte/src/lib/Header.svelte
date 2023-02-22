<script lang="ts">
	// import { page } from '$app/stores';
	import logo from '@/assets/images/svelte-logo.svg';
	import github from '@/assets/images/github.svg';
	import login from '@/assets/images/login.png';
    import { push, replace } from 'svelte-spa-router';
    import { onMount } from 'svelte';

	export let TagUser:Boolean = true
	export let TagFile:Boolean = false
	export let TagChat:Boolean = false

	onMount(()=>{
		if (location.hash === '#/userManager'){
			TagUser = true;
			TagFile = false;
			TagChat = false
		}else if (location.hash === '#/fileManager'){
			TagChat = false;
			TagFile = true;
			TagUser = false;
		}else if (location.hash === '#/chat'){
			TagChat = true;
			TagFile = false;
			TagUser = false;
		}
	})

	const Navigate = (locationx:any)=>{
		replace(locationx)
		if (locationx === '/userManager'){
			TagUser = true;
			TagFile = false;
			TagChat = false
		}else if (locationx === '/fileManager'){
			TagChat = false;
			TagFile = true;
			TagUser = false;
		}else if (locationx === '/chat'){
			TagChat = true;
			TagFile = false;
			TagUser = false;
		}
	}
</script>

<header>
	<div class="corner">
		<a href="https://kit.svelte.dev">
			<img src={logo} alt="SvelteKit" />
		</a>
	</div>

	<nav>
		<svg viewBox="0 0 2 3" aria-hidden="true">
			<path d="M0,0 L1,2 C1.5,3 1.5,3 2,3 L2,0 Z" />
		</svg>
		<ul>
			<li class:active={TagUser}>
				<div on:click={()=>Navigate('/userManager')}>User</div>
			</li>
			<li class:active={TagChat}>
				<!-- <a href="/#/chat">Chat</a> -->
				<div on:click={()=>Navigate('/chat')}>Chat</div>
			</li>
			<li class:active={TagFile}>
				<!-- <a href="/#/fileManager">File</a> -->
				<div on:click={()=>Navigate('/fileManager')}>File</div>
			</li>
		</ul>
		<svg viewBox="0 0 2 3" aria-hidden="true">
			<path d="M0,0 L0,3 C0.5,3 0.5,3 1,2 L2,0 Z" />
		</svg>
	</nav>

	<div class="corner">
		<a href="/#/login">
			<img src={login} alt="Login" />
			<p>登 录</p>
		</a>
	</div>
</header>

<style>
	header {
		display: flex;
		justify-content: space-between;
	}

	.corner {
		width: 5em;
		height: 5em;
	}

	.corner a {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
		text-align: center;
		color: black;
		/* margin-left: -20px; */
	}

	.corner img {
		width: 2em;
		height: 2em;
		object-fit: contain;
	}
	.corner p{
		width: 3em;
		height: 1em;
	}

	nav {
		display: flex;
		justify-content: center;
		--background: rgba(255, 255, 255, 0.7);
	}

	svg {
		width: 2em;
		height: 3em;
		display: block;
	}

	path {
		fill: var(--background);
	}

	ul {
		position: relative;
		padding: 0;
		margin: 0;
		height: 3em;
		display: flex;
		justify-content: center;
		align-items: center;
		list-style: none;
		background: var(--background);
		background-size: contain;
	}

	li {
		position: relative;
		height: 100%;
	}

	li.active::before {
		--size: 6px;
		content: '';
		width: 0;
		height: 0;
		position: absolute;
		top: 0;
		left: calc(50% - var(--size));
		border: var(--size) solid transparent;
		border-top: var(--size) solid var(--color-theme-1);
	}

	nav div {
		display: flex;
		height: 100%;
		align-items: center;
		padding: 0 0.5rem;
		color: var(--color-text);
		font-weight: 700;
		font-size: 0.8rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		text-decoration: none;
		transition: color 0.2s linear;
	}

	a:hover {
		color: var(--color-theme-1);
	}
</style>
