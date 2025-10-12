document.addEventListener('DOMContentLoaded', async () => {
    const auth = new AuthBase('/auth/');

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
            document.body.textContent = '';
            const h1 = document.createElement('h1');
            h1.textContent = 'アクセス権限がありません。';
            document.body.appendChild(h1);
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
        tableBody.textContent = '';

        const gameOptions = games ? games.map(game => {
            const option = document.createElement('option');
            option.value = game.game_id;
            option.textContent = `${game.name} (開始: ${new Date(game.start_time * 1000).toLocaleString()})`;
            return option;
        }) : [];

        users.forEach((user, index) => {
            const joinedGames = joinedGamesResults[index];
            const joinedGamesText = (joinedGames && joinedGames.length > 0) 
                ? joinedGames.map(g => g.gameID).join(', ') 
                : 'なし';

            const tr = document.createElement('tr');

            const td1 = document.createElement('td');
            td1.textContent = user.user_id;
            tr.appendChild(td1);

            const td2 = document.createElement('td');
            td2.textContent = user.name;
            tr.appendChild(td2);

            const td3 = document.createElement('td');
            td3.textContent = user.comment;
            tr.appendChild(td3);

            const td4 = document.createElement('td');
            const a = document.createElement('a');
            a.href = `/statics/admin/walking_history/?userId=${user.user_id}`;
            a.textContent = '歩行履歴を見る';
            td4.appendChild(a);
            tr.appendChild(td4);

            const td5 = document.createElement('td');
            const select = document.createElement('select');
            select.className = 'game-id-select';
            const defaultOption = document.createElement('option');
            defaultOption.value = '';
            defaultOption.textContent = 'ゲームを選択';
            select.appendChild(defaultOption);
            gameOptions.forEach(option => {
                select.appendChild(option.cloneNode(true));
            });
            td5.appendChild(select);
            const button = document.createElement('button');
            button.className = 'add-to-game-btn';
            button.textContent = '追加';
            td5.appendChild(button);
            tr.appendChild(td5);

            const td6 = document.createElement('td');
            td6.textContent = joinedGamesText;
            tr.appendChild(td6);

            tableBody.appendChild(tr);
        });

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

    } catch (error) {
        console.error('Failed to fetch users:', error);
        const p = document.createElement('p');
        p.textContent = 'ユーザー一覧の取得に失敗しました。';
        document.body.appendChild(p);
    }
});
