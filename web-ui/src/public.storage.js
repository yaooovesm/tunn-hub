let User = {
    isLogin: false,
    token: "",
    info: undefined,
    password: "",
}

/**
 *
 * @constructor
 */
function Load() {
    let lo = localStorage.getItem("tunnel_server_user")
    if (lo !== "" && lo !== undefined && lo !== null) {
        let loUsr = JSON.parse(lo)
        User.isLogin = true
        User.info = loUsr.info
        User.token = loUsr.token
        User.password = loUsr.password
    }
}

/**
 *
 * @returns {boolean}
 * @constructor
 */
function IsAdmin() {
    return User.info.auth === "admin"
}

export default {
    User,
    IsAdmin,
    Load
}
