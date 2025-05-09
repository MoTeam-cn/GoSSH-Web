// WebSocket连接管理
let ws = null;
let sessionId = null;

// 连接状态
const STATUS = {
    CONNECTED: 'connected',
    DISCONNECTED: 'disconnected'
};

// 初始化终端
function initTerminal() {
    const term = new Terminal({
        cursorBlink: true,
        fontSize: 14,
        fontFamily: 'Menlo, Monaco, "Courier New", monospace',
        theme: {
            background: '#1e1e1e',
            foreground: '#ffffff'
        }
    });

    return term;
}

// 建立SSH连接
async function connect(host, port, username, password) {
    try {
        const response = await fetch('/api/auth', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ host, port, username, password })
        });

        const data = await response.json();
        if (data.success) {
            sessionId = data.sessionId;
            connectWebSocket();
            updateStatus(STATUS.CONNECTED);
            return true;
        } else {
            throw new Error(data.message || '连接失败');
        }
    } catch (error) {
        console.error('连接错误:', error);
        updateStatus(STATUS.DISCONNECTED);
        return false;
    }
}

// 建立WebSocket连接
function connectWebSocket() {
    if (!sessionId) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/ws/${sessionId}`;
    
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('WebSocket连接已建立');
    };

    ws.onmessage = (event) => {
        term.write(event.data);
    };

    ws.onclose = () => {
        console.log('WebSocket连接已关闭');
        updateStatus(STATUS.DISCONNECTED);
    };

    ws.onerror = (error) => {
        console.error('WebSocket错误:', error);
        updateStatus(STATUS.DISCONNECTED);
    };
}

// 更新连接状态显示
function updateStatus(status) {
    const statusElement = document.getElementById('connection-status');
    statusElement.className = `status ${status}`;
    statusElement.textContent = status === STATUS.CONNECTED ? '已连接' : '未连接';
}

// 断开连接
async function disconnect() {
    if (sessionId) {
        try {
            await fetch('/api/disconnect', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ sessionId })
            });
        } catch (error) {
            console.error('断开连接错误:', error);
        }
    }

    if (ws) {
        ws.close();
    }

    sessionId = null;
    updateStatus(STATUS.DISCONNECTED);
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    const term = initTerminal();
    term.open(document.getElementById('terminal'));

    // 连接表单提交处理
    document.getElementById('connect-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const success = await connect(
            formData.get('host'),
            formData.get('port'),
            formData.get('username'),
            formData.get('password')
        );

        if (success) {
            term.focus();
        }
    });

    // 断开连接按钮处理
    document.getElementById('disconnect-btn').addEventListener('click', disconnect);
}); 