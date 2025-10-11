/**
 * 全てのプロファイルデータをAPIから取得する。
 * AuthorizationヘッダーにBearerなしのトークンを設定する。
 * @returns {Promise<Object|null>} プロファイルデータオブジェクト、またはnull
 */
async function GetAllProfiles() {
    const accessToken = await auth.getToken();

    const req = await fetch("/user/admin/profiles", {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            // Bearer不要
            ...(accessToken && { "Authorization": accessToken }) 
        }
    });

    if (req.ok) {
        return await req.json();
    }
    
    console.error("プロファイル取得失敗:", req.status, req.statusText);
    return null;
}

/**
 * 新しいプロファイルをAPIで作成する。
 * AuthorizationヘッダーにBearerなしのトークンを設定する。
 * @param {string} userId - 作成するユーザーのID
 * @param {string} name - ユーザー名
 * @param {string} comment - コメント
 * @param {string} regionId - 地域ID
 * @returns {Promise<boolean>} 成功/失敗
 */
async function CreateProfile(userId, name, comment, regionId) {
    const accessToken = await auth.getToken();

    const requestBody = {
        user_id: userId,
        name: name,
        comment: comment,
        region_id: regionId,
        system_game_id: "", 
        admin_game_id: ""  
    };

    const req = await fetch("/user/admin/profile", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            "Authorization": accessToken 
        },
        body: JSON.stringify(requestBody)
    });

    if (req.ok) {
        return true;
    } else {
        const errorText = await req.text();
        console.error("プロファイル作成APIエラー:", req.status, req.statusText, errorText);
        return false;
    }
}

/**
 * 指定されたユーザーIDのプロファイルをAPIで削除する。
 * UserIDヘッダーに削除対象IDを設定し、AuthorizationヘッダーにBearerなしのトークンを設定する。
 * @param {string} userId - 削除対象のユーザーID
 * @returns {Promise<boolean>} 成功/失敗
 */
async function DeleteProfile(userId) {
    const accessToken = await auth.getToken();

    if (!accessToken) {
        console.error("削除操作：アクセストークンが取得できませんでした。");
        return false;
    }

    const req = await fetch("/user/admin/profile", {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'UserID': userId,
            'Authorization': accessToken
        }
    });

    if (req.ok) {
        return true;
    } else {
        const errorText = await req.text();
        console.error("プロファイル削除APIエラー:", req.status, req.statusText, errorText);
        return false;
    }
}

/**
 * 地域一覧データをAPIから取得する。
 * @returns {Promise<Array<Object>|null>} 地域データの配列、またはnull
 */
async function GetRegions() {
    const accessToken = await auth.getToken();

    const req = await fetch("/user/admin/regions", {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            // Bearer不要
            ...(accessToken && { "Authorization": accessToken })
        }
    });

    if (req.ok) {
        return await req.json();
    }
    
    console.error("地域データ取得失敗:", req.status, req.statusText);
    return null;
}
