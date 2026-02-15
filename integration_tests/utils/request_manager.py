from dataclasses import dataclass
import requests


@dataclass
class User:
    login: str
    password: str
    token: str
    headers: dict[str, str]


class RequestManager:
    def __init__(self, url):
        self.url = url
        self._user_id = 0

    def post(self, page, json=None, data=None, user=None):
        headers = user.headers if user else None
        return requests.post(f'{self.url}/{page}', json=json, data=data, headers=headers)

    def get(self, page, user=None):
        headers = user.headers if user else None
        return requests.get(f'{self.url}/{page}', headers=headers)

    def put(self, page, json=None, data=None, user=None):
        headers = user.headers if user else None
        return requests.put(f'{self.url}/{page}', json=json, data=data, headers=headers)

    def delete(self, page, json=None, data=None, user=None):
        headers = user.headers if user else None
        return requests.delete(f'{self.url}/{page}', json=json, data=data, headers=headers)

    def register(self, login, password):
        r = self.post('register', json={'login': login, 'password': password})
        if r.status_code != 200:
            raise Exception('wrong status code')

        token = r.json()['token']
        headers = {'Authorization': f'Bearer {token}'}
        return User(login=login, password=password, token=token, headers=headers)

    def login(self, login, password):
        r = self.post('login', json={'login': login, 'password': password})
        if r.status_code != 200:
            raise Exception('wrong status code')

        token = r.json()['token']
        headers = {'Authorization': f'Bearer {token}'}
        return User(login=login, password=password, token=token, headers=headers)

    def get_new_user(self):
        self._user_id += 1
        login = f'user_{self._user_id}'
        password = login
        return self.register(login, password)