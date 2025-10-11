// グローバル変数 (home.jsで定義) を使用
// let currentSort;
// let allProfiles;

/**
 * プロファイル一覧をソートし、テーブルに描画する。
 * @param {Array<Object>} profiles - プロファイルデータの配列
 */
function sortAndDisplayProfiles(profiles) {
    const { key, direction } = currentSort;

    profiles.sort((a, b) => {
        const valA = a[key];
        const valB = b[key];

        if (typeof valA === 'number' && typeof valB === 'number') {
            return direction === 'asc' ? valA - valB : valB - valA;
        }

        if (valA < valB) {
            return direction === 'asc' ? -1 : 1;
        }
        if (valA > valB) {
            return direction === 'asc' ? 1 : -1;
        }
        return 0;
    });

    displayProfiles(profiles);
    updateTableHeaderClasses(key, direction);
}

/**
 * プロファイル一覧をテーブルに描画する。
 * @param {Array<Object>} profiles - プロファイルデータの配列
 */
function displayProfiles(profiles) {
    const tableBody = document.getElementById('profiles-table-body');
    tableBody.innerHTML = '';

    if (!profiles || profiles.length === 0) {
        tableBody.innerHTML = '<tr><td colspan="7" class="no-data">表示するプロファイルがありません。</td></tr>';
        return;
    }

    profiles.forEach(profile => {
        const row = tableBody.insertRow();
        
        const columns = [
            profile.name,
            profile.user_id,
            profile.region_id,
            profile.coin,
            profile.latitude,
            profile.longitude
        ];

        columns.forEach(text => {
            const cell = row.insertCell();
            cell.textContent = text;
        });

        const actionCell = row.insertCell();
        actionCell.className = 'action-cell';

        const editButton = document.createElement('button');
        editButton.textContent = '編集';
        editButton.className = 'action-button edit-btn';
        editButton.onclick = () => handleEdit(profile.user_id);

        const deleteButton = document.createElement('button');
        deleteButton.textContent = '削除';
        deleteButton.className = 'action-button delete-btn';
        deleteButton.onclick = () => handleDelete(profile.user_id, profile.name);

        actionCell.appendChild(editButton);
        actionCell.appendChild(deleteButton);
    });
}

/**
 * テーブルヘッダーにソート機能のイベントリスナーを設定する。
 */
function setupSorting() {
    const headers = document.querySelectorAll('#profiles-table th.header-sortable');
    headers.forEach(header => {
        header.addEventListener('click', () => {
            const key = header.getAttribute('data-key');
            if (currentSort.key === key) {
                currentSort.direction = currentSort.direction === 'asc' ? 'desc' : 'asc';
            } else {
                currentSort.key = key;
                currentSort.direction = 'asc';
            }
            sortAndDisplayProfiles(allProfiles);
        });
    });
}

/**
 * ソート状態に合わせてヘッダーのスタイルを更新する。
 * @param {string} key - 現在ソート中のキー
 * @param {string} direction - 現在のソート方向 ('asc' or 'desc')
 */
function updateTableHeaderClasses(key, direction) {
    const headers = document.querySelectorAll('#profiles-table th.header-sortable');
    headers.forEach(header => {
        header.classList.remove('asc', 'desc');
        if (header.getAttribute('data-key') === key) {
            header.classList.add(direction);
        }
    });
}

/**
 * 地域一覧を取得し、ドロップダウンに設定する。
 */
async function fetchAndPopulateRegions() {
    const regionSelect = document.getElementById('region_id');
    try {
        const regions = await GetRegions();
        
        regionSelect.innerHTML = '';
        
        if (regions && regions.length > 0) {
            const defaultOption = document.createElement('option');
            defaultOption.value = '';
            defaultOption.textContent = '地域を選択してください';
            defaultOption.disabled = true;
            defaultOption.selected = true;
            regionSelect.appendChild(defaultOption);
            
            regions.forEach(region => {
                const option = document.createElement('option');
                // 💡 修正: option.valueにはRegionID、option.textContentにはregionNameを設定
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
