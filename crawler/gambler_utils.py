from requests import get, post, Session
from urllib.parse import urlparse, parse_qs, urlencode
from bs4 import BeautifulSoup
from datetime import datetime
from time import sleep
from collections import defaultdict
import re
import numpy as np
import pandas as pd
from concurrent.futures import ThreadPoolExecutor
import pickle
from http_request_randomizer.requests.proxy.requestProxy import RequestProxy
import time


MAIN_URL = 'https://www.playelephant.com'
LOGIN = urlencode({'gloginname': 'Цезарион', 'gloginpsw': 'yPKkhXDi15'})
HEADERS = {"Host": "www.playelephant.com",
           "Content-Type": "application/x-www-form-urlencoded",
           "Content-Length": "148"}
GAME_TYPES = {'preferance': 9, 'races': 37, 'bandit': 36, 'gusarik': 42, 'almaty': 84}
TEXT_GAME_TYPES = {'preferance': 'Преферанс', 'races': 'Скачки', 'bandit': 'Разбойник',
                   'gusarik': 'Гусарик', 'almaty': 'Алматинка'}
GAMENO = '/php/protocol?{}'
USER_INFO = '/user/info?{}'
PREF_ON_FRIDAYS = 'https://www.playelephant.com/team/transfers?tid=9902'
MONTHS = {u'янв': '01', u'фев': '02', u'мар': '03', u'апр': '04',
          u'май': '05', u'июн': '06', u'июл': '07', u'авг': '08',
          u'сен': '09', u'окт': '10', u'ноя': '11', u'дек': '12'}
TRANSFER_UIN = 611386
TOP = '/php/top?{}'
transfer_columns = ['uin', 'nick', '#transfers', 'transfer timestamps']
user_columns = ['uin', 'nick', 'rating', 'games_played']
USER_TABLE_COLUMNS = ['uin', 'nick', 'reg_date']
TABLE_SIZES = {'preferance': 3000, 'races': 300, 'bandit': 7600, 'gusarik': 5200, 'almaty': 0}
ROWS_LIMIT = 100
BACKWARD = '<<<'
FORWARD = '>>>'
PREFERANCE_TYPES = ('ЖП', 'Питер', 'Сочи', 'Ростов', 'ВМК')
WRAPPERS = (' 100', ' 500', ' 1000', ' 5000')
TIMESTAMP_LOWER = 1328072400
GAMENO_LOWER = 170000000
skipped_filename = 'preferance_skipped'
archive_filename = 'preferance_archive.csv'
users_filename = 'preferance_table.csv'


def submit_workers(func, sess, iterable, n_workers=10):
    executor = ThreadPoolExecutor(max_workers=n_workers)
    futures = []
    for item in iterable:
        future_result = executor.submit(func, sess, item)
        futures.append(future_result)
    results = [f.result() for f in futures]
    return results


def sort_list_by_another(lst, sorting_lst):
    return [x for _, x in sorted(zip(sorting_lst, lst), reverse=True)]


def timestamp2gameno(timestamp):
    return int(0.884 * timestamp - 1002003970)
