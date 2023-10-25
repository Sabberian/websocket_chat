var socket = new WebSocket("ws://localhost:8080/ws");
var savedMessages = localStorage.getItem("savedMessages");

var messages = [];

if (savedMessages) {
    console.log(savedMessages);
    messages = JSON.parse(savedMessages);
    displayMessages(messages);
}

function displayMessages(messages) {
    var chat = document.getElementById("chat");
    chat.innerHTML = "";

    for (var i = 0; i < messages.length; i++) {
        var message = messages[i];
        var messageElement = document.createElement("div");
        messageElement.innerHTML = "<strong>" + message.username + ":</strong> " + message.message;
        chat.appendChild(messageElement);
    }
}

function sendMessage() {
    var username = document.getElementById("username").value;
    var message = document.getElementById("message").value;

    if (username === "" || message === "") {
        alert("Введите имя пользователя и сообщение.");
        return;
    }

    var data = {
        username: username,
        message: message
    };

    socket.send(JSON.stringify(data));
    document.getElementById("message").value = "";
}

function updateChat(message) {
    var chat = document.getElementById("chat");
    var messageElement = document.createElement("div");
    messageElement.innerHTML = "<strong>" + message.username + ":</strong> " + message.message;
    chat.appendChild(messageElement);
}


socket.onmessage = function (event) {
    var message = JSON.parse(event.data);
    messages.push(message);

    localStorage.setItem("savedMessages", JSON.stringify(messages));

    updateChat(message);
}
