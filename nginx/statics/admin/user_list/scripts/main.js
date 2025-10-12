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
        // 1. 全ゲームリストと全ユーザーリストを並行して取得
        const [gameData, users] = await Promise.all([
            auth.Get("/game/list", {}),
            auth.Get("/user/admin/users", {})
        ]);

        if (!users) {
            throw new Error('ユーザー一覧の取得に失敗しました。');
        }

        const games = gameData ? gameData.Data : [];

        // 2. 各ユーザーの参加ゲームリストを並行して取得
        const joinedGamesPromises = users.map(user => 
            auth.Get(`/game/admin/users/${user.user_id}/games`, {})
        );
        const joinedGamesResults = await Promise.all(joinedGamesPromises);

        // 3. HTMLを文字列として構築
        const tableBody = document.querySelector('#user-table tbody');
        let tableContent = '';

        let gameOptions = '<option value="">ゲームを選択</option>';
        if (games) {
            games.forEach(game => {
                gameOptions += `<option value="${game.game_id}">${game.name} (開始: ${new Date(game.start_time * 1000).toLocaleString()})</option>`;
            });
        }

        users.forEach((user, index) => {
            const joinedGames = joinedGamesResults[index];
            const joinedGamesText = (joinedGames && joinedGames.length > 0) 
                ? joinedGames.map(g => g.gameID).join(', ') 
                : 'なし';

            tableContent += `
                <tr>
                    <td>${user.user_id}</td>
                    <td>${user.name}</td>
                    <td>${user.comment}</td>
                    <td><a href="/statics/admin/walking_history/?userId=${user.user_id}">歩行履歴を見る</a></td>
                    <td>
                        <select class="game-id-select">${gameOptions}</select>
                        <button class="add-to-game-btn">追加</button>
                    </td>
                    <td>${joinedGamesText}</td>
                </tr>
            `;
        });

        // 4. DOMを一括で更新
        tableBody.innerHTML = tableContent;

        // 5. イベントリスナーをまとめて追加
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
                        // 表示を更新
                        const gamesCell = row.querySelector('td:last-child');
                        const currentGames = gamesCell.textContent;
                        if (currentGames === 'なし') {
                            gamesCell.textContent = gameId;
                        } else {
                            gamesCell.textContent += `, ${gameId}`;
                        }
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
