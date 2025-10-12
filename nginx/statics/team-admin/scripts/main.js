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
        teamsContainer.textContent = '';
        if (teams.length === 0) {
            const p = document.createElement('p');
            p.className = 'text-gray-500 text-center';
            p.textContent = 'このゲームにはチームがありません。';
            teamsContainer.appendChild(p);
            return;
        }

        teams.forEach(team => {
            const teamCard = document.createElement('div');
            teamCard.className = 'bg-white rounded-lg shadow-md overflow-hidden';

            const teamCard = document.createElement('div');
            teamCard.className = 'bg-white rounded-lg shadow-md overflow-hidden';

            let totalPoints = 0;

            const div1 = document.createElement('div');
            div1.className = 'p-4 bg-gray-50 border-b border-gray-200 flex justify-between items-center';
            
            const div2 = document.createElement('div');
            const h3 = document.createElement('h3');
            h3.className = 'text-lg font-bold text-gray-800';
            h3.textContent = team.teamID;
            div2.appendChild(h3);
            const p = document.createElement('p');
            p.className = 'text-sm text-gray-600';
            p.textContent = '合計ポイント: ';
            const span = document.createElement('span');
            span.className = 'font-semibold';
            p.appendChild(span);
            div2.appendChild(p);
            div1.appendChild(div2);

            const deleteButton = document.createElement('button');
            deleteButton.className = 'delete-team-btn bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-lg text-sm';
            deleteButton.dataset.gameId = gameId;
            deleteButton.dataset.teamId = team.teamID;
            deleteButton.textContent = 'チーム削除';
            div1.appendChild(deleteButton);
            teamCard.appendChild(div1);

            const div3 = document.createElement('div');
            div3.className = 'p-0';
            const table = document.createElement('table');
            table.className = 'min-w-full';
            const thead = document.createElement('thead');
            thead.className = 'bg-gray-50';
            const tr1 = document.createElement('tr');
            const th1 = document.createElement('th');
            th1.className = 'p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider';
            th1.textContent = 'メンバーID';
            tr1.appendChild(th1);
            const th2 = document.createElement('th');
            th2.className = 'p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider';
            th2.textContent = 'ポイント';
            tr1.appendChild(th2);
            const th3 = document.createElement('th');
            th3.className = 'p-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider';
            th3.textContent = '操作';
            tr1.appendChild(th3);
            thead.appendChild(tr1);
            table.appendChild(thead);

            const tbody = document.createElement('tbody');
            tbody.className = 'bg-white divide-y divide-gray-200';

            if (team.members && team.members.length > 0) {
                team.members.forEach(member => {
                    totalPoints += member.points;
                    const tr = document.createElement('tr');
                    tr.className = 'hover:bg-gray-50';

                    const td1 = document.createElement('td');
                    td1.className = 'p-3 text-sm text-gray-700';
                    td1.textContent = member.userID;
                    tr.appendChild(td1);

                    const td2 = document.createElement('td');
                    td2.className = 'p-3 text-sm text-gray-700';
                    td2.textContent = `${member.points} pts`;
                    tr.appendChild(td2);

                    const td3 = document.createElement('td');
                    td3.className = 'p-3 text-right';
                    const a = document.createElement('a');
                    a.href = `/statics/admin/walking_history/?userId=${member.userID}&gameId=${gameId}`;
                    a.className = 'text-blue-500 hover:text-blue-700 text-sm font-semibold mr-4';
                    a.textContent = '行動履歴';
                    td3.appendChild(a);
                    const button = document.createElement('button');
                    button.className = 'delete-member-btn text-red-500 hover:text-red-700 text-sm font-semibold';
                    button.dataset.gameId = gameId;
                    button.dataset.userId = member.userID;
                    button.textContent = '削除';
                    td3.appendChild(button);
                    tr.appendChild(td3);

                    tbody.appendChild(tr);
                });
            } else {
                const tr = document.createElement('tr');
                const td = document.createElement('td');
                td.colSpan = 3;
                td.className = 'p-3 text-center text-gray-500';
                td.textContent = 'メンバーがいません。';
                tr.appendChild(td);
                tbody.appendChild(tr);
            }

            span.textContent = totalPoints;
            table.appendChild(tbody);
            div3.appendChild(table);
            teamCard.appendChild(div3);
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
