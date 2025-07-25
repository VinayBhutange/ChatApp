document.addEventListener("DOMContentLoaded", () => {
    const messages = document.getElementById("messages");
    const messageInput = document.getElementById("message-input");
    const sendButton = document.getElementById("send-button");

    // Establish WebSocket connection
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
        console.log("WebSocket connection established");
        addMessage("System", "Connected to the chat!");
    };

    socket.onmessage = (event) => {
        // For this simple app, we assume the message is just plain text
        // In a real app, you'd likely parse a JSON object with sender info
        addMessage("Anonymous", event.data);
    };

    socket.onclose = () => {
        console.log("WebSocket connection closed");
        addMessage("System", "Connection closed.");
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
        addMessage("System", "An error occurred with the connection.");
    };

    const sendMessage = () => {
        const message = messageInput.value.trim();
        if (message !== "") {
            socket.send(message);
            messageInput.value = "";
        }
    };

    sendButton.addEventListener("click", sendMessage);

    messageInput.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
            sendMessage();
        }
    });

    function addMessage(sender, text) {
        const messageElement = document.createElement("div");
        messageElement.innerHTML = `<strong>${sender}:</strong> ${text}`;
        messages.appendChild(messageElement);
        messages.scrollTop = messages.scrollHeight; // Auto-scroll to the bottom
    }
});
