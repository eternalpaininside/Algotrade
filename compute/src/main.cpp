#include "forecast.h"
#include <algorithm>
#include <cmath>
#include <iomanip>
#include <iostream>
#include <optional>
#include <vector>

using std::abs;
using std::cin;
using std::cout;
using std::endl;
using std::vector;

int main() {
  setlocale(LC_ALL, "RU");
  std::ios_base::sync_with_stdio(false);
  cin.tie(nullptr);

  vector<double> prices{
      100, 102, 101, 104, 105, 100, 101, 101, 103, 101.5
  };

  size_t window = 4;
  const double alpha = 0.5;

  const vector<ForecastResult> momentum_results =
    Calculate_Momentum(prices, window);

  const vector<double> sma_predictions =
    Calculate_SMA(prices, window);

  const vector<double> ema_predictions =
    Calculate_EMA(prices, alpha);

  const vector<double> trend_predictions =
    Calculate_linear_trend(prices, window);

  vector<ForecastResult> sma_results;
  vector<ForecastResult> ema_results;
  vector<ForecastResult> trend_results;

  for (size_t i = window; i < prices.size(); ++i) {
    const double sma_price = sma_predictions[i - window];

    sma_results.push_back({
        i,
        sma_price,
        prices[i],
        sma_price - prices[i]
      });
  }

  for (size_t i = 1; i < prices.size(); ++i) {
    const double ema_price = ema_predictions[i - 1];

    ema_results.push_back({
        i,
        ema_price,
        prices[i],
        ema_price - prices[i]
      });
  }

  for (size_t i = window; i < prices.size(); ++i) {
    const double trend_price = trend_predictions[i - window];

    trend_results.push_back({
        i,
        trend_price,
        prices[i],
        trend_price - prices[i]
      });
  }

  cout << "Window = " << window << '\n';
  cout << "Alpha = " << alpha << '\n';
  cout << "------------\n";

  for (const auto& i : prices) {
    cout << i << " ";
  }

  cout << "\n------------\n";

  cout << std::left << std::setw(10) << "Target"
    << std::right << std::setw(10) << "Actual"
    << std::setw(12) << "Momentum"
    << std::setw(10) << "SMA"
    << std::setw(10) << "EMA"
    << std::setw(10) << "Trend"
    << std::setw(17) << "Momentum error"
    << std::setw(12) << "SMA error"
    << std::setw(12) << "EMA error"
    << std::setw(13) << "Trend error" << '\n';

  cout << std::fixed << std::setprecision(2);

  for (size_t i = 0; i < prices.size(); ++i) {
    cout << std::left << std::setw(10) << i
      << std::right << std::setw(10) << prices[i];

    if (i >= window) {
      cout << std::setw(12)
        << momentum_results[i - window].predicted_price;
    } else {
      cout << std::setw(12) << "";
    }

    if (i >= window) {
      cout << std::setw(10)
        << sma_results[i - window].predicted_price;
    } else {
      cout << std::setw(10) << "";
    }

    if (i >= 1) {
      cout << std::setw(10)
        << ema_results[i - 1].predicted_price;
    } else {
      cout << std::setw(10) << "";
    }

    if (i >= window) {
      cout << std::setw(10)
        << trend_results[i - window].predicted_price;
    } else {
      cout << std::setw(10) << "";
    }

    if (i >= window) {
      cout << std::setw(17)
        << momentum_results[i - window].error;
    } else {
      cout << std::setw(17) << "";
    }

    if (i >= window) {
      cout << std::setw(12)
        << sma_results[i - window].error;
    } else {
      cout << std::setw(12) << "";
    }

    if (i >= 1) {
      cout << std::setw(12)
        << ema_results[i - 1].error;
    } else {
      cout << std::setw(12) << "";
    }

    if (i >= window) {
      cout << std::setw(13)
        << trend_results[i - window].error;
    } else {
      cout << std::setw(13) << "";
    }

    cout << '\n';
  }

  const auto momentum_mae = Calculate_MAE(momentum_results, window);
  const auto sma_mae = Calculate_MAE(sma_results, window);
  const auto ema_mae = Calculate_MAE(ema_results, window);
  const auto trend_mae = Calculate_MAE(trend_results, window);

  if (momentum_mae.has_value()) {
    cout << "MAE Momentum = " << *momentum_mae << '\n';
  } else {
    cout << "MAE Momentum посчитать невозможно\n";
  }

  if (sma_mae.has_value()) {
    cout << "MAE SMA = " << *sma_mae << '\n';
  } else {
    cout << "MAE SMA посчитать невозможно\n";
  }

  if (ema_mae.has_value()) {
    cout << "MAE EMA = " << *ema_mae << '\n';
  } else {
    cout << "MAE EMA посчитать невозможно\n";
  }

  if (trend_mae.has_value()) {
    cout << "MAE Trend = " << *trend_mae << '\n';
  } else {
    cout << "MAE Trend посчитать невозможно\n";
  }

  cout << "\n------------\n";
}