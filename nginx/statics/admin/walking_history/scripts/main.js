document.addEventListener('DOMContentLoaded', async () => {
    const auth = new AuthBase('/api/auth');
    let map = null; // マップオブジェクトを保持する変数
    let currentLayerGroup = null; // 現在のレイヤーグループを保持する変数

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

    const urlParams = new URLSearchParams(window.location.search);
    const userId = urlParams.get('userId');
    const gameId = urlParams.get('gameId'); // gameIdをURLパラメータから取得

    if (!userId) {
        document.body.innerHTML = '<h1>ユーザーIDが指定されていません。</h1>';
        return;
    }
    document.getElementById('user-id').textContent = userId;

    // 地図の初期化
    map = L.map('map').setView([35.681236, 139.767125], 13); // 初期中心を東京駅に設定
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);
    currentLayerGroup = L.layerGroup().addTo(map);

    // URLにgameIdがあるかどうかで処理を分岐
    if (gameId) {
        // gameIdがある場合、プルダウンを非表示にし、特定のゲームのログを直接表示
        document.querySelector('.game-selector').style.display = 'none';
        fetchAndDisplayLogs(userId, gameId);
    } else {
        // gameIdがない場合、従来通りゲーム選択のプルダウンをセットアップ
        setupGameSelection(userId);
    }

    // 特定のゲームの行動ログを取得して表示する関数
    async function fetchAndDisplayLogs(uid, gid) {
        try {
            const logs = await auth.Get(`/game/admin/games/${gid}/movement_logs/${uid}`, {});
            updateMap(logs);
        } catch (error) {
            console.error('Failed to fetch movement logs:', error);
            alert('行動ログの取得に失敗しました。');
        }
    }

    // ゲーム選択プルダウンをセットアップする関数
            async function setupGameSelection(uid) {
            try {
                const games = await auth.Get(`/game/admin/users/${uid}/games`, {});            const gameSelect = document.getElementById('game-select');
            if (games && games.length > 0) {
                games.forEach(game => {
                    const option = document.createElement('option');
                    option.value = game.gameID;
                    const startTime = new Date(game.startTime).toLocaleString('ja-JP');
                    option.textContent = `${game.gameID} (開始: ${startTime})`;
                    gameSelect.appendChild(option);
                });
            } else {
                const option = document.createElement('option');
                option.textContent = '参加しているゲームはありません';
                option.disabled = true;
                gameSelect.appendChild(option);
            }
        } catch (error) {
            console.error('Failed to fetch games:', error);
        }

        // ゲーム選択時のイベントリスナー
        document.getElementById('game-select').addEventListener('change', (event) => {
            const selectedGameId = event.target.value;
            if (selectedGameId) {
                fetchAndDisplayLogs(uid, selectedGameId);
            }
        });
    }

    // 地図を更新する関数
    function updateMap(logs) {
        // 既存のレイヤーをクリア
        currentLayerGroup.clearLayers();

        if (!logs || logs.length === 0) {
            alert('このゲームの行動ログはありません。');
            return;
        }

        const latlngs = logs.map(log => [log.latitude, log.longitude]);

        // ポリラインを作成して地図に追加
        const polyline = L.polyline(latlngs, { color: 'blue' }).addTo(currentLayerGroup);

        // 地図の表示範囲をポリラインに合わせる
        map.fitBounds(polyline.getBounds());

        // マーカーを一定間隔で追加
        const markerInterval = Math.max(1, Math.floor(logs.length / 20)); // 最大20個のマーカー
        logs.forEach((log, index) => {
            if (index % markerInterval === 0) {
                const marker = L.marker([log.latitude, log.longitude])
                    .addTo(currentLayerGroup)
                    .bindPopup(`<b>歩数:</b> ${log.steps}<br><b>時刻:</b> ${new Date(log.timeStamp * 1000).toLocaleString('ja-JP')}`);
            }
        });
    }
});
