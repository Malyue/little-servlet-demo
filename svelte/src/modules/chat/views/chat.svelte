<script lang="ts">
    import { onMount,beforeUpdate } from "svelte";
    import { dataset_dev, onDestroy } from "svelte/internal";
    import Header from '../../../lib/Header.svelte'
    import Button, { Label } from '@smui/button'


    export let selectUser:number = -1
    export let selectUserName:string = "未选择"
    export let inputContent:string = ""
    export let websocket:any = "";
    export let messageList:any[] = [];
    export let onlineStatus:boolean = false
    //显示
    export let message:any = [];
    export let userList:any = [];
    export let userIdList:any = [];
    // export let messageList:any = []; 

    const sendMessage = ()=>{
        if (selectUser === -1) {
            alert("请选择聊天对象")
        }else{
            //聊天对象已离线
            if(onlineStatus === false){
                alert("对方已离线，请切换其他用户")
                return
            }
            const data = {
                type: 1,
                content: inputContent,
                time: new Date(),
                receiver:selectUser
            };
            const mymessage = {
                content: inputContent,
                type:1,
                userName:["我"],
                receiver:[selectUser],
                userid:[-1]
            } 
            messageList.push(mymessage)
            messageList = messageList
            message.push(mymessage)
            message = message
            websocket.send(JSON.stringify(data))
        }
        inputContent=""
    }
    const selectUserToChat = (index:number)=>{
        onlineStatus = true
        message = []
        selectUserName = userList[index]
        selectUser = userIdList[index]
        // console.log(me)
        for(let i=0;i<messageList.length;i++){
            console.log(messageList[i].receiver[0])
            if(messageList[i].userid[0]===selectUser||messageList[i].receiver[0]===selectUser){
                message.push(messageList[i])
                message = message
            }
        }
    }

    window.onbeforeunload = function() {
        websocket.close()
    }

    onDestroy(()=>{
        websocket.close()
    })

    onMount(()=>{
        if (location.hash === '#/chat'){
            let url = "ws://114.132.232.3:9090/api/chat" 
            var token = localStorage.getItem("userToken") || ''
            websocket = new WebSocket(url+'?token='+token)
            // System.setProperty("ws",websocket);
            websocket.opopen = () => {
                console.log("ws connect")
            }
            websocket.onmessage = (res:any) => {
                var data:any
                // console.log(JSON.parse(res.data)) 
                data = JSON.parse(res.data)
                console.log(data)
                //添加在綫用戶
                if (data.type === 2){
                    for(let i = 0 ;i<data.userid.length;++i){
                        userList.push(data.userName[i]);
                        userIdList.push(data.userid[i]);
                        userIdList=userIdList
                        userList=userList
                    }
                }else if (data.type === 3){
                    console.log("离线处理")
                    //離綫處理
                    //从数组中删除
                    var index = userIdList.indexOf(data.userid[0]);
                    console.log(index,data.userid[0]);
                    if (index > -1) {
                        userIdList.splice(index,1);
                        userList.splice(index,1);
                        userIdList = userIdList;
                        userList = userList
                    }
                    //如果为当前聊天对象
                    if (selectUser === data.userid[0]){
                        alert("当前聊天对象已离线...")
                        onlineStatus = false
                    }
                    //如果为此时选中的目标，则显示对方已离线
                }else if (data.type===1){
                    //接受消息
                    messageList.push(data)
                    messageList = messageList
                    //如果是发给该用户或者该用户收到的消息
                    if(selectUser === -1){
                        return
                    }
                    if(data.userid[0] === selectUser){
                        message.push(data)
                        message = message
                    }
                }
            }
            websocket.onclose = (err:any) => {
                console.log(err)
            }
            websocket.onerror = (err:any) => {
                console.log(err)
            }
        }
    })
</script>


<Header />
<div class="chatBox">
    <div class="chat">
        <div class="head">{selectUserName}</div>
        <div class="message">
            {#each message as m}
            <div class="messageItem">
                <div class="messageContent"><p class="name">{m.userName}:  </p><p class="content">{m.content}</p></div>
            </div>
            {/each}
        </div>
        <div class="send">
            <input bind:value={inputContent} />
            <div class="sendMessage" on:click={sendMessage}>发 送</div>
        </div>
    </div>
    <div class="user">
        <div class="title" >在 綫 用 戶</div>
        {#each userList as user,index}
            <div class="userCount" on:click={()=>selectUserToChat(index)}>{user}</div>
        {/each}
    </div>
    <Button variant="unelevated">
        <Label>Raised</Label>
    </Button>
</div>

<style>
    .chatBox{
        display: flex;
        margin-left: 10vw;
        margin-top: 40px;
        /* height:; */
    }
    .chat{
        display: flex;
        flex-direction: column;
        /* background-color: rgb(238,245,253); */
        background-color: white;
        width: 50vw;
        height: 85vh;
        border-radius: 7px;
        /* border: solid 0.5px; */
    }
    .chat .message{
        display: flex;
        width: 100%;
        height: 85%;
        flex-direction: column;
    }
    .chat .head{
        display: flex;
        width: 100%;
        height: 8%;
        justify-content: center;
        align-items: center;
        text-align: center;
        border-bottom: 2px solid rgb(237,240,248);
        color: rgb(71,73,74);
    }
    .chat .send{
        display: flex;
        background-color: rgb(225,238,242);
        height: 7%;
        width: 100%;
    }
    .chat .send input{
        display: flex;
        height: 60%;
        justify-content: center;
        width: 75%;
        margin-left: 2%;
        margin-top: 0.7%;
    }
    .chat .send .sendMessage{
        display: flex;
        justify-content: center;
        align-items: center;
        text-align: center;
        margin-left: 3%;
        height: 60%;
        width: 15%;
        margin-top: 1.2%;
        background-color: rgb(77,168,220);
        border-radius: 5px;
    }
    .chat .send .sendMessage:hover{
        background-color: rgb(42, 144, 199);
        /* border: 2px; */
    }
    .user{
        display: flex;
        flex-direction: column;
        /* background-color: rgb(238,245,253); */
        background-color: white;
        /* border: solid 0.5px; */
        margin-left: 3vw;
        width: 25vw;
        height: 85vh;
        border-radius: 7px;
    }
    .user .title{
        display: flex;
        width: 100%;
        height: 10%;
        justify-content: center;
        align-items: center;
        text-align: center;
        border-bottom: 3px solid rgb(208,224,232);
        color:rgb(93,135,255)
    }
    .user .userCount{
        display: flex;
        width: 100%;
        height: 5%;
        border-bottom: solid 2px rgb(208,224,232);
        text-align: center;
        align-items: center;
        justify-content: center;
        /* line-height: 5%; */
    }
    .user .userCount:hover{
        background-color: rgb(244,244,244);
    }
    .message .messageItem{
        display: flex;
        width: 100%;
        height: 5%;
        /* background-color: aqua; */
    }
    .message .messageItem .messageContent{
        display: flex;
        width: 100%;
        height: 100%;
        /* background-color: aquamarine; */
    }
    .name{
        display: flex;
        color:skyblue
    }
    .content{
        display: flex;
        margin-left: 5px;
    }
</style>


