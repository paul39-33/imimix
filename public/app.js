document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const pass = document.getElementById('password').value;
    const errorMsg = document.getElementById('error-msg');
    const submitBtn = this.querySelector('button[type="submit"]');

    // Reset error
    errorMsg.textContent = '';
    errorMsg.classList.remove('visible');

    // Loading state
    const originalBtnText = submitBtn.textContent;
    submitBtn.textContent = 'Logging in...';
    submitBtn.disabled = true;

    try {
        console.log('Attempting login with:', { username, pass });
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, pass })
        });

        const data = await response.json();

        if (response.ok) {
            // Success
            localStorage.setItem('token', data.access_token);
            localStorage.setItem('user', JSON.stringify(data.user));

            // Redirect or show success (for now alert/log as I don't have a dashboard page yet)
            // window.location.href = '/app/dashboard.html'; 
            submitBtn.textContent = 'Success!';
            submitBtn.style.background = '#10b981'; // Green
            setTimeout(() => {
                // alert('Login Successful!');
                window.location.href = 'dashboard.html';
            }, 500);

        } else {
            // Error
            throw new Error(data.error || 'Login failed');
        }
    } catch (error) {
        console.error('Login Error:', error);
        errorMsg.textContent = error.message;
        errorMsg.classList.add('visible');

        // Reset button
        submitBtn.textContent = originalBtnText;
        submitBtn.disabled = false;

        // Shake animation
        const card = document.querySelector('.login-card');
        card.style.transform = 'translateX(5px)';
        setTimeout(() => {
            card.style.transform = 'translateX(-5px)';
            setTimeout(() => {
                card.style.transform = 'translateX(0)';
            }, 100);
        }, 100);
    }
});
