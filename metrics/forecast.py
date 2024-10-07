import pandas as pd
from metrics import indicators


def precise_forecast(data: pd.DataFrame, span: int = 10, threshold: float = 0.01) -> pd.Series:
    '''
    Checks forecast at the end of the span, the extrema will occur at the end of the span.
    1 for good buy, -1 for good sell, 0 for none
    '''
    close_values = indicators.close(data)

    res = []

    for i in range(len(data) - span):
        next_close = close_values.iloc[i+span]

        res.append(next_close)

    res = pd.Series(res)
    res.name = f"precise_forecast_{span}"

    return res
