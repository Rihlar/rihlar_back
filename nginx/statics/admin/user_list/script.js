document.addEventListener('DOMContentLoaded', () => {
    const userListElement = document.getElementById('user-list');

    // 仮のユーザーデータ (実際にはAPIから取得)
    const users = [
        { id: 'user123', name: 'Alice' },
        { id: 'user456', name: 'Bob' },
        { id: 'user789', name: 'Charlie' },
    ];

    users.forEach(user => {
        const listItem = document.createElement('li');
        listItem.innerHTML = `
            <span>${user.name} (ID: ${user.id})</span>
            <a href="/statics/admin/walking_history/index.html?user_id=${user.id}">歩行履歴を見る</a>
        `;
        userListElement.appendChild(listItem);
    });
});