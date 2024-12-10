let socket; // WebSocket variable
let currentUsername; // Store the current user's username
let selectedUser; // Store the selected user for messaging
let messageQueue = []; // Queue to store messages before WebSocket is open
let messagesReceived = {}; // save the messages from each user

function setupWebSocket(username) {
    if (!username) {
        console.error("Username is required to establish WebSocket connection");
        return;
    }
    currentUsername = username; // Set the current user's username
    console.log(`Setting up WebSocket connection for user: ${username}`);
    socket = new WebSocket(`ws://localhost:8080/ws?username=${username}`); // Create a new WebSocket connection

    socket.onopen = function () {
        console.log("WebSocket connection established");
        // Send any messages that were queued before the connection was open
        while (messageQueue.length > 0) {
            const msg = messageQueue.shift();
            console.log("Sending queued message:", msg);
            socket.send(JSON.stringify(msg));
        }
    };

socket.onmessage = function (event) {
    console.log("Message received from server:", event.data);
    const message = JSON.parse(event.data);
    
    if (message.newUser) {
        if (message.newUser !== currentUsername) {
            if (!onlineUsers.includes(message.newUser)) {
                onlineUsers.push(message.newUser);
                updateUserLists();
            }
            addNotification(message.newUser, "logged in!");
        }
    } else if (message.removeUser) {
        onlineUsers = onlineUsers.filter(user => user !== message.removeUser);
        updateUserLists();
        addNotification(message.removeUser, "logged out!", "status");
    } else {
        // Handle regular messages
        handleNewMessage(message);
    }
};

    socket.onerror = function (error) {
        console.error("WebSocket error:", error);
    };

    socket.onclose = function () {
        console.log("WebSocket connection closed");
    };
}





function handleNewMessage(message) {
    const messageTime = new Date(message.timeSent).getTime();
    
    // Update message history immediately
    if (!userMessageHistory[message.senderUsername]) {
        userMessageHistory[message.senderUsername] = {};
    }
    if (!userMessageHistory[message.receiverUsername]) {
        userMessageHistory[message.receiverUsername] = {};
    }

    // Update timestamps and force immediate sort
    userMessageHistory[message.senderUsername][message.receiverUsername] = messageTime;
    userMessageHistory[message.receiverUsername][message.senderUsername] = messageTime;
    
    // Force immediate re-sort and render
    requestAnimationFrame(() => {
        updateUserLists();
    });

    // Handle notifications for new messages
    if (message.receiverUsername === currentUsername && 
        message.senderUsername !== currentUsername) {
        addNotification(message.senderUsername, message.messageText);
        updateUserLists();
    }

    // Update chat if it's the current conversation
    if ((message.senderUsername === selectedUser && message.receiverUsername === currentUsername) || 
        (message.senderUsername === currentUsername && message.receiverUsername === selectedUser)) {
        appendNewMessage(message);
    }
}




// Update the user lists function to handle all states

function updateUserLists() {
    // Separate users into online and offline
    const offlineUsers = messageUsers.filter(user => !onlineUsers.includes(user));
    
    // Sort both lists
    const sortedOnlineUsers = sortUsersByLastMessage([...onlineUsers]);
    const sortedOfflineUsers = sortUsersByLastMessage([...offlineUsers]);

    // Render the lists
    renderUserLists(sortedOnlineUsers, sortedOfflineUsers);

    

    // Disable the chat input if no users are online or if no user is selected for messaging
    const chatInput = document.getElementById('chatInput');

    if (!selectedUser || selectedUser === '') {
        chatInput.disabled = true;
        chatInput.placeholder = 'Select a user to start';
    } else {
        chatInput.disabled = false;
        chatInput.placeholder = 'Type a message...';
    }
}


function updateOnlineUsers(newUser) {
    if (!onlineUsers.includes(newUser)) {
        onlineUsers.push(newUser);
        updateUserLists();  // Changed from renderOnlineUsers to updateUserLists
    }
}

function sendMessage(senderUsername, receiverUsername, messageText) {
    const msg = {
        senderUsername: senderUsername,
        receiverUsername: receiverUsername,
        messageText: messageText,
        timeSent: new Date().toISOString()
    };
    
    // Update message history immediately for real-time sorting
    const messageTime = new Date(msg.timeSent).getTime();
    if (!userMessageHistory[senderUsername]) {
        userMessageHistory[senderUsername] = {};
    }
    if (!userMessageHistory[receiverUsername]) {
        userMessageHistory[receiverUsername] = {};
    }
    
    // Update both sides of the conversation
    userMessageHistory[senderUsername][receiverUsername] = messageTime;
    userMessageHistory[receiverUsername][senderUsername] = messageTime;
    
    // Update UI immediately
    updateUserLists();
    
    // Send message through WebSocket
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(msg));
    } else {
        console.error("WebSocket is not open, queuing message");
        messageQueue.push(msg);
    }
}

function sendChatMessage() {
    const chatInput = document.getElementById("chatInput");
    const messageText = chatInput.value;
    if (messageText.trim() !== "") {
        sendMessage(currentUsername, selectedUser, messageText);
        chatInput.value = ""; // Clear the input field after sending the message
    }
}