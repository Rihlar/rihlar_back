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
        { RegionID: "c161edb9-6aff-4244-8749-707bff2fa3be", RegionName: "関西" },
        { RegionID: "regionId-005344a4-523d-4aab-9d1e-d232322cf54e", RegionName: "中国" },
        { RegionID: "regionId-3501902d-8bab-40cd-926f-30c53d80efc5", RegionName: "四国" },
        { RegionID: "regionId-ef5aa179-53e0-481d-b64d-ae7654049a88", RegionName: "九州" },
    ];

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
            fetchGamesAPI: () => {
                console.log("API: Fetching all games...");
                return [...backendData.games]; // データのコピーを返す
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
                fetch('/game/create', { method: 'POST', headers: { 'Content-Type': 'application/json', 'Authorization': `${token}` }, body: JSON.stringify(apiPayload) })
                    .then(response => response.json())
                    .then(data => {
                        console.log(data);
                    });

                // シミュレーションのため、直接データを追加
                const gameWithId = { ...newGameData, id: Date.now().toString() };
                backendData.games.push(gameWithId);
                saveBackendData();
                return gameWithId; // UIに返すデータは、追加されたゲームオブジェクト
            },

            /**
             * ゲームを削除するAPIシミュレーション
             * @param {string} gameId - 削除するゲームのID
             */
            deleteGameAPI: (gameId) => {
                console.log("API: Deleting game...", gameId);
                const initialLength = backendData.games.length;
                backendData.games = backendData.games.filter(game => game.id !== gameId);
                saveBackendData();
                if (backendData.games.length === initialLength) {
                    throw new Error("Game not found for deletion.");
                }
            },

            /**
             * ゲームのステータスを更新するAPIシミュレーション
             * @param {string} gameId - 更新するゲームのID
             * @param {string} newStatus - 新しいステータス ('予定', '開催中', '終了')
             * @returns {Object} 更新されたゲームのデータ
             */
            updateGameStatusAPI: (gameId, newStatus) => {
                console.log("API: Updating game status...", gameId, newStatus);
                const gameIndex = backendData.games.findIndex(game => game.id === gameId);
                if (gameIndex !== -1) {
                    backendData.games[gameIndex].status = newStatus;
                    saveBackendData();
                    return backendData.games[gameIndex];
                } else {
                    throw new Error("Game not found for status update.");
                }
            },

            /**
             * 特定のゲームのチーム一覧を取得するAPIシミュレーション
             * @param {string} gameId - チーム一覧を取得するゲームのID
             * @returns {Array} チームデータの配列
             */
            fetchTeamsAPI: (gameId) => {
                console.log("API: Fetching teams for game...", gameId);
                const game = backendData.games.find(g => g.id === gameId);
                if (game) {
                    return [...(game.teams || [])]; // データのコピーを返す
                } else {
                    throw new Error("Game not found for fetching teams.");
                }
            },

            /**
             * 特定のゲームにメンバーを追加するAPIシミュレーション
             * @param {string} gameId - メンバーを追加するゲームのID
             * @param {string} memberName - 追加するメンバーの名前
             * @returns {Array} 更新されたメンバーリスト
             */
            addMemberAPI: (gameId, memberName) => {
                console.log("API: Adding member to game...", gameId, memberName);
                const gameIndex = backendData.games.findIndex(game => game.id === gameId);
                if (gameIndex !== -1) {
                    if (!backendData.games[gameIndex].members.includes(memberName)) {
                        backendData.games[gameIndex].members.push(memberName);
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
                const gameIndex = backendData.games.findIndex(game => game.id === gameId);
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
                const gameIndex = backendData.games.findIndex(game => game.id === gameId);
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
    const loadGames = () => {
        try {
            games = backendService.fetchGamesAPI(); // Promiseではなく直接データを取得
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
            gameCard.dataset.id = game.id; // データIDを設定して詳細表示に利用

            gameCard.innerHTML = `
                <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4">
                    <h3 class="text-2xl font-bold text-gray-900 mb-2 sm:mb-0">${game.name}</h3>
                    <span class="px-3 py-1 rounded-full text-sm font-semibold ${game.status === '予定' ? 'status-予定' :
                    game.status === '開催中' ? 'status-開催中' :
                        'status-終了'
                }">${game.status}</span>
                </div>
                <p class="text-gray-700 mb-2"><strong>リージョン:</strong> ${getRegionNameById(game.region)}</p>
                <p class="text-gray-700 mb-2"><strong>開催時間:</strong> ${new Date(game.time).toLocaleString('ja-JP', {
                    year: 'numeric', month: 'numeric', day: 'numeric',
                    hour: '2-digit', minute: '2-digit', hour12: false
                })}</p>
                <p class="text-gray-700 mb-2"><strong>開催期間:</strong> ${game.duration} 日間</p>
                <p class="text-gray-700 mb-4"><strong>ゲームタイプ:</strong> ${game.gameType === 'personal' ? '個人のゲーム' : 'チームのゲーム'}</p>
                <div class="flex flex-wrap gap-3 mt-4">
                    <button data-id="${game.id}" class="delete-btn bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-200">削除</button>
                    <button data-id="${game.id}" class="members-btn bg-gray-600 hover:bg-gray-700 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-200">メンバー (${game.members ? game.members.length : 0})</button>
                    ${game.status === '予定' || game.status === '終了' ? `<button data-id="${game.id}" class="start-btn bg-gray-700 hover:bg-gray-800 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-200">開始</button>` : ''}
                    ${game.status === '開催中' ? `<button data-id="${game.id}" class="end-btn bg-gray-600 hover:bg-gray-700 text-white font-semibold py-2 px-4 rounded-lg shadow-sm transition duration-300">終了</button>` : ''}
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
            button.onclick = (e) => {
                openConfirmModal('本当にこのゲームを削除しますか？', () => {
                    try {
                        backendService.deleteGameAPI(e.target.dataset.id);
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
            button.onclick = (e) => {
                try {
                    backendService.updateGameStatusAPI(e.target.dataset.id, '開催中');
                    loadGames(); // ステータス更新後にリストを再ロード
                } catch (error) {
                    console.error("Failed to start game:", error);
                    openConfirmModal('ゲームの開始に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                }
            };
        });
        document.querySelectorAll('.end-btn').forEach(button => {
            button.onclick = (e) => {
                try {
                    backendService.updateGameStatusAPI(e.target.dataset.id, '終了');
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
    gameForm.addEventListener('submit', (e) => {
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
            backendService.addGameAPI(newGameData); // API経由で新しいゲームを追加
            loadGames(); // 追加後にリストを再ロード
            hideModal(gameModal, gameModalContent);
        } catch (error) {
            console.error("Failed to add game:", error);
            openConfirmModal('ゲームの作成に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
        }
    });

    // ゲーム詳細モーダルを開く
    const openGameDetailsModal = (id) => {
        const game = games.find(g => g.id === id);
        if (game) {
            currentSelectedGameId = id;
            detailsModalGameName.textContent = game.name;
            detailsModalRegion.textContent = getRegionNameById(game.region);
            detailsModalTime.textContent = new Date(game.time).toLocaleString('ja-JP', {
                year: 'numeric', month: 'numeric', day: 'numeric',
                hour: '2-digit', minute: '2-digit', hour12: false
            });
            detailsModalDuration.textContent = game.duration;
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
                    <span class="text-gray-800 font-medium">チームID: ${team.teamId}</span>
                    <span class="text-gray-600">ポイント: ${team.points}</span>
                    <button data-index="${index}" class="remove-team-btn text-red-500 hover:text-red-700 font-semibold text-sm">削除</button>
                `;
                teamListDiv.appendChild(teamItem);
            });
            // チーム削除ボタンにイベントリスナーを追加
            document.querySelectorAll('.remove-team-btn').forEach(button => {
                button.onclick = (e) => {
                    const teamIndex = parseInt(e.target.dataset.index);
                    try {
                        backendService.removeTeamAPI(currentSelectedGameId, teamIndex);
                        // 削除後、詳細モーダルを閉じずにチームリストのみ更新
                        const updatedTeams = backendService.fetchTeamsAPI(currentSelectedGameId);
                        renderTeams(updatedTeams);
                    } catch (error) {
                        console.error("Failed to remove team:", error);
                        openConfirmModal('チームの削除に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                    }
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
        // メンバーリストはゲームオブジェクトに直接含まれているため、API呼び出しは不要
        const game = games.find(g => g.id === currentMembersGameId);
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
                    <span class="text-gray-800">${member}</span>
                    <button data-index="${index}" class="remove-member-btn text-red-500 hover:text-red-700 font-semibold text-sm">削除</button>
                `;
                memberListDiv.appendChild(memberItem);
            });
            // メンバー削除ボタンにイベントリスナーを追加
            document.querySelectorAll('.remove-member-btn').forEach(button => {
                button.onclick = (e) => {
                    const memberIndex = parseInt(e.target.dataset.index);
                    try {
                        backendService.removeMemberAPI(currentMembersGameId, memberIndex);
                        // 削除後、メンバーリストとゲームリストを再ロード
                        loadGames();
                        const game = games.find(g => g.id === currentMembersGameId);
                        renderMembers(game ? game.members : []);
                    } catch (error) {
                        console.error("Failed to remove member:", error);
                        openConfirmModal('メンバーの削除に失敗しました。', () => hideModal(confirmModal, confirmModalContent), 'OK');
                    }
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