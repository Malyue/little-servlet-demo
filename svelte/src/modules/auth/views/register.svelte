<script lang="ts">
    
</script>

<div class="registerBg">
    <div class="registerbox">
        <div class="registerhead">注 册</div>
        <!-- <div class="registerContent"></div> -->
    </div>
</div>

<style>
    .registerBg{
        /* display: flex; */
        width:100vw;
        height:100vh;
        background: url(../../../lib/images/svelte-welcome.png) no-repeat;
        background-size: 100%,100%;
        background-attachment: fixed;
    }
    .registerbox{
        position:absolute;
        top:30%;
        bottom:0;
        left:0;
        right: 0;
        margin: auto;
        width: 50%;
        height: 50%;
        border-radius: 5px;
        background: rgba(203,217,229,0.3);
        border:2px skyblue;
        display: flex;
        flex-flow: column;
    }
    .registerhead{
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 15%;
        font-size: x-large;
        font-weight: 600;
    }
</style>