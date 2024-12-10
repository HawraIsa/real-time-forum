document.addEventListener("DOMContentLoaded", function() {
    const app = document.getElementById("app"); // Get the main app container

    const session = sessionStorage.getItem("session"); // Check if there is a session stored in sessionStorage
    console.log(session); // Log the session to the console
    if (window.location.pathname == "/404") {
        loadTemplate("404");
    } else  if (session) { // If a session exists
        fetch("/authSession?session=" + session).then((response) => { // Fetch the session authentication endpoint
            response.json().then((json) => { // Parse the response as JSON
                if (json.error) { // If there is an error in the response
                    console.error(json.error); // Log the error to the console
                    loadTemplate("login"); // Load the login template
                } else { // If the session is valid
                    username = json.username; // Set the username from the response
                    currentUsername = username; // Set the current username
                    loadTemplate("index"); // Load the main index template
                }
            });
        });
    } else { // If no session exists
        loadTemplate("login"); // Load the login template
    }

    // Expose functions to global scope
    window.setLikes = setLikes; // Expose setLikes function
    window.setCommentLikes = setCommentLikes; // Expose setCommentLikes function
    window.sendChatMessage = sendChatMessage; // Expose sendChatMessage function
    window.handleLogout = handleLogout; // Expose handleLogout function
    window.loadPrivateMessages = loadPrivateMessages; // Expose loadPrivateMessages function
});


const path = window.location.pathname;


// Function to load private messages for a user
function loadPrivateMessages(username) {
    setupWebSocket(username); // Set up the WebSocket connection with the current user's username
}