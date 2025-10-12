const auth = new AuthBase('/auth/');

document.addEventListener('DOMContentLoaded', () => {
    const toggleButton = document.getElementById('toggle-create-form-btn');
    const createFormSection = document.getElementById('create-game-section');

    toggleButton.addEventListener('click', () => {
        createFormSection.classList.toggle('hidden');
    });
    const gamesTableBody = document.getElementById('games-table-body');
    const createForm = document.getElementById('create-game-form');
    const regionSelect = document.getElementById('region_id');
    const errorMessage = document.getElementById('error-message');

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

    const statusMap = {
        0: '予定',
        1: '開催中',
        2: '終了'
    };

    async function apiRequest(url, options) {
        const token = await auth.getToken();
        const headers = {
            'Content-Type': 'application/json',
            'Authorization': `${token}`,
            ...options.headers,
        };
        const response = await fetch(url, { ...options, headers });
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ message: 'An unknown error occurred' }));
            throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
        }
        return response.json();
    }

    function populateRegions() {
        regionSelect.textContent = '';
        const option = document.createElement('option');
        option.value = '';
        option.disabled = true;
        option.selected = true;
        option.textContent = '地域を選択';
        regionSelect.appendChild(option);
        regions.forEach(region => {
            const option = document.createElement('option');
            option.value = region.RegionID;
            option.textContent = region.RegionName;
            regionSelect.appendChild(option);
        });
    }

    function getRegionNameById(id) {
        const region = regions.find(r => r.RegionID === id);
        return region ? region.RegionName : 'N/A';
    }

    async function loadGames() {
        try {
            gamesTableBody.textContent = '';
            const tr = document.createElement('tr');
            const td = document.createElement('td');
            td.colSpan = 7;
            td.className = 'loading-message';
            td.textContent = 'データを読み込み中...';
            tr.appendChild(td);
            gamesTableBody.appendChild(tr);
            const data = await apiRequest('/game/list', { method: 'GET' });
            renderGames(data.Data || []);
        } catch (error) {
            errorMessage.textContent = `ゲームの読み込みに失敗しました: ${error.message}`;
            errorMessage.style.display = 'block';
        }
    }

    function renderGames(games) {
        gamesTableBody.textContent = '';
        if (games.length === 0) {
            gamesTableBody.textContent = '';
            const tr = document.createElement('tr');
            const td = document.createElement('td');
            td.colSpan = 7;
            td.textContent = 'ゲームが見つかりません。';
            tr.appendChild(td);
            gamesTableBody.appendChild(tr);
            return;
        }

        games.forEach(game => {
            const row = document.createElement('tr');
            const td1 = document.createElement('td');
            td1.textContent = game.name;
            row.appendChild(td1);

            const td2 = document.createElement('td');
            td2.textContent = game.game_id;
            row.appendChild(td2);

            const td3 = document.createElement('td');
            td3.textContent = getRegionNameById(game.region_id);
            row.appendChild(td3);

            const td4 = document.createElement('td');
            td4.textContent = statusMap[game.status] || '不明';
            row.appendChild(td4);

            const td5 = document.createElement('td');
            td5.textContent = new Date(game.start_time * 1000).toLocaleString();
            row.appendChild(td5);

            const td6 = document.createElement('td');
            td6.textContent = game.dulation_date;
            row.appendChild(td6);

            const td7 = document.createElement('td');
            td7.className = 'actions';
            const startButton = document.createElement('button');
            startButton.className = 'start-btn';
            startButton.dataset.id = game.game_id;
            if (game.status === 1) {
                startButton.disabled = true;
            }
            startButton.textContent = '開始';
            td7.appendChild(startButton);

            const endButton = document.createElement('button');
            endButton.className = 'end-btn';
            endButton.dataset.id = game.game_id;
            if (game.status !== 1) {
                endButton.disabled = true;
            }
            endButton.textContent = '終了';
            td7.appendChild(endButton);

            const deleteButton = document.createElement('button');
            deleteButton.className = 'delete-btn';
            deleteButton.dataset.id = game.game_id;
            deleteButton.textContent = '削除';
            td7.appendChild(deleteButton);

            const a = document.createElement('a');
            a.href = `../team-admin/?game_id=${game.game_id}`;
            a.className = 'manage-teams-btn';
            a.textContent = 'チーム管理';
            td7.appendChild(a);
            row.appendChild(td7);
            gamesTableBody.appendChild(row);
        });
    }

    createForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(createForm);
        const gameData = {
            name: formData.get('name'),
            region_id: formData.get('region_id'),
            start_time: new Date(formData.get('start_time')).getTime() / 1000,
            dulation_date: parseInt(formData.get('dulation_date'), 10),
        };

        try {
            await apiRequest('/game/create', {
                method: 'POST',
                body: JSON.stringify(gameData),
            });
            createForm.reset();
            loadGames();
        } catch (error) {
            errorMessage.textContent = `ゲームの作成に失敗しました: ${error.message}`;
            errorMessage.style.display = 'block';
        }
    });

    gamesTableBody.addEventListener('click', async (e) => {
        const target = e.target;
        const id = target.dataset.id;

        if (target.classList.contains('delete-btn')) {
            if (confirm('本当にこのゲームを削除しますか？')) {
                try {
                    await apiRequest('/game/delete', { 
                        method: 'DELETE',
                        headers: { 'GameID': id }
                    });
                    loadGames();
                } catch (error) {
                    errorMessage.textContent = `ゲームの削除に失敗しました: ${error.message}`;
                    errorMessage.style.display = 'block';
                }
            }
        } else if (target.classList.contains('start-btn')) {
            try {
                await apiRequest('/game/start', { 
                    method: 'PATCH',
                    body: JSON.stringify({ game_id: id })
                });
                loadGames();
            } catch (error) {
                errorMessage.textContent = `ゲームの開始に失敗しました: ${error.message}`;
                errorMessage.style.display = 'block';
            }
        } else if (target.classList.contains('end-btn')) {
            try {
                await apiRequest('/game/end', { 
                    method: 'PATCH',
                    body: JSON.stringify({ game_id: id })
                });
                loadGames();
            } catch (error) {
                errorMessage.textContent = `ゲームの終了に失敗しました: ${error.message}`;
                errorMessage.style.display = 'block';
            }
        }
    });

    async function init() {
        try {
            const userData = await auth.GetInfo();
            if (!userData) {
                window.location.href = '../login.html';
                return;
            }
            populateRegions();
            loadGames();
        } catch (error) {
            window.location.href = '../login.html';
        }
    }

    init();
});
