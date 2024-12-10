// Get the container element where notifications will be displayed
let notificationsContainer = document.getElementById("notifications");

// Initialize a unique ID counter for each notification
let id = 0;

// Function to add a new notification
function addNotification(user, message) {
    // Create a new div element for the notification
    let newNotification = document.createElement("div");
    
    // Add the "toast" class to the new notification for styling
    newNotification.classList.add("toast");
    
    // Set a unique ID for the new notification
    newNotification.id = `not_${id}`;
    
    // Set the inner HTML of the notification, truncating the message if it's too long
    newNotification.innerHTML = `💬 ${user}: ${message.length < 20 ? message : `${message.substring(0,20)}...`}`;
    
    // Append the new notification to the notifications container
    notificationsContainer.appendChild(newNotification);
    
    // Set a timeout to remove the notification after 5 seconds
    setTimeout(() => {
        // Remove the first child of the notifications container (the oldest notification)
        notificationsContainer.removeChild(notificationsContainer.children[0]);
    }, 5000);
    
    // Increment the unique ID counter for the next notification
    id++;
}