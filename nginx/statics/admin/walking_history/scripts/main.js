// scripts/main.js

import { RouteMapper } from './route-mapper.js';

document.addEventListener('DOMContentLoaded', async () => {
    // 認証ベースのインスタンスは、AuthBaseがグローバルに読み込まれていることを前提とする
    const auth = new AuthBase('/api/auth');
    let map = null;
    let currentLayerGroup = null;

    // --- DOM要素 ---
    const gameSelect = document.getElementById('game-select');
    const modeDisplayBtn = document.getElementById('mode-display');
    const modeAddBtn = document.getElementById('mode-add');
    const routeInfoArea = document.getElementById('route-info-area');
    // -------------------------

    let currentMode = 'display'; // 'display' or 'add'
    let routeMapper = null; // RouteMapperインスタンスを保持

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

    const urlParams = new URLSearchParams(window.location.search);
    const userId = urlParams.get('userId');
    const gameId = urlParams.get('gameId');

    if (!userId) {
        document.body.innerHTML = '<h1>ユーザーIDが指定されていません。</h1>';
        return;
    }
    document.getElementById('user-id').textContent = userId;

    // 地図の初期化
    // 【重要】ダブルクリックでピンを打つため、Leaflet標準のダブルクリックズームを無効化
    map = L.map('map', { doubleClickZoom: false }).setView([35.681236, 139.767125], 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);
    currentLayerGroup = L.layerGroup().addTo(map);

    // RouteMapperの初期化
    routeMapper = new RouteMapper(map, auth);

    // 初期化後の処理
    if (gameId) {
        document.querySelector('.game-selector').style.display = 'none';
        fetchAndDisplayLogs(userId, gameId);
    } else {
        setupGameSelection(userId);
    }
    
    // 【削除】元の map.on('click') イベントは削除
    
    // モード切り替えボタンのイベントリスナー
    modeDisplayBtn.addEventListener('click', () => setMode('display'));
    modeAddBtn.addEventListener('click', () => setMode('add'));
    
    // ゲーム選択時のイベントリスナー
    gameSelect.addEventListener('change', (event) => {
        const selectedGameId = event.target.value;
        if (selectedGameId && currentMode === 'display') {
            fetchAndDisplayLogs(userId, selectedGameId);
        }
    });

    /**
     * モードを切り替える関数
     */
    function setMode(mode) {
        if (currentMode === mode) return;

        currentMode = mode;
        modeDisplayBtn.classList.remove('active');
        modeAddBtn.classList.remove('active');
        
        // RouteMapperのステータスを更新
        routeMapper.setActive(mode === 'add');

        if (mode === 'display') {
            modeDisplayBtn.classList.add('active');
            // 履歴表示モードのUI
            routeInfoArea.style.display = 'none';
            currentLayerGroup.addTo(map); // 履歴レイヤーを表示
            // 現在選択されているゲームの履歴を再表示
            if (gameSelect.value) {
                fetchAndDisplayLogs(userId, gameSelect.value);
            } else {
                 // 履歴がない場合、レイヤーはクリアされたまま
                currentLayerGroup.clearLayers();
            }
        } else {
            modeAddBtn.classList.add('active');
            // ログ追加モードのUI
            routeInfoArea.style.display = 'block';
            currentLayerGroup.clearLayers(); // 履歴レイヤーを非表示
        }
    }


    /**
     * 特定のゲームの行動ログを取得して表示する関数
     */
    async function fetchAndDisplayLogs(uid, gid) {
        try {
            const logs = await auth.Get(`/game/admin/games/${gid}/movement_logs/${uid}`, {});
            updateMap(logs);
        } catch (error) {
            console.error('Failed to fetch movement logs:', error);
            alert('行動ログの取得に失敗しました。');
        }
    }

    /**
     * ゲーム選択プルダウンをセットアップする関数
     */
    async function setupGameSelection(uid) {
        try {
            const games = await auth.Get(`/game/admin/users/${uid}/games`, {});
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
    }

    /**
     * 地図を更新する関数
     */
    function updateMap(logs) {
        // 既存のレイヤーをクリア
        currentLayerGroup.clearLayers();
        routeMapper.clearRoute(); // 追加ルートもクリア

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
