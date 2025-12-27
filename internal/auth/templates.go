package auth

const setupTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Beeper CLI Setup</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-deep: #050508;
            --bg-card: rgba(12, 12, 18, 0.8);
            --bg-input: #0a0a10;
            --border: rgba(255, 255, 255, 0.06);
            --border-focus: #6953f2;
            --text: #f0f0f5;
            --text-muted: #8888a0;
            --text-dim: #4a4a5c;
            --accent: #6953f2;
            --accent-blue: #0c52f9;
            --beeper-purple: #6953f2;
            --gradient: linear-gradient(225deg, #6953f2, #0c52f9);
            --accent-glow: rgba(105, 83, 242, 0.2);
            --accent-hover: #7a66f5;
            --success: #10b981;
            --success-glow: rgba(16, 185, 129, 0.15);
            --error: #f43f5e;
            --error-glow: rgba(244, 63, 94, 0.15);
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'DM Sans', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg-deep);
            color: var(--text);
            height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 1.5rem;
            position: relative;
            overflow: hidden;
        }

        /* Gradient mesh background */
        body::before {
            content: '';
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background:
                radial-gradient(ellipse 80% 50% at 20% -20%, rgba(105, 83, 242, 0.15) 0%, transparent 50%),
                radial-gradient(ellipse 60% 40% at 80% 120%, rgba(12, 82, 249, 0.12) 0%, transparent 50%);
            pointer-events: none;
            z-index: 0;
        }

        /* Subtle noise texture overlay */
        body::after {
            content: '';
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
            opacity: 0.03;
            pointer-events: none;
            z-index: 0;
        }

        .container {
            width: 100%;
            max-width: 520px;
            position: relative;
            z-index: 1;
        }

        /* Terminal header */
        .terminal-header {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            margin-bottom: 1.25rem;
            padding-bottom: 1rem;
            border-bottom: 1px solid var(--border);
        }

        .terminal-prompt {
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.875rem;
            color: var(--text-muted);
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .terminal-prompt::before {
            content: '$';
            color: #2D52F6;
        }

        /* Logo */
        .logo-section {
            text-align: center;
            margin-bottom: 1.5rem;
        }

        .logo {
            width: 56px;
            height: 56px;
            margin-bottom: 1rem;
            display: inline-block;
        }

        h1 {
            font-size: 1.5rem;
            font-weight: 600;
            letter-spacing: -0.02em;
            margin-bottom: 0.375rem;
        }

        .subtitle {
            color: var(--text-muted);
            font-size: 0.9375rem;
        }

        /* Card */
        .card {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 16px;
            padding: 1.5rem;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
            backdrop-filter: blur(20px);
            -webkit-backdrop-filter: blur(20px);
        }

        /* Form */
        .form-group {
            margin-bottom: 1.25rem;
        }

        label {
            display: block;
            font-size: 0.8125rem;
            font-weight: 500;
            color: var(--text-muted);
            margin-bottom: 0.5rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        input {
            width: 100%;
            padding: 0.875rem 1rem;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.9375rem;
            background: var(--bg-input);
            border: 1px solid var(--border);
            border-radius: 12px;
            color: var(--text);
            transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
        }

        input::placeholder {
            color: var(--text-dim);
        }

        input:focus {
            outline: none;
            border-color: var(--border-focus);
            box-shadow: 0 0 0 4px var(--accent-glow), 0 0 20px rgba(105, 83, 242, 0.15);
        }

        input:hover:not(:focus) {
            border-color: rgba(255, 255, 255, 0.1);
        }

        .input-hint {
            font-size: 0.75rem;
            color: var(--text-muted);
            margin-top: 0.375rem;
            font-family: 'JetBrains Mono', monospace;
            transition: color 0.2s;
        }

        .input-hint.error {
            color: var(--error);
        }

        input.error {
            border-color: var(--error);
            box-shadow: 0 0 0 3px var(--error-glow);
        }

        /* Buttons */
        .btn-group {
            display: flex;
            gap: 0.75rem;
            margin-top: 1.5rem;
        }

        button {
            flex: 1;
            padding: 0.875rem 1.5rem;
            font-family: 'DM Sans', sans-serif;
            font-size: 0.9375rem;
            font-weight: 600;
            border-radius: 12px;
            cursor: pointer;
            transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
            border: none;
        }

        .btn-secondary {
            background: rgba(255, 255, 255, 0.03);
            border: 1px solid var(--border);
            color: var(--text-muted);
        }

        .btn-secondary:hover {
            background: rgba(255, 255, 255, 0.06);
            border-color: rgba(255, 255, 255, 0.12);
            color: var(--text);
        }

        .btn-primary {
            background: var(--gradient);
            color: white;
            box-shadow: 0 4px 20px rgba(105, 83, 242, 0.35), 0 0 0 1px rgba(255, 255, 255, 0.1) inset;
            position: relative;
            overflow: hidden;
        }

        .btn-primary::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: linear-gradient(225deg, rgba(255,255,255,0.1), transparent 50%);
            pointer-events: none;
        }

        .btn-primary:hover {
            box-shadow: 0 6px 28px rgba(105, 83, 242, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.15) inset;
        }

        .btn-primary:active {
            box-shadow: 0 4px 16px rgba(105, 83, 242, 0.35);
        }

        button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            transform: none !important;
        }

        /* Status messages */
        .status {
            margin-top: 1.5rem;
            padding: 1rem;
            border-radius: 10px;
            font-size: 0.875rem;
            display: none;
            align-items: center;
            gap: 0.75rem;
            font-family: 'JetBrains Mono', monospace;
        }

        .status.show {
            display: flex;
        }

        .status.loading {
            background: var(--accent-glow);
            border: 1px solid rgba(106, 75, 229, 0.2);
            color: var(--beeper-purple);
        }

        .status.success {
            background: var(--success-glow);
            border: 1px solid rgba(34, 197, 94, 0.2);
            color: var(--success);
        }

        .status.error {
            background: var(--error-glow);
            border: 1px solid rgba(239, 68, 68, 0.2);
            color: var(--error);
        }

        .spinner {
            width: 16px;
            height: 16px;
            border: 2px solid currentColor;
            border-top-color: transparent;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        /* Help section */
        .help-section {
            margin-top: 1.25rem;
            padding-top: 1rem;
            border-top: 1px solid var(--border);
        }

        .help-title {
            font-size: 0.6875rem;
            font-weight: 500;
            color: var(--text-dim);
            text-transform: uppercase;
            letter-spacing: 0.08em;
            margin-bottom: 0.75rem;
        }

        .help-item {
            display: flex;
            align-items: flex-start;
            gap: 0.625rem;
            margin-bottom: 0.5rem;
            font-size: 0.75rem;
            color: var(--text-muted);
        }

        .help-item:last-child {
            margin-bottom: 0;
        }

        .help-icon {
            flex-shrink: 0;
            width: 20px;
            height: 20px;
            background: var(--bg-input);
            border-radius: 5px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.625rem;
            color: var(--text-dim);
        }

        .help-item code {
            font-family: 'JetBrains Mono', monospace;
            background: var(--bg-input);
            padding: 0.125rem 0.375rem;
            border-radius: 4px;
            font-size: 0.6875rem;
            color: #2D52F6;
        }

        /* Footer */
        .footer {
            text-align: center;
            margin-top: 1.25rem;
            font-size: 0.75rem;
            color: var(--text-dim);
        }

        .footer a {
            color: var(--text-muted);
            text-decoration: none;
            transition: color 0.2s;
        }

        .footer a:hover {
            color: #2D52F6;
        }

        .github-link {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
        }

        .github-link svg {
            opacity: 0.7;
            transition: opacity 0.2s;
        }

        .github-link:hover svg {
            opacity: 1;
        }

        /* Animations */
        .fade-in {
            animation: fadeIn 0.5s ease forwards;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .card { animation-delay: 0.1s; opacity: 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="terminal-header">
            <div class="terminal-prompt">
                beeper auth login
            </div>
        </div>

        <div class="logo-section fade-in">
            <div class="logo">
                <svg viewBox="0 0 42 54" xmlns="http://www.w3.org/2000/svg">
                    <defs>
                        <linearGradient id="beeperGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                            <stop offset="0%" style="stop-color:#6953f2"/>
                            <stop offset="100%" style="stop-color:#0c52f9"/>
                        </linearGradient>
                    </defs>
                    <path fill-rule="evenodd" clip-rule="evenodd" d="M6.3877 27.4864C6.3877 21.7626 11.0065 17.0532 16.8442 17.0532H25.4918C31.2667 17.0532 35.9482 21.7114 35.9482 27.4864C35.9414 30.5966 34.583 33.2757 32.6924 35.0576L30.7113 37.0319C33.5964 38.7151 35.9243 41.356 35.9243 44.9371C35.9243 50.2996 31.5772 54.1488 26.2147 54.1488H11.6043V36.5248C8.51598 34.7225 6.3877 31.3519 6.3877 27.4864ZM16.8442 19.9513C12.6186 19.9513 9.28579 23.3475 9.28579 27.4864C9.28579 31.6252 12.6186 35.0214 16.8442 35.0214H20.4444L18.4419 45.1466L30.6569 32.995C32.0562 31.6908 33.0501 29.7181 33.0501 27.4864C33.0501 23.2841 29.6087 19.906 25.4918 19.9513H16.8442Z" fill="url(#beeperGrad)"/>
                </svg>
            </div>
            <h1>Connect to Beeper</h1>
            <p class="subtitle">Configure your CLI to interact with Beeper Desktop</p>
        </div>

        <div class="card fade-in">
            <form id="setupForm" autocomplete="off">
                <div class="form-group">
                    <label for="apiToken">Bearer Token</label>
                    <input
                        type="password"
                        id="apiToken"
                        name="apiToken"
                        placeholder="Enter your Bearer token"
                        required
                    >
                    <div class="input-hint">
                        Beeper Desktop → Settings → Developers → Click <code>+</code>
                    </div>
                </div>

                <div class="btn-group">
                    <button type="button" id="testBtn" class="btn-secondary">Test Connection</button>
                    <button type="submit" id="submitBtn" class="btn-primary">Save & Connect</button>
                </div>

                <div id="status" class="status"></div>
            </form>

            <div class="help-section">
                <div class="help-title">How to get your token</div>
                <div class="help-item">
                    <span class="help-icon">1</span>
                    <span>Open Beeper Desktop app on this computer</span>
                </div>
                <div class="help-item">
                    <span class="help-icon">2</span>
                    <span>Go to Settings → Developers</span>
                </div>
                <div class="help-item">
                    <span class="help-icon">3</span>
                    <span>Click the <code>+</code> button to create a new token</span>
                </div>
                <div class="help-item">
                    <span class="help-icon">4</span>
                    <span>Copy the token and paste it above</span>
                </div>
            </div>
        </div>

        <div class="footer fade-in" style="animation-delay: 0.2s; opacity: 0;">
            <a href="https://github.com/salmonumbrella/beeper-cli" target="_blank" class="github-link">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
                    <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
                </svg>
                View on GitHub
            </a>
        </div>
    </div>

    <script>
        const form = document.getElementById('setupForm');
        const testBtn = document.getElementById('testBtn');
        const submitBtn = document.getElementById('submitBtn');
        const status = document.getElementById('status');
        const csrfToken = '{{.CSRFToken}}';

        const tokenInput = document.getElementById('apiToken');
        const inputHint = document.querySelector('.input-hint');
        const originalHint = inputHint.innerHTML;

        function showStatus(type, message) {
            status.className = 'status show ' + type;
            if (type === 'loading') {
                status.innerHTML = '<div class="spinner"></div><span>' + message + '</span>';
            } else {
                const icon = type === 'success' ? '&#10003;' : '&#10007;';
                status.innerHTML = '<span>' + icon + '</span><span>' + message + '</span>';
            }
        }

        function showInputError(message) {
            tokenInput.classList.add('error');
            inputHint.classList.add('error');
            inputHint.textContent = message;
            status.className = 'status';
        }

        function clearInputError() {
            tokenInput.classList.remove('error');
            inputHint.classList.remove('error');
            inputHint.innerHTML = originalHint;
        }

        function hideStatus() {
            status.className = 'status';
        }

        tokenInput.addEventListener('input', clearInputError);

        function getFormData() {
            return {
                token: document.getElementById('apiToken').value.trim()
            };
        }

        testBtn.addEventListener('click', async () => {
            const data = getFormData();

            if (!data.token) {
                showStatus('error', 'Please enter your token');
                return;
            }

            testBtn.disabled = true;
            submitBtn.disabled = true;
            showStatus('loading', 'Testing connection...');

            try {
                const response = await fetch('/validate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-Token': csrfToken
                    },
                    body: JSON.stringify(data)
                });

                const result = await response.json();

                if (result.success) {
                    showStatus('success', result.message);
                } else {
                    if (result.error && (result.error.includes('connection refused') || result.error.includes('Connection failed'))) {
                        showInputError('Beeper API not running — enable it in Settings → Developers');
                    } else if (result.error && result.error.includes('Invalid token')) {
                        showInputError('Invalid token — check your token and try again');
                    } else {
                        showInputError(result.error || 'Connection failed');
                    }
                }
            } catch (err) {
                showInputError('Request failed: ' + err.message);
            } finally {
                testBtn.disabled = false;
                submitBtn.disabled = false;
            }
        });

        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const data = getFormData();

            if (!data.token) {
                showStatus('error', 'Please enter your token');
                return;
            }

            testBtn.disabled = true;
            submitBtn.disabled = true;
            showStatus('loading', 'Saving credentials...');

            try {
                const response = await fetch('/submit', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-Token': csrfToken
                    },
                    body: JSON.stringify(data)
                });

                const result = await response.json();

                if (result.success) {
                    showStatus('success', 'Credentials saved! Redirecting...');
                    setTimeout(() => {
                        window.location.href = '/success';
                    }, 1000);
                } else {
                    if (result.error && (result.error.includes('connection refused') || result.error.includes('Connection failed'))) {
                        showInputError('Beeper API not running — enable it in Settings → Developers');
                    } else if (result.error && result.error.includes('Invalid token')) {
                        showInputError('Invalid token — check your token and try again');
                    } else {
                        showInputError(result.error || 'Connection failed');
                    }
                    testBtn.disabled = false;
                    submitBtn.disabled = false;
                }
            } catch (err) {
                showInputError('Request failed: ' + err.message);
                testBtn.disabled = false;
                submitBtn.disabled = false;
            }
        });
    </script>
