document.getElementById('registerForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const job = document.getElementById('job').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm_password').value;
    const errorMsg = document.getElementById('error-msg');
    const submitBtn = this.querySelector('button[type="submit"]');

    // Reset error
    errorMsg.textContent = '';
    errorMsg.classList.remove('visible');

    // Client-side validation
    if (password !== confirmPassword) {
        errorMsg.textContent = "Passwords do not match!";
        errorMsg.classList.add('visible');
        return;
    }

    // Loading state
    const originalBtnText = submitBtn.textContent;
    submitBtn.textContent = 'Creating Account...';
    submitBtn.disabled = true;

    try {
        const response = await fetch('/api/create_user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username,
                job,
                password,
                confirm_password: confirmPassword
            })
        });

        const data = await response.json();

        if (response.ok) {
            // Success
            submitBtn.textContent = 'Success!';
            submitBtn.style.background = '#10b981'; // Green
            setTimeout(() => {
                alert('Account Created Successfully! Please Log In.');
                window.location.href = 'index.html';
            }, 500);

        } else {
            // Error
            throw new Error(data.error || 'Registration failed');
        }
    } catch (error) {
        console.error('Registration Error:', error);
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
