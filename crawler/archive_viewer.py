from gambler_utils import *
general_columns = ['gameno', 'timestamp', 'duration', 'deals', 'game_id',
                   'game_type', 'tournament', 'wrappers', 'timeout']
result_columns = ['uin1', 'result1', 'rate1', 'change_rate1',
                  'uin2', 'result2', 'rate2', 'change_rate2',
                  'uin3', 'result3', 'rate3', 'change_rate3',
                  'uin4', 'result4', 'rate4', 'change_rate4']


def extract_datetime_and_duration(string):
    dt, dur = string.split(',')
    dt, dur = dt.strip(), dur.strip()

    if not dt[1].isdigit():
        dt = '0' + dt
    month = dt[3:6]
    dt = dt.replace(month, MONTHS[month])

    # current year
    if '-' not in dt:
        dt = dt[:5] + ' ' + str(datetime.now().year) + dt[5:]
    else:
        dt = dt.replace('-', ' ')

    digits = re.findall(r'\d\d?', dur)
    if len(digits) == 1:
        duration = int(digits[0])
    elif len(digits) == 2:
        duration = int(digits[1]) + 60 * int(digits[0])
    else:
        raise ValueError('Bad duration!')

    return dt, duration


class UserViewer:
    def __init__(self, filename=None):
        self.uins = set()
        if filename is not None:
            saved_table = pd.read_csv(filename)
            print("Table with {} rows was loaded from {}".format(saved_table.shape[0], filename))
            self.uins.update(set(saved_table['uin']))
        self.bad_users = set()
        self.result = []

    def get_register_date(self, sess, uin):
        url = MAIN_URL + USER_INFO.format(urlencode({'user': uin}))
        while True:
            try:
                r = sess.get(url)
            except Exception as e:
                print("Error {} while getting {} occured, sleep for 3 second...".format(e, url))
                sleep(3)
                continue
            if r.ok:
                r.encoding = 'utf-8'
                soup = BeautifulSoup(r.text, "lxml")
                try:
                    table = soup.find('table', attrs={'class': 'gtborder', 'align': 'center'})
                    reg = table.find('span', title=re.compile(r'\d{2}:\d{2}'))
                    nick = soup.find('h2', attrs={'style': 'margin-bottom: 1pt'}).text
                except Exception as e:
                    print('Error {} in uin={}'.format(e, uin))
                    self.bad_users.add(uin)
                    break
                self.uins.add(uin)
                dt = reg.text + ' ' + reg['title']
                self.result.append([uin, nick, dt])
                break      
            else:
                raise ConnectionError("Error while getting url={}, code= {}".format(url, r.status_code))

    def view_datetimes(self, uins, max_workers=20):
        with Session() as s:
            r = s.post(MAIN_URL, data=LOGIN, headers=HEADERS)
            if not r.ok:
                raise ConnectionError('Login post failed, code {}'.format(r.status_code))
            uins_to_handle = uins - self.uins
            user_dt = submit_workers(self.get_register_date, s, uins_to_handle, n_workers=max_workers)
            return self.result


