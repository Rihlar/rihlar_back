document.addEventListener('DOMContentLoaded', async () => {
    const auth = new AuthBase('/api/auth');

    // 認証チェック
    const token = await auth.getToken();
    if (!token) {
        window.location.href = '/login.html';
        return;
    }

    // 管理者権限チェック
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        if (!payload.labels || !payload.labels.includes('admin')) {
            document.body.innerHTML = '<h1>アクセス権限がありません。</h1>';
            return;
        }
    } catch (e) {
        console.error('Token decode error:', e);
        window.location.href = '/login.html';
        return;
    }

    // ユーザー一覧の取得と表示
    try {
        const users = await auth.Get("/user/admin/users", {});
        if (users) {
            const tableBody = document.querySelector('#user-table tbody');
            users.forEach(user => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${user.user_id}</td>
                    <td>${user.name}</td>
                    <td>${user.comment}</td>
                    <td><a href="/statics/admin/walking_history/?userId=${user.user_id}">歩行履歴を見る</a></td>
                `;
                tableBody.appendChild(row);
            });
        }
    } catch (error) {
        console.error('Failed to fetch users:', error);
        document.body.innerHTML += '<p>ユーザー一覧の取得に失敗しました。</p>';
    }
});
