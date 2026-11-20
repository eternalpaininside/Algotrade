#include "forecast.h"

vector<double> Calculate_SMA(const vector<double>& prices, const size_t window) {
  if (window == 0 || prices.size() <= window) {
    return vector<double>{};
  }

  vector<double> sma;

  for (size_t i = window; i < prices.size(); ++i) {
    double sum = 0;
    for (size_t j = 1; j <= window; ++j) {
      sum += prices[i - j];
    }
    sma.push_back(sum / window);
  }
  return sma;
}

std::optional<double> Calculate_MAE(const vector<ForecastResult>& results,
                                    size_t first_target_index) {
  double total_error = 0;
  size_t count = 0;

  for (size_t i = 0; i < results.size(); ++i) {
    if (results[i].target_index < first_target_index) {
      continue;
    }

    total_error += std::abs(results[i].error);
    ++count;
  }

  if (count == 0) {
    return std::nullopt;
  }

  double mae = total_error / static_cast<double>(count);
  return mae;
}

vector<ForecastResult> Calculate_Momentum(const vector<double>& prices, const size_t window) {
  if (window < 2 || prices.size() <= window) {
    return vector<ForecastResult>{};
  }

  vector<ForecastResult> results;

  for (size_t i = window; i < prices.size(); ++i) {
    double change = 0;
    for (size_t j = i - window + 1; j < i; ++j) {
      change += prices[j] - prices[j - 1];
    }
    results.push_back({
        i,
        prices[i - 1] + change / (window - 1),
        prices[i],
        prices[i - 1] + change / (window - 1) - prices[i]
      });
  }
  return results;
}

vector<double> Calculate_EMA(const vector<double>& prices, const double alpha) {
  if (prices.empty() || alpha <= 0 || alpha > 1) {
    return vector<double>{};
  }
  vector<double> ema{ prices[0] };
  for (size_t i = 1; i < prices.size(); ++i) {
    double forecast = ema[i - 1];
    double error = forecast - prices[i];
    ema.push_back(alpha * prices[i] + (1 - alpha) * ema[i - 1]);
  }
  return ema;
}

vector<double> Calculate_linear_trend(const vector<double>& prices, const size_t window) {
  if (window < 2 || prices.size() <= window) {
    return vector<double>{};
  }

  vector<double> linear_trend;

  for (size_t i = window; i < prices.size(); ++i) {
    double meanX = 0;
    double meanY = 0;

    for (size_t j = 1; j <= window; ++j) {
      meanY += prices[i - j];
      meanX += i - j;
    }
    meanY /= window;
    meanX /= window;

    double numerator = 0;
    double denominator = 0;

    for (size_t j = 1; j <= window; ++j) {
      const double x = static_cast<double>(i - j);
      const double y = prices[i - j];

      numerator += (x - meanX) * (y - meanY);
      denominator += (x - meanX) * (x - meanX);
    }

    const double k = numerator / denominator;
    const double b = meanY - k * meanX;

    linear_trend.push_back(k * static_cast<double>(i) + b);
  }
  return linear_trend;
}