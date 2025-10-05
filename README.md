# Hype Detector
This repository contains the implementation for the research paper "A Hype Detector based on LSTM, Stock Prices from Market Sentiment and Technical Indicators". The project develops a neural network model that integrates market sentiment from Indian financial news sources with historical stock prices and technical indicators to improve stock price prediction accuracy.

The model addresses limitations in traditional time series approaches (e.g., LSTM, ARIMA, GRU) by incorporating sentiment analysis from news articles. Data is sourced from Indian financial news websites (Live Mint, Economic Times, Money Control) and historical stock data from Yahoo Finance. Key features include:
- LSTM layers for processing stock price sequences.
- Dense layers for technical indicators (e.g., engulfing candles, Marubozu candles).
- Embedding layers with LSTM for sentiment analysis using VADER.
- Achieved a scaled Root Mean Squared Error (RMSE) of 0.1010 (unscaled RMSE of 1.36 INR) on an average stock price of 243.73 INR, outperforming baselines.

## Components
The repository is structured around three main components:
 1. [News Data Collector](https://github.com/papaya147/hype-detector/tree/main/data-collector/news):
-   A Go-based script for web scraping financial news articles
-   Scrapes real-time market-related articles from Live Mint, Economic Times, and Money Control
-   Stores articles in a centralized GitHub repository for preprocessing and sentiment analysis
-   Dependencies: Managed via [`go.mod`](https://github.com/papaya147/hype-detector/blob/main/data-collector/news/go.mod)
 2. [Stock Prices Data Collector](https://github.com/papaya147/hype-detector/tree/main/data-collector/stock-prices)
-   A Python-based script for fetching minute-wise stock data
-   Pulls data (open, close, high, low, volume) from Yahoo Finance using stock symbols listed in [`stock-list.txt`](https://github.com/papaya147/hype-detector/blob/main/stock%20list.txt)
-   Enables high-frequency analysis to capture real-time impacts of news sentiment on stock fluctuations
-   Dependencies: Listed in [`requirements.txt`](https://github.com/papaya147/hype-detector/blob/main/data-collector/stock-prices/requirements.txt).
 3. Prediction Model
 -   The core component that uses collected news and stock data to predict stock prices following news article publications
-   Preprocesses data with z-score normalization, generates sequences, computes technical indicators, and performs sentiment analysis (using VADER)
-   Integrates inputs into a unified neural network for forecasting
-   Focuses on detecting "hype" through sentiment-driven price movements in chaotic financial markets

## Installation
### Prerequisites
- Go (for news collection)
- Conda (for environment management)

### Steps
1. Clone the repository
```bash
git clone https://github.com/papaya147/hype-detector.git
cd hype-detector
conda create --prefix ./env python=3.11
conda activate ./env
```
2. News collection (optional)
```bash
cd data-collector/news
go run .
```
3. Stock price collection (optional)
```bash
cd data-collector/stock-prices
conda create --prefix ./env python=3.11
conda activate ./env
pip install -r requirements.txt
python main.py
```
