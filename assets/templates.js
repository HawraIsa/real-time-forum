let pendingPostId = null;
const loginDialog = document.getElementById("loginDialog"); // Get the login dialog element
const registerDialog = document.getElementById("registerDialog"); // Get the register dialog element
const addPostDialog = document.getElementById("addPostDialog"); // Get the add post dialog element
//const errDialog = document.getElementById("404Dialog"); // Get the 404 dialog element

// Function to load different templates based on the template name
function loadTemplate(templateName, data = {}) {
    loginDialog.close(); // Close the login dialog
    registerDialog.close(); // Close the register dialog
    //errDialog.close();

    document.getElementById("loginForm").reset(); // Reset the login form
    document.getElementById("registerForm").reset(); // Reset the register form

    // Check for pending post when loading index
    if (templateName === "index") {
        renderIndexTemplate();
        if (pendingPostId) {
            // If there was a pending post, fetch and show it
            fetchAndFilterPosts().then(() => {
                showSinglePost(pendingPostId);
                pendingPostId = null;
            });
        } else {
            fetchAndFilterPosts();
        }
        loadPrivateMessages(currentUsername);
        fetchOnlineUsers();
    } else if (templateName === "post") {
        fetchPostData(data.postID); // Fetch data for a specific post
    } else if (templateName === "login") {
        app.innerHTML = ""; // Clear the app container
        loginDialog.showModal(); // Show the login dialog
    } else if (templateName === "register") {
        app.innerHTML = ""; // Clear the app container
        registerDialog.showModal(); // Show the register dialog
    } else if (templateName === "404" ) {
        show404();
    }

}


// function to show err page 404
function show404() {
    const contentDiv = document.getElementById('app');
    contentDiv.innerHTML = `
        <div class="404div">
        <h1>404</h1>
        <p>Oops! The page you're looking for doesn't exist.</p>
        <a href="/">Go Back Home</a>
        </div>
    `;
}

// Function to render the index template
function renderIndexTemplate() {
    const indexHTML = `
            <nav>
                <div><a href="/">&#128510; Home</a></div>
                <div>
                <a onclick="addpost();">New Post</a>
                <a onclick="logout();">Logout</a>
                </div>
            </nav>
            <main id="postsContainer">
                <!-- Posts will be dynamically inserted here -->
            </main>
            <div id="post">
                <div style="text-align: center;">
                <br><br><br>
                <svg height="80px" width="80px" version="1.1" id="Capa_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" 
                        viewBox="0 0 177.218 177.218" xml:space="preserve">
                    <g>
                        <g>
                            <polygon style="fill:#ffffff;" points="88.008,68.431 120.493,95.945 123.685,93.247 91.196,65.736 122.866,38.913 
                                123.685,38.219 120.493,35.524 88.008,63.038 55.522,35.524 52.334,38.219 84.829,65.736 53.153,92.56 52.334,93.247 
                                55.522,95.945 		"/>
                            <path style="fill:#ffffff;" d="M164.77,5.819l-0.086-0.812H12.526L0.1,119.609L0,120.614h0.319v51.596h176.566v-51.596h0.333
                                L164.77,5.819z M17.722,10.797h141.765l11.281,104.026h-48.014v2.895c0,18.828-15.317,34.149-34.149,34.149
                                c-18.835,0-34.153-15.321-34.153-34.149v-2.895H6.442L17.722,10.797z M171.094,166.427H6.109v-45.813h42.656
                                c1.492,20.804,18.9,37.041,39.84,37.041c20.936,0,38.34-16.234,39.825-37.041h42.66v45.813H171.094z"/>
                        </g>
                    </g>
                    </svg>
                <br><br><br> No Post Selected </div>
            </div>
            <div>
                <div id="privateMessagesContainer" class="private-messages-container">
                    <div id="onlineUsersContainer" class="online-users-container"></div>
                    <div class="chat-container">
                        <h2>Private Messages</h2>
                        <div id="chatMessages" class="chat-messages"></div>
                        <div class="chat-input-container">
                            <input type="text" id="chatInput" class="chat-input">
                            <button onclick="sendChatMessage()" class="chat-send-button">Send</button>
                        </div>
                    </div>
                </div>
            </div>
        `;
    app.innerHTML = indexHTML; // Insert the index HTML into the app container
    document.getElementById("chatMessages").addEventListener("scroll", throttle(loadMore, 500)); // Add a scroll event listener to the chat messages container
}