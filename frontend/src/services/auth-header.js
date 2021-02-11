
export default function authHeader() {
  let user = JSON.parse(localStorage.getItem('user'));
  if(!user) {
    user = JSON.parse(localStorage.getItem('user2fa'));
  }

  if (user && user.token) {
    return { Authorization: 'Bearer ' + user.token };
  } else {
    return {};
  }
}
