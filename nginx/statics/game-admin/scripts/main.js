const auth = new AuthBase('/auth/');

document.addEventListener('DOMContentLoaded', () => {
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
        regionSelect.innerHTML = '<option value="" disabled selected>地域を選択</option>';
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
            gamesTableBody.innerHTML = '<tr><td colspan="7" class="loading-message">データを読み込み中...</td></tr>';
            const data = await apiRequest('/game/list', { method: 'GET' });
            renderGames(data.Data || []);
        } catch (error) {
            errorMessage.textContent = `ゲームの読み込みに失敗しました: ${error.message}`;
            errorMessage.style.display = 'block';
        }
    }

    function renderGames(games) {
        gamesTableBody.innerHTML = '';
        if (games.length === 0) {
            gamesTableBody.innerHTML = '<tr><td colspan="7">ゲームが見つかりません。</td></tr>';
            return;
        }

        games.forEach(game => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${game.name}</td>
                <td>${game.game_id}</td>
                <td>${getRegionNameById(game.region_id)}</td>
                <td>${statusMap[game.status] || '不明'}</td>
                <td>${new Date(game.start_time * 1000).toLocaleString()}</td>
                <td>${game.dulation_date}</td>
                <td class="actions">
                    <button class="start-btn" data-id="${game.game_id}" ${game.status !== 0 ? 'disabled' : ''}>開始</button>
                    <button class="end-btn" data-id="${game.game_id}" ${game.status !== 1 ? 'disabled' : ''}>終了</button>
                    <button class="delete-btn" data-id="${game.game_id}">削除</button>
                </td>
            `;
            gamesTableBody.appendChild(row);
        });
    }

    createForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(createForm);
        const gameData = {
            name: formData.get('name'),
            region_id: formData.get('region_id'),
            start_time: new Date(formData.get('start_time')).getTime(),
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
