#!/usr/bin/env python3

import sys
import json
from os.path import basename

from requests import Request, Session
from requests.exceptions import ConnectionError, Timeout, TooManyRedirects

API_URL = 'https://sandbox-api.coinmarketcap.com/v1'
# API_URL = 'https://pro-api.coinmarketcap.com/v1'

API_KEY = 'b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c'

def cook_headers(api_key):
    return {
        'Accepts': 'application/json',
        'X-CMC_PRO_API_KEY': api_key
    }

def process_error(data):
    err_msg = data['status']['error_message']
    err_list = err_msg.split(":")
    if len(err_list) == 2:
        if err_list[0] in ('Invalid value for "convert"', 'Invalid value for "symbol"'):
            print(f'Unknown currency symbol: {err_list[1]}')
        else:
            print(err_msg)
    else:
        print(err_msg)

def get_price(data, from_currency, to_currency):
    if from_currency in data['data']:
        # Sandbox API
        return data['data'][from_currency]['quote'][to_currency]['price']
    else:
        # Pro API 
        return data['data']['quote'][to_currency]['price']

def convert(amount, from_currency, to_currency):
    session = Session()
    session.headers.update(cook_headers(API_KEY))
    
    url = f'{API_URL}/tools/price-conversion'
    params = {
        'symbol': from_currency,
        'convert': to_currency,
        'amount': amount,
    }
    
    try:
        response = session.get(url, params=params)
        data = json.loads(response.text)
        print(json.dumps(data, indent=4))
        if data['status']['error_code'] == 0:
            price = get_price(data, from_currency, to_currency)
            print(f'{amount} {from_currency} is {price} {to_currency}')
        else:
            process_error(data)
    except (ConnectionError, Timeout, TooManyRedirects) as e:
        print(e)

def print_usage(name):
    msg = """
Yet Another Crypto Currency Converter

Usage:
{} <AMOUNT> <FROM_CURRENCY_SYMBOL> <TO_CURRENCY_SYMBOL>

For example:
{} 0.5 BTC USD
"""
    print(msg.format(name, name))

def do():
    if len(sys.argv) < 4:
        print_usage(basename(sys.argv[0]))
        sys.exit(4)

    try:
        amount = float(sys.argv[1])
    except ValueError as e:
        print(f'Amount should be integer of float, not "{sys.argv[1]}"')
        sys.exit(4)
        
    convert(amount, sys.argv[2].upper(), sys.argv[3].upper())

if __name__ == '__main__':
    do()