</body>
</html>`

const successTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Setup Complete - Beeper CLI</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-deep: #050508;
            --bg-card: rgba(12, 12, 18, 0.8);
            --bg-input: #0a0a10;
            --border: rgba(255, 255, 255, 0.06);
            --text: #f0f0f5;
            --text-muted: #8888a0;
            --text-dim: #4a4a5c;
            --beeper-purple: #6953f2;
            --beeper-glow: rgba(105, 83, 242, 0.25);
            --gradient: linear-gradient(225deg, #6953f2, #0c52f9);
            --success: #10b981;
            --success-glow: rgba(16, 185, 129, 0.15);
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        body {
            font-family: 'DM Sans', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg-deep);
            color: var(--text);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 2rem;
            position: relative;
        }

        body::before {
            content: '';
            position: fixed;
            top: 0; left: 0; right: 0; bottom: 0;
            background:
                radial-gradient(ellipse 80% 50% at 50% -10%, rgba(105, 83, 242, 0.2) 0%, transparent 60%),
                radial-gradient(ellipse 60% 40% at 50% 110%, rgba(12, 82, 249, 0.15) 0%, transparent 50%);
            pointer-events: none;
        }

        body::after {
            content: '';
            position: fixed;
            top: 0; left: 0; right: 0; bottom: 0;
            background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
            opacity: 0.03;
            pointer-events: none;
        }

        .container {
            width: 100%;
            max-width: 560px;
            position: relative;
            z-index: 1;
            text-align: center;
        }

        .logo {
            width: 80px;
            height: 80px;
            margin: 0 auto 2rem;
            display: block;
            animation: scaleIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
            filter: drop-shadow(0 8px 40px var(--beeper-glow));
        }

        @keyframes scaleIn {
            from { transform: scale(0); }
            to { transform: scale(1); }
        }

        h1 {
            font-size: 2rem;
            font-weight: 600;
            letter-spacing: -0.02em;
            margin-bottom: 0.5rem;
            animation: fadeUp 0.5s ease 0.2s both;
        }

        .subtitle {
            color: var(--text-muted);
            font-size: 1rem;
            margin-bottom: 2.5rem;
            animation: fadeUp 0.5s ease 0.3s both;
        }

        @keyframes fadeUp {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Terminal card */
        .terminal {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 16px;
            overflow: hidden;
            text-align: left;
            animation: fadeUp 0.5s ease 0.4s both;
            box-shadow: 0 4px 32px rgba(0, 0, 0, 0.4);
        }

        .terminal-bar {
            background: var(--bg-input);
            padding: 0.75rem 1rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
            border-bottom: 1px solid var(--border);
        }

        .terminal-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
        }

        .terminal-dot.red { background: #ff5f57; }
        .terminal-dot.yellow { background: #febc2e; }
        .terminal-dot.green { background: #28c840; }

        .terminal-title {
            flex: 1;
            text-align: center;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
            color: var(--text-dim);
        }

        .terminal-body {
            padding: 1.5rem;
        }

        .terminal-line {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.875rem;
            margin-bottom: 1rem;
        }

        .terminal-line:last-child {
            margin-bottom: 0;
        }

        .terminal-prompt {
            color: var(--beeper-purple);
            user-select: none;
        }

        .terminal-text {
            color: var(--text);
        }

        .terminal-cursor {
            display: inline-block;
            width: 10px;
            height: 20px;
            background: var(--beeper-purple);
            animation: cursorBlink 1.2s step-end infinite;
            margin-left: 2px;
            vertical-align: middle;
        }

        @keyframes cursorBlink {
            0%, 50% { opacity: 1; }
            50.01%, 100% { opacity: 0; }
        }

        .terminal-output {
            color: var(--success);
            padding-left: 1.25rem;
            margin-top: -0.5rem;
            margin-bottom: 1rem;
        }

        .terminal-comment {
            color: var(--text-dim);
            font-style: italic;
        }

        /* Message */
        .message {
            margin-top: 2rem;
            padding: 1.25rem;
            background: rgba(106, 75, 229, 0.08);
            border: 1px solid rgba(106, 75, 229, 0.15);
            border-radius: 12px;
            animation: fadeUp 0.5s ease 0.5s both;
        }

        .message-icon {
            font-size: 1.5rem;
            margin-bottom: 0.5rem;
        }

        .message-title {
            font-weight: 600;
            margin-bottom: 0.25rem;
            color: var(--text);
        }

        .message-text {
            font-size: 0.875rem;
            color: var(--text-muted);
        }

        .footer {
            margin-top: 2rem;
            font-size: 0.8125rem;
            color: var(--text-dim);
            animation: fadeUp 0.5s ease 0.6s both;
        }
    </style>
</head>
<body>
    <div class="container">
        <svg class="logo" viewBox="0 0 42 54" xmlns="http://www.w3.org/2000/svg">
            <defs>
                <linearGradient id="beeperGradSuccess" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#6953f2"/>
                    <stop offset="100%" style="stop-color:#0c52f9"/>
                </linearGradient>
            </defs>
            <path fill-rule="evenodd" clip-rule="evenodd" d="M6.3877 27.4864C6.3877 21.7626 11.0065 17.0532 16.8442 17.0532H25.4918C31.2667 17.0532 35.9482 21.7114 35.9482 27.4864C35.9414 30.5966 34.583 33.2757 32.6924 35.0576L30.7113 37.0319C33.5964 38.7151 35.9243 41.356 35.9243 44.9371C35.9243 50.2996 31.5772 54.1488 26.2147 54.1488H11.6043V36.5248C8.51598 34.7225 6.3877 31.3519 6.3877 27.4864ZM16.8442 19.9513C12.6186 19.9513 9.28579 23.3475 9.28579 27.4864C9.28579 31.6252 12.6186 35.0214 16.8442 35.0214H20.4444L18.4419 45.1466L30.6569 32.995C32.0562 31.6908 33.0501 29.7181 33.0501 27.4864C33.0501 23.2841 29.6087 19.906 25.4918 19.9513H16.8442Z" fill="url(#beeperGradSuccess)"/>
        </svg>

        <h1>You're all set!</h1>
        <p class="subtitle">Beeper CLI is now connected and ready to use</p>

        <div class="terminal">
            <div class="terminal-bar">
                <span class="terminal-dot red"></span>
                <span class="terminal-dot yellow"></span>
                <span class="terminal-dot green"></span>
                <span class="terminal-title"></span>
            </div>
            <div class="terminal-body">
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-text">beeper chats list</span>
                </div>
                <div class="terminal-output">Fetching chats...</div>
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-text">beeper reminders list</span>
                </div>
                <div class="terminal-output">Found 2 reminders</div>
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-cursor"></span>
                </div>
            </div>
        </div>

        <div class="message">
            <div class="message-icon">&#8592;</div>
            <div class="message-title">Return to your terminal</div>
            <div class="message-text">You can close this window and start using the CLI. Try running <strong>beeper --help</strong> to see all available commands.</div>
        </div>

        <p class="footer">This window will close automatically.</p>
    </div>

    <script>
        // Signal completion to server
        fetch('/complete', { method: 'POST' }).catch(() => {});
    </script>
</body>
</html>`
