import time

import psycopg
import psycopg.rows


class DbManager:
    def __init__(self, dsn: str):
        self._dsn = dsn
        self._conn = None

    def _ensure_connection(self):
        if self._conn is None or self._conn.closed:
            self._conn = psycopg.connect(self._dsn)

    def close(self):
        if self._conn and not self._conn.closed:
            self._conn.close()
        self._conn = None

    def execute(self, query: str, params: tuple = ()) -> list[dict]:
        try:
            return self._try_execute(query, params)
        except Exception:
            time.sleep(5)
            self._conn = None
            return self._try_execute(query, params)

    def _try_execute(self, query: str, params: tuple) -> list[dict]:
        self._ensure_connection()
        with self._conn.cursor(row_factory=psycopg.rows.dict_row) as cur:
            cur.execute(query, params)
            return cur.fetchall()