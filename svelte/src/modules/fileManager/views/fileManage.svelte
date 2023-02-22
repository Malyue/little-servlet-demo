<script lang="ts">
    import { onMount } from "svelte";
    import Header from '../../../lib/Header.svelte'


    export let fileList:any[] = []
    // export let fileIdList:any[] = []
    // export let uploadfile:any
    export let formData:FormData = new FormData()
    export let fileFormdata:any[] = []

    export function getList(){
        let token =localStorage.getItem("userToken")||''
        let url = "/api/getFile" +"?page="+1+"&count="+10
        fetch(url,{
            method:"GET",
            mode:"cors",
            cache:'no-cache',
            credentials:'same-origin',
            headers:{
                'Content-Type':'application/json',
                'Authorization':token
            },
        }).then((res)=>{
            return res.json()
        }).then((res)=>{
            if (res.status === 0){
                fileList = []
                if (res.data != null){
                    fileList = fileList.concat(res.data)
                }
            }else{
                alert("查找错误")
            }
        })
    };

    export function deleteFile(index:any){
        let token =localStorage.getItem("userToken")||''
        let url = "/api/deleteFile" +"?"+"fid="+fileList[index].fid
        fetch(url,{
            method:"GET",
            mode:"cors",
            cache:'no-cache',
            credentials:'same-origin',
            headers:{
                'Content-Type':'application/json',
                'Authorization':token
            },
        }).then((res)=>{
            return res.json()
        }).then((res)=>{
            if (res.status === 0){
                alert("删除成功")
                getList()
            }else{
                alert("删除错误")
            }
        })
    }
    export function downloadFile(index:any){
        let url = "/api/download" +"?fid="+fileList[index].fid
        // location.href = url
        window.open(url)
    }
    export function fileChange(e:any){
        const input = e.target;
        const files = e.target.files;
        formData = new FormData;
        if (files[0]) {
            const file =files[0]
            if(file.size > 1024*1024*3){
                alert("文件过大");
                input.value = '';
                return false
            }else{
                formData.append("file",file);
                formData = formData
                // fileFormdata
                fileFormdata = []
                fileFormdata.push(file["name"]);
                fileFormdata = fileFormdata;
            }
        }
    }
    export function uploadFile(){
        if (formData){
            let url = "/api/upload";
            let token =localStorage.getItem("userToken")||''
            fetch(url,{
            method:"POST",
            mode:"cors",
            cache:'no-cache',
            credentials:'same-origin',
            headers:{
                // 'Content-Type':'multipart/form-data',
                // 'Content-Type':'application/x-www-form-urlencoded',
                'Authorization':token
            },
            body:formData
        }).then((res)=>{
            return res.json()
        }).then((res)=>{
            if (res.status === 0){
                alert("上传成功")
                formData = new FormData()
                let fileTag = document.getElementById("file")||''
                // fileTag.value = '';
                fileTag = '';
                getList()
            }else{
                alert("查找错误")
                formData = new FormData()
            }
        })
        }else{
            alert("请选择文件")
        }
    };

    onMount(()=>{
        getList()
    })
</script>


<Header />
<div class="fileBox">
    <div class="title"><div>全部文件</div></div>
    <div class="fileList">
        <div class="fileCount text">
            <td>文件名</td>
            <td>操作</td>
        </div>
        {#each fileList as file,index}
            <div class="fileCount">
                <td>{file.filename}</td>
                <td>
                    <span class="download" on:click={()=>downloadFile(index)} >下载</span>
                    <span class="delete" on:click={()=>deleteFile(index)}>删除</span>
                </td>
            </div>
        {/each}
    </div>
    <div class="fileFunction">
        <div class="file">
            {#each fileFormdata as file,index}
                <div class="fileCount">{file}</div>
            {/each}
        </div>
        <div class="function">
            <div class="input">
                <input id="file" type="file" name="file" on:change={fileChange}/>
                <p>选择文件</p>
            </div>
            <div class="subFunction" on:click={uploadFile}>
                上传
            </div>
        </div>
    </div>
</div>



<style>
    .fileBox{
        display: flex;
        width:80vw;
        height: 85vh;
        background-color:rgb(238,245,253);
        box-shadow: 3px 3px 4px 3px rgba(0,0,0,0.3);
        margin-top: 5vh;
        margin-left: 10vw;
        flex-direction: column;
    }
    .title{
        display: flex;
        height: 10%;
        width: 100%;
        border-bottom: solid 2px rgb(186,198,209);
    }
    .title div{
        display: flex;
        font-size: medium;
        /* background-color: aqua; */
        height: 100%;
        width: 180px;
        justify-content: center;
        align-items: center;
        text-align: center;
        letter-spacing: 1px;
        font-weight: 550;
    }
    .fileList{
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 80%;
    }
    .fileCount{
        display: flex;
        height: 10%;
        justify-content: center;
        text-align: center;
        align-items: center;
        width: 100%;
        flex-wrap: wrap;
        border-bottom: solid 2px rgb(208,224,232);
    }
    .text{
        color: rgb(77,166,224);
    }
    .fileCount td{
        flex:1;
    }
    .fileFunction{
        height: 10%;
        width:100%;
        display: flex;
        /* background-color: aqua; */
    }
    .fileFunction .function{
        display: flex;
        /* display: inline-block; */
        flex-direction: row;
        width: 100%;
        height: 100%;
        /* background-color: aqua; */
    }
    .fileFunction .function .subFunction{
        display: flex;
        width: 100px;
        height: 25px;
        border-radius: 7px;
        background-color: rgb(27,183,255);
        justify-content: center;
        align-items: center;
        text-align: center;
        margin-left: 80%;
    }
    .fileFunction .function .subFunction:hover{
        background-color: rgb(64,193,243);
    }
    .function .input{
        display: flex;
        border-radius: 7px;
        background-color: rgb(27,183,255);
        width: 100px;
        height: 25px;
        justify-content: center;
        align-items: center;
        text-align: center;
        margin-left:20px;
    }
    .function .input p{
        width: 100%;
        height: 100%;
        margin-left: -100px;
        display: flex;
        text-align: center;
        align-items: center;
        justify-content: center;
    }
    .function .input input{
        display: flex;
        opacity: 0;
        width: 100%;
        height: 100%;
    }
    .function .file .fileCount{
        display: flex;
        /* display: inline-block; */
        margin-left: 70px;
        justify-content: center;
        align-items: center;
        text-align: center;
        height: 25px;
        width: 250px;
    }
    .download{
        color: skyblue;
    }
    .delete{
        color:red;
        margin-left: 20px;
    }
</style>