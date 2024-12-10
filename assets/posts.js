let postsJson = []; // Array to hold the fetched posts

function getPostIdFromURL() {
    // Check for /post/ path without ID
    if (window.location.pathname === '/post/' || window.location.pathname === '/post') {
        show404();
    }

    // Check for /post/{id} format
    const pathMatch = window.location.pathname.match(/\/post\/(\d+)/);
    if (pathMatch) {
        // Ensure the ID is a valid number
        const id = parseInt(pathMatch[1]);
        if (isNaN(id) || id <= 0) {
            return null;
        }
        return pathMatch[1];
    }

    // Fallback to query parameter
    const urlParams = new URLSearchParams(window.location.search);
    const queryId = urlParams.get('post');
    if (!queryId || isNaN(parseInt(queryId)) || parseInt(queryId) <= 0) {
        return null;
    }
    return queryId;
}

// Function to fetch posts from the backend and update visibility based on filter criteria
async function fetchAndFilterPosts() {
    await fetch("/getposts?username=" + currentUsername) // Fetch posts from the backend with the current username
        .then(response => response.json()) // Parse the response as JSON
        .then(posts => {
            postsJson = posts; // Store the fetched posts in the postsJson array
            renderPosts(posts); // Render the fetched posts

            const postId = getPostIdFromURL();
            if (postId) {
                const post = posts.find(p => p.post_id == postId);
                if (post) {
                    showSinglePost(postId);
                }
            }
            
        })
        .catch(error => console.error("Error fetching and filtering posts:", error)); // Log any errors
}

// Function to render posts dynamically
function renderPosts(posts) {
    const postsContainer = document.getElementById("postsContainer"); // Get the posts container element
    postsContainer.innerHTML = ""; // Clear existing posts
    posts.forEach(post => {
        // Create HTML for each post
        const postHTML = `
            <div class="card" id="${post.user_id}">
                <div onclick="showSinglePost(${post.post_id})">
                    <h3><a>${post.post_title}</a></h3>
                    <span>
                        by ${post.user_name}
                        ${post.categories.map(category => `<div class="chip">${category.cat_name}</div>`).join('')}
                    </span>
                </div>
                <div class="stats">
                    <div ${post.liked_by_user ? "" : 'class="grey"'}>
                        <div id="like_${post.post_id}">${post.likes}</div>
                        <div onclick="setLikes(${post.post_id}, true);">👍</div>
                    </div>
                    <div ${post.disliked_by_user ? "" : 'class="grey"'}>
                        <div id="dislike_${post.post_id}">${post.dislikes}</div>
                        <div onclick="setLikes(${post.post_id}, false);">👎</div>
                    </div>
                </div>
            </div>
        `;
        // Append the post HTML to the container
        postsContainer.innerHTML += postHTML;
    });
}


