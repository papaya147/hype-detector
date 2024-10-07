import yfinance as yf
from datetime import datetime, timedelta
import os
import pandas as pd


def save_stock_data_1m() -> None:
    with open("../stock list.txt", "r") as file:
        data = file.read()
    stock_list = list(set(data.split()))

    stock_tickers = yf.Tickers(" ".join(stock_list))

    for stock in stock_list:
        stock_ticker = stock_tickers.tickers[stock]

        end_date = datetime.now().strftime('%Y-%m-%d')
        start_date = (datetime.now() - timedelta(days=7)).strftime('%Y-%m-%d')
        exists = False

        # check if the csv exists already, if it does, make start_date from the latest date
        try:
            file_name = f"../stock-prices/{stock}.csv"
            data = pd.read_csv(file_name)
            latest_date_record = data["Datetime"][len(data) - 1]
            next_date = (datetime.strptime(latest_date_record,
                                           "%Y-%m-%d %H:%M:%S%z") + timedelta(days=1))
            start_date = next_date.strftime('%Y-%m-%d')
            exists = True
        except Exception as e:
            pass

        hist = stock_ticker.history(start=start_date, end=end_date,
                                    period="max", interval="1m")
        hist.to_csv(f"../stock-prices/{stock}.csv",
                    mode="a" if exists else "w", header=not exists)


def save_stock_data_1h() -> None:
    with open("../stock list.txt", "r") as file:
        data = file.read()
    stock_list = list(set(data.split()))

    stock_tickers = yf.Tickers(" ".join(stock_list))

    for stock in stock_list:
        stock_ticker = stock_tickers.tickers[stock]

        end_date = datetime.now().strftime('%Y-%m-%d')
        start_date = (datetime.now() -
                      timedelta(days=729)).strftime('%Y-%m-%d')
        exists = False

        # check if the csv exists already, if it does, make start_date from the latest date
        try:
            file_name = f"../stock-prices-1h/{stock}.csv"
            data = pd.read_csv(file_name)
            latest_date_record = data["Datetime"][len(data) - 1]
            next_date = (datetime.strptime(latest_date_record,
                                           "%Y-%m-%d %H:%M:%S%z") + timedelta(days=1))
            start_date = next_date.strftime('%Y-%m-%d')
            exists = True
        except Exception as e:
            pass

        hist = stock_ticker.history(start=start_date, end=end_date,
                                    period="max", interval="1h")
        hist.to_csv(f"../stock-prices-1h/{stock}.csv",
                    mode="a" if exists else "w", header=not exists)


save_stock_data_1m()
save_stock_data_1h()
