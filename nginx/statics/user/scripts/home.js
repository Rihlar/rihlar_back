const auth = new AuthBase('/auth/');

let currentSort = { key: 'name', direction: 'asc' };
let allProfiles = [];

// DOM読み込み後にフォームイベントを設定
document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('create-profile-form');
    if (form) {
        form.addEventListener('submit', handleCreateProfileSubmit);
    }

    const toggleButton = document.getElementById('toggle-create-form-btn');
    const createFormSection = document.getElementById('create-profile-section');

    if (toggleButton && createFormSection) {
        toggleButton.addEventListener('click', () => {
            createFormSection.classList.toggle('hidden');
        });
    }

    const editModal = document.getElementById('edit-profile-modal');
    const closeButton = editModal.querySelector('.close-button');
    closeButton.addEventListener('click', () => {
        editModal.style.display = 'none';
    });

    const editForm = document.getElementById('edit-profile-form');
    editForm.addEventListener('submit', handleUpdateProfileSubmit);
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

async function handleEdit(userId) {
    const modal = document.getElementById('edit-profile-modal');
    modal.style.display = 'block';

    const profile = await GetProfile(userId);
    if (profile) {
        document.getElementById('edit-user-id').value = profile.user_id;
        document.getElementById('edit-name').value = profile.name;
        document.getElementById('edit-comment').value = profile.comment;
        
        const regionSelect = document.getElementById('edit-region_id');
        await fetchAndPopulateRegionsForEdit(regionSelect);
        regionSelect.value = profile.region_id;
    }
}

async function handleUpdateProfileSubmit(event) {
    event.preventDefault();

    const form = event.target;
    const formData = new FormData(form);
    const userId = formData.get('user_id');
    const data = {
        name: formData.get('name'),
        comment: formData.get('comment'),
        region_id: formData.get('region_id'),
    };

    const statusMessage = document.getElementById('update-status-message');
    const updateBtn = document.getElementById('update-btn');

    updateBtn.disabled = true;
    statusMessage.textContent = '更新中...';
    statusMessage.style.color = 'black';

    try {
        const success = await UpdateProfile(userId, data);
        if (success) {
            statusMessage.textContent = '✅ プロファイルの更新に成功しました！';
            statusMessage.style.color = 'green';
            await fetchAndDisplayProfiles();
            setTimeout(() => {
                document.getElementById('edit-profile-modal').style.display = 'none';
                statusMessage.textContent = '';
            }, 2000);
        } else {
            statusMessage.textContent = '❌ プロファイルの更新に失敗しました。';
            statusMessage.style.color = 'red';
        }
    } catch (error) {
        console.error("プロファイル更新中にエラー:", error);
        statusMessage.textContent = `❌ エラーが発生しました: ${error.message}`;
        statusMessage.style.color = 'red';
    } finally {
        updateBtn.disabled = false;
    }
}

async function fetchAndPopulateRegionsForEdit(regionSelect) {
    try {
        const regions = await GetRegions();
        
        regionSelect.innerHTML = '';
        
        if (regions && regions.length > 0) {
            regions.forEach(region => {
                const option = document.createElement('option');
                option.value = region.RegionID;
                option.textContent = region.regionName;
                regionSelect.appendChild(option);
            });
        } else {
            const noDataOption = document.createElement('option');
            noDataOption.value = '';
            noDataOption.textContent = '地域データがありません';
            noDataOption.disabled = true;
            noDataOption.selected = true;
            regionSelect.appendChild(noDataOption);
            console.warn("地域データが見つかりませんでした。");
        }
    } catch (error) {
        console.error("地域データの取得に失敗しました:", error);
        regionSelect.innerHTML = '<option value="" disabled selected>地域の読み込みに失敗しました</option>';
    }
}


// 初期化とソート機能のセットアップを実行
Init().then(() => {
    setupSorting();
});