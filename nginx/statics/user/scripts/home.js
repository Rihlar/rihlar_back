const auth = new AuthBase('/auth/');

let currentSort = { key: 'name', direction: 'asc' };
let allProfiles = [];

// DOM読み込み後にフォームイベントを設定
document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('create-profile-form');
    if (form) {
        form.addEventListener('submit', handleCreateProfileSubmit);
    }
});

async function Init() {    
    const tableBody = document.getElementById('profiles-table-body');
    const errorMessage = document.getElementById('error-message');
    errorMessage.style.display = 'none';

    try {
        const userData = await auth.GetInfo();
        if (userData == null) {
            window.location.href = './login.html';
            return;
        }
        console.log("ユーザーデータ:", userData);
        
        await fetchAndPopulateRegions();
        await fetchAndDisplayProfiles();
        
        const accessToken = await auth.getToken();
        console.log("アクセストークン:", accessToken);

    } catch (error) {
        console.error("初期化エラー:", error);
        errorMessage.textContent = '初期化中にエラーが発生しました。';
        errorMessage.style.display = 'block';
        tableBody.innerHTML = '<tr><td colspan="7" class="no-data">エラーによりデータを取得できませんでした。</td></tr>';
    }
}

async function fetchAndDisplayProfiles() {
    const profilesResult = await GetAllProfiles();
    if (profilesResult && profilesResult.profiles) {
        allProfiles = profilesResult.profiles;
        sortAndDisplayProfiles(allProfiles); 
    } else {
        document.getElementById('profiles-table-body').innerHTML = '<tr><td colspan="7" class="no-data">プロファイルデータが見つかりませんでした。</td></tr>';
    }
}

// プロファイル作成フォーム送信ハンドラ
async function handleCreateProfileSubmit(event) {
    event.preventDefault();
    
    const form = event.target;
    const formData = new FormData(form);
    
    const name = formData.get('name');
    const comment = formData.get('comment') || '';
    const regionId = formData.get('region_id');

    const userId = `userid-${generateUUID()}`;

    const statusMessage = document.getElementById('create-status-message');
    const createBtn = document.getElementById('create-btn');
    
    createBtn.disabled = true;
    statusMessage.textContent = '作成中...';
    statusMessage.style.color = 'black';

    try {
        const success = await CreateProfile(userId, name, comment, regionId);
        
        if (success) {
            statusMessage.textContent = '✅ プロファイルの作成に成功しました！';
            statusMessage.style.color = 'green';
            form.reset();
            await fetchAndDisplayProfiles();
            const regionSelect = document.getElementById('region_id');
            // ドロップダウンを「地域を選択してください」に戻す
            regionSelect.value = ''; 
        } else {
            statusMessage.textContent = '❌ プロファイルの作成に失敗しました。';
            statusMessage.style.color = 'red';
        }
    } catch (error) {
        console.error("プロファイル作成中にエラー:", error);
        statusMessage.textContent = `❌ エラーが発生しました: ${error.message}`;
        statusMessage.style.color = 'red';
    } finally {
        createBtn.disabled = false;
        setTimeout(() => { statusMessage.textContent = ''; }, 5000);
    }
}

// 削除ボタンがクリックされたときの処理
async function handleDelete(userId, name) {
    const confirmMessage = `ユーザー名: ${name} (ID: ${userId}) を削除しますか？`;
    
    if (!confirm(confirmMessage)) {
        return;
    }

    try {
        const success = await DeleteProfile(userId);
        if (success) {
            alert(`✅ ユーザーID: ${userId} の削除に成功しました。`);
            await fetchAndDisplayProfiles();
        } else {
            alert(`❌ ユーザーID: ${userId} の削除に失敗しました。`);
        }
    } catch (error) {
        console.error("削除処理中にエラー:", error);
        alert(`❌ 削除中に予期せぬエラーが発生しました: ${error.message}`);
    }
}

function handleEdit(userId) {
    alert(`ユーザーID: ${userId} の編集画面へ遷移します。（未実装）`);
}


// 初期化とソート機能のセットアップを実行
Init().then(() => {
    setupSorting();
});
