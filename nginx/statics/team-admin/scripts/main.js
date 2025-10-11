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
        teamsContainer.innerHTML = '';
        if (teams.length === 0) {
            teamsContainer.innerHTML = '<p>このゲームにはチームがありません。</p>';
            return;
        }

        teams.forEach(team => {
            const teamElement = document.createElement('div');
            teamElement.classList.add('team-card');

            let totalPoints = 0;
            const membersList = team.members ? team.members.map(member => {
                totalPoints += member.points;
                return `<li>${member.userID} - ${member.points} pts</li>`;
            }).join('') : '<li>メンバーがいません。</li>';

            teamElement.innerHTML = `
                <h3>チームID: ${team.teamID}</h3>
                <p>合計ポイント: ${totalPoints}</p>
                <h4>メンバー:</h4>
                <ul>
                    ${membersList}
                </ul>
            `;
            teamsContainer.appendChild(teamElement);
        });
    }

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
