const auth = new AuthBase('/auth/');

async function Init() {
    try {
        // 情報を取得
        const userData = await auth.GetInfo();

        if (userData == null) {
            // ログインにリダイレクト
            window.location.href = '../login.html';
            return;
        }

        console.log("ユーザーデータ:", userData);
    } catch (error) {
        console.error("認証エラー:", error);
        // ログインにリダイレクト
        window.location.href = '../login.html';
    }
}

// DOMContentLoadedイベントでスクリプトが実行されるようにする
document.addEventListener('DOMContentLoaded', () => {
    // Init関数をページのロード時に実行
    Init();

    // DOM要素の取得
    const addGameBtn = document.getElementById('addGameBtn');
    const gameModal = document.getElementById('gameModal');
    const gameModalContent = document.getElementById('gameModalContent');
    const cancelGameModalBtn = document.getElementById('cancelGameModalBtn');
    const gameForm = document.getElementById('gameForm');
    const modalTitle = document.getElementById('modalTitle');
    const gameNameInput = document.getElementById('gameName');
    const gameRegionSelect = document.getElementById('gameRegion');
    const gameTimeInput = document.getElementById('gameTime');
    const gameDurationInput = document.getElementById('gameDuration');
    const gameTypeSelect = document.getElementById('gameTypeSelect'); // セレクトボックスのIDを取得
    const gameListDiv = document.getElementById('gameList');

    const gameDetailsModal = document.getElementById('gameDetailsModal');
    const gameDetailsModalContent = document.getElementById('gameDetailsModalContent');
    const detailsModalGameName = document.getElementById('detailsModalGameName');
    const detailsModalRegion = document.getElementById('detailsModalRegion');
    const detailsModalTime = document.getElementById('detailsModalTime');
    const detailsModalDuration = document.getElementById('detailsModalDuration');
    const detailsModalGameType = document.getElementById('detailsModalGameType'); // ゲームタイプ表示要素
    const teamListDiv = document.getElementById('teamList');
    const noTeamsMessage = document.getElementById('noTeamsMessage');
    const closeDetailsModalBtn = document.getElementById('closeDetailsModalBtn');

    const membersModal = document.getElementById('membersModal');
    const membersModalContent = document.getElementById('membersModalContent');
    const memberListDiv = document.getElementById('memberList'); // 正しいIDで取得
    // メンバー追加関連のDOM要素は削除
    const closeMembersModalBtn = document.getElementById('closeMembersModalBtn');

    const confirmModal = document.getElementById('confirmModal');
    const confirmModalContent = document.getElementById('confirmModalContent');
    const confirmMessage = document.getElementById('confirmMessage');
    const cancelConfirmBtn = document.getElementById('cancelConfirmBtn');
    const confirmActionBtn = document.getElementById('confirmActionBtn');

    let games = []; // ゲームデータを格納する配列 (クライアントサイドのキャッシュ)
    let currentSelectedGameId = null; // 現在選択されているゲームのID (詳細表示用)
    let currentMembersGameId = null; // メンバー一覧表示中のゲームID
    let confirmCallback = null; // 確認モーダルのコールバック関数

    // 定義済みのリージョンデータ
    const regions = [
        { RegionID: "regionId-a3d3dcd1-7a73-4a31-9908-b9cab944280d", RegionName: "北海道" },
        { RegionID: "regionId-65051f6a-9e94-439d-a7f6-5c127ad0c885", RegionName: "東北" },
        { RegionID: "regionId-fb145c05-e0e5-4f22-86e1-9f40326faf31", RegionName: "関東" },
        { RegionID: "regionId-16d687a3-8eab-4b36-8563-b62514823fe8", RegionName: "中部" },
        { RegionID: "regionId-c161edb9-6aff-4244-8749-707bff2fa3be", RegionName: "関西" },
        { RegionID: "regionId-005344a4-523d-4aab-9d1e-d232322cf54e", RegionName: "中国" },
        { RegionID: "regionId-3501902d-8bab-40cd-926f-30c53d80efc5", RegionName: "四国" },
        { RegionID: "regionId-ef5aa179-53e0-481d-b64d-ae7654049a88", RegionName: "九州" },
    ];

    // ステータス数値と表示文字列のマッピング
    const statusMap = {
        0: '予定',      // 開始前
        1: '開催中',    // 開催中
        2: '終了'       // 終了済
    };

    // ====================================================================
    // バックエンドAPIシミュレーション関数群 (Promiseではなく実体を返す)
    // 実際のアプリケーションでは、ここにfetch()などを使ったAPI呼び出しを記述します。
    // ====================================================================

    const backendService = (() => {
        // 擬似的なバックエンドデータストア
        let backendData = {
            games: JSON.parse(localStorage.getItem('games_backend_sim')) || []
        };

        // バックエンドデータをローカルストレージに保存 (シミュレーション用)
        const saveBackendData = () => {
            localStorage.setItem('games_backend_sim', JSON.stringify(backendData.games));
        };

        return {
            /**
             * 全てのゲームを取得するAPIシミュレーション
             * @returns {Array} ゲームデータの配列
             */
            fetchGamesAPI: async () => { // async を追加
                console.log("API: Fetching all games from /game/list...");
                const token = await auth.getToken(); // トークンを取得

                try {
                    const response = await fetch('/game/list', { // await を追加
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': `${token}` // Authorization ヘッダーを追加
                        },
                    });

                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`);
                    }

                    const data = await response.json(); // await を追加
                    console.log("Fetched game list:", data);

                    // APIのレスポンス形式に合わせてデータを整形
                    const fetchedGames = data.Data.map(game => ({
                        game_id: game.game_id,
                        name: game.name,
                        region_id: game.region_id,
                        status: statusMap[game.status], // 数値を文字列に変換
                        start_time: game.start_time, // Unix timestamp in seconds
                        dulation_date: game.dulation_date,
                        members: game.members || [], // membersはオブジェクトの配列
                        teams: game.teams || [],     // teamsはオブジェクトの配列
                        // gameTypeはバックエンドから提供されないため、フロントエンドで推測
                        gameType: (game.teams && game.teams.length > 0) ? 'team' : 'personal' // チームの有無で判断
                    }));

                    backendData.games = fetchedGames; // 取得したデータをローカルストレージに保存
                    saveBackendData();
                    return [...backendData.games]; // データのコピーを返す
                } catch (error) {
                    console.error("Error fetching games:", error);
                    throw error; // エラーを再スローして呼び出し元で処理させる
                }
            },

            /**
             * 新しいゲームを追加するAPIシミュレーション
             * @param {Object} newGameData - 追加するゲームのデータ
             * @returns {Object} 追加されたゲームのデータ
             */
            addGameAPI: async (newGameData) => {
                console.log("API: Adding new game to /gcore/create (simulated POST)...", newGameData);
                // GoのCreateGameArgs構造体に合わせたデータを作成
                const apiPayload = {
                    name: newGameData.name,
                    region_id: newGameData.region,
                    start_time: new Date(newGameData.time).getTime(), // JavaScriptのミリ秒単位のUnixタイムスタンプ
                    dulation_date: newGameData.duration,
                };

                // トークンを取得
                const token = await auth.getToken();

                // ここで実際のfetch()呼び出しをシミュレート
                const response = await fetch('/game/create', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `${token}`
                    },
                    body: JSON.stringify(apiPayload)
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.json();
                console.log(data);

                // シミュレーションのため、直接データを追加
                // 実際のAPIでは、APIが返すデータにIDなどが含まれるため、そちらを使用するのが望ましい
                const gameWithId = {
                    game_id: data.game_id || Date.now().toString(), // APIが返すIDを使用
                    name: newGameData.name,
                    region_id: newGameData.region,
                    start_time: new Date(newGameData.time).getTime(),
                    dulation_date: newGameData.duration,
                    gameType: newGameData.gameType,
                    status: '予定', // 新規作成時は「予定」
                    members: [], // 新規作成時は空
                    teams: []    // 新規作成時は空
                };
                backendData.games.push(gameWithId);
                saveBackendData();
                return gameWithId; // UIに返すデータは、追加されたゲームオブジェクト
            },

            /**
             * ゲームを削除するAPIシミュレーション
             * @param {string} gameId - 削除するゲームのID
             */
            deleteGameAPI: async (gameId) => { // async を追加
                console.log("API: Deleting game...", gameId);
                const token = await auth.getToken(); // トークンを取得

                try {
                    const response = await fetch('/game/delete', { // DELETE リクエスト
                        method: 'DELETE',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': `${token}`,
                            'GameID': gameId // GameID ヘッダーを追加
                        },
                    });

                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`);
                    }

                    // 成功レスポンスをログに出力
                    const data = await response.json();
                    console.log("Game deletion response:", data);

                    // バックエンドで削除されるため、クライアントサイドのデータ操作は不要
                    // loadGames() が呼び出し元でリストを再ロードする
                } catch (error) {
                    console.error("Failed to delete game via API:", error);
                    throw error;
                }
            },

            /**
             * ゲームのステータスを更新するAPIシミュレーション
             * @param {string} gameId - 更新するゲームのID
             * @param {number} newStatusValue - 新しいステータス値 (0:予定, 1:開催中, 2:終了)
             * @returns {Object} 更新されたゲームのデータ
             */
            updateGameStatusAPI: async (gameId, newStatusValue) => { // async を追加
                console.log("API: Updating game status...", gameId, newStatusValue);
                const token = await auth.getToken(); // トークンを取得

                try {
                    // 仮のステータス更新エンドポイントへのPOSTリクエスト
                    const response = await fetch('/game/update_status', { // 仮のAPIエンドポイント
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': `${token}`
                        },
                        body: JSON.stringify({
                            game_id: gameId,
                            status: newStatusValue
                        })
                    });

                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`);
                    }

                    const data = await response.json();
                    console.log("Status update response:", data);

                    // クライアントサイドのデータを更新 (APIレスポンスから直接更新)
                    const gameIndex = backendData.games.findIndex(game => game.game_id === gameId);
                    if (gameIndex !== -1) {
                        backendData.games[gameIndex].status = statusMap[newStatusValue]; // 数値を文字列に変換して更新
                        saveBackendData();
                        return backendData.games[gameIndex];
                    } else {
                        throw new Error("Game not found for status update (client-side update failed).");
                    }
                } catch (error) {
                    console.error("Failed to update game status via API:", error);
                    throw error;
                }
            },

            /**
             * 特定のゲームのチーム一覧を取得するAPIシミュレーション
             * @param {string} gameId - チーム一覧を取得するゲームのID
             * @returns {Array} チームデータの配列
             */
            fetchTeamsAPI: (gameId) => {
                console.log("API: Fetching teams for game...", gameId);
                const game = backendData.games.find(g => g.game_id === gameId); // game_id を使用
                if (game) {
                    return [...(game.teams || [])]; // データのコピーを返す
                } else {
                    throw new Error("Game not found for fetching teams.");
                }
            },

            /**
             * 特定のゲームにメンバーを追加するAPIシミュレーション
             * @param {string} gameId - メンバーを追加するゲームのID
             * @param {Object} memberData - 追加するメンバーのデータ (user_id, user_name, points)
             * @returns {Array} 更新されたメンバーリスト
             */
            addMemberAPI: (gameId, memberData) => {
                console.log("API: Adding member to game...", gameId, memberData);
                const gameIndex = backendData.games.findIndex(game => game.game_id === gameId);
                if (gameIndex !== -1) {
                    // user_id で重複チェック
                    if (!backendData.games[gameIndex].members.some(m => m.user_id === memberData.user_id)) {
                        backendData.games[gameIndex].members.push(memberData);
                        saveBackendData();
                        return [...backendData.games[gameIndex].members];
                    } else {
                        throw new Error("Member already exists.");
                    }
                } else {
                    throw new Error("Game not found for adding member.");
                }
            },

            /**
             * 特定のゲームからメンバーを削除するAPIシミュレーション
             * @param {string} gameId - メンバーを削除するゲームのID
             * @param {number} memberIndex - 削除するメンバーのインデックス
             * @returns {Array} 更新されたメンバーリスト
             */
            removeMemberAPI: (gameId, memberIndex) => {
                console.log("API: Removing member from game...", gameId, memberIndex);
                const gameIndex = backendData.games.findIndex(game => game.game_id === gameId); // game_id を使用
                if (gameIndex !== -1 && backendData.games[gameIndex].members) {
                    if (memberIndex >= 0 && memberIndex < backendData.games[gameIndex].members.length) {
                        backendData.games[gameIndex].members.splice(memberIndex, 1);
                        saveBackendData();
                        return [...backendData.games[gameIndex].members];
                    } else {
                        throw new Error("Member index out of bounds.");
                    }
                } else {
                    throw new Error("Game not found or no members for removal.");
                }
            },

            /**
             * 特定のゲームからチームを削除するAPIシミュレーション
             * @param {string} gameId - チームを削除するゲームのID
             * @param {number} teamIndex - 削除するチームのインデックス
             * @returns {Array} 更新されたチームリスト
             */
            removeTeamAPI: (gameId, teamIndex) => {
                console.log("API: Removing team from game...", gameId, teamIndex);
                const gameIndex = backendData.games.findIndex(game => game.game_id === gameId); // game_id を使用
                if (gameIndex !== -1 && backendData.games[gameIndex].teams) {
                    if (teamIndex >= 0 && teamIndex < backendData.games[gameIndex].teams.length) {
                        backendData.games[gameIndex].teams.splice(teamIndex, 1);
                        saveBackendData();
                        return [...backendData.games[gameIndex].teams];
                    } else {
                        throw new Error("Team index out of bounds.");
                    }
                } else {
                    throw new Error("Game not found or no teams for removal.");
                }
            }
        };
    })(); // 即時実行関数でbackendServiceを定義

    // ====================================================================
    // UIロジック
    // ====================================================================

    // リージョン選択ドロップダウンを初期化
    const populateRegionsDropdown = () => {
        gameRegionSelect.innerHTML = '<option value="">リージョンを選択してください</option>'; // デフォルトオプション
        regions.forEach(region => {
            const option = document.createElement('option');
            option.value = region.RegionID;
            option.textContent = region.RegionName;
            gameRegionSelect.appendChild(option);
        });
    };

    // RegionIDからRegionNameを取得するヘルパー関数
    const getRegionNameById = (id) => {
        const region = regions.find(r => r.RegionID === id);
        return region ? region.RegionName : '不明なリージョン';
    };

    // ゲームデータをロード (API経由)
    const loadGames = async () => { // async を追加
        try {
            games = await backendService.fetchGamesAPI(); // await を追加
            renderGames(); // ゲームリストをレンダリング
        } catch (error) {
            console.error("Failed to load games:", error);
            openConfirmModal('ゲームの読み込みに失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
        }
    };

    // モーダル表示/非表示のヘルパー関数
    const showModal = (modal, content) => {
        modal.classList.remove('hidden');
        setTimeout(() => {
            content.classList.remove('opacity-0', 'scale-95');
            content.classList.add('opacity-100', 'scale-100');
        }, 10); // 小さな遅延でトランジションをトリガー
    };

    const hideModal = (modal, content) => {
        content.classList.remove('opacity-100', 'scale-100');
        content.classList.add('opacity-0', 'scale-95');
        setTimeout(() => {
            modal.classList.add('hidden');
        }, 300); // トランジションが完了するまで待つ
    };

    // ゲームリストのレンダリング
    const renderGames = () => {
        gameListDiv.innerHTML = ''; // リストをクリア

        if (games.length === 0) {
            gameListDiv.innerHTML = '<p class="text-center text-gray-500">ゲームがありません。新しいゲームを作成してください。</p>';
            return;
        }

        games.forEach(game => {
            const gameCard = document.createElement('div');
            // ゲームカード全体をクリック可能にする
            gameCard.className = `bg-white rounded-xl shadow-md p-6 border-l-4 game-card-clickable ${game.status === '予定' ? 'border-gray-500' :
                game.status === '開催中' ? 'border-gray-700' :
                    'border-gray-900'
                }`;
            gameCard.dataset.id = game.game_id; // game_id を使用

            gameCard.innerHTML = `
                <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4">
                    <h3 class="text-2xl font-bold text-gray-900 mb-2 sm:mb-0">${game.name}</h3>
                    <span class="px-3 py-1 rounded-full text-sm font-semibold ${game.status === '予定' ? 'status-予定' :
                    game.status === '開催中' ? 'status-開催中' :
                        'status-終了'
                }">${game.status}</span>
                </div>
                <p class="text-gray-700 mb-2"><strong>リージョン:</strong> ${getRegionNameById(game.region_id)}</p>
                <p class="text-gray-700 mb-2"><strong>開催時間:</strong> ${new Date(game.start_time * 1000).toLocaleString('ja-JP', { // 秒をミリ秒に変換
                    year: 'numeric', month: 'numeric', day: 'numeric',
                    hour: '2-digit', minute: '2-digit', hour12: false
                })}</p>
                <p class="text-gray-700 mb-2"><strong>開催期間:</strong> ${game.dulation_date} 日間</p>
                <p class="text-gray-700 mb-4"><strong>ゲームタイプ:</strong> ${game.gameType === 'personal' ? '個人のゲーム' : 'チームのゲーム'}</p>
                <div class="flex flex-wrap gap-3 mt-4">
                    <button data-id="${game.game_id}" class="delete-btn bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-200">削除</button>
                    <button data-id="${game.game_id}" class="members-btn bg-gray-600 hover:bg-gray-700 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-200">メンバー (${game.members ? game.members.length : 0})</button>
                    ${game.status === '予定' || game.status === '終了' ? `<button data-id="${game.game_id}" class="start-btn bg-gray-700 hover:bg-gray-800 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-200">開始</button>` : ''}
                    ${game.status === '開催中' ? `<button data-id="${game.game_id}" class="end-btn bg-gray-600 hover:bg-gray-700 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-300">終了</button>` : ''}
                </div>
            `;
            gameListDiv.appendChild(gameCard);
        });

        // イベントリスナーを再割り当て
        addEventListenersToGameButtons();
    };

    // ゲームボタンにイベントリスナーを追加
    const addEventListenersToGameButtons = () => {
        // ゲームカードクリックで詳細モーダルを開く
        document.querySelectorAll('.game-card-clickable').forEach(card => {
            card.onclick = (e) => {
                // ボタンクリックイベントがカード全体に伝播しないようにする
                if (e.target.tagName === 'BUTTON') {
                    return;
                }
                openGameDetailsModal(card.dataset.id);
            };
        });
        document.querySelectorAll('.delete-btn').forEach(button => {
            button.onclick = async (e) => { // async を追加
                openConfirmModal('本当にこのゲームを削除しますか？', async () => { // async を追加
                    try {
                        await backendService.deleteGameAPI(e.target.dataset.id); // await を追加
                        loadGames(); // 削除後にリストを再ロード
                        hideModal(confirmModal, confirmModalContent);
                    } catch (error) {
                        console.error("Failed to delete game:", error);
                        openConfirmModal('ゲームの削除に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                    }
                });
            };
        });
        document.querySelectorAll('.members-btn').forEach(button => {
            button.onclick = (e) => openMembersModal(e.target.dataset.id);
        });
        document.querySelectorAll('.start-btn').forEach(button => {
            button.onclick = async (e) => { // async を追加
                try {
                    await backendService.updateGameStatusAPI(e.target.dataset.id, 1); // ステータスを「開催中」(1)に更新
                    loadGames(); // ステータス更新後にリストを再ロード
                } catch (error) {
                    console.error("Failed to start game:", error);
                    openConfirmModal('ゲームの開始に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                }
            };
        });
        document.querySelectorAll('.end-btn').forEach(button => {
            button.onclick = async (e) => { // async を追加
                try {
                    await backendService.updateGameStatusAPI(e.target.dataset.id, 2); // ステータスを「終了」(2)に更新
                    loadGames(); // ステータス更新後にリストを再ロード
                } catch (error) {
                    console.error("Failed to end game:", error);
                    openConfirmModal('ゲームの終了に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                }
            };
        });
    };

    // ゲーム作成モーダルを開く
    addGameBtn.addEventListener('click', () => {
        gameForm.reset(); // フォームをリセット
        gameRegionSelect.value = ""; // リージョン選択をリセット
        gameDurationInput.value = 1; // 期間をデフォルト値にリセット
        gameTypeSelect.value = 'personal'; // ゲームタイプをデフォルトに設定
        showModal(gameModal, gameModalContent);
    });

    // ゲームモーダルを閉じる
    cancelGameModalBtn.addEventListener('click', () => {
        hideModal(gameModal, gameModalContent);
    });

    // ゲームフォームの送信処理 (作成のみ)
    gameForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const gameName = gameNameInput.value;
        const gameRegion = gameRegionSelect.value; // 選択されたRegionIDを取得
        const gameTime = gameTimeInput.value;
        const gameDuration = parseInt(gameDurationInput.value, 10); // 期間を数値として取得
        const selectedGameType = gameTypeSelect.value; // 選択されたゲームタイプを取得

        if (!gameRegion) {
            openConfirmModal('リージョンを選択してください。', () => hideModal(confirmModal, confirmModalContent), 'OK');
            return;
        }
        if (isNaN(gameDuration) || gameDuration < 1) {
            openConfirmModal('開催期間は1以上の数値を入力してください。', () => hideModal(confirmModal, confirmModalContent), 'OK');
            return;
        }

        const newGameData = {
            name: gameName,
            region: gameRegion,
            time: gameTime,
            duration: gameDuration,
            gameType: selectedGameType, // ゲームタイプを追加
            status: '予定', // 初期ステータスは「予定」
            members: [], // 初期メンバーは空
            teams: [] // 初期チームは空
        };

        try {
            await backendService.addGameAPI(newGameData);
            loadGames(); // 追加後にリストを再ロード
            hideModal(gameModal, gameModalContent);
        } catch (error) {
            console.error("Failed to add game:", error);
            openConfirmModal('ゲームの作成に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
        }
    });

    // ゲーム詳細モーダルを開く
    const openGameDetailsModal = (id) => {
        const game = games.find(g => g.game_id === id); // game_id を使用
        if (game) {
            currentSelectedGameId = id;
            detailsModalGameName.textContent = game.name;
            detailsModalRegion.textContent = getRegionNameById(game.region_id); // region_id を使用
            detailsModalTime.textContent = new Date(game.start_time * 1000).toLocaleString('ja-JP', { // 秒をミリ秒に変換
                year: 'numeric', month: 'numeric', day: 'numeric',
                hour: '2-digit', minute: '2-digit', hour12: false
            });
            detailsModalDuration.textContent = game.dulation_date; // dulation_date を使用
            detailsModalGameType.textContent = game.gameType === 'personal' ? '個人のゲーム' : 'チームのゲーム'; // ゲームタイプを表示

            try {
                const teamsData = backendService.fetchTeamsAPI(id);
                renderTeams(teamsData); // チームリストをレンダリング
            }
            catch (error) {
                console.error("Failed to fetch teams:", error);
                openConfirmModal('チーム情報の取得に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                renderTeams([]); // エラー時は空のリストを表示
            }
            showModal(gameDetailsModal, gameDetailsModalContent);
        }
    };

    // チームリストをレンダリング
    const renderTeams = (teams) => {
        teamListDiv.innerHTML = ''; // リストをクリア
        if (teams && teams.length > 0) {
            noTeamsMessage.classList.add('hidden');
            teams.forEach((team, index) => {
                const teamItem = document.createElement('div');
                teamItem.className = 'flex justify-between items-center bg-gray-50 p-3 rounded-md shadow-sm';
                teamItem.innerHTML = `
                    <span class="text-gray-800 font-medium">チームID: ${team.team_id}</span>
                    <span class="text-gray-600">ポイント: ${team.points}</span>
                    <button data-index="${index}" class="remove-team-btn text-red-500 hover:text-red-700 font-semibold text-sm">削除</button>
                `;
                teamListDiv.appendChild(teamItem);
            });
            // チーム削除ボタンにイベントリスナーを追加
            document.querySelectorAll('.remove-team-btn').forEach(button => {
                button.onclick = (e) => {
                    const teamIndex = parseInt(e.target.dataset.index);
                    openConfirmModal('本当にこのチームを削除しますか？', () => { // Added confirmation for team deletion
                        try {
                            backendService.removeTeamAPI(currentSelectedGameId, teamIndex);
                            // 削除後、詳細モーダルを閉じずにチームリストのみ更新
                            const updatedTeams = backendService.fetchTeamsAPI(currentSelectedGameId);
                            renderTeams(updatedTeams);
                            hideModal(confirmModal, confirmModalContent); // Hide confirmation modal
                        } catch (error) {
                            console.error("Failed to remove team:", error);
                            openConfirmModal('チームの削除に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                        }
                    });
                };
            });
        } else {
            noTeamsMessage.classList.remove('hidden');
        }
    };

    // ゲーム詳細モーダルを閉じる
    closeDetailsModalBtn.addEventListener('click', () => {
        hideModal(gameDetailsModal, gameDetailsModalContent);
        currentSelectedGameId = null;
    });


    // メンバー一覧モーダルを開く
    const openMembersModal = (id) => {
        currentMembersGameId = id;
        const game = games.find(g => g.game_id === currentMembersGameId); // game_id を使用
        renderMembers(game ? game.members : []); // メンバーリストをレンダリング
        showModal(membersModal, membersModalContent);
    };

    // メンバー一覧をレンダリング
    const renderMembers = (members) => {
        memberListDiv.innerHTML = '';
        if (members && members.length > 0) {
            members.forEach((member, index) => {
                const memberItem = document.createElement('div');
                memberItem.className = 'flex justify-between items-center bg-gray-50 p-3 rounded-md shadow-sm';
                memberItem.innerHTML = `
                    <span class="text-gray-800">${member.user_name || member.user_id}</span> <!-- Display user_name or user_id -->
                    <button data-index="${index}" class="remove-member-btn text-red-500 hover:text-red-700 font-semibold text-sm">削除</button>
                `;
                memberListDiv.appendChild(memberItem);
            });
            // メンバー削除ボタンにイベントリスナーを追加
            document.querySelectorAll('.remove-member-btn').forEach(button => {
                button.onclick = (e) => {
                    const memberIndex = parseInt(e.target.dataset.index);
                    openConfirmModal('本当にこのメンバーを削除しますか？', () => { // Added confirmation for member deletion
                        try {
                            backendService.removeMemberAPI(currentMembersGameId, memberIndex);
                            // 削除後、メンバーリストとゲームリストを再ロード
                            loadGames(); // Reloads all games, which includes updated member lists
                            const game = games.find(g => g.game_id === currentMembersGameId); // Fetch updated game object
                            renderMembers(game ? game.members : []); // Re-render members for the current game
                            hideModal(confirmModal, confirmModalContent); // Hide confirmation modal
                        } catch (error) {
                            console.error("Failed to remove member:", error);
                            openConfirmModal('メンバーの削除に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                        }
                    });
                };
            });
        } else {
            memberListDiv.innerHTML = '<p class="text-center text-gray-500">このゲームにはまだメンバーがいません。</p>';
        }
    };

    // メンバー追加機能は削除されたため、addMemberBtnのイベントリスナーも削除

    // メンバーモーダルを閉じる
    closeMembersModalBtn.addEventListener('click', () => {
        hideModal(membersModal, membersModalContent);
        currentMembersGameId = null;
    });

    // 確認モーダルを開く
    const openConfirmModal = (message, callback, confirmButtonText = '確認', cancelButtonText = 'キャンセル') => {
        confirmMessage.textContent = message;
        confirmActionBtn.textContent = confirmButtonText;
        cancelConfirmBtn.textContent = cancelButtonText;

        confirmCallback = callback; // コールバック関数を保存
        showModal(confirmModal, confirmModalContent);

        // OKボタンのみの場合、キャンセルボタンを非表示にする
        if (confirmButtonText === 'OK' && cancelButtonText === 'キャンセル') {
            cancelConfirmBtn.classList.add('hidden');
            confirmActionBtn.classList.remove('w-auto'); // 幅を自動調整
            confirmActionBtn.classList.add('w-full'); // 幅をフルに
            confirmActionBtn.classList.remove('bg-red-500', 'hover:bg-red-600'); // 赤色を削除
            confirmActionBtn.classList.add('bg-gray-700', 'hover:bg-gray-800'); // モノクロの色を追加
        } else {
            cancelConfirmBtn.classList.remove('hidden');
            confirmActionBtn.classList.remove('w-full');
            confirmActionBtn.classList.add('w-auto');
            confirmActionBtn.classList.add('bg-red-500', 'hover:bg-red-600'); // 赤色を維持
            confirmActionBtn.classList.remove('bg-gray-700', 'hover:bg-gray-800'); // モノクロの色を削除
        }
    };

    // 確認モーダルのキャンセルボタン
    cancelConfirmBtn.addEventListener('click', () => {
        hideModal(confirmModal, confirmModalContent);
        confirmCallback = null; // コールバックをクリア
    });

    // 確認モーダルの確認ボタン
    confirmActionBtn.addEventListener('click', () => {
        if (confirmCallback) {
            confirmCallback(); // 保存されたコールバックを実行
        }
        // OKボタンのみの場合は、クリック後にモーダルを閉じる
        if (confirmActionBtn.textContent === 'OK') {
            hideModal(confirmModal, confirmModalContent);
        }
    });

    // アプリケーション起動時にリージョンとゲームをロード
    populateRegionsDropdown(); // リージョンドロップダウンを初期化
    loadGames(); // API経由でゲームをロード
});
