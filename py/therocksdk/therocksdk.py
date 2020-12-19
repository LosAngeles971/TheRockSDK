"""
This module hosts the wrapper to The Rock Trading API v1.

The module does not wrap all functionalities exposed by API, but:



"""
import json
import requests
import hmac
import hashlib
import time
from remida.exchange import exchange
from remida.financial import wallets
import remida.system.logger as l


API_BASEURL = 'https://api.therocktrading.com/v1/'

def getheaders(url, key, secret):
    if key is None:
        return {
            'User-Agent': 'PyRock v1', 
            'content-type': 'application/json'
        }
    else:
        nonce   = str(int(time.time() * 1e6))
        message = str(nonce) + url
        return {
            'User-Agent': 'PyRock v1', 
            'content-type': 'application/json', 
            'X-TRT-KEY': key, 
            'X-TRT-SIGN': hmac.new(secret.encode(), msg = message.encode(), digestmod = hashlib.sha512).hexdigest(), 
            'X-TRT-NONCE': nonce
        }

def get(url, key, secret, timeout = 30):
    r = requests.get(url, headers = getheaders(url, key, secret), timeout = timeout)
    return r.json()

class TheRockTrading:

    FUNDS       = None
    PRECISIONS  = {}
    ro_key      = None
    ro_secret   = None
    rw_key      = None
    rw_secret   = None

    def __init__(self,  ro_key, ro_secret, rw_key, rw_secret):
        self.ro_key    = ro_key
        self.ro_secret = ro_secret
        self.rw_key    = rw_key
        self.rw_secret = rw_secret

    def __get_funds(self):
        """
        This method provides the list of supported funds, with the characteristics of each funds.
        All funds are cached into memory.

        Response's format is:
        [
            {
                "id":"BTCEUR",
                "description":"Trade Bitcoin with Euro",
                "type":"currency",
                "base_currency":"EUR",
                "trade_currency":"BTC",
                "buy_fee":0.2,
                "sell_fee":0.2,
                "minimum_price_offer":0.01,
                "minimum_quantity_offer":0.0005,
                "base_currency_decimals":2,
                "trade_currency_decimals":4,
                "leverages":[]
            },
            ...
        ]
        """
        if self.FUNDS is None:
            self.FUNDS = get(API_BASEURL + 'funds', self.ro_key, self.ro_secret)['funds']
        return self.FUNDS

    def __get_ticker_from_json(self, t):
        ticker = exchange.Ticker()
        ticker.asset = t['fund_id'][:3]
        ticker.currency = t['fund_id'][3:]
        ticker.date = t['date']
        ticker.last = t['last']
        ticker.bid = t['bid']
        ticker.ask = t['ask']
        ticker.open = t['open']
        ticker.close = t['close']
        ticker.high = t['high']
        ticker.low = t['low']
        ticker.volume = t['volume']
        return ticker

    def get_precision(self, asset):
        if self.PRECISIONS.get(asset, None) is not None:
            return self.PRECISIONS[asset], 'down'
        j = get(API_BASEURL + 'currencies/' + asset, None, None)
        self.PRECISIONS[asset] = j['decimals']
        return self.PRECISIONS[asset], 'down'

    def get_ticker(self, asset, currency):
        url = API_BASEURL + 'funds/' + asset + currency + '/ticker'
        j = requests.get(url, headers = getheaders(url, None, None), timeout = 30).json()
        l.debug(j)
        return self.__get_ticker_from_json(j)

    def get_market(self):
        """
        Get the current status of the market for all supported trades
        """
        j = get(API_BASEURL + 'funds/tickers', self.ro_key, self.ro_secret)
        tickers = []
        for t in j['tickers']: tickers.append(self.__get_ticker_from_json(t))
        return tickers

    def check_trade_validity(self, asset, currency):
        fundid = asset + currency
        r_fundid = currency + asset
        for f in self.__get_funds():
            if f['id'] == fundid or f['id'] == r_fundid:
                return True
        return False

    def __match_to_findid(self, asset, currency):
        """
        This method is necessary to the complicated world of the FUCKING Exchanges.

        In the fucking world of TRT, there are fundid like BTCEUR.
        BTCEUR means buying BTC with EUR.

        So if the asset is EUR and currency is BTC, you have to convert the buy intent into a
        sell order.
        """
        fundid = asset + currency
        r_fundid = currency + asset
        for f in self.__get_funds():
            if f['id'] == fundid:
                return 'match'
            if f['id'] == r_fundid:
                return 'reverse'
        return 'no'

    def get_wallets(self):
        """
        Respose is:
        {
            "balances": [
                {"currency":"LTC","balance":6.50884835,"trading_balance":2.30884835},
                {"currency":"BTC","balance":3.50884835,"trading_balance":1.30884835}
            ]
        }
        """
        aw = wallets.create(self.credential.username)
        url = API_BASEURL + 'balances'
        r = requests.get(url, headers = getheaders(url, self.credential.ro_apikey, self.credential.ro_apisecret), timeout = 30)
        for b in r.json()['balances']:
            w = wallets.Wallet()
            w.asset = b['currency']
            w.amount = b['balance']
            w.available = b['trading_balance']
            aw.add(w)
        return aw

    def get_trade_fee(self, order: exchange.Order):        
        return order.amount_out * 0.02, order.c_out

    def trt_buy(self, fund, amount, price):
        """
        ATTENTION:
        - amount    : the amount you want to Buy/Sell, so the amount you want to BUY is amount_in
        - price     : is the price of your order to be filled. If price is 0 (zero) a market order will be placed.
        """
        url = API_BASEURL + 'funds/' + fund.upper() + '/orders'
        values = { 
            'fund_id' : fund.upper(),
            'side' : 'buy',
            'amount' : str(amount),
            'price' : str(price)
        }
        r = requests.post(url, data = json.dumps(values), headers = getheaders(url, self.credential.ro_apikey, self.credential.ro_apisecret), timeout = 30)
        return r.json()

    def trt_sell(self, fund, amount, price):
        """
        ATTENTION:
        - amount    : the amount you want to Buy/Sell, so the amount you want to BUY is amount_out
        - price     : is the price of your order to be filled. If price is 0 (zero) a market order will be placed.
        """
        url = API_BASEURL + 'funds/' + fund.upper() + '/orders'
        values = { 
        'fund_id' : fund.upper(),
        'side' : 'sell',
        'amount' : amount,
        'price' : price
        }
        r = requests.post(url, data = json.dumps(values), headers = getheaders(url, self.credential.ro_apikey, self.credential.ro_apisecret), timeout = 30)
        return r.json()

    def __convert_order_status(self, label):
        if label == 'active':
            return exchange.Order.STATUS_RUNNING
        elif label == 'conditional':
            return exchange.Order.STATUS_RUNNING
        elif label == 'executed':
            return exchange.Order.STATUS_COMPLETED
        elif label == 'deleted':
            return exchange.Order.STATUS_FAILED
        else:
            raise Exception('Unrecognized status: ' + str(label))

    def prepare_buy(self, asset, currency, amount_out, price, expire = 365):
        m = self.__match_to_findid(asset, currency)
        if m == 'match':
            return self.goaskalice_buy(asset, currency, amount_out, price, expire = expire)
        elif m == 'reverse':
            l.warning('Indirect buy for asset: ' + asset + ' with currency: ' + currency)
            """
            I would BUY asset (c_in) PAYING with currency (c_out) but TRT misses the fundid asset+currency,
            while TRT supports the fundid is currency+asset
            Thus I need the operation:
            I want to SELL currency (c_out) to HAVE asset (c_in) 
            """
            return self.goaskalice_sell(currency, asset, amount_out, price, expire = expire)
        else:
            raise Exception('This trade is not supported, asset: ' + str(asset) + ' currency: ' + str(currency))

    def prepare_sell(self, asset, currency, amount_out, price, expire = 365):
        m = self.__match_to_findid(asset, currency)
        if m == 'match':
            return self.goaskalice_sell(asset, currency, amount_out, price, expire = expire)
        elif m == 'reverse':
            l.warning('Indirect sell for asset: ' + asset + ' with currency: ' + currency)
            """
            I would SELL asset (c_out) to HAVE currency (c_in) but TRT misses the fundid asset+currency,
            while TRT supports the fundid is currency+asset
            Thus I need the operation:
            I want to BUY currency (c_in) PAYING with asset (c_out) 
            """
            return self.goaskalice_buy(currency, asset, amount_out, price, expire = expire)
        else:
            raise Exception('This trade is not supported, asset: ' + str(asset) + ' currency: ' + str(currency))

    def place(self, order: exchange.Order):
        """
        Response of TRT API is:
        HTTP/1.1 200 OK
        {
            "id": 4325578,
            "fund_id":"BTCEUR",
            "side":"buy",
            "type":"limit",
            "status":"executed",
            "price":0.0102,
            "amount": 50.0,
            "amount_unfilled": 0.0,
            "conditional_type": null,
            "conditional_price": null,
            "date":"2015-06-03T00:49:48.000Z",
            "close_on": nil,
            "leverage": 1.0,
            "position_id": null,
            "trades": [
                { 
                "id":237338,
                "fund_id":"BTCEUR",
                "amount":50,
                "price":0.0102,
                "side":"buy",
                "dark":false,
                "date":"2015-06-03T00:49:49.000Z"
                }
            ]
        }
        """
        if self.check_trade_validity(order.asset, order.currency) is False:
            raise Exception('This trade is not supported, asset: ' + str(order.asset) + ' currency: ' + str(order.currency))
        if order.operation == exchange.Order.TYPE_BUY:
            j = self.trt_buy(order.asset + order.currency, order.amount_in, order.price)
        elif order.operation == exchange.Order.TYPE_SELL:
            j = self.trt_sell(order.asset + order.currency, order.amount_out, order.price)
        else:
            raise Exception('Unrecognized order type: ' + str(order.operation))
        order.exchange_id = j['id']
        order.status = self.__convert_order_status(j['status'])
        order.price = j['price']
        order.sprice = str(j['price'])
    
    def get_order_status(self, order: exchange.Order):
        """
        """
        if self.check_trade_validity(order.asset, order.currency) is False:
            raise Exception('This trade is not supported, asset: ' + str(order.asset) + ' currency: ' + str(order.currency))
        url = API_BASEURL + 'funds/' + order.asset + order.currency +'/orders/' + str(order.exchange_id)
        j = requests.get(url, headers = getheaders(url, self.credential.ro_apikey, self.credential.ro_apisecret), timeout = 30).json()
        order.exchange_id = j['id']
        order.status = self.__convert_order_status(j['status'])
        order.price = j['price']
        order.sprice = str(j['price'])
        return order.status
