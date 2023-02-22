<script lang="ts">
    import { onMount } from "svelte";
    import Header from '../../../lib/Header.svelte'

    export let num = 0
    // [{id:"0",account:"12314",name:"libai",role:"管理员",note:"111111111111"},{id:"1",account:"Mike",role:"用户",note:""}]
    export let peopleList:any = []
    export let checked:Array<boolean> = []
    export let checkedAll = false
    export let currentPage = 1
    export let addshow = false
    export let inputSearch:string = ''
    export let inputAccount:string = ''
    export let inputPassword:string = ''
    export let inputUserName:string = ''


    export function check(index:number){
        console.log(index)
        if (checked[index]===false){
            checked[index] = true
        }else{
            checked[index] = false
        }
        console.log(checked)
    }

    export function checkAll(){
        if (checkedAll === true) {
            //取消全选
            for (let i:number=0;i<checked.length;i++){
                checked[i] = false;
            }
            checkedAll = false
        }else if (checkedAll === false){
            //全选
            for (let i:number = 0;i<checked.length;i++){
                checked[i] = true;
            }
            checkedAll = true
        }
    }
    export function getLastUser(){
        if (currentPage !== 1){
            currentPage = currentPage - 1;
            getUser();
        }
    }
    export function getNextUser(){
        if (currentPage === Math.ceil(num/10)){
            return
        }
        //如果不是最后一页
        if (currentPage !== Math.ceil(num/10)){
            currentPage = currentPage + 1;
            getUser();
        }
    }
    export function getUser(){
        checkedAll = false
        checked = []
        let token =localStorage.getItem("userToken")||''
        let url = "/api/getUser"
        fetch(url,{
            method:"POST",
            mode:"cors",
            cache:'no-cache',
            credentials:'same-origin',
            headers:{
                'Content-Type':'application/json',
                'Authorization':token
            },
            body:JSON.stringify({
                "page":currentPage,
                "count":10
            })
        }).then((res)=>{
            return res.json()
        }).then((res)=>{
            peopleList = []
            num = res.data.count;
            if (res.data.userInfo.length != 0){
                peopleList = peopleList.concat(res.data.userInfo)
            }
            for (let i =0 ;i<peopleList.length;i++){
                checked.push(false)
                checked = checked
            }
        })
    };

    export function deleteUser(){
        let deleteUserList:any[] = []
        for(let i =0;i<checked.length;++i){
            if(checked[i] === true){
                // console.log(i)
                deleteUserList = [...deleteUserList,peopleList[i].id]
            }
        }
        if (deleteUserList.length === 0){
            alert("请选择要删除用户")
            return
        }
        let token =localStorage.getItem("userToken")||''
        let url = "/api/deleteUser"
        fetch(url,{
            method:"POST",
            mode:"cors",
            cache:'no-cache',
            credentials:'same-origin',
            headers:{
                'Content-Type':'application/json',
                'Authorization':token
            },
            body:JSON.stringify({
                "deleteUserList":deleteUserList,
                "count":10
            })
        }).then((res)=>{
            return res.json()
        }).then((res)=>{
            if (res.status == 0){
                currentPage = 1
                getUser();
                alert("删除成功")
            }else{
                alert("删除失败")
            }
        })
    };

    export function addShow(){
        addshow = true
    }

    export function searchUser(){
        if (inputSearch === ''){
            getUser()
            return
        }
        let token =localStorage.getItem("userToken")||''
        let url = "/api/getUser"
        fetch(url,{
            method:"POST",
            mode:"cors",
            cache:'no-cache',
            credentials:'same-origin',
            headers:{
                'Content-Type':'application/json',
                'Authorization':token
            },
            body:JSON.stringify({
                "page":1,
                "count":10,
                "search":inputSearch,
            })
        }).then((res)=>{
            return res.json()
        }).then((res)=>{
           if (res.status === 0){
                peopleList = []
                num = res.data.count;
                peopleList = peopleList.concat(res.data.userInfo)
                for (let i =0 ;i<peopleList.length;i++){
                    checked.push(false)
                    checked = checked
                }
                inputSearch = ''
           }
        })
    }

    export function addUserToDb(){
        if (inputAccount!==''&&inputPassword!==''&&inputUserName!==''){
            let token =localStorage.getItem("userToken")||''
            let url = "/api/register"
            fetch(url,{
                method:"POST",
                mode:"cors",
                cache:'no-cache',
                credentials:'same-origin',
                headers:{
                    'Content-Type':'application/json',
                    'Authorization':token
                },
                body:JSON.stringify({
                    "userName":inputUserName,
                    "account":inputAccount,
                    "password":inputPassword,
                })
            }).then((res)=>{
                return res.json()
            }).then((res)=>{
                if (res.status === 0){
                    alert("添加成功")
                    inputAccount = ''
                    inputPassword = ''
                    inputUserName = ''
                }else if (res.status === 1005){
                    alert("账号重复")
                }
            })
        }else{
            alert("请输入账号/密码/用户名")
        }
    }

    onMount(()=>{
        getUser()
    })
</script>

