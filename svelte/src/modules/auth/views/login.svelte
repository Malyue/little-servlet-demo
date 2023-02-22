<script lang="ts" type="module">
    import { push } from 'svelte-spa-router';

    import Textfield from '@smui/textfield';
    import Icon from '@smui/textfield/icon';
    import HelperText from '@smui/textfield/helper-text';
    import Button, { Label } from '@smui/button';
    import Snackbar, { Actions } from '@smui/snackbar';
    import IconButton from '@smui/icon-button';

    export let account:string = "";
    export let password:string = "";
    let snackbarWithClose: Snackbar;
    let snackbarSuccess: Snackbar;
    let snackbarWarning: Snackbar;
    const login = ()=>{
        if(account!=""&&password!=""){
            let url = "/api/login";
            //发送请求
            fetch(url,{
                method:"POST",
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'same-origin',
                headers:{
                    'Content-Type':'application/json',
                    Authorization:``
                },
                body:JSON.stringify({
                    "account":account,
                    "password":password,
                })
            }).then((res)=>{
                return res.json();
            }).then((res) =>{
                if (res.status == 0){
                    localStorage.setItem("userToken","")
                    //设置缓存
                    localStorage.setItem("userToken",res.data)
                    // console.log(res.data.token)
                    // console.log(localStorage.getItem("userToken")) 
                    snackbarSuccess.open()
                    push('/userManager')                                                                                                                                                                                                                                                                          
                }else{
                    if (res.status == 1003){
                        snackbarWarning.open()
                    }
                }
            })
        }else{
            // let id = $page.url.searchParams.get("id")
                    // console.log(id)
            snackbarWarning.open()
        }
    }
    const register = ()=>{
        push('/register')
        // goto('/auth/register')
    }
</script>

<svelte:head>
    <title>登录</title>
</svelte:head>

<Snackbar bind:this={snackbarSuccess} class="demo-success">
    <Label
      >登录成功!</Label
    >
    <Actions>
      <IconButton class="material-icons" title="Dismiss">close</IconButton>
    </Actions>
</Snackbar>


<Snackbar bind:this={snackbarWarning} class="demo-warning">
    <Label>请输入账号密码</Label>
    <Actions>
      <IconButton class="material-icons" title="Dismiss">close</IconButton>
    </Actions>
  </Snackbar>

<div class="loginBg">
    <div class="loginbox">
        <div class="loginhead"><span>登 录</span></div>
        <div class="loginUser">
            <!-- <input placeholder="   Please input your account..." bind:value="{account}" />  -->
            <Textfield variant="outlined" bind:value={account} label="Account" style="width:100%;background-color:rgba(255,255,255);!important">
                <HelperText slot="helper">Helper Text</HelperText>
            </Textfield>
        </div>
        <div class="loginPassword">
            <!-- <input placeholder="   Please input your password..." bind:value="{password}" type="password" /> -->
            <Textfield variant="outlined" bind:value={password} label="Password" style="width:100%;background-color:rgba(255,255,255)">
                <HelperText slot="helper">Helper Text</HelperText>
            </Textfield>
        </div>
        <div class="function">
            <!-- <div class="login subFunction" on:click={login}>登 录</div> -->
            <Button on:click={login} style="background-color:skyblue;width:30%;height:70%;color:black">
                <Label>登 录</Label>
            </Button>
            <!-- <div class="register subFunction" on:click={register}>注 册</div> -->
            <Button on:click={() => register} style="background-color:skyblue;width:30%;height:70%;color:black;margin-left:10%">
                <Label>注 册</Label>
            </Button>
        </div>
    </div>
</div>

<style>
    .loginBg{
        /* display: flex; */
        width:100vw;
        height:100vh;
        background: url('../../../assets/images/svelte-welcome.png') no-repeat;
        background-size: 100%,100%;
        background-attachment: fixed;
    }
    .loginbox{
        position:absolute;
        top:30%;
        bottom:0;
        left:0;
        right: 0;
        margin: auto;
        width: 50%;
        height: 50%;
        border-radius: 5px;
        /* background:rgba(255,255,255,0.3); */
        background: rgba(203,217,229,0.3);
        /* box-shadow: 3px 3px 6px 3px rgba(0,0,0,0.3); */
        border:2px skyblue;
        display: flex;
        flex-flow: column;
    }
    .loginhead {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 15%;
        font-size: x-large;
        font-weight: 600;
        /* background-color: aqua; */
    }
    .loginUser {
        width: 60%;
        height: 10%;
        margin-left: 20%;
        margin-top: 6%;
        flex:1;
    }
    .loginUserInput{
        width:100%;
        height: 100%;
        border-radius: 10px;
    }
    .loginPassword {
        width: 60%;
        height: 10%;
        margin-left: 20%;
        margin-top: 10%;
    }
    .loginPassword input{
        width:100%;
        height: 100%;
        border-radius: 10px;
    }
    .function{
        display: flex;
        flex-direction: row;
        /* background-color: aqua; */
        width: 70%;
        height: 20%;
        margin-left: 26%;
        margin-top: 7%;
    }
    .function .subFunction{
        display: flex;
        background-color: skyblue;
        width: 30%;
        height: 70%;
        justify-content: center;
        text-align: center;
        border-radius: 6px;
        align-items:center ;
    }
    .function .subFunction:hover{
        background-color: rgb(135,230,250);
    } 
    .function .register{
        margin-left: 10%;
    }
</style>

