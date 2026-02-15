import pytest

from utils.request_manager import RequestManager
from utils.db_manager import DbManager


def pytest_addoption(parser):
    parser.addoption(
        "--url",
        action="store",
        # default="http://localhost:8080",
        help="Base URL of the server (e.g. http://localhost:8080)",
    )
    parser.addoption(
        "--db",
        action="store",
        # default="postgresql://postgres:postgres@localhost:5432/yuki_buy_log",
        help="PostgreSQL DSN (e.g. postgresql://user:pass@host:port/dbname)",
    )


@pytest.fixture(scope="session")
def db(request):
    base_url = request.config.getoption("--db")
    manager = DbManager(base_url)
    yield manager
    manager.close()


@pytest.fixture(scope="session")
def req(request):
    url = request.config.getoption("--url")
    manager = RequestManager(url)
    yield manager


