function createClass(){

let name = document.getElementById("createName").value

if(name === ""){
alert("Enter your name")
return
}

localStorage.setItem("username", name)
localStorage.setItem("creatorName", name)

let code = Math.random().toString(36).substring(2,8).toUpperCase()

localStorage.setItem("classCode", code)

alert("Classroom Code: " + code)

window.location.href = "/editor"

}

function login(){
window.location.href = "/dashboard"
}

//function signup(){
window.location.href = "/dashboard"
//}

function joinClass(){

let name = document.getElementById("joinName").value
let code = document.getElementById("classCode").value

if(name === "" || code === ""){
alert("Enter name and code")
return
}

localStorage.setItem("username", name)
localStorage.setItem("classCode", code)

window.location.href = "/editor"

}

function leave(){
window.location.href = "/dashboard"
}

function sendMessage(){

let input = document.getElementById("chatInput")
let msg = input.value

if(msg === "") return

socket.send(msg)

input.value = ""

}