// Function to show a single post in detail
async function showSinglePost(id) {
    // Validate id parameter
    if (!id || isNaN(parseInt(id)) || parseInt(id) <= 0) {
        show404();
    }

    // Update URL to use /post/{id} format
    const newPath = `/post/${id}`;
    if (window.location.pathname !== newPath) {
        window.history.pushState({}, '', newPath);
    }

    // Clear existing post content
    const postDialog = document.getElementById("post");
    postDialog.innerHTML = "";

    try {
        // First try to find post in existing posts
        let post = postsJson.find(post => post.post_id == id);

        // If post not found in current posts array, fetch it
        if (!post) {
            const response = await fetch(`/getposts?username=${currentUsername}`);
            if (!response.ok) {
                throw new Error('Failed to fetch posts');
            }
            const posts = await response.json();
            post = posts.find(p => p.post_id == id);
            
            if (!post) {
                show404();
            }
            
            // Update posts array and render
            postsJson = posts;
            renderPosts(posts);
        }

        // Create HTML for the single post
        let postHTML = `
            <div class="card" id="${post.user_id}">
                <div>
                    <h3>${post.post_title}</h3>
                    <span>
                        by ${post.user_name}
                        ${post.categories.map(category => `<div class="chip">${category.cat_name}</div>`).join('')}
                    </span>
                </div>
                <div class="stats">
                    <div ${post.liked_by_user ? "" : 'class="grey"'}>
                        <div id="like_${post.post_id}">${post.likes}</div>
                        <div onclick="setLikes(${post.post_id}, true);">👍</div>
                    </div>
                    <div ${post.disliked_by_user ? "" : 'class="grey"'}>
                        <div id="dislike_${post.post_id}">${post.dislikes}</div>
                        <div onclick="setLikes(${post.post_id}, false);">👎</div>
                    </div>
                </div>
            </div>
            <div class="card">
                <p>${post.post_content}</p>
            </div>
        `;

        // Add comment form if user is logged in
        if (currentUsername) {
            postHTML += `
                <br>
                <div class="card">
                    <form id="commentForm">
                        <label>
                            enter comment text
                            <textarea required name="comment" rows="3"></textarea>
                        </label>
                        <br>
                        <br>
                        <input type="hidden" name="id" value="${post.post_id}">
                        <input type="submit" value="comment">
                    </form>
                </div>
            `;
        }

        // Add comments section if there are comments
        if (post.comments && post.comments.length > 0) {
            postHTML += "<h2>Comments</h2>";
            post.comments.forEach(comment => {
                postHTML += `
                    <div class="card comment">
                        <div>
                            <h3>${comment.comment_content}</h3>
                            <span>by ${comment.user_name}</span>
                        </div>
                        <div class="stats">
                            <div id="like_comment_wrapper_${comment.comment_id}" class="${comment.liked_by_user ? '' : 'grey'}">
                                <div id="like_c_${comment.comment_id}">
                                    ${comment.likes}
                                </div>
                                <div onclick="setCommentLikes(${comment.comment_id}, true);">👍</div>
                        </div>
                    <!-- Dislike Section for Comment -->
                        <div id="dislike_comment_wrapper_${comment.comment_id}" class="${comment.disliked_by_user ? '' : 'grey'}">
                            <div id="dislike_c_${comment.comment_id}">
                                ${comment.dislikes}
                            </div>
                            <div onclick="setCommentLikes(${comment.comment_id}, false);">👎</div>
                        </div>
                        </div>
                    </div>
                `;
            });
        } else if (post.comments && post.comments.length === 0) {
            postHTML += `
                <div class="card">
                    <p>No comments yet. Be the first to comment!</p>
                </div>
            `;
        }

        // Set the post content
        postDialog.innerHTML = postHTML;

        // Add event listener for comment form if it exists
        const commentForm = document.getElementById("commentForm");
        if (commentForm) {
            commentForm.addEventListener("submit", async (e) => {
                e.preventDefault();
                const success = await addComment(e);
                if (success) {
                    // Refresh posts and show updated post
                    await fetchAndFilterPosts();
                    showSinglePost(post.post_id);
                }
            });
        }

    } catch (error) {
        console.error("Error showing post:", error);
        show404();
    }
}

// Function to add a comment to a post
async function addComment(event) {
    event.preventDefault();
    try {
        await fetch("/addComment?username=" + currentUsername, {
            method: "POST",
            body: new FormData(event.target),
        }).then((response) => {
            response.json().then((json) => {
                console.log(json);
            }).catch((e) => console.error(e));
        }).catch((e) => console.error(e));
        return true;
    } catch (error) {
        alert(error);
    }
    return false;
}

// Function to set likes or dislikes for a post
async function setLikes(postID, isLike) {
    const response = await fetch("/likePost/?postID=" + postID + "&isLike=" + isLike + "&username=" + currentUsername);
    if (response.ok) {
        const post = postsJson.find(post => post.post_id == postID);
        const body = await response.json();
        const likeDiv = document.getElementById("like_" + postID);
        const dislikeDiv = document.getElementById("dislike_" + postID);

        likeDiv.innerHTML = body.likes;
        dislikeDiv.innerHTML = body.dislikes;

        if (body.actionTaken == 1) {
            likeDiv.parentElement.classList.add("grey");
            dislikeDiv.parentElement.classList.remove("grey");
        } else if (body.actionTaken == 0) {
            likeDiv.parentElement.classList.remove("grey");
            dislikeDiv.parentElement.classList.add("grey");
        } else {
            likeDiv.parentElement.classList.add("grey");
            dislikeDiv.parentElement.classList.add("grey");
        }
        fetchAndFilterPosts();
    } else {
        alert(await response.text());
    }
}

