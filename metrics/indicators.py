import pandas as pd
import numpy as np


def open(data: pd.DataFrame) -> pd.Series:
    assert "Open" in data.columns, "Open not present in data"
    return data["Open"]


def close(data: pd.DataFrame) -> pd.Series:
    assert "Close" in data.columns, "Close not present in data"
    return data["Close"]


def high(data: pd.DataFrame) -> pd.Series:
    assert "High" in data.columns, "High not present in data"
    return data["High"]


def low(data: pd.DataFrame) -> pd.Series:
    assert "Low" in data.columns, "Low not present in data"
    return data["Low"]


def volume(data: pd.DataFrame) -> pd.Series:
    assert "Volume" in data.columns, "Volume not present in data"
    return data["Volume"]


def timestamp(data: pd.DataFrame) -> pd.Series:
    assert "Datetime" in data.columns, "Datetime not present in data"
    return pd.to_datetime(data["Datetime"])


def dividends(data: pd.DataFrame) -> pd.Series:
    assert "Dividends" in data.columns, "Dividends not present in data"
    return data["Dividends"]


def engulfing(data: pd.DataFrame) -> pd.Series:
    open_prices = open(data)
    close_prices = close(data)

    res = [np.nan]

    for i in range(1, len(data)):
        prev_open = open_prices.iloc[i - 1]
        prev_close = close_prices.iloc[i - 1]
        curr_open = open_prices.iloc[i]
        curr_close = close_prices.iloc[i]

        if curr_close > curr_open and curr_close > prev_open and curr_open < prev_close:
            res.append("bullish")
        elif curr_close < curr_open and curr_close < prev_open and curr_open > prev_close:
            res.append("bearish")
        else:
            res.append("none")

    res = pd.Series(res)
    res.name = "engulfing"

    return res


def marubozo(data: pd.DataFrame) -> pd.Series:
    open_prices = open(data)
    close_prices = close(data)
    high_prices = high(data)
    low_prices = low(data)

    res = []

    for i in range(len(data)):
        curr_open = open_prices.iloc[i]
        curr_close = close_prices.iloc[i]
        curr_high = high_prices.iloc[i]
        curr_low = low_prices.iloc[i]

        if curr_open == curr_low and curr_close == curr_high:
            res.append("bullish")
        elif curr_open == curr_high and curr_close == curr_low:
            res.append("bearish")
        else:
            res.append("none")

    res = pd.Series(res)
    res.name = "marubozo"

    return res


def doji(data: pd.DataFrame, epsilon: float = 0.01) -> pd.Series:
    '''
    1 for doji, 0 for not
    '''
    open_prices = data['Open']
    close_prices = data['Close']
    high_prices = data['High']
    low_prices = data['Low']

    res = []

    for i in range(len(data)):
        curr_open = open_prices.iloc[i]
        curr_close = close_prices.iloc[i]

        if abs(curr_open - curr_close) <= epsilon * (high_prices.iloc[i] - low_prices.iloc[i]):
            res.append(1)
        else:
            res.append(0)

    res = pd.Series(res)
    res.name = "doji"

    return res


def hammer(data: pd.DataFrame, body_ratio: float = 0.3, shadow_ratio: float = 2.0) -> pd.Series:
    '''
    1 for hammer, 0 for not
    '''
    open_prices = data['Open']
    close_prices = data['Close']
    high_prices = data['High']
    low_prices = data['Low']

    res = []

    for i in range(len(data)):
        curr_open = open_prices.iloc[i]
        curr_close = close_prices.iloc[i]
        curr_high = high_prices.iloc[i]
        curr_low = low_prices.iloc[i]

        body_length = abs(curr_close - curr_open)
        lower_shadow_length = min(curr_open, curr_close) - curr_low
        upper_shadow_length = curr_high - max(curr_open, curr_close)

        if (body_length <= body_ratio * (curr_high - curr_low)) and \
           (lower_shadow_length >= shadow_ratio * body_length) and \
           (upper_shadow_length <= body_length):
            res.append(1)
        else:
            res.append(0)

    res = pd.Series(res)
    res.name = "hammer"

    return res


def inverted_hammer(data: pd.DataFrame, body_ratio: float = 0.3, shadow_ratio: float = 2.0) -> pd.Series:
    '''
    1 for inverted hammer, 0 for not
    '''
    open_prices = data['Open']
    close_prices = data['Close']
    high_prices = data['High']
    low_prices = data['Low']

    res = []

    for i in range(len(data)):
        curr_open = open_prices.iloc[i]
        curr_close = close_prices.iloc[i]
        curr_high = high_prices.iloc[i]
        curr_low = low_prices.iloc[i]

        body_length = abs(curr_close - curr_open)
        upper_shadow_length = curr_high - max(curr_open, curr_close)
        lower_shadow_length = min(curr_open, curr_close) - curr_low

        if (body_length <= body_ratio * (curr_high - curr_low)) and \
           (upper_shadow_length >= shadow_ratio * body_length) and \
           (lower_shadow_length <= body_length):
            res.append(1)
        else:
            res.append(0)

    res = pd.Series(res)
    res.name = "inverted_hammer"

    return res


def macd(data: pd.DataFrame, short_span: int = 12, long_span: int = 26, signal_span: int = 9) -> pd.DataFrame:
    short = close(data).ewm(span=short_span, adjust=False).mean()
    long = close(data).ewm(span=long_span, adjust=False).mean()

    macd = short - long
    signal = macd.ewm(span=signal_span, adjust=False).mean()
    histogram = macd - signal

    return pd.DataFrame({
        "macd": macd,
        "macd_signal": signal,
        "macd_histogram": histogram
    })


def ewma(data: pd.DataFrame, span: int) -> pd.Series:
    ewma = close(data).ewm(span=span, adjust=False).mean()
    return ewma


def bollinger_bands(data: pd.DataFrame, period: int = 20, num_std_dev: float = 2) -> pd.DataFrame:
    close_prices = close(data)

    middle_band = close_prices.rolling(window=period).mean()

    std_dev = close_prices.rolling(window=period).std()

    upper_band = middle_band + (std_dev * num_std_dev)
    lower_band = middle_band - (std_dev * num_std_dev)

    return pd.DataFrame({
        "bollinger_middle_band": middle_band / close_prices,
        "bollinger_upper_band": upper_band / close_prices,
        "bollinger_lower_band": lower_band / close_prices
    })


def rsi(data: pd.DataFrame, period: int = 14) -> pd.Series:
    close_prices = close(data)

    delta = close_prices.diff()

    gain = (delta.where(delta > 0, 0))
    loss = (-delta.where(delta < 0, 0))

    avg_gain = gain.rolling(window=period, min_periods=1).mean()
    avg_loss = loss.rolling(window=period, min_periods=1).mean()

    rs = avg_gain / avg_loss

    rsi = 1 - (1 / (1 + rs))
    rsi.name = "rsi"

    return rsi
