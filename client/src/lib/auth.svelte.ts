const TOKEN_KEY = 'token';

let token = $state<string | null>(localStorage.getItem(TOKEN_KEY));

export const auth = {
  get token() {
    return token;
  },
  get isAuthenticated() {
    return token !== null;
  },
  login(newToken: string) {
    token = newToken;
    localStorage.setItem(TOKEN_KEY, newToken);
  },
  logout() {
    token = null;
    localStorage.removeItem(TOKEN_KEY);
  },
};