const auth = new AuthBase('/auth/');

document.addEventListener('DOMContentLoaded', () => {
    const gameIdHeading = document.getElementById('game-id-heading');
    const teamsContainer = document.getElementById('teams-container');
    const errorMessage = document.getElementById('error-message');

    const urlParams = new URLSearchParams(window.location.search);
    const gameId = urlParams.get('game_id');

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

    async function loadTeams() {
        if (!gameId) {
            errorMessage.textContent = 'Game IDが指定されていません。';
            errorMessage.style.display = 'block';
            return;
        }

        gameIdHeading.textContent = `ゲームID: ${gameId}`;

        try {
            const teams = await apiRequest(`/game/game/${gameId}/teams`, { method: 'GET' });
            renderTeams(teams || []);
        } catch (error) {
            errorMessage.textContent = `チームの読み込みに失敗しました: ${error.message}`;
            errorMessage.style.display = 'block';
        }
    }

    function renderTeams(teams) {
        const teamsContainer = document.getElementById('teams-container');
        teamsContainer.innerHTML = '';
        if (teams.length === 0) {
            teamsContainer.innerHTML = '<p class="text-gray-500 text-center">このゲームにはチームがありません。</p>';
            return;
        }

        teams.forEach(team => {
            const teamCard = document.createElement('div');
            teamCard.className = 'bg-white rounded-lg shadow-md overflow-hidden';

            let totalPoints = 0;
            const memberRows = team.members ? team.members.map(member => {
                totalPoints += member.points;
                return `
                    <tr class="hover:bg-gray-50">
                        <td class="p-3 text-sm text-gray-700">${member.userID}</td>
                        <td class="p-3 text-sm text-gray-700">${member.points} pts</td>
                        <td class="p-3 text-right">
                            <a href="/statics/admin/walking_history/?userId=${member.userID}&gameId=${gameId}" class="text-blue-500 hover:text-blue-700 text-sm font-semibold mr-4">行動履歴</a>
                            <button class="delete-member-btn text-red-500 hover:text-red-700 text-sm font-semibold" data-game-id="${gameId}" data-user-id="${member.userID}">削除</button>
                        </td>
                    </tr>
                `;
            }).join('') : '<tr><td colspan="3" class="p-3 text-center text-gray-500">メンバーがいません。</td></tr>';

            teamCard.innerHTML = `
                <div class="p-4 bg-gray-50 border-b border-gray-200 flex justify-between items-center">
                    <div>
                        <h3 class="text-lg font-bold text-gray-800">${team.teamID}</h3>
                        <p class="text-sm text-gray-600">合計ポイント: <span class="font-semibold">${totalPoints}</span></p>
                    </div>
                    <button class="delete-team-btn bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-lg text-sm" data-game-id="${gameId}" data-team-id="${team.teamID}">チーム削除</button>
                </div>
                <div class="p-0">
                    <table class="min-w-full">
                        <thead class="bg-gray-50">
                            <tr>
                                <th class="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">メンバーID</th>
                                <th class="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ポイント</th>
                                <th class="p-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            ${memberRows}
                        </tbody>
                    </table>
                </div>
            `;
            teamsContainer.appendChild(teamCard);
        });
    }

    document.getElementById('teams-container').addEventListener('click', async (e) => {
        const target = e.target;
        if (target.closest('.delete-team-btn')) {
            const button = target.closest('.delete-team-btn');
            const teamId = button.dataset.teamId;
            const gameId = button.dataset.gameId;
            if (confirm(`本当にチーム ${teamId} を削除しますか？`)) {
                try {
                    await apiRequest('/game/team/delete', { 
                        method: 'DELETE',
                        headers: { 'GameID': gameId, 'TeamID': teamId }
                    });
                    loadTeams();
                } catch (error) {
                    errorMessage.textContent = `チームの削除に失敗しました: ${error.message}`;
                    errorMessage.classList.remove('hidden');
                }
            }
        } else if (target.closest('.delete-member-btn')) {
            const button = target.closest('.delete-member-btn');
            const userId = button.dataset.userId;
            const gameId = button.dataset.gameId;
            if (confirm(`本当にメンバー ${userId} を削除しますか？`)) {
                try {
                    await apiRequest('/game/member/delete', { 
                        method: 'DELETE',
                        headers: { 'GameID': gameId, 'UserID': userId }
                    });
                    loadTeams();
                } catch (error) {
                    errorMessage.textContent = `メンバーの削除に失敗しました: ${error.message}`;
                    errorMessage.classList.remove('hidden');
                }
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
            loadTeams();
        } catch (error) {
            window.location.href = '../login.html';
        }
    }

    init();
});