// Function to set likes or dislikes for a comment
async function setCommentLikes(commentID, isLike) {
    const response = await fetch("/likeComment/?id=" + commentID + "&isLike=" + isLike + "&username=" + currentUsername);
    if (response.ok) {
        const body = await response.json();

        // Locate DOM elements for like/dislike buttons
        const likeWrapper = document.getElementById(`like_comment_wrapper_${commentID}`);
        const dislikeWrapper = document.getElementById(`dislike_comment_wrapper_${commentID}`);
        const likeCount = document.getElementById(`like_c_${commentID}`);
        const dislikeCount = document.getElementById(`dislike_c_${commentID}`);
        const likeIcon = likeWrapper.querySelector('div:last-child');  // 👍 icon
        const dislikeIcon = dislikeWrapper.querySelector('div:last-child');  // 👎 icon

        // Update counts
        likeCount.innerHTML = body.likes;
        dislikeCount.innerHTML = body.dislikes;

        // Reset the active and grey classes
        likeWrapper.classList.add("grey");
        dislikeWrapper.classList.add("grey");

        if (body.actionTaken == 1) { // if liked
            likeWrapper.classList.remove("grey");
            //likeIcon.parentElement.classList.remove("grey");
            //dislikeWrapper.parentElement.classList.add("grey");
            //dislikeIcon.parentElement.classList.add("grey");
        } else if (body.actionTaken == 0) { // if disliked
            likeWrapper.parentElement.classList.remove("grey");
            likeIcon.parentElement.classList.remove("grey");
            dislikeWrapper.parentElement.classList.add("grey");
            dislikeIcon.parentElement.classList.add("grey");
        } else {
            likeWrapper.parentElement.classList.add("grey");
            likeIcon.parentElement.classList.add("grey");
            dislikeWrapper.parentElement.classList.add("grey");
            dislikeIcon.parentElement.classList.add("grey");
        }
        // fetchAndFilterPosts();
    } else {
        alert(await response.text());
    }
}

// Function to show the add post dialog
function addpost() {
    addPostDialog.showModal();
}

// Function to handle post submission
function postSubmitted(event) {
    event.preventDefault();
    let formData = new FormData(event.target);
    let categories = [];
    document.getElementsByName("cat").forEach((checkbox) => {
        if (checkbox.checked) {
            categories.push(checkbox.value);
        }
    });
    console.log(formData.set("cat", categories.join(",")));
    fetch("/addpost?username=" + currentUsername, { method: "post", body: formData, }).then((response) => {
        console.log(response);
        if (response.status == 200) {
            addPostDialog.close();
            addPostDialog.children[1].reset();
            fetchAndFilterPosts();
        } else {
            response.json().then((json) => {
                console.error(json);
            });
        }
    }).catch((error) => {
        console.error(error);
    });
}

//route handling function
function handleRoute() {
    const path = window.location.pathname;

    // Only allow root path and valid post URLs
    if (path === '/') {
        document.getElementById("post").innerHTML = "";
        document.getElementById("postsContainer").style.display = "block";
        return;
    }

    // Check if it's a post URL
    if (path.startsWith('/post/')) {
        const postId = getPostIdFromURL();
        // Show 404 if no postId or invalid format
        if (!postId) {
            show404();
            return;
        }

        // Check if post exists in our data
        const post = postsJson.find(p => p.post_id == postId);
        if (!post) {
            show404();
            return;
        }

        document.getElementById("postsContainer").style.display = "none";
        showSinglePost(postId);
        return;
    }

    // Show 404 for any other path
    show404();
}

// Update popstate event listener
window.addEventListener('popstate', handleRoute);

// Update the DOMContentLoaded event listener
document.addEventListener('DOMContentLoaded', function() {
    const path = window.location.pathname;
    
    // Handle initial route
    if (path === '/') {
        fetchAndFilterPosts();
        return;
    }

    if (path.startsWith('/post/')) {
        const postId = getPostIdFromURL();
        if (!postId) {
            show404();
        }

        if (!currentUsername) {
            pendingPostId = postId;
            loadTemplate('login');
        } else {
            fetchAndFilterPosts(); // This will also handle showing the post if it exists
        }
    }

    // Show 404 for any other path
    show404();
});

function returnToHome() {
    // Update URL
    window.history.pushState({}, '', '/');
    
    // Reset app content
    document.getElementById("app").innerHTML = `
        <div id="postsContainer"></div>
        <div id="post"></div>
    `;
    
    // Show posts container and fetch posts
    document.getElementById("postsContainer").style.display = "block";
    document.getElementById("post").innerHTML = "";
    fetchAndFilterPosts();
}

function handleRoute() {
    const path = window.location.pathname;

    // Only allow root path and valid post URLs
    if (path === '/') {
        document.getElementById("post").innerHTML = "";
        document.getElementById("postsContainer").style.display = "block";
        return;
    }

    // Check if it's a post URL with a valid ID
    if (path.startsWith('/post/')) {
        const postId = getPostIdFromURL();
        // Show 404 if no postId or invalid format
        if (!postId) {
            show404();
        }

        // Check if post exists in our data
        const post = postsJson.find(p => p.post_id == postId);
        if (!post) {
            show404();
        }

        document.getElementById("postsContainer").style.display = "none";
        showSinglePost(postId);
        return;
    }

    // Show 404 for any other path
    show404();
}

function loadHome() {
    // Reset app content
    document.getElementById("app").innerHTML = `
        <div id="postsContainer"></div>
        <div id="post"></div>
    `;
    
    // Show posts container and fetch posts
    document.getElementById("postsContainer").style.display = "block";
    document.getElementById("post").innerHTML = "";
    fetchAndFilterPosts();
}



