document.addEventListener('DOMContentLoaded', async () => {
    const auth = new AuthBase('/api/auth');

    // 認証チェック
    const token = await auth.getToken();
    if (!token) {
        window.location.href = '/statics/login.html';
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
        window.location.href = '/statics/login.html';
        return;
    }

    // ユーザー一覧の取得と表示
    try {
        const gameData = await auth.Get("/game/list", {});
        const games = gameData.Data;

        const users = await auth.Get("/user/admin/users", {});
        if (users) {
            const tableBody = document.querySelector('#user-table tbody');
            
            let gameOptions = '<option value="">ゲームを選択</option>';
            if (games) {
                games.forEach(game => {
                    gameOptions += `<option value="${game.game_id}">${game.name} (開始: ${new Date(game.start_time * 1000).toLocaleString()})</option>`;
                });
            }

            users.forEach(user => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${user.user_id}</td>
                    <td>${user.name}</td>
                    <td>${user.comment}</td>
                    <td><a href="/statics/admin/walking_history/?userId=${user.user_id}">歩行履歴を見る</a></td>
                    <td>
                        <select class="game-id-select">${gameOptions}</select>
                        <button class="add-to-game-btn">追加</button>
                    </td>
                `;
                tableBody.appendChild(row);
            });

            // イベントリスナーを追加
            document.querySelectorAll('.add-to-game-btn').forEach(button => {
                button.addEventListener('click', async (event) => {
                    const row = event.target.closest('tr');
                    const userId = row.querySelector('td:first-child').textContent;
                    const gameId = row.querySelector('.game-id-select').value;

                    if (!gameId) {
                        alert('ゲームを選択してください。');
                        return;
                    }

                    try {
                        const result = await auth.Post('/game/admin/member/join', {
                            'Content-Type': 'application/json'
                        }, JSON.stringify({ user_id: userId, game_id: gameId }));

                        if (result) {
                            alert('ユーザーをゲームに追加しました。');
                        } else {
                            alert('ユーザーの追加に失敗しました。');
                        }
                    } catch (error) {
                        console.error('Failed to add user to game:', error);
                        alert('エラーが発生しました。');
                    }
                });
            });
        }
    } catch (error) {
        console.error('Failed to fetch users:', error);
        document.body.innerHTML += '<p>ユーザー一覧の取得に失敗しました。</p>';
    }
});
