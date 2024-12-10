function handleRegister(event) {
    event.preventDefault(); // Prevent default form submission
    console.log("Form submission prevented");

    const formData = new FormData(event.target); // Get form data from the form element
    console.log("Form data collected");

    // Log form data for debugging
    for (let [key, value] of formData.entries()) {
        console.log(`${key}: ${value}`); // Log each key-value pair in the form data
    }

    fetch("/register", {
        method: "POST", // Specify the request method as POST
        body: formData // Attach the form data to the request body
    })
        .then(response => {
            console.log("Received response from server", response);
            if (!response.ok) {
                console.log("Response not OK, status:", response.status);
                response.text().then(err => {
                    alert(err);
                    throw new Error(err);
                }); // Parse and throw error if response is not OK
            }
            return response.json(); // Parse the JSON response
        })
        .then(data => {
            console.log("Parsed JSON data", data);
            if (data.status && data.status == "ok") {
                console.log("Registration successful");
                loadTemplate("login"); // Load login template on successful registration
            } else {
                console.log("Registration failed:", data.message);
                alert("Registration failed: " + data.message); // Show error message on failed registration
            }
        })
        .catch(error => {
            console.error("Error registering:", error); // Log error if any
        });
}

// Function to handle login form submission
async function handleLogin(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const username = formData.get("username");
    console.log(username);

    try {
        const response = await fetch("/login", {
            method: "POST",
            body: formData
        });

        // console.log("Received response from server", response);

        // Check if the response is JSON
        const contentType = response.headers.get("content-type");
        if (!contentType) {
            throw new Error("Unexpected response format: " + contentType);
        }
        if (!response.ok) {
            console.log("Response not OK, status:", response.status);
            const err = await response.json();
            throw new Error(err.message); // Parse and throw error if response is not OK
            return;
        }

        const data = await response.json(); // Parse the JSON response

        console.log("Parsed JSON data", data);

        if (data.status == "ok") {
            
            sessionStorage.setItem("session", data.session);
            currentUsername = data.username;
            // Check if there was a pending post
            if (pendingPostId) {
                loadTemplate("index");
            } else {
                loadTemplate("index");
            }
        } else {
            alert("Login failed: " + data.message);
        }
    } catch (error) {
        console.error("Error logging in:", error); // Log error if any
        alert("An error occurred while logging in. Please try again."); // Show a generic error message
    }
}

// Function to handle logout
function handleLogout() {
    // Hide the main content
    document.getElementById('app').style.display = 'none';
    // Immediately show the login dialog
    document.getElementById('loginDialog').showModal();

    fetch("/logout", {
        method: "POST"
    })
        .then(response => response.json())
        .then(data => {
            if (data.status == "ok") {
                loadTemplate("login");
            } else {
                alert("Logout failed: " + data.message);
            }
        })
        .catch(error => console.error("Error logging out:", error));
}

function logout() {
    fetch("/logout?username=" + currentUsername).then(response => {
        try {
            response.json().then(json => console.log(json));
            sessionStorage.removeItem("session");
            loadTemplate('login');
            if (socket) {
                socket.close(); // Close the WebSocket connection
            }
        } catch (error) {
            alert(error);
            console.error(error);
        }
    });
}


