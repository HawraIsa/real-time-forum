function sendChatMessage() {
    const chatInput = document.getElementById("chatInput"); // Get the chat input element
    const messageText = chatInput.value; // Get the message text from the input
    if (messageText.trim() !== "") { // If the message text is not empty
        sendMessage(currentUsername, selectedUser, messageText); // Call sendMessage from websocket.js
        chatInput.value = ""; // Clear the input field after sending the message
    }
}

// Add this new function to render the sorted lists
function renderUserLists(onlineUsers = [], offlineUsers = []) {
    const createUserListHTML = (users, isOnline) => {
        return users
            .filter(user => user !== currentUsername)
            .map(user => {
                const lastMessageTime = userMessageHistory[currentUsername]?.[user] || 
                                      userMessageHistory[user]?.[currentUsername];
                
                const timeHtml = lastMessageTime ? 
                    `<span class="last-message-time" data-timestamp="${lastMessageTime}"></span>` : '';
                const selectedClass = user === selectedUser ? 'selected' : '';
                
                return `<li class="user-item ${isOnline ? 'online' : 'offline'} ${selectedClass}" 
                           onclick="selectUser('${user}')">
                    <div class="user-status-dot"></div>
                    <div class="user-info">
                        <span class="username">${user}</span>
                        ${timeHtml}
                    </div>
                    ${selectedClass ? '<div class="active-indicator"></div>' : ''}
                </li>`;
            })
            .join('');
    };

    const onlineUsersHTML = `
        <div id="onlineUsers" class="user-section">
            <h2>Online — ${onlineUsers.filter(user => user !== currentUsername).length}</h2>
            <ul class="user-list">
                ${createUserListHTML(onlineUsers, true)}
            </ul>
        </div>
        <div id="offlineUsers" class="user-section">
            <h2>Offline — ${offlineUsers.length}</h2>
            <ul class="user-list">
                ${createUserListHTML(offlineUsers, false)}
            </ul>
        </div>
    `;
    
    document.getElementById("onlineUsersContainer").innerHTML = onlineUsersHTML;
}

let onlineUsers = []; // Array to hold online users
let messageUsers = []; // Array to hold message users
let userMessageHistory = {}; // Object to track message history for each user

async function fetchOnlineUsers() {
    try {
        // Fetch all necessary data in parallel for better performance
        const [messageUsersResponse, onlineUsersResponse, allMessagesResponse] = await Promise.all([
            fetch("/getMessageUsers?username=" + currentUsername),
            fetch("/getOnlineUsers"),
            fetch("/getAllUserMessages?username=" + currentUsername)
        ]);

        // Process the responses
        messageUsers = await messageUsersResponse.json();
        onlineUsers = [...new Set(await onlineUsersResponse.json())];
        const allMessages = await allMessagesResponse.json();

        // Initialize message history from all messages
        allMessages.forEach(message => {
            const messageTime = new Date(message.timeSent).getTime();
            
            // Initialize nested objects if needed
            if (!userMessageHistory[message.senderUsername]) {
                userMessageHistory[message.senderUsername] = {};
            }
            if (!userMessageHistory[message.receiverUsername]) {
                userMessageHistory[message.receiverUsername] = {};
            }

            // Update message history for both users
            // Only update if this is a more recent message
            const currentTime = userMessageHistory[message.senderUsername][message.receiverUsername];
            if (!currentTime || messageTime > currentTime) {
                userMessageHistory[message.senderUsername][message.receiverUsername] = messageTime;
                userMessageHistory[message.receiverUsername][message.senderUsername] = messageTime;
            }
        });

        // Update the lists with complete message history
        updateUserLists();
        
    } catch (error) {
        console.error("Error initializing chat:", error);
    }
}


function selectUser(username) {
    currentPage = 1; // Reset the current page
    selectedUser = (selectedUser === username) ? null : username; // Set the selected user for messaging

      // Force UI update immediately
     updateUserLists(); // This will re-render the user lists with the new selection


    if (!socket || socket.readyState !== WebSocket.OPEN) { // If the WebSocket is not open
        setupWebSocket(currentUsername); // Set up the WebSocket connection if not already open
    }
    fetchMessages(currentPage); // Load the first page of messages
}

function formatTimestamp(utcTimestamp) {
    const date = new Date(utcTimestamp); // Create a Date object from the UTC timestamp

    // Extract date components
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are zero-indexed
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    // Format the date to 'YYYY-MM-DD HH:mm:ss'
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}


function createMessageElement(message) {
    const messageElement = document.createElement("div");
    const formattedTime = formatTimestamp(message.timeSent);
    messageElement.textContent = `${formattedTime} - ${message.senderUsername}: ${message.messageText}`;
    return messageElement;
}


function appendNewMessage(message) {
    const chatMessages = document.getElementById("chatMessages");
    const messageElement = createMessageElement(message);
    chatMessages.appendChild(messageElement); // Append to the end
    chatMessages.scrollTop = chatMessages.scrollHeight; // Scroll to the bottom
}