class ArchiveViewer:
    def __init__(self, game, user_table_filename=users_filename, archive_table_filename=archive_filename):
        self.archive_table = pd.read_csv(archive_table_filename)
        print("Archive table with {} games has been loaded".format(self.archive_table.shape[0]))

        self.user_table = pd.read_csv(user_table_filename)
        self.unviewed_user_table = self.user_table[self.user_table['viewed'] == 0]
        print("Table with {} users has been loaded, {} unviewed users among them".format(
            self.user_table.shape[0], self.unviewed_user_table.shape[0]
        ))
        self.all_games = set(self.archive_table['gameno'])
        self.viewed_nicks = dict()
        self.viewed_reg_dates = dict()
        
        with open(skipped_filename, 'rb') as f:
            self.skipped_games = pickle.load(f)
        self.game = game
        self.games_to_add = list()

        start = time.time()
        self.req_proxy = RequestProxy()
        print("Initialization took: {0} sec".format((time.time() - start)))
        print("Size: {0}".format(len(self.req_proxy.get_proxy_list())))
        print("ALL = {0} ".format(list(map(lambda x: x.get_address(), self.req_proxy.get_proxy_list()))))

    def get_viewed_players(self, size):
        self.games_to_add = list()
        subtable = self.unviewed_user_table[:min(size, self.unviewed_user_table.shape[0])]
        self.viewed_nicks = dict(zip(subtable['uin'], subtable['nick']))
        self.viewed_reg_dates = dict(zip(subtable['uin'], subtable['reg_datetime']))
        return list(subtable['uin'])

    @staticmethod
    def continue_viewing(link):
        if link is None:
            return False
        query_dict = parse_qs(urlparse(link).query)
        if 'fintime' in query_dict:
            return TIMESTAMP_LOWER < int(query_dict['fintime'][0])
        if 'togameno' in query_dict:
            return GAMENO_LOWER < int(query_dict['togameno'][0])
        return True

    def handle_exception(self, game_no):
        if game_no == 194921494:
            general_data = [game_no, int(datetime(2012, 11, 2, 19, 58).timestamp()), 40, 27,
                            GAME_TYPES[self.game], 'Преферанс', 1, 0, 1]
            result_data = [1008710, 234, 162, 5, 439224, -34, -2, 679118, -200, 128, -3]
            return dict(zip(general_columns + result_columns, general_data + result_data))

    def extract_games(self, uin, rows):
        archive_addon = pd.DataFrame(columns=general_columns + result_columns)
        for row in rows:
            players = []
            results = []
            rates = []
            change_rates = []
            cells = row.findAll('td')

            title_a = cells[1].findAll('a')
            if not title_a:
                # actually there was no play
                continue

            game_no = re.findall(re.compile(r'gameno=(\d{9})'), title_a[0]['href'])
            if not game_no:
                # no gameno - no game
                continue

            game_no = int(game_no[0][-9:])
            if game_no in self.all_games:
                # skip game if it has been handled already
                continue

            if game_no == 194921494:
                archive_addon = archive_addon.append(self.handle_exception(game_no), ignore_index=True)
                continue

            if self.game == 'preferance':
                text_type = 'Преферанс'
                for tt in TEXT_GAME_TYPES:
                    if tt in title_a[0].text:
                        text_type = tt
                        break
            else:
                text_type = TEXT_GAME_TYPES[self.game]

            is_wrappers = False
            if self.game == 'preferance':
                for wrap in WRAPPERS:
                    if wrap in title_a[0].text:
                        is_wrappers = True
                        break

            is_tournament = False
            if len(title_a) > 1 and 'турнир' in title_a[1].text:
                is_tournament = True

            dt_dur = cells[0].find('span', attrs={'class': 'comment'})
            # bug in https://www.playelephant.com/user/archive?uin=895916&game=9&fromgameno=253861827
            if dt_dur.text == 'играется':
                continue
            try:
                dt, dur = extract_datetime_and_duration(dt_dur.text)
            except Exception as e:
                print('Error {} in gameno={}, dt_dur={}'.format(e, game_no, dt_dur.text))
                self.skipped_games.add(game_no)
                continue

            date_time = datetime.strptime(dt, '%d %m %Y %H:%M')

            try:
                deals = int(re.findall(re.compile(r'\d+'), str(cells[1]['title']))[0])
            except Exception as e:
                print('Error {} in gameno={}, cells_1={}'.format(e, game_no, cells[1]['title']))
                self.skipped_games.add(game_no)
                continue

            try:
                players.append(uin)
                results.append(int(cells[2].text))
            except Exception as e:
                print('Error {} in gameno={}, cells_2={}'.format(e, game_no, cells[2].text))
                self.skipped_games.add(game_no)
                continue

            span_4 = cells[4].find('span')
            if span_4 is None:
                # no rate info - probably no game
                print("No rate info in game {}".format(game_no))
                self.skipped_games.add(game_no)
                continue

            if 'title' in span_4.attrs:
                title_text = span_4['title']
                if title_text != 'не на рейтинг':
                    try:
                        rates.append(int(re.findall(re.compile(r'[+-]?\d+'), span_4['title'])[0]))
                        change_rates.append(int(span_4.text))
                    except Exception as e:
                        print('Error {} in gameno={}, span_4_title={}, span_4_text={}'.format(
                            e, game_no, span_4['title'], span_4.text)
                        )
                        continue
            else:
                try:
                    change_rates.append(int(span_4.text))
                    rates.append(int(re.findall(re.compile(r'\(([+-]?\d+)\)'), cells[4].text)[0]))
                except Exception as e:
                    print('Error {} in gameno={}, span_4={}, cells_4={}'.format(
                        e, game_no, span_4.text, cells[4].text)
                    )
                    continue

            timeout = False
            if 'тайм-аут' in str(cells[4]):
                timeout = True

            opps = cells[3].findAll('span')
            for opp in opps:
                if rates:
                    numbers = re.findall(re.compile(r'[+-]?\d+'), opp['title'])
                    assert (len(numbers) == 3)
                    results.append(int(numbers[0]))
                    rates.append(int(numbers[1]))
                    change_rates.append(int(numbers[2]))
                else:
                    results.append(int(opp['title']))
                opp_uin = int(re.findall(re.compile(r'uin=(\d+)'), str(opp.find('a')['href']))[0])
                players.append(opp_uin)
                # self.all_players[opp_uin] = opp.find('a').text
                # if opp_uin not in self.new_players:
                #     self.new_players[opp_uin] = opp.find('a').text

            players = sort_list_by_another(players, results)
            if rates:
                rates = sort_list_by_another(rates, results)
                change_rates = sort_list_by_another(change_rates, results)
            results = sorted(results, reverse=True)

            self.all_games.add(game_no)

            general_data = [game_no, int(date_time.timestamp()), dur, deals, GAME_TYPES[self.game],
                            text_type, int(is_tournament), int(is_wrappers), int(timeout)]
            result_data = []
            added_columns = list(general_columns)
            for i in range(len(players)):
                result_data += [players[i], results[i]]
                added_columns += result_columns[4 * i:4 * i + 2]
                if rates:
                    result_data += [rates[i], change_rates[i]]
                    added_columns += result_columns[4 * i + 2:4 * i + 4]

            archive_addon = archive_addon.append(
                dict(zip(added_columns, general_data + result_data)), ignore_index=True
            )
        archive_addon[result_columns] = archive_addon[result_columns].astype(object)
        self.games_to_add.append(archive_addon)
        return archive_addon.shape[0]

    def view_archive(self, sess, uin):
        print("Начинается просмотр архива игрока {}...".format(self.viewed_nicks[uin]))
        direction = FORWARD
        begin = max(datetime.strptime(self.viewed_reg_dates[uin], '%d.%m.%y %H:%M').timestamp(), TIMESTAMP_LOWER)
        previous_link = MAIN_URL + '/user/archive?{}'.format(
                urlencode({'uin': uin, 'game': GAME_TYPES[self.game],
                           'fromgameno': timestamp2gameno(begin)})
            )
        n_games_viewed = 0
        n_games_added = 0
        while self.continue_viewing(previous_link):
            try:
                r = sess.get(previous_link, proxies={"http": self.req_proxy.randomize_proxy()})
            except Exception as e:
                print("Error {} while getting {} occured, sleep for 5 second...".format(e, previous_link))
                sleep(5)
                continue
            if r.ok:
                soup = BeautifulSoup(r.text, 'lxml')
                arrows = soup.find('a', string=direction)
                previous_link = MAIN_URL + '/user/' + str(arrows['href']) if arrows else None
                rows = soup.findAll('tr', attrs={'valign': 'middle'})
                n_games_viewed += len(rows)
                if rows:
                    n_games_added += self.extract_games(uin, rows)
        print("Архив игрока {} просмотрен, {} / {} игр добавлено".format(
            self.viewed_nicks[uin], n_games_added, n_games_viewed
        ))
        self.user_table.loc[self.user_table['uin'] == uin, 'viewed'] = 1
        return [n_games_added, n_games_viewed]

    def view_user_archives(self, user_size=20, coef=1.0):
        users_to_view = self.get_viewed_players(user_size)
        with Session() as s:
            r = s.post(MAIN_URL, data=LOGIN, headers=HEADERS,
                       proxies={"http": self.req_proxy.randomize_proxy()})
            if not r.ok:
                raise ConnectionError('Login post failed, code {}'.format(r.status_code))
            n_game_pairs = submit_workers(
                self.view_archive, s, users_to_view, int(coef * len(users_to_view))
                )
            n_added, n_viewed = list(np.sum(n_game_pairs, axis=0))
            percentage = 100.0 * n_added / n_viewed if n_viewed else 0.0
            print("Добавлено {} / {} = {}% игр".format(n_added, n_viewed, percentage))
            with open(skipped_filename, 'wb') as f:
                pickle.dump(self.skipped_games, f)
            added_archive = pd.concat(self.games_to_add, ignore_index=True)
            self.archive_table = self.archive_table.append(added_archive, ignore_index=True)
            self.archive_table.to_csv(archive_filename, index=False)
            self.user_table.to_csv(users_filename, index=False)
            self.unviewed_user_table = self.user_table[self.user_table['viewed'] == 0]
