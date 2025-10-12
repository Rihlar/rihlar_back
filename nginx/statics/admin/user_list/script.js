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
        const span = document.createElement('span');
        span.textContent = `${user.name} (ID: ${user.id})`;
        listItem.appendChild(span);

        const a = document.createElement('a');
        a.href = `/statics/admin/walking_history/index.html?user_id=${user.id}`;
        a.textContent = '歩行履歴を見る';
        listItem.appendChild(a);
        userListElement.appendChild(listItem);
    });
});