function prependMessages(messages) {
    const chatMessages = document.getElementById("chatMessages");
    let initialScrollHeight = chatMessages.scrollHeight;
    messages.forEach(message => {
        const messageElement = createMessageElement(message);
        chatMessages.insertBefore(messageElement, chatMessages.firstChild); // Insert at the top
    });
    chatMessages.scrollTop += chatMessages.scrollHeight - initialScrollHeight; // Maintain view
}

function renderMessages(messages) {
    const chatMessages = document.getElementById("chatMessages");
    chatMessages.innerHTML = ''; // Clear the chat messages container
    if (messages.error) {
        console.error(messages.error);
    } else {
        messages.forEach(message => {
            const messageElement = createMessageElement(message);
            chatMessages.appendChild(messageElement); // Append to the end
        });
        chatMessages.scrollTop = chatMessages.scrollHeight; // Scroll to the bottom
    }
}


//the sorting function to consider only direct message interactions
function sortUsersByLastMessage(users) {
    return users.sort((a, b) => {
        // Get last message times for both users
        const aTime = userMessageHistory[currentUsername]?.[a] || 
                     userMessageHistory[a]?.[currentUsername] || 0;
        const bTime = userMessageHistory[currentUsername]?.[b] || 
                     userMessageHistory[b]?.[currentUsername] || 0;
        
        // If both users have message history
        if (aTime && bTime) {
            return bTime - aTime; // Most recent first
        }
        
        // If only one has message history
        if (aTime) return -1;
        if (bTime) return 1;
        
        // If neither has message history, sort by online status
        const aIsOnline = onlineUsers.includes(a);
        const bIsOnline = onlineUsers.includes(b);
        if (aIsOnline !== bIsOnline) {
            return aIsOnline ? -1 : 1;
        }
        
        // If online status is the same, sort alphabetically
        return a.localeCompare(b, undefined, { sensitivity: 'accent' });
    });
}


// Update the fetchMessages function to  handle message history
function fetchMessages(page) {
    fetch(`/getPrivateMessages?senderUsername=${currentUsername}&receiverUsername=${selectedUser}&page=${page}`)
        .then(response => response.json())
        .then(messages => {
            // Update message history only for the selected conversation
            if (messages.length > 0) {
                const latestMessage = messages[0];
                const messageTime = new Date(latestMessage.timeSent).getTime();
                
                // Initialize nested objects if needed
                if (!userMessageHistory[currentUsername]) {
                    userMessageHistory[currentUsername] = {};
                }
                if (!userMessageHistory[selectedUser]) {
                    userMessageHistory[selectedUser] = {};
                }

                // Update message history for both users
                userMessageHistory[currentUsername][selectedUser] = messageTime;
                userMessageHistory[selectedUser][currentUsername] = messageTime;
                
                // Update the lists to reflect the new interaction
                updateUserLists();
            }
            
            // Handle message display
            if (page === 1) {
                renderMessages(messages.reverse());
            } else {
                prependMessages(messages);
            }
            currentPage++;
        })
        .catch(error => { 
            console.error("Error fetching messages:", error); 
            isFetching = false; 
        });
}

function appendMessages(messages, toEnd = false) {
    const chatMessages = document.getElementById("chatMessages");
    messages.forEach(message => {
        const messageElement = createMessageElement(message);
        chatMessages.appendChild(messageElement); // Append to the end
    });
    if (toEnd) {
        chatMessages.scrollTop = chatMessages.scrollHeight; // Scroll to the bottom
    }
}



let currentPage = 1; // Variable to hold the current page number
function loadMore() {
    let messageContainer = document.getElementById("chatMessages"); // Get the chat messages container
    if (messageContainer.scrollTop <= 10) { // If the scroll position is near the top
        fetchMessages(currentPage); // Fetch more messages
    }
}

function throttle(mainFunction, delay) {
    let start = Date.now(); // Get the current time

    return (...args) => { // Return a throttled function
        if (Date.now() - start > delay) { // If the delay has passed
            mainFunction(...args); // Call the main function
            start = Date.now(); // Reset the start time
        }
    };
}

function formatLastMessageTime(timestamp) {
    const now = new Date();
    const messageDate = new Date(timestamp);
    const diff = now - messageDate;
    
    // Less than a minute
    if (diff < 60000) {
        return 'just now';
    }
    // Less than an hour
    if (diff < 3600000) {
        const minutes = Math.floor(diff / 60000);
        return `${minutes}m ago`;
    }
    // Less than a day
    if (diff < 86400000) {
        const hours = Math.floor(diff / 3600000);
        return `${hours}h ago`;
    }
    // Less than a week
    if (diff < 604800000) {
        const days = Math.floor(diff / 86400000);
        return `${days}d ago`;
    }
    // Otherwise show date
    return messageDate.toLocaleDateString();
}