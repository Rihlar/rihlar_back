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
        const tableBody = document.getElementById('teams-table-body');
        tableBody.innerHTML = '';
        if (teams.length === 0) {
            tableBody.innerHTML = '<tr><td colspan="4">このゲームにはチームがありません。</td></tr>';
            return;
        }

        teams.forEach(team => {
            const row = document.createElement('tr');

            let totalPoints = 0;
            const membersHtml = team.members ? team.members.map(member => {
                totalPoints += member.points;
                return `<li>${member.userID} - ${member.points} pts <button class="delete-member-btn" data-game-id="${gameId}" data-user-id="${member.userID}">削除</button></li>`;
            }).join('') : '<li>メンバーがいません。</li>';

            row.innerHTML = `
                <td>${team.teamID}</td>
                <td>${totalPoints}</td>
                <td><ul>${membersHtml}</ul></td>
                <td class="actions">
                    <button class="delete-team-btn" data-game-id="${gameId}" data-team-id="${team.teamID}">チーム削除</button>
                </td>
            `;
            tableBody.appendChild(row);
        });
    }

    document.getElementById('teams-table-body').addEventListener('click', async (e) => {
        const target = e.target;
        if (target.classList.contains('delete-team-btn')) {
            const teamId = target.dataset.teamId;
            const gameId = target.dataset.gameId;
            if (confirm(`本当にチーム ${teamId} を削除しますか？`)) {
                try {
                    await apiRequest('/game/team/delete', { 
                        method: 'DELETE',
                        headers: { 'GameID': gameId, 'TeamID': teamId }
                    });
                    loadTeams();
                } catch (error) {
                    errorMessage.textContent = `チームの削除に失敗しました: ${error.message}`;
                    errorMessage.style.display = 'block';
                }
            }
        } else if (target.classList.contains('delete-member-btn')) {
            const userId = target.dataset.userId;
            const gameId = target.dataset.gameId;
            if (confirm(`本当にメンバー ${userId} を削除しますか？`)) {
                try {
                    await apiRequest('/game/member/delete', { 
                        method: 'DELETE',
                        headers: { 'GameID': gameId, 'UserID': userId }
                    });
                    loadTeams();
                } catch (error) {
                    errorMessage.textContent = `メンバーの削除に失敗しました: ${error.message}`;
                    errorMessage.style.display = 'block';
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