<Header />
{#if addshow === true}
    <div class="AddUser" on:click={()=>{addshow = false}}>
        <div class="AddUserBox" on:click|stopPropagation={()=>{}}>
            <div class="header">增加用户</div>
            <div class="addUserName addinput">
                <input placeholder="   请输入用户名..." bind:value={inputUserName}/>
            </div>
            <div class="addUserAccount addinput">
                <input placeholder="   请输入账号..." bind:value={inputAccount}/>
            </div>
            <div class="addUserPassword addinput">
                <input placeholder="   请输入密码..." bind:value={inputPassword}/>
            </div>
            <div class="function">
                <div class="confirm" on:click={addUserToDb}>确 定</div>
            </div>
        </div>
    </div>
{/if}
<div class="userManageBox">
    <div class="title">
        <div>全部用户<span class="num">(共{num}条)</span></div>
        <div class="search">
            <input bind:value={inputSearch} class="input" placeholder="请输入用户名" />
            <div class="searchbtn" on:click={searchUser}>搜索</div>
        </div>
    </div>
    <div class="userList">
        <div class="peopleCount text">
            <input type=checkbox checked={checkedAll} on:click={checkAll} />
            <td>账号</td>
            <td>姓名</td>
            <td>角色</td>
            <td>注册时间</td>
        </div>
        {#each peopleList as people,index}
            <div class="peopleCount">
                <input type=checkbox checked={checked[index]} on:click={()=>check(index)} />
                <td>{people.account}</td>
                <td>{people.name}</td>
                <td>{people.role}</td>
                <td>{people.time}</td>
            </div>
        {/each}
    </div>
    <div class="userFunction">
        <div class="function">
            <div class="subFunction" on:click={addShow}>增 加</div>
            <div class="subFunction delete" on:click={deleteUser}>删 除</div>
        </div>
        <div class="pages">
            {#if currentPage !== 1 }
            <div class="last page"  on:click={getLastUser}>上一页</div>
            {/if}
            {#if currentPage === 1}
                <div class="last page">上一页</div>
            {/if}
            {#if currentPage !== Math.ceil(num/10)}
            <div class="next page" on:click={getNextUser}>下一页</div>
            {/if}
            {#if currentPage === Math.ceil(num/10)}
            <div class="next page">下一页</div>
            {/if}
         </div>   
    </div>
</div>


<style>
    .userManageBox{
        display: flex;
        width:80vw;
        height: 85vh;
        background-color:rgb(238,245,253);
        box-shadow: 3px 3px 4px 3px rgba(0,0,0,0.3);
        margin-top: 5vh;
        margin-left: 10vw;
        flex-direction: column;
    }
    .userManageBox .title{
        display: flex;
        height: 10%;
        width: 100%;
        border-bottom: solid 2px rgb(186,198,209);
    }
    .title div{
        display: flex;
        font-size: medium;
        height: 100%;
        width: 180px;
        justify-content: center;
        align-items: center;
        text-align: center;
        letter-spacing: 1px;
        font-weight: 550;
    }
    .title .num{
        font-size: smaller;
        margin-left: 5px;
        margin-top: 2px;
        font-weight: 400;
    }
    .userManageBox .userList{
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 80%;
    }
    .peopleCount{
        display: flex;
        height: 10%;
        justify-content: center;
        text-align: center;
        align-items: center;
        width: 100%;
        flex-wrap: wrap;
        border-bottom: solid 2px rgb(208,224,232) ;
    }
    .peopleCount input{
        display: flex;
        /* margin-left: 50%; */
        justify-content: center;
        align-items: center;
        text-align: center;
        width: 100px;
    }
    .text{
        color: rgb(77,166,224);
    }
    .peopleCount td{
        flex:1;
    }
    .userFunction{
        height: 10%;
        width:100%;
        display: flex;
    }
    .userFunction .function{
        display: flex;
        flex-direction: row;
        width: 40%;
        height: 100%;
    }
    .userFunction .function .subFunction{
        display: flex;
        width: 100px;
        height: 25px;
        border-radius: 7px;
        background-color: rgb(27,183,255);
        justify-content: center;
        align-items: center;
        text-align: center;
        margin-left: 10%;
    }
    .userFunction .function .subFunction:hover{
        background-color: rgb(64,193,243);
    }
    .userFunction .function .delete{
        margin-left: 20px;
    }
    .pages{
        display: flex;
        width: 60%;
        height: 100%;
        flex-direction: row;
        margin-left: 50%;
    }
    .pages .page{
        display: flex;
        width: 30%;
        height: 50%;
        border-radius: 7px;
        background-color: rgb(27,183,255);
        justify-content: center;
        align-items: center;
        text-align: center;
        margin-left: 20%;
    }
    .pages .page:hover{
        background-color: rgb(52, 174, 226);
    }
    .pages .next{
        margin-left: 20px;
    }
    .search{
        display: flex;
        height:100%;
        margin-left: 67%;
    }
    .search .input{
        display: flex;
        height: 40%;
        width: 800px;
        border: solid 1px ;
    }
    .search .searchbtn{
        display: flex;
        width: 1000px;
        height: 40%;
        margin-left: 10px;
        background-color: rgb(27,183,255);
    }
    .AddUser{
        position: absolute;
        width:100vw;
        height: 100vh;
        background: rgba(0,0,0,0.2);
    }
    .AddUser .AddUserBox{
        display: flex;
        margin:0 auto;
        margin-top: 50px;
        width: 40vw;
        height: 70vh;
        background-color: aliceblue;
        flex-direction: column;
        text-align: center;
        justify-content: center;
        align-items: center;
    }
    .AddUser .AddUserBox .header{
        display: flex;
        width:100%;
        height: 10%;
        font-size:large;
        font-weight: 600;
        justify-content: center;
        margin-top: -30%;
        align-items: center;
    }
    .AddUser .AddUserBox .addinput{
        display: flex;
        width: 80%;
        height: 12%;
        margin:0 auto;
        margin-top: 20px;
    }
    .AddUser .AddUserBox .addinput input{
        width: 100%;
        height: 100%;
        border-radius: 10px;
    }
    .AddUser .AddUserBox .function{
        display: flex;
        width: 100%;
    }
    .AddUser .AddUserBox .function .confirm{
        position: relative;
        left:80%;
        top:150px;
        height: 30px;
        width: 100px;
        border-radius:7px ;
        background-color: skyblue;
        justify-content: center;
        align-items: center;
        text-align: center;
    }
</style>