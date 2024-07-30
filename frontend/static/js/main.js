document.addEventListener('DOMContentLoaded', () => {
    const loginPopup = document.getElementById('login-popup');
    const registerPopup = document.getElementById('register-popup');
    const loginBtn = document.getElementById('login-btn');
    const registerBtn = document.getElementById('register-btn');
    const closeLogin = document.getElementById('close-login');
    const closeRegister = document.getElementById('close-register');
    const logoutBtn = document.getElementById('logout-btn');
    const userInfoContainer = document.getElementById('user-info');
    const userNameElement = document.getElementById('user-name');
    const userEmailElement = document.getElementById('user-email');
    const filterButton = document.getElementById('filter-button');
    const categoryFilter = document.getElementById('category-filter');

    // for profile user info
    const userInfoHeader = document.getElementById('user-info-header');
    const userUsernameHeaderElement = document.getElementById('user-username-header');
    const userEmailHeaderElement = document.getElementById('user-email-header');

    // Define profileLink
    const profileLink = document.getElementById('profile-link');

    const togglePopup = (popup, action) => {
        popup.style.display = action === 'open' ? 'block' : 'none';
    };

    loginBtn.addEventListener('click', (event) => {
        event.preventDefault();
        togglePopup(loginPopup, 'open');
    });

    registerBtn.addEventListener('click', (event) => {
        event.preventDefault();
        togglePopup(registerPopup, 'open');
    });

    closeLogin.addEventListener('click', () => {
        togglePopup(loginPopup, 'close');
    });

    closeRegister.addEventListener('click', () => {
        togglePopup(registerPopup, 'close');
    });

    window.addEventListener('click', (event) => {
        if (event.target === loginPopup) {
            togglePopup(loginPopup, 'close');
        }
        if (event.target === registerPopup) {
            togglePopup(registerPopup, 'close');
        }
    });

    const handleFormSubmit = (form, url) => {
        if (form) {
            form.addEventListener('submit', async (event) => {
                event.preventDefault();
                const formData = new FormData(form);
                const data = Object.fromEntries(formData);

                try {
                    const response = await fetch(url, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(data)
                    });

                    if (response.ok) {
                        const responseData = await response.json();
                        if (form.id === 'login-form') {
                            localStorage.setItem('user', JSON.stringify(responseData));
                            loadUser(); // Reload user information
                            togglePopup(loginPopup, 'close');
                        } else if (form.id === 'register-form') {
                            alert("Registered successfully");
                            togglePopup(registerPopup, 'close');
                        }
                    } else {
                        const responseData = await response.json();
                        alert(responseData.error);
                    }
                } catch (error) {
                    console.error(`${form.id} request failed:`, error);
                }
            });
        }
    };

    handleFormSubmit(document.getElementById('login-form'), '/login');
    handleFormSubmit(document.getElementById('register-form'), '/register');

    const toggleUserUI = (isLoggedIn) => {
        loginBtn.style.display = isLoggedIn ? 'none' : 'inline';
        registerBtn.style.display = isLoggedIn ? 'none' : 'inline';
        logoutBtn.style.display = isLoggedIn ? 'inline' : 'none';
        userInfoHeader.style.display = isLoggedIn ? 'inline' : 'none'; // Updated line

        // Show or hide the profile link
        if (profileLink) {
            profileLink.style.display = isLoggedIn ? 'inline' : 'none';
        }
    };

    const loadUser = async () => {
        try {
            const response = await fetch('/models/user');
            if (!response.ok) {
                throw new Error('Failed to fetch user data');
            }
            const user = await response.json();
            toggleUserUI(true);
            if (userNameElement) {
                userNameElement.textContent = `${user.name} ${user.surname}`;
            }
            if (userEmailElement) {
                userEmailElement.textContent = user.email;
            }
            if (userUsernameHeaderElement) {
                userUsernameHeaderElement.textContent = user.username;
            }
            if (userEmailHeaderElement) {
                userEmailHeaderElement.textContent = user.email;
            }
            if (window.location.pathname === '/profile.html') {
                document.getElementById('profile-name').textContent = `${user.name} ${user.surname}`;
                document.getElementById('profile-email').textContent = user.email;
            }
        } catch (error) {
            console.error('Error loading user:', error);
            toggleUserUI(false);
        }
    };

    loadUser();

    // Fetch threads only if on the forum page
    if (window.location.pathname === '/forum') {
        fetchThreads();
    }

    if (filterButton) {
        filterButton.addEventListener('click', () => {
            const category = categoryFilter.value;
            fetchThreads(category);
        });
    }

    logoutBtn.addEventListener('click', async () => {
        try {
            const response = await fetch('/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (response.ok) {
                const responseData = await response.json();
                alert(responseData.message);
                localStorage.removeItem('user');
                toggleUserUI(false);
                window.location.href = "/";  // Redirect to homepage after logout
            } else {
                const responseData = await response.json();
                alert(responseData.error);
            }
        } catch (error) {
            console.error('Logout request failed:', error);
        }
    });
});
